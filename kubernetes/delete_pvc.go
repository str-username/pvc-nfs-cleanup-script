package kubernetes

import (
	"context"
	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeletePvc : delete pvc by name \ namespace
func DeletePvc(pvcName string, pvcNamespace string) error {
	c, err := Client()

	if err != nil {
		log.Fatal().Err(err)
		panic(err)
	}

	log.Info().Str("ns", pvcNamespace).Str("pvc", pvcName).Msg("delete")

	return c.CoreV1().PersistentVolumeClaims(pvcNamespace).Delete(context.TODO(), pvcName, metav1.DeleteOptions{})
}
