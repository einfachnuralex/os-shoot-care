package workloadtopology

import (
	"context"
	gardenv1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/stackitcloud/gophercloud-wrapper/pkg/openstack"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

type TopologyBuilder struct {
	server       openstack.ServerClient
	gardenClient client.Client
}

//   --Shoot--
//   |		|
// Pod <=> Node <=> Hypervisor

// Pod <=> Hypervisor

func NewTopologyBuilder(serverClient openstack.ServerClient, gardenClient client.Client) *TopologyBuilder {
	return &TopologyBuilder{server: serverClient, gardenClient: gardenClient}
}

func (r TopologyBuilder) PrintTree(projectID string) {

}

func (r *TopologyBuilder) getShootByNode(ctx context.Context, nodeName string) (*gardenv1beta1.Shoot, error) {
	nodeNameParts := strings.Split(nodeName, "--")
	namespace := nodeNameParts[1]
	shootName := nodeNameParts[2]

	var shoot gardenv1beta1.Shoot
	if err := r.gardenClient.Get(ctx, client.ObjectKey{Name: shootName, Namespace: namespace}, &shoot); err != nil {
		return nil, err
	}

	return &shoot, nil
}

func (r TopologyBuilder) GetHypervisorByNode(projectID string) (TreeNode, error) {
	serverList, err := r.server.ListExtended(servers.ListOpts{
		AllTenants: true,
		TenantID:   projectID,
	})

	os := OpenStack{
		BaseNode{
			name: projectID,
		},
	}

	for _, server := range serverList {
		hypervisor := os.GetChildWithName(server.HypervisorHostname)
		if hypervisor == nil {
			hypervisor = &Hypervisor{
				BaseNode{
					name: server.HypervisorHostname,
				},
			}
			os.AddChild(hypervisor)
		}

		hypervisor.AddChild(&Node{BaseNode{name: server.Name}})
	}

	return &os, err
}
