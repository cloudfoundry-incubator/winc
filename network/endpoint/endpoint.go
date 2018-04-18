package endpoint

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

	"code.cloudfoundry.org/winc/network"
	"code.cloudfoundry.org/winc/network/firewall"
	"code.cloudfoundry.org/winc/network/netinterface"
	"github.com/Microsoft/hcsshim"
	"github.com/sirupsen/logrus"
)

//go:generate counterfeiter -o fakes/hcs_client.go --fake-name HCSClient . HCSClient
type HCSClient interface {
	GetHNSNetworkByName(string) (*hcsshim.HNSNetwork, error)
	CreateEndpoint(*hcsshim.HNSEndpoint) (*hcsshim.HNSEndpoint, error)
	UpdateEndpoint(*hcsshim.HNSEndpoint) (*hcsshim.HNSEndpoint, error)
	GetHNSEndpointByID(string) (*hcsshim.HNSEndpoint, error)
	GetHNSEndpointByName(string) (*hcsshim.HNSEndpoint, error)
	DeleteEndpoint(*hcsshim.HNSEndpoint) (*hcsshim.HNSEndpoint, error)
	HotAttachEndpoint(containerID string, endpointID string, endpointReady func() (bool, error)) error
	HotDetachEndpoint(containerID string, endpointID string) error
}

//go:generate counterfeiter -o fakes/firewall.go --fake-name Firewall . Firewall
type Firewall interface {
	DeleteRule(string) error
	RuleExists(string) (bool, error)
}

type EndpointManager struct {
	hcsClient   HCSClient
	firewall    Firewall
	containerId string
	config      network.Config
}

func NewEndpointManager(hcsClient HCSClient, firewall Firewall, containerId string, config network.Config) *EndpointManager {
	return &EndpointManager{
		hcsClient:   hcsClient,
		firewall:    firewall,
		containerId: containerId,
		config:      config,
	}
}

func (e *EndpointManager) Create() (hcsshim.HNSEndpoint, error) {
	network, err := e.hcsClient.GetHNSNetworkByName(e.config.NetworkName)
	if err != nil {
		return hcsshim.HNSEndpoint{}, err
	}

	providerAddress := network.ManagementIP
	paPolicy, err := json.Marshal(hcsshim.PaPolicy{
		Type: hcsshim.PA,
		PA:   providerAddress,
	})

	if err != nil {
		return hcsshim.HNSEndpoint{}, err
	}

	natPolicy, err := json.Marshal(hcsshim.OutboundNatPolicy{
		Policy: hcsshim.Policy{Type: hcsshim.OutboundNat},
	})

	if err != nil {
		return hcsshim.HNSEndpoint{}, err
	}

	policies := []json.RawMessage{paPolicy, natPolicy}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	endpointIP := strings.TrimSuffix(network.Subnets[0].GatewayAddress, ".1") + fmt.Sprintf(".%d", r.Intn(256))

	endpoint := &hcsshim.HNSEndpoint{
		VirtualNetwork: network.Id,
		Name:           e.containerId,
		Policies:       policies,
		IPAddress:      net.ParseIP(endpointIP),
		GatewayAddress: network.Subnets[0].GatewayAddress,
	}

	if e.config.MaximumOutgoingBandwidth != 0 {
		policy, err := json.Marshal(hcsshim.QosPolicy{
			Type: hcsshim.QOS,
			MaximumOutgoingBandwidthInBytes: uint64(e.config.MaximumOutgoingBandwidth),
		})
		if err != nil {
			return hcsshim.HNSEndpoint{}, err
		}

		endpoint.Policies = append(endpoint.Policies, policy)
	}

	if len(e.config.DNSServers) > 0 {
		endpoint.DNSServerList = strings.Join(e.config.DNSServers, ",")
	}

	createdEndpoint, err := e.createEndpoint(endpoint)
	if err != nil {
		return hcsshim.HNSEndpoint{}, err
	}

	attachedEndpoint, err := e.attachEndpoint(createdEndpoint)
	if err != nil {
		if _, err := e.hcsClient.DeleteEndpoint(createdEndpoint); err != nil {
			logrus.Error(fmt.Sprintf("Error deleting endpoint %s: %s", endpoint.Id, err.Error()))
		}

		return hcsshim.HNSEndpoint{}, err
	}

	return *attachedEndpoint, nil
}

func (e *EndpointManager) attachEndpoint(endpoint *hcsshim.HNSEndpoint) (*hcsshim.HNSEndpoint, error) {
	endpointReady := func() (bool, error) {
		interfaceAlias := fmt.Sprintf("vEthernet (%s)", e.containerId)
		return netinterface.InterfaceExists(interfaceAlias)
	}

	if err := e.hcsClient.HotAttachEndpoint(e.containerId, endpoint.Id, endpointReady); err != nil {
		return nil, err
	}

	allocatedEndpoint, err := e.hcsClient.GetHNSEndpointByID(endpoint.Id)
	if err != nil {
		logrus.Error(fmt.Sprintf("Unable to load updated endpoint %s: %s", endpoint.Id, err.Error()))
		return nil, err
	}

	var compartmentId uint32
	var endpointPortGuid string

	for _, a := range allocatedEndpoint.Resources.Allocators {
		if a.Type == hcsshim.EndpointPortType {
			compartmentId = a.CompartmentId
			endpointPortGuid = a.EndpointPortGuid
			break
		}
	}

	if compartmentId == 0 || endpointPortGuid == "" {
		return nil, fmt.Errorf("invalid endpoint %s allocators: %+v", endpoint.Id, allocatedEndpoint.Resources.Allocators)
	}

	//	ruleName := fmt.Sprintf("Compartment %d - %s", compartmentId, endpointPortGuid)
	//	if err := e.deleteFirewallRule(ruleName); err != nil {
	//		return nil, err
	//	}

	return allocatedEndpoint, nil

}

