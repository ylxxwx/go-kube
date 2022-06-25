package main

import (
	"context"
	"flag"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig string
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "/home/kmt/.kube/config", "kubeconfig file")
}

func showPods(cs *kubernetes.Clientset) {
	pods, err := cs.CoreV1().Pods(apiv1.NamespaceDefault).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Get Pods failed.")
		return
	}
	fmt.Printf("%-20s\t%-20s\n", "namespace", "name")
	for _, pod := range pods.Items {
		fmt.Printf("%-20s\t%-20s\n", pod.Namespace, pod.Name)
	}
}

func deployPod(cs *kubernetes.Clientset) {
	pod := apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-pod",
		},
		Spec: apiv1.PodSpec{
			Containers: []apiv1.Container{
				{
					Name:  "c-test",
					Image: "nginx:1.14.2",
				},
			},
		},
	}
	_, err := cs.CoreV1().Pods(metav1.NamespaceDefault).Create(context.Background(), &pod, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err.Error())
	}
}

func deletePod(cs *kubernetes.Clientset) {
	cs.CoreV1().Pods(metav1.NamespaceDefault).Delete(context.TODO(), "test-pod", metav1.DeleteOptions{})
}
func main() {
	fmt.Println("enter main.")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	showPods(clientSet)

	//deployPod(clientSet)
	deletePod(clientSet)

}
