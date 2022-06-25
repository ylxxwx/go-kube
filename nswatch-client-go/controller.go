package main

import (
	corev1informer "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
)

type Controller struct {
	clientSet *kubernetes.Clientset
	informer  corev1informer.NamespaceInformer
}

func NewController(cs *kubernetes.Clientset, informer corev1informer.NamespaceInformer) *Controller {
	return &Controller{
		clientSet: cs,
		informer:  informer,
	}
}