func (e *EndpointManager) deleteFirewallRule(ruleName string) error {
	ruleCreated := false

	for i := 0; i < 3; i++ {
		var err error
		time.Sleep(time.Millisecond * 200 * time.Duration(i))

		ruleCreated, err = e.firewall.RuleExists(ruleName)
		if err != nil {
			return err
		}

		if ruleCreated {
			if err := e.firewall.DeleteRule(ruleName); err != nil {
				logrus.Error(fmt.Sprintf("Unable to delete generated firewall rule %s: %s", ruleName, err.Error()))
				return err
			}

			break
		}
	}

	if !ruleCreated {
		return fmt.Errorf("firewall rule %s not generated in time", ruleName)
	}

	return nil
}

func (e *EndpointManager) ApplyMappings(endpoint hcsshim.HNSEndpoint, mappings []hcsshim.NatPolicy, acls []hcsshim.ACLPolicy) (hcsshim.HNSEndpoint, error) {
	var policies []json.RawMessage

	for _, mapping := range mappings {
		policy, err := json.Marshal(mapping)
		if err != nil {
			return hcsshim.HNSEndpoint{}, err
		}
		policies = append(policies, policy)
	}

	if len(acls) == 0 {
		// make sure everything's blocked if no netout rules present
		acls = []hcsshim.ACLPolicy{
			{
				Type:      hcsshim.ACL,
				RuleType:  hcsshim.Switch,
				Action:    hcsshim.Block,
				Direction: hcsshim.Out,
				Protocol:  uint16(firewall.NET_FW_IP_PROTOCOL_ANY),
			},
			{
				Type:      hcsshim.ACL,
				RuleType:  hcsshim.Switch,
				Action:    hcsshim.Block,
				Direction: hcsshim.In,
				Protocol:  uint16(firewall.NET_FW_IP_PROTOCOL_ANY),
			},
		}
	}

	// these rules are necessary to get throught the host firewall
	// Priority: 100 is some sort of magic number that makes this work...
	acls = append(acls, hcsshim.ACLPolicy{
		Type:      hcsshim.ACL,
		RuleType:  hcsshim.Host,
		Action:    hcsshim.Allow,
		Direction: hcsshim.In,
		Priority:  100,
		Protocol:  uint16(firewall.NET_FW_IP_PROTOCOL_ANY),
	},
		hcsshim.ACLPolicy{
			Type:      hcsshim.ACL,
			RuleType:  hcsshim.Host,
			Action:    hcsshim.Allow,
			Direction: hcsshim.Out,
			Priority:  100,
			Protocol:  uint16(firewall.NET_FW_IP_PROTOCOL_ANY),
		},
	)

	for _, acl := range acls {
		policy, err := json.Marshal(acl)
		if err != nil {
			return hcsshim.HNSEndpoint{}, err
		}
		policies = append(policies, policy)
	}

	endpoint.Policies = append(endpoint.Policies, policies...)
	updatedEndpoint, err := e.hcsClient.UpdateEndpoint(&endpoint)
	if err != nil {
		fmt.Println("are we really here?")
		return hcsshim.HNSEndpoint{}, err
	}

	//id := updatedEndpoint.Id
	//var natAllocated bool
	//var allocatedEndpoint *hcsshim.HNSEndpoint

	//for i := 0; i < 10; i++ {
	//	natAllocated = false
	//	allocatedEndpoint, err = e.hcsClient.GetHNSEndpointByID(id)

	//	for _, a := range allocatedEndpoint.Resources.Allocators {
	//		if a.Type == hcsshim.NATPolicyType {
	//			natAllocated = true
	//			break
	//		}
	//	}

	//	if natAllocated {
	//		break
	//	}

	//	time.Sleep(200 * time.Millisecond)
	//}

	//if !natAllocated {
	//	return hcsshim.HNSEndpoint{}, errors.New("NAT not initialized in time")
	//}

	//return *allocatedEndpoint, nil
	return *updatedEndpoint, nil
}

func (e *EndpointManager) Delete() error {
	endpoint, err := e.hcsClient.GetHNSEndpointByName(e.containerId)
	if err != nil {
		if _, ok := err.(hcsshim.EndpointNotFoundError); ok {
			return nil
		}

		return err
	}

	var detachErr error
	err = e.hcsClient.HotDetachEndpoint(e.containerId, endpoint.Id)
	if err != hcsshim.ErrComputeSystemDoesNotExist {
		detachErr = err
	}

	_, deleteErr := e.hcsClient.DeleteEndpoint(endpoint)

	if detachErr != nil && deleteErr != nil {
		return fmt.Errorf("%s, %s", detachErr.Error(), deleteErr.Error())
	}

	if detachErr != nil {
		return detachErr
	}

	if deleteErr != nil {
		return deleteErr
	}

	return nil
}

func (e *EndpointManager) createEndpoint(endpoint *hcsshim.HNSEndpoint) (*hcsshim.HNSEndpoint, error) {
	var createErr error
	var createdEndpoint *hcsshim.HNSEndpoint
	for i := 0; i < 3 && createdEndpoint == nil; i++ {
		createdEndpoint, createErr = e.hcsClient.CreateEndpoint(endpoint)
		if createErr != nil {
			if createErr.Error() != "HNS failed with error : Unspecified error" {
				return nil, createErr
			}
		}
	}

	if createdEndpoint == nil {
		return nil, createErr
	}

	return createdEndpoint, nil
}
