package shootcare

import (
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v2/volumes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"strings"
)

type Gandalf struct {
	Network *gophercloud.ServiceClient
	Compute *gophercloud.ServiceClient
	Storage *gophercloud.ServiceClient
}

func (g *Gandalf) GetNetworkByName(name string, projectId string) (networks.Network, error) {
	listOpts := networks.ListOpts{
		TenantID: projectId,
		Name:     name,
	}

	allPages, err := networks.List(g.Network, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allNetworks, err := networks.ExtractNetworks(allPages)
	if err != nil {
		panic(err)
	}

	if len(allNetworks) == 0 {
		return networks.Network{}, nil
	}

	if len(allNetworks) > 1 {
		return networks.Network{}, fmt.Errorf("more than one network found")
	}

	return allNetworks[0], nil
}

func (g *Gandalf) GetInstancesByName(name string, projectId string) ([]servers.Server, error) {
	listOpts := servers.ListOpts{
		Name: name + ".*",
		//AllTenants:   true,
		TenantID: projectId,
	}

	allPages, err := servers.List(g.Compute, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allServers, err := servers.ExtractServers(allPages)
	if err != nil {
		panic(err)
	}

	return allServers, nil
}

func (g *Gandalf) GetInstancesByNetwork(net string, projectId string) ([]servers.Server, error) {
	listOpts := ports.ListOpts{
		NetworkID: net,
		ProjectID: projectId,
	}

	allPages, err := ports.List(g.Network, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allPorts, err := ports.ExtractPorts(allPages)
	if err != nil {
		panic(err)
	}

	var allServers []servers.Server

	for _, port := range allPorts {
		if strings.HasPrefix(port.DeviceOwner, "compute:") {
			srv, err := g.getServer(port.DeviceID)
			if err == nil {
				allServers = append(allServers, srv)
			}
		}
	}

	return allServers, nil
}

func (g *Gandalf) getServer(id string) (servers.Server, error) {
	srv, err := servers.Get(g.Compute, id).Extract()
	if err != nil {
		return servers.Server{}, err
	}
	return *srv, nil
}

func (g *Gandalf) GetPortsByNetwork(net string, project string) ([]ports.Port, error) {
	listOpts := ports.ListOpts{
		NetworkID: net,
		ProjectID: project,
	}

	allPages, err := ports.List(g.Network, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allPorts, err := ports.ExtractPorts(allPages)
	if err != nil {
		panic(err)
	}

	return allPorts, nil
}

func (g *Gandalf) GetServerLostVolumes(serverID string) ([]volumes.Volume, error) {
	var faultyVolumes []volumes.Volume
	srv, err := servers.Get(g.Compute, serverID).Extract()
	if err != nil {
		return []volumes.Volume{}, err
	}
	vols := srv.AttachedVolumes
	for _, vol := range vols {
		gvol, err := volumes.Get(g.Storage, vol.ID).Extract()

		if err == nil {
			f := gvol.Attachments
			for _, g := range f {
				if g.ServerID != srv.ID {
					faultyVolumes = append(faultyVolumes, *gvol)
				}

			}
		}
	}

	return faultyVolumes, nil
}
