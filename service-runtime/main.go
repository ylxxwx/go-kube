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

func showService(cl client.Client) {
	svcList := &corev1.ServiceList{}
	cl.List(context.Background(), svcList, client.InNamespace("default"))
	fmt.Printf("%-15s\t%-10s\t%-20s\n", "NAMESPACE", "STATUS", "NAME")
	for _, svc := range svcList.Items {
		fmt.Printf("%-15s\t%v\n", svc.Namespace, svc.Status)
	}
}

func deployService(cl client.Client) (func() error, error) {
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-svc",
			Namespace: "default",
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "def-port",
					Protocol: "TCP",
					Port:     21,
				},
			},
		},
	}

	if err := cl.Create(context.Background(), &svc, &client.CreateOptions{}); err != nil {
		return nil, err
	}
	return nil, nil
}

func patchService(cl client.Client) (func() error, error) {
	svc := corev1.Service{}

	if err := cl.Get(context.Background(), client.ObjectKey{
		Namespace: "default",
		Name:      "test-svc",
	}, &svc); err != nil {
		return nil, err
	}

	patch := client.MergeFrom(svc.DeepCopy())
	svc.Spec.Ports = append(svc.Spec.Ports,
		corev1.ServicePort{
			Name:     "http-service",
			Protocol: "TCP",
			Port:     80},
	)
	err := cl.Patch(context.Background(), &svc, patch)
	if err != nil {
		fmt.Printf("Patch failed. %s\n", err.Error())
	}

	return nil, nil
}

func main() {
	cl, err := client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		fmt.Println("failed to create client")
		os.Exit(1)
	}
	_, err = deployService(cl)
	if err != nil {
		fmt.Println("Deploy pod failed", err.Error())
		return
	}
	time.Sleep(1 * time.Second)
	showService(cl)

	patchService(cl)

	time.Sleep(1 * time.Second)
	showService(cl)

}
