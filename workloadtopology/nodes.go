package workloadtopology

import (
	"github.com/stackitcloud/gophercloud-wrapper/pkg/openstack"
	corev1 "k8s.io/api/core/v1"
)

type TreeNode interface {
	Type() string
	Children() []TreeNode
	Name() string
	GetChildWithName(name string) TreeNode
	AddChild(TreeNode)
}

type BaseNode struct {
	nodes []TreeNode
	name  string
}

func (o *BaseNode) Children() []TreeNode {
	return o.nodes
}

func (o *BaseNode) Name() string {
	return o.name
}

func (o *BaseNode) AddChild(node TreeNode) {
	o.nodes = append(o.nodes, node)
}

func (o *BaseNode) GetChildWithName(name string) TreeNode {
	for _, node := range o.Children() {
		if node.Name() == name {
			return node
		}
	}
	return nil
}

type OpenStack struct {
	BaseNode
}

func (o *OpenStack) Type() string {
	return "Openstack"
}

type Hypervisor struct {
	BaseNode
}

func (h *Hypervisor) Type() string {
	return "Hypervisor"
}

type Node struct {
	BaseNode
	server openstack.ExtendedServer
}

func (n *Node) Type() string {
	return "Node"
}

type Pod struct {
	BaseNode
	pod corev1.Pod
}

func (n *Pod) Name() string {
	return n.pod.ObjectMeta.Name
}

func (n *Pod) Type() string {
	return "Pod"
}
