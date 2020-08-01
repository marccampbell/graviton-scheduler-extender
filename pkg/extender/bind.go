package extender

import (
	"k8s.io/apimachinery/pkg/types"
	extenderapi "k8s.io/kube-scheduler/extender/v1"
)

type Bind struct {
	Func func(podName string, podNamespace string, podUID types.UID, node string) error
}

func (b Bind) Handler(args extenderapi.ExtenderBindingArgs) *extenderapi.ExtenderBindingResult {
	err := b.Func(args.PodName, args.PodNamespace, args.PodUID, args.Node)
	return &extenderapi.ExtenderBindingResult{
		Error: err.Error(),
	}
}
