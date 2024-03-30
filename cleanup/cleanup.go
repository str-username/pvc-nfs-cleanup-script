package cleanup

import (
	"github.com/rs/zerolog/log"
	"os"
	"pvc-cleanup/filesystem"
	"pvc-cleanup/kubernetes"
	"strings"
)

func RemovePvcs(basePath string) error {
	folders, err := filesystem.CheckDirectory(basePath)
	if err != nil {
		panic(err)
	}

	claims, err := kubernetes.ListUnusedPVCs()
	if err != nil {
		panic(err)
	}

	envDryRun := os.Getenv("DRY_RUN")

	if envDryRun == "" {
		envDryRun = "true"
	}

	if err != nil {
		panic(err)
	}

	for _, folderName := range folders {
		for namespaceName, pvcs := range claims {
			folderNameWithoutNamespace := strings.Replace(folderName, namespaceName+"-", "", 1)
			for _, pvcName := range pvcs {
				if folderNameWithoutNamespace == pvcName {
					if envDryRun == "true" {
						log.Info().Str("ns", namespaceName).Str("pvc", pvcName).Msg("will be delete")
					} else {
						err := kubernetes.DeletePvc(pvcName, namespaceName)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}
