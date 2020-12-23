package scheduler

import (
	"k8s.io/api/core/v1"
	extender "k8s.io/kube-scheduler/extender/v1"
)

type Predicate struct {
	Name string
	Func func(pod v1.Pod, node v1.Node) (bool, error)
}

func (p Predicate) Handler(args extender.ExtenderArgs) *extender.ExtenderFilterResult {
	pod := args.Pod
	canSchedule := make([]v1.Node, 0, len(args.Nodes.Items))
	canNotSchedule := make(map[string]string)

	for _, node := range args.Nodes.Items {
		result, err := p.Func(*pod, node)
		if err != nil {
			canNotSchedule[node.Name] = err.Error()
		} else {
			if result {
				canSchedule = append(canSchedule, node)
			}
		}
	}

	result := extender.ExtenderFilterResult{
		Nodes: &v1.NodeList{
			Items: canSchedule,
		},
		FailedNodes: canNotSchedule,
		Error:       "",
	}

	return &result
}