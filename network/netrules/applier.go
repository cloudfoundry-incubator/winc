package netrules

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"code.cloudfoundry.org/localip"
	"code.cloudfoundry.org/winc/network/firewall"
	"github.com/Microsoft/hcsshim"
)

//go:generate counterfeiter -o fakes/netsh_runner.go --fake-name NetShRunner . NetShRunner
type NetShRunner interface {
	RunContainer([]string) error
}

//go:generate counterfeiter -o fakes/port_allocator.go --fake-name PortAllocator . PortAllocator
type PortAllocator interface {
	AllocatePort(handle string, port int) (int, error)
	ReleaseAllPorts(handle string) error
}

//go:generate counterfeiter -o fakes/netinterface.go --fake-name NetInterface . NetInterface
type NetInterface interface {
	ByName(string) (*net.Interface, error)
	ByIP(string) (*net.Interface, error)
	SetMTU(string, int) error
}

//go:generate counterfeiter -o fakes/firewall.go --fake-name Firewall . Firewall
type Firewall interface {
	CreateRule(firewall.Rule) error
	DeleteRule(string) error
	RuleExists(string) (bool, error)
}

type Applier struct {
	netSh         NetShRunner
	containerId   string
	networkName   string
	portAllocator PortAllocator
	netInterface  NetInterface
	firewall      Firewall
}

func NewApplier(netSh NetShRunner, containerId string, networkName string, portAllocator PortAllocator, netInterface NetInterface, firewall Firewall) *Applier {
	return &Applier{
		netSh:         netSh,
		containerId:   containerId,
		networkName:   networkName,
		portAllocator: portAllocator,
		netInterface:  netInterface,
		firewall:      firewall,
	}
}

func (a *Applier) In(rule NetIn, containerIP string) (hcsshim.NatPolicy, hcsshim.ACLPolicy, error) {
	externalPort := rule.HostPort

	if externalPort == 0 {
		allocatedPort, err := a.portAllocator.AllocatePort(a.containerId, 0)
		if err != nil {
			return hcsshim.NatPolicy{}, hcsshim.ACLPolicy{}, err
		}
		externalPort = uint32(allocatedPort)
	}

	if err := a.netInModifyHostVM(rule, containerIP); err != nil {
		return hcsshim.NatPolicy{}, hcsshim.ACLPolicy{}, err
	}

	if err := a.openPort(rule.ContainerPort); err != nil {
		return hcsshim.NatPolicy{}, hcsshim.ACLPolicy{}, err
	}

	return hcsshim.NatPolicy{
			Type:         hcsshim.Nat,
			Protocol:     "TCP",
			ExternalPort: uint16(externalPort),
			InternalPort: uint16(rule.ContainerPort),
		}, hcsshim.ACLPolicy{
			Type:           hcsshim.ACL,
			Action:         hcsshim.Allow,
			Direction:      hcsshim.In,
			Protocol:       uint16(firewall.NET_FW_IP_PROTOCOL_TCP),
			LocalAddresses: containerIP,
			LocalPort:      strconv.FormatUint(uint64(rule.ContainerPort), 10),
		}, nil
}

func (a *Applier) netInModifyHostVM(rule NetIn, containerIP string) error {
	fr := firewall.Rule{
		Name:           a.containerId,
		Action:         firewall.NET_FW_ACTION_ALLOW,
		Direction:      firewall.NET_FW_RULE_DIR_IN,
		Protocol:       firewall.NET_FW_IP_PROTOCOL_TCP,
		LocalAddresses: containerIP,
		LocalPorts:     strconv.FormatUint(uint64(rule.ContainerPort), 10),
	}

	return a.firewall.CreateRule(fr)
}

