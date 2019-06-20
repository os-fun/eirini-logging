package main

import (
	"fmt"
	"os"

	eirinix "github.com/SUSE/eirinix"
)

func main() {

	ns := "default"
	x := eirinix.NewManager(
		eirinix.ManagerOptions{
			Namespace:           ns,
			Host:                "10.0.2.2",
			Port:                3000,
			KubeConfig:          os.Getenv("KUBECONFIG"),
			FilterEiriniApps:    false,
			OperatorFingerprint: "eirini-app-logging",
		})

	x.AddExtension(&Extension{Namespace: ns})
	fmt.Println(x.Start())
}
