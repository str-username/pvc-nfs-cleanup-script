package filesystem

import (
	"github.com/rs/zerolog/log"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// containsFiles: check that folder don't have a files
func containsFiles(basePath string) (bool, error) {
	var contains bool

	err := filepath.WalkDir(basePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal().Str("WalkDir", basePath).Msg("can't walk in dir")
			return err
		}
		if !d.IsDir() {
			contains = true
			return filepath.SkipDir
		}
		return nil
	})

	if err != nil {
		log.Fatal().Str("WalkDir", basePath).Msg("can't walk in dir")
		return false, err
	}

	return contains, nil
}

// CheckDirectory : recursive check any folders
func CheckDirectory(basePath string) ([]string, error) {
	var empty []string

	directories, err := os.ReadDir(basePath)

	if err != nil {
		log.Fatal().Str("ReadDir", basePath).Msg("can't read dir")
		return nil, err
	}

	for _, directory := range directories {
		if directory.IsDir() {
			path := filepath.Join(basePath, directory.Name())
			containsFiles, err := containsFiles(path)
			if err != nil {
				log.Warn().Err(err).Str("dir", path).Msg("error checking directory")
				continue
			}
			if !containsFiles {
				removeBasePath := strings.Replace(path, basePath, "", 1)
				re := regexp.MustCompile(`^(.*?-pvc)[^-]*`)
				pvcName := re.FindString(removeBasePath)
				empty = append(empty, pvcName)
				log.Info().Str("dir", pvcName).Msg("empty directory")
			}
		}
	}
	return empty, nil
}
