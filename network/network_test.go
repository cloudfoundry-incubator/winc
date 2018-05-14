package network_test

import (
	"errors"
	"io/ioutil"
	"net"

	"code.cloudfoundry.org/winc/network"
	"code.cloudfoundry.org/winc/network/fakes"
	"code.cloudfoundry.org/winc/network/netrules"
	"github.com/Microsoft/hcsshim"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("NetworkManager", func() {
	const containerId = "some-container-id"

	var (
		networkManager  *network.NetworkManager
		netRuleApplier  *fakes.NetRuleApplier
		hcsClient       *fakes.HCSClient
		endpointManager *fakes.EndpointManager
		hnsNetwork      *hcsshim.HNSNetwork
	)

	BeforeEach(func() {
		hcsClient = &fakes.HCSClient{}
		netRuleApplier = &fakes.NetRuleApplier{}
		endpointManager = &fakes.EndpointManager{}
		config := network.Config{
			MTU:            1434,
			SubnetRange:    "123.45.0.0/67",
			GatewayAddress: "123.45.0.1",
			NetworkName:    "unit-test-name",
		}

		networkManager = network.NewNetworkManager(hcsClient, netRuleApplier, endpointManager, containerId, config)

		logrus.SetOutput(ioutil.Discard)
	})

	Describe("CreateHostNATNetwork", func() {
		BeforeEach(func() {
			hcsClient.GetHNSNetworkByNameReturns(nil, hcsshim.NetworkNotFoundError{NetworkName: "unit-test-name"})
		})

		It("creates the network with the correct values", func() {
			Expect(networkManager.CreateHostNATNetwork()).To(Succeed())

			Expect(hcsClient.GetHNSNetworkByNameCallCount()).To(Equal(1))
			Expect(hcsClient.GetHNSNetworkByNameArgsForCall(0)).To(Equal("unit-test-name"))

			Expect(hcsClient.CreateNetworkCallCount()).To(Equal(1))
			net, _ := hcsClient.CreateNetworkArgsForCall(0)
			Expect(net.Name).To(Equal("unit-test-name"))
			Expect(net.Subnets).To(ConsistOf(hcsshim.Subnet{AddressPrefix: "123.45.0.0/67", GatewayAddress: "123.45.0.1"}))

			Expect(netRuleApplier.NatMTUCallCount()).To(Equal(1))
			Expect(netRuleApplier.NatMTUArgsForCall(0)).To(Equal(1434))
		})

		Context("the network already exists with the correct values", func() {
			BeforeEach(func() {
				hnsNetwork = &hcsshim.HNSNetwork{
					Name:    "unit-test-name",
					Subnets: []hcsshim.Subnet{{AddressPrefix: "123.45.0.0/67", GatewayAddress: "123.45.0.1"}},
				}
				hcsClient.GetHNSNetworkByNameReturns(hnsNetwork, nil)
			})

			It("does not create the network", func() {
				Expect(networkManager.CreateHostNATNetwork()).To(Succeed())

				Expect(hcsClient.GetHNSNetworkByNameCallCount()).To(Equal(1))
				Expect(hcsClient.GetHNSNetworkByNameArgsForCall(0)).To(Equal("unit-test-name"))
				Expect(hcsClient.CreateNetworkCallCount()).To(Equal(0))
			})
		})

		Context("the network already exists with an incorrect address prefix", func() {
			BeforeEach(func() {
				hnsNetwork = &hcsshim.HNSNetwork{
					Name:    "unit-test-name",
					Subnets: []hcsshim.Subnet{{AddressPrefix: "123.89.0.0/67", GatewayAddress: "123.45.0.1"}},
				}
				hcsClient.GetHNSNetworkByNameReturns(hnsNetwork, nil)
			})

			It("returns an error", func() {
				err := networkManager.CreateHostNATNetwork()
				Expect(err).To(BeAssignableToTypeOf(&network.SameNATNetworkNameError{}))
			})
		})

		Context("the network already exists with an incorrect gateway address", func() {
			BeforeEach(func() {
				hnsNetwork = &hcsshim.HNSNetwork{
					Name:    "unit-test-name",
					Subnets: []hcsshim.Subnet{{AddressPrefix: "123.45.0.0/67", GatewayAddress: "123.45.67.89"}},
				}
				hcsClient.GetHNSNetworkByNameReturns(hnsNetwork, nil)
			})

			It("returns an error", func() {
				err := networkManager.CreateHostNATNetwork()
				Expect(err).To(BeAssignableToTypeOf(&network.SameNATNetworkNameError{}))
			})
		})

		Context("GetHNSNetwork returns a non network not found error", func() {
			BeforeEach(func() {
				hcsClient.GetHNSNetworkByNameReturns(nil, errors.New("some HNS error"))
			})

			It("returns an error", func() {
				err := networkManager.CreateHostNATNetwork()
				Expect(err).To(HaveOccurred())
			})
		})

		Context("CreateNetwork returns an error", func() {
			BeforeEach(func() {
				hcsClient.CreateNetworkReturns(nil, errors.New("couldn't create HNS network"))
			})

			It("returns an error", func() {
				err := networkManager.CreateHostNATNetwork()
				Expect(err).To(HaveOccurred())
			})
		})

		Context("NatMTU returns an error", func() {
			BeforeEach(func() {
				netRuleApplier.NatMTUReturns(errors.New("couldn't set MTU on NAT network"))
			})

			It("returns an error", func() {
				err := networkManager.CreateHostNATNetwork()
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("DeleteHostNATNetwork", func() {
		BeforeEach(func() {
			hnsNetwork = &hcsshim.HNSNetwork{Name: "unit-test-name"}
			hcsClient.GetHNSNetworkByNameReturnsOnCall(0, hnsNetwork, nil)
		})

		It("deletes the network", func() {
			Expect(networkManager.DeleteHostNATNetwork()).To(Succeed())

			Expect(hcsClient.GetHNSNetworkByNameCallCount()).To(Equal(1))
			Expect(hcsClient.GetHNSNetworkByNameArgsForCall(0)).To(Equal("unit-test-name"))

			Expect(hcsClient.DeleteNetworkCallCount()).To(Equal(1))
			Expect(hcsClient.DeleteNetworkArgsForCall(0)).To(Equal(hnsNetwork))
		})

		Context("the network does not exist", func() {
			BeforeEach(func() {
				hcsClient.GetHNSNetworkByNameReturnsOnCall(0, nil, hcsshim.NetworkNotFoundError{NetworkName: "unit-test-name"})
			})

			It("returns success", func() {
				Expect(networkManager.DeleteHostNATNetwork()).To(Succeed())

				Expect(hcsClient.GetHNSNetworkByNameCallCount()).To(Equal(1))
				Expect(hcsClient.GetHNSNetworkByNameArgsForCall(0)).To(Equal("unit-test-name"))
			})
		})

		Context("GetHNSNetwork returns a non network not found error", func() {
			BeforeEach(func() {
				hcsClient.GetHNSNetworkByNameReturns(nil, errors.New("some HNS error"))
			})

			It("returns an error", func() {
				err := networkManager.CreateHostNATNetwork()
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Up", func() {
		var (
			inputs          network.UpInputs
			createdEndpoint hcsshim.HNSEndpoint
			containerIP     net.IP
			nat1            hcsshim.NatPolicy
			nat2            hcsshim.NatPolicy
			inAcl1          hcsshim.ACLPolicy
			inAcl2          hcsshim.ACLPolicy
			outAcl1         hcsshim.ACLPolicy
			outAcl2         hcsshim.ACLPolicy
		)

		BeforeEach(func() {
			containerIP = net.ParseIP("111.222.33.44")

			createdEndpoint = hcsshim.HNSEndpoint{
				IPAddress: containerIP,
			}

			inputs = network.UpInputs{
				Pid: 1234,
				NetIn: []netrules.NetIn{
					{HostPort: 0, ContainerPort: 666},
					{HostPort: 0, ContainerPort: 888},
				},
				NetOut: []netrules.NetOut{
					{Protocol: 6},
					{Protocol: 17},
				},
			}

			nat1 = hcsshim.NatPolicy{
				Type:         hcsshim.Nat,
				Protocol:     "TCP",
				ExternalPort: 111,
				InternalPort: 666,
			}

			nat2 = hcsshim.NatPolicy{
				Type:         hcsshim.Nat,
				Protocol:     "TCP",
				ExternalPort: 222,
				InternalPort: 888,
			}

			inAcl1 = hcsshim.ACLPolicy{
				Type:      hcsshim.ACL,
				LocalPort: "666",
				Direction: hcsshim.In,
				Action:    hcsshim.Allow,
			}

			inAcl2 = hcsshim.ACLPolicy{
				Type:      hcsshim.ACL,
				LocalPort: "888",
				Direction: hcsshim.In,
				Action:    hcsshim.Allow,
			}

			outAcl1 = hcsshim.ACLPolicy{
				Type:      hcsshim.ACL,
				Direction: hcsshim.Out,
				Action:    hcsshim.Allow,
				Protocol:  6,
			}

			outAcl2 = hcsshim.ACLPolicy{
				Type:      hcsshim.ACL,
				Direction: hcsshim.In,
				Action:    hcsshim.Allow,
				Protocol:  17,
			}

			netRuleApplier.InReturnsOnCall(0, nat1, inAcl1, nil)
			netRuleApplier.InReturnsOnCall(1, nat2, inAcl2, nil)

			netRuleApplier.OutReturnsOnCall(0, outAcl1, nil)
			netRuleApplier.OutReturnsOnCall(1, outAcl2, nil)

			endpointManager.CreateReturns(createdEndpoint, nil)
		})

		It("creates an endpoint, applies net in, applies net out, handles mtu, and returns the up outputs", func() {
			output, err := networkManager.Up(inputs)
			Expect(err).NotTo(HaveOccurred())

			Expect(output.Properties.ContainerIP).To(Equal(containerIP.String()))
			Expect(output.Properties.DeprecatedHostIP).To(Equal("255.255.255.255"))
			Expect(output.Properties.MappedPorts).To(Equal(`[{"HostPort":111,"ContainerPort":666},{"HostPort":222,"ContainerPort":888}]`))

			Expect(endpointManager.CreateCallCount()).To(Equal(1))

			Expect(netRuleApplier.InCallCount()).To(Equal(2))
			inRule, ip := netRuleApplier.InArgsForCall(0)
			Expect(inRule).To(Equal(netrules.NetIn{HostPort: 0, ContainerPort: 666}))
			Expect(ip).To(Equal(containerIP.String()))

			inRule, ip = netRuleApplier.InArgsForCall(1)
			Expect(inRule).To(Equal(netrules.NetIn{HostPort: 0, ContainerPort: 888}))
			Expect(ip).To(Equal(containerIP.String()))

			Expect(netRuleApplier.OutCallCount()).To(Equal(2))
			outRule, ip := netRuleApplier.OutArgsForCall(0)
			Expect(outRule).To(Equal(netrules.NetOut{Protocol: 7}))
			Expect(ip).To(Equal(containerIP.String()))

			outRule, ip = netRuleApplier.OutArgsForCall(1)
			Expect(outRule).To(Equal(netrules.NetOut{Protocol: 8}))
			Expect(ip).To(Equal(containerIP.String()))

			Expect(endpointManager.ApplyPoliciesCallCount()).To(Equal(1))
			ep, nats, acls := endpointManager.ApplyPoliciesArgsForCall(0)
			Expect(ep).To(Equal(createdEndpoint))
			Expect(nats).To(Equal([]hcsshim.NatPolicy{nat1, nat2}))
			Expect(acls).To(Equal([]hcsshim.ACLPolicy{inAcl1, inAcl2, outAcl1, outAcl2}))

			Expect(netRuleApplier.ContainerMTUCallCount()).To(Equal(1))
			mtu := netRuleApplier.ContainerMTUArgsForCall(0)
			Expect(mtu).To(Equal(1434))
		})

		Context("when the config specifies DNS servers", func() {
			BeforeEach(func() {
				config := network.Config{
					DNSServers: []string{"1.1.1.1", "2.2.2.2"},
				}
				networkManager = network.NewNetworkManager(hcsClient, netRuleApplier, endpointManager, containerId, config)
				inputs.NetOut = []netrules.NetOut{}
			})

			It("creates netout rules for the servers", func() {
				_, err := networkManager.Up(inputs)
				Expect(err).NotTo(HaveOccurred())

				dnsServer1 := net.ParseIP("1.1.1.1")
				dnsServer2 := net.ParseIP("2.2.2.2")
				Expect(netRuleApplier.OutCallCount()).To(Equal(4))

				outRule, ip := netRuleApplier.OutArgsForCall(0)
				Expect(outRule).To(Equal(netrules.NetOut{
					Protocol: netrules.ProtocolTCP,
					Networks: []netrules.IPRange{{Start: dnsServer1, End: dnsServer1}},
					Ports:    []netrules.PortRange{{Start: 53, End: 53}},
				}))
				Expect(ip).To(Equal(containerIP.String()))

				outRule, ip = netRuleApplier.OutArgsForCall(1)
				Expect(outRule).To(Equal(netrules.NetOut{
					Protocol: netrules.ProtocolUDP,
					Networks: []netrules.IPRange{{Start: dnsServer1, End: dnsServer1}},
					Ports:    []netrules.PortRange{{Start: 53, End: 53}},
				}))
				Expect(ip).To(Equal(containerIP.String()))

				outRule, ip = netRuleApplier.OutArgsForCall(2)
				Expect(outRule).To(Equal(netrules.NetOut{
					Protocol: netrules.ProtocolTCP,
					Networks: []netrules.IPRange{{Start: dnsServer2, End: dnsServer2}},
					Ports:    []netrules.PortRange{{Start: 53, End: 53}},
				}))
				Expect(ip).To(Equal(containerIP.String()))

				outRule, ip = netRuleApplier.OutArgsForCall(3)
				Expect(outRule).To(Equal(netrules.NetOut{
					Protocol: netrules.ProtocolUDP,
					Networks: []netrules.IPRange{{Start: dnsServer2, End: dnsServer2}},
					Ports:    []netrules.PortRange{{Start: 53, End: 53}},
				}))
				Expect(ip).To(Equal(containerIP.String()))
			})
		})

		Context("net in fails", func() {
			BeforeEach(func() {
				netRuleApplier.InReturnsOnCall(0, netrules.PortMapping{}, errors.New("couldn't allocate port"))
			})

			It("cleans up allocated ports", func() {
				_, err := networkManager.Up(inputs)
				Expect(err).To(MatchError("couldn't allocate port"))
				Expect(netRuleApplier.CleanupCallCount()).To(Equal(1))
			})
		})

		Context("endpoint create fails", func() {
			BeforeEach(func() {
				endpointManager.CreateReturns(hcsshim.HNSEndpoint{}, errors.New("couldn't create endpoint"))
			})

			It("cleans up allocated ports", func() {
				_, err := networkManager.Up(inputs)
				Expect(err).To(MatchError("couldn't create endpoint"))
				Expect(netRuleApplier.CleanupCallCount()).To(Equal(1))
			})
		})

		Context("net out fails", func() {
			BeforeEach(func() {
				netRuleApplier.OutReturns(errors.New("couldn't set firewall rules"))
			})

			It("cleans up allocated ports, firewall rules and deletes the endpoint", func() {
				_, err := networkManager.Up(inputs)
				Expect(err).To(MatchError("couldn't set firewall rules"))
				Expect(netRuleApplier.CleanupCallCount()).To(Equal(1))
				Expect(endpointManager.DeleteCallCount()).To(Equal(1))
			})
		})

		Context("MTU fails", func() {
			BeforeEach(func() {
				netRuleApplier.ContainerMTUReturns(errors.New("couldn't set MTU"))
			})

			It("cleans up allocated ports, firewall rules and deletes the endpoint", func() {
				_, err := networkManager.Up(inputs)
				Expect(err).To(MatchError("couldn't set MTU"))
				Expect(netRuleApplier.CleanupCallCount()).To(Equal(1))
				Expect(endpointManager.DeleteCallCount()).To(Equal(1))
			})
		})
	})

	Describe("Down", func() {
		It("deletes the endpoint and cleans up the ports and firewall rules", func() {
			Expect(networkManager.Down()).To(Succeed())
			Expect(endpointManager.DeleteCallCount()).To(Equal(1))
			Expect(netRuleApplier.CleanupCallCount()).To(Equal(1))
		})

		Context("endpoint delete fails", func() {
			BeforeEach(func() {
				endpointManager.DeleteReturns(errors.New("couldn't delete endpoint"))
			})

			It("cleans up allocated ports, firewall rules but returns an error", func() {
				Expect(networkManager.Down()).To(MatchError("couldn't delete endpoint"))
				Expect(netRuleApplier.CleanupCallCount()).To(Equal(1))
				Expect(endpointManager.DeleteCallCount()).To(Equal(1))
			})
		})

		Context("host cleanup fails", func() {
			BeforeEach(func() {
				netRuleApplier.CleanupReturns(errors.New("couldn't remove firewall rules"))
			})

			It("deletes the endpoint but returns an error", func() {
				Expect(networkManager.Down()).To(MatchError("couldn't remove firewall rules"))
				Expect(netRuleApplier.CleanupCallCount()).To(Equal(1))
				Expect(endpointManager.DeleteCallCount()).To(Equal(1))
			})
		})

		Context("host cleanup + endpoint delete fail", func() {
			BeforeEach(func() {
				endpointManager.DeleteReturns(errors.New("couldn't delete endpoint"))
				netRuleApplier.CleanupReturns(errors.New("couldn't remove firewall rules"))
			})

			It("deletes the endpoint but returns an error", func() {
				Expect(networkManager.Down()).To(MatchError("couldn't delete endpoint, couldn't remove firewall rules"))
				Expect(netRuleApplier.CleanupCallCount()).To(Equal(1))
				Expect(endpointManager.DeleteCallCount()).To(Equal(1))
			})
		})
	})
})