func (a *Applier) Out(rule NetOut, containerIP string) (hcsshim.ACLPolicy, error) {
	lAddrs := []string{}

	for _, ipr := range rule.Networks {
		lAddrs = append(lAddrs, IPRangeToCIDRs(ipr)...)
	}

	acl := hcsshim.ACLPolicy{
		Type:            hcsshim.ACL,
		Action:          hcsshim.Allow,
		Direction:       hcsshim.Out,
		LocalAddresses:  containerIP,
		RemoteAddresses: strings.Join(lAddrs, ","),
	}

	switch rule.Protocol {
	case ProtocolTCP:
		acl.RemotePort = firewallRulePortRange(rule.Ports)
		acl.Protocol = uint16(firewall.NET_FW_IP_PROTOCOL_TCP)
	case ProtocolUDP:
		acl.RemotePort = firewallRulePortRange(rule.Ports)
		acl.Protocol = uint16(firewall.NET_FW_IP_PROTOCOL_UDP)
	case ProtocolICMP:
		acl.Protocol = uint16(firewall.NET_FW_IP_PROTOCOL_ICMP)
	case ProtocolAll:
		acl.Protocol = uint16(firewall.NET_FW_IP_PROTOCOL_ANY)
	default:
		return hcsshim.ACLPolicy{}, fmt.Errorf("invalid protocol: %d", rule.Protocol)
	}

	if err := a.netOutModifyHostVM(rule, containerIP); err != nil {
		return hcsshim.ACLPolicy{}, err
	}

	return acl, nil
}

func (a *Applier) netOutModifyHostVM(rule NetOut, containerIP string) error {
	fr := firewall.Rule{
		Name:            a.containerId,
		Action:          firewall.NET_FW_ACTION_ALLOW,
		Direction:       firewall.NET_FW_RULE_DIR_OUT,
		LocalAddresses:  containerIP,
		RemoteAddresses: firewallRuleIPRange(rule.Networks),
	}

	switch rule.Protocol {
	case ProtocolTCP:
		fr.RemotePorts = firewallRulePortRange(rule.Ports)
		fr.Protocol = firewall.NET_FW_IP_PROTOCOL_TCP
	case ProtocolUDP:
		fr.RemotePorts = firewallRulePortRange(rule.Ports)
		fr.Protocol = firewall.NET_FW_IP_PROTOCOL_UDP
	case ProtocolICMP:
		fr.Protocol = firewall.NET_FW_IP_PROTOCOL_ICMP
	case ProtocolAll:
		fr.Protocol = firewall.NET_FW_IP_PROTOCOL_ANY
	default:
		return fmt.Errorf("invalid protocol: %d", rule.Protocol)
	}

	return a.firewall.CreateRule(fr)
}

func (a *Applier) ContainerMTU(mtu int) error {
	if mtu == 0 {
		iface, err := a.netInterface.ByName(fmt.Sprintf("vEthernet (%s)", a.networkName))
		if err != nil {
			return err
		}
		mtu = iface.MTU
	}

	interfaceAlias := fmt.Sprintf("vEthernet (%s)", a.containerId)
	return a.netInterface.SetMTU(interfaceAlias, mtu)
}

func (a *Applier) NatMTU(mtu int) error {
	if mtu == 0 {
		hostIP, err := localip.LocalIP()
		if err != nil {
			return err
		}
		iface, err := a.netInterface.ByIP(hostIP)
		if err != nil {
			return err
		}
		mtu = iface.MTU
	}

	interfaceId := fmt.Sprintf("vEthernet (%s)", a.networkName)
	return a.netInterface.SetMTU(interfaceId, mtu)
}

func (a *Applier) openPort(port uint32) error {
	args := []string{"http", "add", "urlacl", fmt.Sprintf("url=http://*:%d/", port), "user=Users"}
	return a.netSh.RunContainer(args)
}

func (a *Applier) Cleanup() error {
	portReleaseErr := a.portAllocator.ReleaseAllPorts(a.containerId)

	// we can just delete the rule here since it will succeed
	// if the rule does not exist
	deleteErr := a.firewall.DeleteRule(a.containerId)

	if portReleaseErr != nil && deleteErr != nil {
		return fmt.Errorf("%s, %s", portReleaseErr.Error(), deleteErr.Error())
	}
	if portReleaseErr != nil {
		return portReleaseErr
	}
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}
