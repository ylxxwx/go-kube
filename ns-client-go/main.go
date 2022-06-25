package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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
	if err != nil {
		panic(err)
	}
	f, err := createNS(cs)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	time.Sleep(5 * time.Second)
	showNS(cs)
	f()
	time.Sleep(5 * time.Second)
	showNS(cs)
}

func showNS(cs *kubernetes.Clientset) {
	nss, err := cs.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("list ns failed.%s\n", err.Error())
		return
	}
	fmt.Printf("%-15s\t%-15s\n", "Name", "CreateTime")
	for _, ns := range nss.Items {
		fmt.Printf("%-15s\t%-15s\n", ns.Name, ns.CreationTimestamp)
	}
}

func createNS(cs *kubernetes.Clientset) (func() error, error) {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-ns-9",
		},
		Spec: corev1.NamespaceSpec{
			Finalizers: []corev1.FinalizerName{},
		},
	}
	_, err := cs.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("create ns failed.%s\n", err.Error())
		return nil, err
	}

	return func() error {
		return cs.CoreV1().Namespaces().Delete(context.TODO(), ns.ObjectMeta.Name, metav1.DeleteOptions{})
	}, nil

}
