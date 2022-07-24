package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/ylxxwx/go-kube/crd_3/pkg/generated/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	//cs, err := kubernetes.NewForConfig(config)
	fcs, err := versioned.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	listFoo(fcs)
}

func listFoo(cs *versioned.Clientset) {
	foos, err := cs.FoogroupV1alpha1().Foos("default").List(context.Background(), v1.ListOptions{})
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("%-15s\t%-15s\t%-15s\n", "Name", "Namespace", "Time")
	for _, foo := range foos.Items {
		fmt.Printf("%-15s\t%-15s\t%-15s\n", foo.Name, foo.Namespace, foo.CreationTimestamp)
	}
}
