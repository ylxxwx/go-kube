package main

import (
	"flag"
	"fmt"
	"time"

	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig string

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "/home/kmt/.kube/config", "kubeconfig file")
}

func main() {
	fmt.Println("Enter main")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	cs, err := kubernetes.NewForConfig(config)
	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(cs, time.Second*30)

	if err != nil {
		panic(err)
	}

	controller := NewController(cs, kubeInformerFactory.Core().V1().Namespaces())

}
