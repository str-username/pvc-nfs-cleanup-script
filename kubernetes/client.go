package kubernetes

import (
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

// Client : create kubernetes client
func Client() (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error

	if _, exists := os.LookupEnv("KUBERNETES_SERVICE_HOST"); exists {
		config, err = rest.InClusterConfig()
	} else {
		configPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", configPath)
	}
	if err != nil {
		log.Fatal().Err(err).Msg("error build client")
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(config)

	if err != nil {
		log.Fatal().Err(err).Msg("error create clientset")
		return nil, err
	}

	return clientSet, nil
}
