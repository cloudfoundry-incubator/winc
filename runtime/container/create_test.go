package container_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/winc/container"
	"code.cloudfoundry.org/winc/container/fakes"
	"code.cloudfoundry.org/winc/hcs"
	hcsfakes "code.cloudfoundry.org/winc/hcs/fakes"
	"github.com/Microsoft/hcsshim"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
)

var _ = Describe("Create", func() {
	const (
		containerVolume = "containervolume"
		hostName        = "some-hostname"
		containerId     = "my-container"
	)

	var (
		layerFolders     []string
		hcsClient        *fakes.HCSClient
		containerManager *container.Manager
		spec             *specs.Spec
	)

	BeforeEach(func() {
		layerFolders = []string{
			"some-layer",
			"some-other-layer",
			"some-rootfs",
		}

		spec = &specs.Spec{
			Version: specs.Version,
			Process: &specs.Process{
				Args: []string{"cmd.exe"},
				Cwd:  "C:\\",
			},
			Root: &specs.Root{
				Path: containerVolume,
			},
			Windows: &specs.Windows{
				LayerFolders: layerFolders,
			},
			Hostname: hostName,
		}

		hcsClient = &fakes.HCSClient{}
		logger := (&logrus.Logger{
			Out: ioutil.Discard,
		}).WithField("test", "create")

		containerManager = container.NewManager(logger, hcsClient, containerId)
	})

	Context("when the specified container does not already exist", func() {
		var (
			expectedHcsshimLayers []hcsshim.Layer
			fakeContainer         hcsfakes.Container
		)

		BeforeEach(func() {
			fakeContainer = hcsfakes.Container{}
			hcsClient.GetContainerPropertiesReturns(hcsshim.ContainerProperties{}, &hcs.NotFoundError{})

			expectedHcsshimLayers = []hcsshim.Layer{}
			for i, l := range layerFolders {
				guid := hcsshim.NewGUID(fmt.Sprintf("layer-%d", i))
				hcsClient.NameToGuidReturnsOnCall(i, *guid, nil)
				expectedHcsshimLayers = append(expectedHcsshimLayers, hcsshim.Layer{
					ID:   guid.ToString(),
					Path: l,
				})
			}

			hcsClient.CreateContainerReturns(&fakeContainer, nil)
			hcsClient.OpenContainerReturns(&fakeContainer, nil)
		})

		It("creates and starts it", func() {
			Expect(containerManager.Create(spec)).To(Succeed())

			Expect(hcsClient.GetContainerPropertiesCallCount()).To(Equal(1))
			Expect(hcsClient.GetContainerPropertiesArgsForCall(0)).To(Equal(containerId))

			Expect(hcsClient.NameToGuidCallCount()).To(Equal(len(layerFolders)))
			for i, l := range layerFolders {
				Expect(hcsClient.NameToGuidArgsForCall(i)).To(Equal(filepath.Base(l)))
			}

			Expect(hcsClient.CreateContainerCallCount()).To(Equal(1))
			actualContainerId, containerConfig := hcsClient.CreateContainerArgsForCall(0)
			Expect(actualContainerId).To(Equal(containerId))
			Expect(containerConfig).To(Equal(&hcsshim.ContainerConfig{
				SystemType:        "Container",
				HostName:          hostName,
				VolumePath:        containerVolume,
				LayerFolderPath:   "ignored",
				Layers:            expectedHcsshimLayers,
				MappedDirectories: []hcsshim.MappedDir{},
			}))

			Expect(fakeContainer.StartCallCount()).To(Equal(1))
		})

		Context("when mounts are specified in the spec", func() {
			var (
				expectedMappedDirs []hcsshim.MappedDir
				mount              string
			)

			BeforeEach(func() {
				var err error
				mount, err = ioutil.TempDir("", "mountdir")
				Expect(err).ToNot(HaveOccurred())

				spec.Mounts = []specs.Mount{
					{Source: mount, Destination: "/bar"},
				}

				expectedMappedDirs = []hcsshim.MappedDir{
					{HostPath: mount, ContainerPath: "C:\\bar", ReadOnly: true},
				}
			})

			AfterEach(func() {
				Expect(os.RemoveAll(mount)).To(Succeed())
			})

			Context("mount options do not specify ro or rw", func() {
				BeforeEach(func() {
					spec.Mounts[0].Options = []string{"bind"}

					expectedMappedDirs[0].ReadOnly = true
				})

				It("creates the container with the specified mounts", func() {
					Expect(containerManager.Create(spec)).To(Succeed())

					Expect(hcsClient.CreateContainerCallCount()).To(Equal(1))
					actualContainerId, containerConfig := hcsClient.CreateContainerArgsForCall(0)
					Expect(actualContainerId).To(Equal(containerId))
					Expect(containerConfig.MappedDirectories).To(ConsistOf(expectedMappedDirs))
				})
			})

			Context("mount options specify ro", func() {
				BeforeEach(func() {
					spec.Mounts[0].Options = []string{"bind", "ro"}

					expectedMappedDirs[0].ReadOnly = true
				})

				It("creates the container with the specified mounts", func() {
					Expect(containerManager.Create(spec)).To(Succeed())

					Expect(hcsClient.CreateContainerCallCount()).To(Equal(1))
					actualContainerId, containerConfig := hcsClient.CreateContainerArgsForCall(0)
					Expect(actualContainerId).To(Equal(containerId))
					Expect(containerConfig.MappedDirectories).To(ConsistOf(expectedMappedDirs))
				})
			})

			Context("mount options specify rw", func() {
				BeforeEach(func() {
					spec.Mounts[0].Options = []string{"bind", "rw"}

					expectedMappedDirs[0].ReadOnly = false
				})

				It("creates the container with the specified mounts", func() {
					Expect(containerManager.Create(spec)).To(Succeed())

					Expect(hcsClient.CreateContainerCallCount()).To(Equal(1))
					actualContainerId, containerConfig := hcsClient.CreateContainerArgsForCall(0)
					Expect(actualContainerId).To(Equal(containerId))
					Expect(containerConfig.MappedDirectories).To(ConsistOf(expectedMappedDirs))
				})
			})

			Context("mount options specify both rw and ro", func() {
				BeforeEach(func() {
					spec.Mounts[0].Options = []string{"bind", "rw", "ro"}
				})

				It("errors", func() {
					err := containerManager.Create(spec)
					Expect(err).To(HaveOccurred())
					Expect(err).To(BeAssignableToTypeOf(&container.InvalidMountOptionsError{}))
				})
			})

			Context("when the mount does not exist", func() {
				BeforeEach(func() {
					Expect(os.RemoveAll(mount)).To(Succeed())
				})

				It("errors", func() {
					err := containerManager.Create(spec)
					Expect(os.IsNotExist(err)).To(BeTrue())
				})
			})

			Context("when a file is specified as a mount", func() {
				var mountFile string

				BeforeEach(func() {
					m, err := ioutil.TempFile("", "mountfile")
					Expect(err).ToNot(HaveOccurred())
					Expect(m.Close()).To(Succeed())
					mountFile = m.Name()

					spec.Mounts = append(spec.Mounts, specs.Mount{
						Source:      mountFile,
						Destination: "foo",
					})
				})

				AfterEach(func() {
					Expect(os.RemoveAll(mountFile)).To(Succeed())
				})

				It("ignores it", func() {
					Expect(containerManager.Create(spec)).To(Succeed())

					Expect(hcsClient.CreateContainerCallCount()).To(Equal(1))
					actualContainerId, containerConfig := hcsClient.CreateContainerArgsForCall(0)
					Expect(actualContainerId).To(Equal(containerId))
					Expect(containerConfig.MappedDirectories).To(ConsistOf(expectedMappedDirs))
				})
			})
		})

		Context("when memory limits are specified in the spec", func() {
			var expectedMemoryMaxinMB uint64

			BeforeEach(func() {
				expectedMemoryMaxinMB = uint64(64)
				expectedMemoryMaxinBytes := expectedMemoryMaxinMB * 1024 * 1024
				spec.Windows.Resources = &specs.WindowsResources{
					Memory: &specs.WindowsMemoryResources{
						Limit: &expectedMemoryMaxinBytes,
					},
				}
			})

			It("creates the container with the specified memory limits", func() {
				Expect(containerManager.Create(spec)).To(Succeed())

				Expect(hcsClient.CreateContainerCallCount()).To(Equal(1))
				_, containerConfig := hcsClient.CreateContainerArgsForCall(0)
				Expect(containerConfig.MemoryMaximumInMB).To(Equal(int64(expectedMemoryMaxinMB)))
			})
		})

		Context("when cpu limits are specified in the spec", func() {
			var expectedCPUShares uint16

			BeforeEach(func() {
				expectedCPUShares = 8080
				spec.Windows.Resources = &specs.WindowsResources{
					CPU: &specs.WindowsCPUResources{
						Shares: &expectedCPUShares,
					},
				}
			})

			It("creates the container with the specified cpu limits", func() {
				Expect(containerManager.Create(spec)).To(Succeed())

				Expect(hcsClient.CreateContainerCallCount()).To(Equal(1))
				_, containerConfig := hcsClient.CreateContainerArgsForCall(0)
				Expect(containerConfig.ProcessorWeight).To(Equal(uint64(expectedCPUShares)))
			})
		})

		Context("when network settings are specified in the spec", func() {
			Context("when NetworkSharedContainerName is specified", func() {
				var (
					networkSharedContainerName string
					sharedEndpointId           string
				)

				BeforeEach(func() {
					networkSharedContainerName = "some-networked-container"
					sharedEndpointId = "some-shared-endpoint-id"
					spec.Windows.Network = &specs.WindowsNetwork{NetworkSharedContainerName: networkSharedContainerName}

					hcsClient.GetHNSEndpointByNameReturns(&hcsshim.HNSEndpoint{Id: sharedEndpointId}, nil)
				})

				It("creates the container with a NetworkSharedContainerName and EndpointList", func() {
					Expect(containerManager.Create(spec)).To(Succeed())

					Expect(hcsClient.CreateContainerCallCount()).To(Equal(1))
					_, containerConfig := hcsClient.CreateContainerArgsForCall(0)
					Expect(containerConfig.NetworkSharedContainerName).To(Equal(networkSharedContainerName))
					Expect(containerConfig.Owner).To(Equal(networkSharedContainerName))
					Expect(containerConfig.EndpointList).To(Equal([]string{sharedEndpointId}))

					Expect(hcsClient.GetHNSEndpointByNameArgsForCall(0)).To(Equal(networkSharedContainerName))
				})

				Context("when getting the endpoint fails", func() {
					BeforeEach(func() {
						hcsClient.GetHNSEndpointByNameReturns(nil, errors.New("couldn't get endpoint"))
					})

					It("returns an error", func() {
						err := containerManager.Create(spec)
						Expect(err).To(MatchError("couldn't get endpoint"))
					})
				})
			})

			Context("when NetworkSharedContainerName is empty", func() {
				BeforeEach(func() {
					spec.Windows.Network = &specs.WindowsNetwork{}
				})

				It("creates a container without a NetworkSharedContainerName or EndpointList", func() {
					Expect(containerManager.Create(spec)).To(Succeed())

					Expect(hcsClient.CreateContainerCallCount()).To(Equal(1))
					_, containerConfig := hcsClient.CreateContainerArgsForCall(0)
					Expect(containerConfig.NetworkSharedContainerName).To(BeEmpty())
					Expect(containerConfig.EndpointList).To(BeEmpty())
				})
			})
		})

		Context("when CreateContainer fails", func() {
			BeforeEach(func() {
				hcsClient.CreateContainerReturns(nil, errors.New("couldn't create"))
			})

			It("returns an error", func() {
				err := containerManager.Create(spec)
				Expect(err).To(MatchError("couldn't create"))
			})
		})

		Context("when container Start fails", func() {
			BeforeEach(func() {
				fakeContainer.StartReturns(errors.New("couldn't start"))
				hcsClient.GetContainerPropertiesReturnsOnCall(1, hcsshim.ContainerProperties{Stopped: true}, nil)
			})

			It("closes but doesn't shutdown or terminate the container", func() {
				err := containerManager.Create(spec)
				Expect(err).To(MatchError("couldn't start"))

				Expect(fakeContainer.CloseCallCount()).To(Equal(1))
				Expect(fakeContainer.ShutdownCallCount()).To(Equal(0))
				Expect(fakeContainer.TerminateCallCount()).To(Equal(0))
			})
		})
	})
})
