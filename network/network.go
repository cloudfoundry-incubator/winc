package network

import (
	"encoding/json"
	"fmt"
	"net"

	"code.cloudfoundry.org/winc/network/netinterface"
	"code.cloudfoundry.org/winc/network/netrules"

	"github.com/Microsoft/hcsshim"
)

//go:generate counterfeiter -o fakes/net_rule_applier.go --fake-name NetRuleApplier . NetRuleApplier
type NetRuleApplier interface {
	In(netrules.NetIn, string) (hcsshim.NatPolicy, hcsshim.ACLPolicy, error)
	Out(netrules.NetOut, string) (hcsshim.ACLPolicy, error)
	NatMTU(int) error
	ContainerMTU(int) error
	Cleanup() error
}

//go:generate counterfeiter -o fakes/endpoint_manager.go --fake-name EndpointManager . EndpointManager
type EndpointManager interface {
	Create() (hcsshim.HNSEndpoint, error)
	Delete() error
	ApplyPolicies(hcsshim.HNSEndpoint, []hcsshim.NatPolicy, []hcsshim.ACLPolicy) (hcsshim.HNSEndpoint, error)
}

//go:generate counterfeiter -o fakes/hcs_client.go --fake-name HCSClient . HCSClient
type HCSClient interface {
	GetHNSNetworkByName(string) (*hcsshim.HNSNetwork, error)
	CreateNetwork(*hcsshim.HNSNetwork, func() (bool, error)) (*hcsshim.HNSNetwork, error)
	DeleteNetwork(*hcsshim.HNSNetwork) (*hcsshim.HNSNetwork, error)
}

type Config struct {
	MTU                      int      `json:"mtu"`
	NetworkName              string   `json:"network_name"`
	SubnetRange              string   `json:"subnet_range"`
	GatewayAddress           string   `json:"gateway_address"`
	DNSServers               []string `json:"dns_servers"`
	MaximumOutgoingBandwidth uint64   `json:"maximum_outgoing_bandwidth"`
}

type UpInputs struct {
	Pid        int
	Properties map[string]interface{}
	NetOut     []netrules.NetOut `json:"netout_rules"`
	NetIn      []netrules.NetIn  `json:"netin"`
}

type UpOutputs struct {
	Properties struct {
		ContainerIP      string `json:"garden.network.container-ip"`
		DeprecatedHostIP string `json:"garden.network.host-ip"`
		MappedPorts      string `json:"garden.network.mapped-ports"`
	} `json:"properties"`
	DNSServers []string `json:"dns_servers,omitempty"`
}

type NetworkManager struct {
	hcsClient       HCSClient
	applier         NetRuleApplier
	endpointManager EndpointManager
	containerId     string
	config          Config
}

func NewNetworkManager(client HCSClient, applier NetRuleApplier, endpointManager EndpointManager, containerId string, config Config) *NetworkManager {
	return &NetworkManager{
		hcsClient:       client,
		applier:         applier,
		endpointManager: endpointManager,
		containerId:     containerId,
		config:          config,
	}
}

func (n *NetworkManager) CreateHostNATNetwork() error {
	existingNetwork, err := n.hcsClient.GetHNSNetworkByName(n.config.NetworkName)
	if err != nil {
		if _, isNotExist := err.(hcsshim.NetworkNotFoundError); !isNotExist {
			return err
		}
	}

	subnets := []hcsshim.Subnet{{AddressPrefix: n.config.SubnetRange, GatewayAddress: n.config.GatewayAddress}}

	if existingNetwork != nil {
		if len(existingNetwork.Subnets) == 1 && subnetsMatch(existingNetwork.Subnets[0], subnets[0]) {
			return nil
		}

		return &SameNATNetworkNameError{Name: n.config.NetworkName, Subnets: existingNetwork.Subnets}
	}

	network := &hcsshim.HNSNetwork{
		Name:    n.config.NetworkName,
		Type:    "nat",
		Subnets: subnets,
	}

	networkReady := func() (bool, error) {
		interfaceAlias := fmt.Sprintf("vEthernet (%s)", network.Name)
		return netinterface.InterfaceExists(interfaceAlias)
	}

	_, err = n.hcsClient.CreateNetwork(network, networkReady)
	if err != nil {
		return err
	}

	return n.applier.NatMTU(n.config.MTU)
}

