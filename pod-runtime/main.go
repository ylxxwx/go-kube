package main

import (
	"context"
	"fmt"
	"os"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func init() {
	//flag.StringVar(&kubeconfig, "kubeconfig", "/home/kmt/.kube/config", "kubeconfig file")
}

func showPods(cl client.Client) {
	podList := &corev1.PodList{}
	cl.List(context.Background(), podList, client.InNamespace("default"))
	fmt.Printf("%-15s\t%-10s\t%-20s\n", "NAMESPACE", "STATUS", "NAME")
	for _, pod := range podList.Items {
		fmt.Printf("%-15s\t%-10s\t%-20s\n", pod.Namespace, pod.Status.Phase, pod.Name)
	}
}

func deployPods(cl client.Client) (func() error, error) {
	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
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

	if err := cl.Create(context.Background(), &pod, &client.CreateOptions{}); err != nil {
		return nil, err
	}
	return func() error {
		return cl.Delete(context.Background(), &pod, &client.DeleteOptions{})
	}, nil
}

func main() {
	cl, err := client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		fmt.Println("failed to create client")
		os.Exit(1)
	}
	f, err := deployPods(cl)
	if err != nil {
		fmt.Println("Deploy pod failed", err.Error())
		return
	}
	time.Sleep(10 * time.Second)
	showPods(cl)

	_ = f()

	time.Sleep(10 * time.Second)
	showPods(cl)

}
