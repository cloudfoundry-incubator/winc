package process

import (
	"time"

	"code.cloudfoundry.org/winc/hcs"
	"github.com/Microsoft/hcsshim"
)

//go:generate counterfeiter -o fakes/hcsclient.go --fake-name HCSClient . HCSClient
type HCSClient interface {
	GetContainerProperties(string) (hcsshim.ContainerProperties, error)
	OpenContainer(string) (hcs.Container, error)
}

type Manager struct {
	hcsClient HCSClient
}

func NewManager(hcsClient HCSClient) *Manager {
	return &Manager{
		hcsClient: hcsClient,
	}
}

func (m *Manager) ContainerPid(id string) (int, error) {
	container, err := m.hcsClient.OpenContainer(id)
	if err != nil {
		return -1, err
	}
	defer container.Close()

	pl, err := container.ProcessList()
	if err != nil {
		return -1, err
	}

	var process hcsshim.ProcessListItem
	oldestTime := time.Now()
	for _, v := range pl {
		if v.ImageName == "wininit.exe" && v.CreateTimestamp.Before(oldestTime) {
			oldestTime = v.CreateTimestamp
			process = v
		}
	}

	return int(process.ProcessId), nil
}
