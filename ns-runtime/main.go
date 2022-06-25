package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {
	fmt.Println("Enter main")
	cc, err := client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		panic(err)
	}
	showNS(cc)
}

func showNS(cc client.Client) {
	nss := &corev1.NamespaceList{}
	err := cc.List(context.TODO(), nss, client.InNamespace(""))
	if err != nil {
		return
	}
	for _, ns := range nss.Items {
		fmt.Printf("%-15s\t%-15s\n", ns.Name, ns.CreationTimestamp)
	}
}
