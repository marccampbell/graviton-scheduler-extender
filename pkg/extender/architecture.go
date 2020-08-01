package extender

import (
	"context"
	"fmt"

	"github.com/containers/image/v5/docker"
	imagetypes "github.com/containers/image/v5/types"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
)

type PodArchitectureCounts struct {
	CountTotal int
	CountAmd64 int
	CountArm64 int
}

func architectureForPod(pod v1.Pod) (*PodArchitectureCounts, error) {
	counts := PodArchitectureCounts{}

	for _, container := range pod.Spec.InitContainers {
		hasArm, hasAmd, err := architectureForContainer(container)
		if err != nil {
			return nil, err
		}

		if hasArm || hasAmd {
			counts.CountTotal++

			if hasArm {
				counts.CountArm64++
			}
			if hasAmd {
				counts.CountAmd64++
			}
		}
	}

	for _, container := range pod.Spec.Containers {
		hasArm, hasAmd, err := architectureForContainer(container)
		if err != nil {
			return nil, err
		}

		if hasArm || hasAmd {
			counts.CountTotal++

			if hasArm {
				counts.CountArm64++
			}
			if hasAmd {
				counts.CountAmd64++
			}
		}
	}

	return &counts, nil
}

func architectureForContainer(container v1.Container) (bool, bool, error) {
	// TODO cache to help with rate limits
	// the images library caches some

	ref, err := docker.ParseReference(fmt.Sprintf("//%s", container.Image))
	if err != nil {
		return false, false, errors.Wrapf(err, "docker parsereference %q", container.Image)
	}

	ctx := context.Background()

	// check all architectures
	hasAmd64, err := hasImageForArch(ctx, ref, "amd64")
	if err != nil {
		return false, false, nil
	}
	hasArm64, err := hasImageForArch(ctx, ref, "arm64")
	if err != nil {
		return false, false, nil
	}

	return hasArm64, hasAmd64, nil
}

func hasImageForArch(ctx context.Context, ref imagetypes.ImageReference, arch string) (bool, error) {
	systemCtx := imagetypes.SystemContext{
		ArchitectureChoice: arch,
	}
	img, err := ref.NewImage(ctx, &systemCtx)
	if err != nil {
		return false, errors.Wrap(err, "newimage")
	}
	defer img.Close()
	b, _, err := img.Manifest(ctx)
	if err != nil {
		return false, errors.Wrap(err, "manifest")
	}

	fmt.Printf("\n\n%s: %s\n", arch, b)

	return false, nil
}

func getNodeArchitecture(node v1.Node) (string, error) {
	nodeArchitecture, ok := node.Labels["kubernetes.io/arch"]
	if ok {
		return nodeArchitecture, nil
	}

	nodeArchitecture, ok = node.Labels["beta.kubernetes.io/arch"]
	if ok {
		return nodeArchitecture, nil
	}

	return "", errors.New("no arch label")
}
