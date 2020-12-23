package scheduler

import (
	"k8s.io/api/core/v1"
	extender "k8s.io/kube-scheduler/extender/v1"
)

type Prioritize struct {
	Name string
	Func func(pod v1.Pod, nodes []v1.Node) (*extender.HostPriorityList, error)
}

func (p Prioritize) Handler(args extender.ExtenderArgs) (*extender.HostPriorityList, error) {
	return p.Func(*args.Pod, args.Nodes.Items)
}
