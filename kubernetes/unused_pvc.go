package kubernetes

import (
	"context"
	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListUnusedPVCs : check if pvc used by pod. Return map { namespace: [ pvc-names ] }
func ListUnusedPVCs() (map[string][]string, error) {
	c, err := Client()

	if err != nil {
		return nil, err
	}

	namespaces := metav1.NamespaceAll
	pvcs, err := c.CoreV1().PersistentVolumeClaims(namespaces).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		log.Fatal().Err(err).Msg("error list pvcs")
		return nil, err
	}

	pods, err := c.CoreV1().Pods(namespaces).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		log.Fatal().Err(err).Msg("error list pods")
		return nil, err
	}

	usedPvcs := make(map[string]bool)

	for _, pod := range pods.Items {
		for _, vol := range pod.Spec.Volumes {
			if pvc := vol.PersistentVolumeClaim; pvc != nil {
				usedPvcs[pvc.ClaimName] = true
			}
		}
	}

	unusedPvcs := make(map[string][]string)

	for _, pvc := range pvcs.Items {
		if _, used := usedPvcs[pvc.Name]; !used {
			unusedPvcs[pvc.Namespace] = append(unusedPvcs[pvc.Namespace], pvc.Name)
		}
	}

	return unusedPvcs, nil
}
