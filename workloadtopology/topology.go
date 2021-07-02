package workloadtopology

import (
	"context"
	"errors"
	"fmt"
	gardenv1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/stackitcloud/gophercloud-wrapper/pkg/openstack"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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

func (r TopologyBuilder) PrintTree(ctx context.Context, projectID string) (TreeNode, error) {
	os, err := r.GetHypervisorByNode(projectID)
	if err != nil {
		return nil, err
	}

	for _, hypervisor := range os.Children() {
		for _, vm := range hypervisor.Children() {
			fmt.Println("fetching pods for vm " + vm.Name())
			pods, err := r.getPodsForNode(ctx, vm.Name())
			if err != nil {
				continue
			}

			for _, pod := range pods {
				vm.AddChild(
					&Pod{pod: pod},
				)
			}
		}
	}

	return os, nil
}

var sasa = map[string]map[string][]corev1.Pod{}

func (r TopologyBuilder) getPodsForNode(ctx context.Context, node string) ([]corev1.Pod, error) {
	shoot, err := r.getShootByNode(ctx, node)
	if err != nil {
		return nil, err
	}

	if sasa[shoot.Name] != nil {
		return sasa[shoot.Name][node], nil
	}

	shootClient, err := r.getKubeClientForShoot(ctx, shoot)
	if err != nil {
		return nil, err
	}

	var podList corev1.PodList
	if err := shootClient.List(ctx, &podList, &client.ListOptions{}); err != nil {
		return nil, err
	}

	if sasa[shoot.Name] == nil {
		sasa[shoot.Name] = map[string][]corev1.Pod{}
	}

	for _, pod := range podList.Items {
		if sasa[shoot.Name][pod.Spec.NodeName] == nil {
			sasa[shoot.Name][pod.Spec.NodeName] = []corev1.Pod{}
		}

		sasa[shoot.Name][pod.Spec.NodeName] = append(sasa[shoot.Name][pod.Spec.NodeName], pod)
	}

	return sasa[shoot.Name][node], nil
}

func (r *TopologyBuilder) getKubeClientForShoot(ctx context.Context, shoot *gardenv1beta1.Shoot) (client.Client, error) {
	var secret corev1.Secret
	err := r.gardenClient.Get(
		ctx,
		client.ObjectKey{Name: fmt.Sprintf("%s.kubeconfig", shoot.Name), Namespace: shoot.Namespace},
		&secret,
	)
	if err != nil {
		return nil, err
	}

	kubeconfigYaml, ok := secret.Data["kubeconfig"]
	if !ok {
		return nil, errors.New("kubeconfig not found")
	}

	clientConf, err := clientcmd.NewClientConfigFromBytes(kubeconfigYaml)
	if err != nil {
		return nil, err
	}

	restConf, err := clientConf.ClientConfig()
	if err != nil {
		return nil, err
	}

	return client.New(restConf, client.Options{})
}

func getConfigFromKubeconfig(kubeconfig string) (*rest.Config, error) {

	restConfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{},
	).ClientConfig()

	if err != nil {
		return nil, err
	}

	return restConfig, nil
}

func (r *TopologyBuilder) getShootByNode(ctx context.Context, nodeName string) (*gardenv1beta1.Shoot, error) {
	nodeNameParts := strings.Split(nodeName, "--")
	if len(nodeNameParts) < 3 {
		return nil, errors.New("node is not a shoot")
	}

	namespace := nodeNameParts[1]
	shootName := strings.Split(nodeNameParts[2], "-")[0]

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

		hypervisor.AddChild(&Node{server: server})
	}

	return &os, err
}
