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

var (
	kubeconfig string
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "/home/kmt/.kube/config", "kubeconfig file")
}

func showPods(cs *kubernetes.Clientset) {
	pods, err := cs.CoreV1().Pods(corev1.NamespaceDefault).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Get Pods failed.")
		return
	}
	fmt.Printf("%-15s\t%-10s\t%-20s\n", "NAMESPACE", "STATUS", "NAME")
	for _, pod := range pods.Items {
		fmt.Printf("%-15s\t%-10s\t%-20s\n", pod.Namespace, pod.Status.Phase, pod.Name)
	}
}

func deployPod(cs *kubernetes.Clientset) (func() error, error) {
	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-pod",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
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

	return func() error {
		return cs.CoreV1().Pods(metav1.NamespaceDefault).Delete(context.TODO(), "test-pod", metav1.DeleteOptions{})
	}, nil
}

func main() {
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	f, err := deployPod(clientSet)
	if err != nil {
		fmt.Printf("Create Pod failed. %s\n", err.Error())
		return
	}
	time.Sleep(10 * time.Second)

	showPods(clientSet)

	_ = f()
	time.Sleep(10 * time.Second)

	showPods(clientSet)
}
