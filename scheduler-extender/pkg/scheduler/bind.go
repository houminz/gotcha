package scheduler

import (
	"k8s.io/apimachinery/pkg/types"
	extender "k8s.io/kube-scheduler/extender/v1"
)

type Bind struct {
	Func func(podName string, podNamespace string, podUID types.UID, node string) error
}

func (b Bind) Handler(args extender.ExtenderBindingArgs) *extender.ExtenderBindingResult {
	err := b.Func(args.PodName, args.PodNamespace, args.PodUID, args.Node)
	return &extender.ExtenderBindingResult{
		Error: err.Error(),
	}
}

