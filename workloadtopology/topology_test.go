package workloadtopology_test

import (
	"fmt"
	"github.com/einfachnuralex/os-shoot-care/openstack"
	"github.com/einfachnuralex/os-shoot-care/workloadtopology"
	wrapper "github.com/stackitcloud/gophercloud-wrapper/pkg/openstack"
	"testing"
)

func TestTopologyBuilder_GetHypervisorByNode(t *testing.T) {
	computeClient, err := openstack.CreateComputeClient()
	if err != nil {
		fmt.Println(err)
	}
	serverClient := wrapper.OSServerClient{}
	serverClient.Configure(computeClient)

	topologyBuilder := workloadtopology.NewTopologyBuilder(&serverClient, nil)

	_, _ = topologyBuilder.GetHypervisorByNode("e52ced9461b1489d94181932f6c393e4")
}
