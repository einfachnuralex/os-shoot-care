package workloadtopology

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
}

func (n *Node) Type() string {
	return "Node"
}
