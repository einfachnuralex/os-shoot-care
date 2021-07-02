package workloadtopology_test

import (
	"context"
	"fmt"
	"github.com/einfachnuralex/os-shoot-care/openstack"
	"github.com/einfachnuralex/os-shoot-care/workloadtopology"
	gardenschema "github.com/gardener/gardener/pkg/apis/core/install"
	wrapper "github.com/stackitcloud/gophercloud-wrapper/pkg/openstack"
	"k8s.io/client-go/kubernetes/scheme"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
)

// AAAAAAAAAAAHHHHHHHHHHH
func TestTopologyBuilder_GetHypervisorByNode(t *testing.T) {
	computeClient, err := openstack.CreateComputeClient()
	if err != nil {
		fmt.Println(err)
	}
	serverClient := wrapper.OSServerClient{}
	serverClient.Configure(computeClient)

	gardenschema.Install(scheme.Scheme)

	client, err := client.New(ctrl.GetConfigOrDie(), client.Options{Scheme: scheme.Scheme})
	if err != nil {
		return
	}

	topologyBuilder := workloadtopology.NewTopologyBuilder(&serverClient, client)

	pipi, err := topologyBuilder.PrintTree(context.Background(), "e52ced9461b1489d94181932f6c393e4")

	file, err := os.Create("/home/xorax/hv-table")
	defer file.Close()

	for _, hypervisor := range pipi.Children() {
		for _, vm := range hypervisor.Children() {
			if vm.Children() == nil || len(vm.Children()) == 0 {
				fmt.Fprintf(file, "%s,%s\n", hypervisor.Name(), vm.Name())
			}

			for _, pod := range vm.Children() {
				fmt.Fprintf(file, "%s,%s,%s\n", hypervisor.Name(), vm.Name(), pod.Name())
			}
		}
	}
}
