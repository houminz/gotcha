package scheduler

import (
	"k8s.io/api/core/v1"
	extender "k8s.io/kube-scheduler/extender/v1"
)

type Preemption struct {
	Func func(
		pod v1.Pod,
		nodeNameToVictims map[string]*extender.Victims,
		nodeNameToMetaVictims map[string]*extender.MetaVictims,
	) map[string]*extender.MetaVictims
}

func (b Preemption) Handler(args extender.ExtenderPreemptionArgs, ) *extender.ExtenderPreemptionResult {
	nodeNameToMetaVictims := b.Func(*args.Pod, args.NodeNameToVictims, args.NodeNameToMetaVictims)
	return &extender.ExtenderPreemptionResult{
		NodeNameToMetaVictims: nodeNameToMetaVictims,
	}
}