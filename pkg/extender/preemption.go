package extender

import (
	v1 "k8s.io/api/core/v1"
	extenderapi "k8s.io/kube-scheduler/extender/v1"
)

type Preemption struct {
	Func func(pod v1.Pod, nodeNameToVictims map[string]*extenderapi.Victims, nodeNameToMetaVictims map[string]*extenderapi.MetaVictims) map[string]*extenderapi.MetaVictims
}

func (b Preemption) Handler(args extenderapi.ExtenderPreemptionArgs) *extenderapi.ExtenderPreemptionResult {
	nodeNameToMetaVictims := b.Func(*args.Pod, args.NodeNameToVictims, args.NodeNameToMetaVictims)
	return &extenderapi.ExtenderPreemptionResult{
		NodeNameToMetaVictims: nodeNameToMetaVictims,
	}
}