func subnetsMatch(a, b hcsshim.Subnet) bool {
	return (a.AddressPrefix == b.AddressPrefix) && (a.GatewayAddress == b.GatewayAddress)
}

func (n *NetworkManager) DeleteHostNATNetwork() error {
	network, err := n.hcsClient.GetHNSNetworkByName(n.config.NetworkName)
	if err != nil {
		if _, ok := err.(hcsshim.NetworkNotFoundError); ok {
			return nil
		}

		return err
	}
	_, err = n.hcsClient.DeleteNetwork(network)
	return err
}

func (n *NetworkManager) Up(inputs UpInputs) (UpOutputs, error) {
	outputs, err := n.up(inputs)
	if err != nil {
		n.applier.Cleanup()
		n.endpointManager.Delete()
	}
	return outputs, err
}

func (n *NetworkManager) up(inputs UpInputs) (UpOutputs, error) {
	outputs := UpOutputs{}

	createdEndpoint, err := n.endpointManager.Create()
	if err != nil {
		return outputs, err
	}

	mappedPorts := []netrules.PortMapping{}
	hnsAcls := []hcsshim.ACLPolicy{}
	hnsNats := []hcsshim.NatPolicy{}

	for _, rule := range inputs.NetIn {
		nat, acl, err := n.applier.In(rule, createdEndpoint.IPAddress.String())
		if err != nil {
			return outputs, err
		}

		mappedPorts = append(mappedPorts, netrules.PortMapping{
			ContainerPort: uint32(nat.InternalPort),
			HostPort:      uint32(nat.ExternalPort),
		})
		hnsAcls = append(hnsAcls, acl)
		hnsNats = append(hnsNats, nat)
	}

	for _, dnsServer := range n.config.DNSServers {
		serverIP := net.ParseIP(dnsServer)
		inputs.NetOut = append(inputs.NetOut,
			netrules.NetOut{
				Protocol: netrules.ProtocolTCP,
				Networks: []netrules.IPRange{{Start: serverIP, End: serverIP}},
				Ports:    []netrules.PortRange{{Start: 53, End: 53}},
			},
			netrules.NetOut{
				Protocol: netrules.ProtocolUDP,
				Networks: []netrules.IPRange{{Start: serverIP, End: serverIP}},
				Ports:    []netrules.PortRange{{Start: 53, End: 53}},
			},
		)
	}

	for _, rule := range inputs.NetOut {
		acl, err := n.applier.Out(rule, createdEndpoint.IPAddress.String())
		if err != nil {
			return outputs, err
		}

		hnsAcls = append(hnsAcls, acl)
	}

	if _, err := n.endpointManager.ApplyPolicies(createdEndpoint, hnsNats, hnsAcls); err != nil {
		return outputs, err
	}

	if err := n.applier.ContainerMTU(n.config.MTU); err != nil {
		return outputs, err
	}

	portBytes, err := json.Marshal(mappedPorts)
	if err != nil {
		return outputs, err
	}

	outputs.Properties.MappedPorts = string(portBytes)
	outputs.Properties.ContainerIP = createdEndpoint.IPAddress.String()
	outputs.Properties.DeprecatedHostIP = "255.255.255.255"

	return outputs, nil
}

func (n *NetworkManager) Down() error {
	deleteErr := n.endpointManager.Delete()
	cleanupErr := n.applier.Cleanup()

	if deleteErr != nil && cleanupErr != nil {
		return fmt.Errorf("%s, %s", deleteErr.Error(), cleanupErr.Error())
	}

	if deleteErr != nil {
		return deleteErr
	}
	if cleanupErr != nil {
		return cleanupErr
	}
	return nil
}
