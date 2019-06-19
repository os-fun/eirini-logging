package main

import (
"fmt"
    "os"

    eirinix "github.com/SUSE/eirinix"
)

func main() {
    x := eirinix.NewManager(
            eirinix.ManagerOptions{
                Namespace:           "default",
                Host:                "10.0.2.2",
                Port:                3000,
                KubeConfig:          os.Getenv("KUBECONFIG"),
                FilterEiriniApps: false,
                OperatorFingerprint: "eirini-app-logging",
            })

    x.AddExtension(&Extension{})
    fmt.Println(x.Start())
}
