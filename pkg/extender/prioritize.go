package extender

import (
	v1 "k8s.io/api/core/v1"
	extenderapi "k8s.io/kube-scheduler/extender/v1"
)

type Prioritize struct {
	Name string
	Func func(pod v1.Pod, nodes []v1.Node) (*extenderapi.HostPriorityList, error)
}

func (p Prioritize) Handler(args extenderapi.ExtenderArgs) (*extenderapi.HostPriorityList, error) {
	return p.Func(*args.Pod, args.Nodes.Items)
}
