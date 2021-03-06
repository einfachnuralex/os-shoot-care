package shoot

import (
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v2/volumes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/pagination"
	"strings"
)

type OSClients struct {
	Network *gophercloud.ServiceClient
	Compute *gophercloud.ServiceClient
	Storage *gophercloud.ServiceClient
}

func (g *OSClients) GetNetworkByName(name string, projectId string) (networks.Network, error) {
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

func (g *OSClients) GetInstancesByName(name string, projectId string) ([]servers.Server, error) {
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

func (g *OSClients) GetInstancesByNetwork(net string, projectId string) ([]servers.Server, error) {
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

func (g *OSClients) getServer(id string) (servers.Server, error) {
	srv, err := servers.Get(g.Compute, id).Extract()
	if err != nil {
		return servers.Server{}, err
	}
	return *srv, nil
}

func (g *OSClients) GetPortsByNetwork(net string, project string) ([]ports.Port, error) {
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

func (g *OSClients) GetServerLostVolumes(serverID string) ([]volumes.Volume, error) {
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

func (g *OSClients) GetVolumeAttachmentsForVolume(volumeId string) ([]volumes.Attachment, []servers.Server, error) {
	attachment, err := g.getVolumeAttachmentsFromVolume(volumeId)
	if err != nil {
		return nil, nil, err
	}

	servers, err := g.getServerWithAttachmentForVolume(volumeId)
	if err != nil {
		return nil, nil, err
	}

	return attachment, servers, nil
}

func (g *OSClients) getVolumeAttachmentsFromVolume(volumeId string) ([]volumes.Attachment, error) {
	vol, err := volumes.Get(g.Storage, volumeId).Extract()
	if err != nil {
		return nil, err
	}
	if vol == nil {
		return nil, fmt.Errorf("nilpointer for volume %s", volumeId)
	}
	return vol.Attachments, nil
}

func (g *OSClients) getServerWithAttachmentForVolume(volumeId string) ([]servers.Server, error) {
	hasVolumeAttachmentFilter := NewServerPageForHasVolumeAttachmentFilter()
	err := g.getAllServerAsPager().EachPage(hasVolumeAttachmentFilter.filter(volumeId))
	if err != nil {
		return nil, err
	}

	return hasVolumeAttachmentFilter.serverList, nil
}

func (g *OSClients) getAllServerAsPager() pagination.Pager {
	return servers.List(g.Compute, servers.ListOpts{})
}

type serverPageForHasVolumeAttachmentFilter struct {
	serverList []servers.Server
}

func NewServerPageForHasVolumeAttachmentFilter() *serverPageForHasVolumeAttachmentFilter {
	return &serverPageForHasVolumeAttachmentFilter{[]servers.Server{}}
}

func (r *serverPageForHasVolumeAttachmentFilter) filter(volumeId string) func(page pagination.Page) (bool, error) {
	return func(page pagination.Page) (bool, error) {
		extractedServers, _ := servers.ExtractServers(page)

		for _, server := range extractedServers {
			for _, attachment := range server.AttachedVolumes {
				if attachment.ID == volumeId {
					r.serverList = append(r.serverList, server)
				}
			}
		}
		return true, nil
	}
}
