package util

import clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

func GetCurrentContext(config clientcmdapi.Config) string {
	return config.CurrentContext
}

func GetCurrentNameSpace(config clientcmdapi.Config) string {
	currentContext := GetCurrentContext(config)
	return config.Contexts[currentContext].Namespace
}
