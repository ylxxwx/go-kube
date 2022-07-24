package main

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	fooV1 "github.com/ylxxwx/go-kube/crd_3/pkg/apis/foogroup/v1alpha1"
	"github.com/ylxxwx/go-kube/crd_3/pkg/generated/clientset/versioned/scheme"
)

func main() {
	fmt.Println("Enter main")
	cc, err := client.New(config.GetConfigOrDie(), client.Options{Scheme: scheme.Scheme})
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	showFoo(cc)
}

func init() {
	fooV1.AddToScheme(scheme.Scheme)
}

func showFoo(cc client.Client) {
	foo := &fooV1.Foo{}
	err := cc.Get(context.Background(), client.ObjectKey{
		Namespace: "default",
		Name:      "example-foo",
	}, foo)
	if err != nil {
		fmt.Println("Get Foo return:", err.Error())
	}
	fmt.Printf("%s\t%s\n", foo.Name, foo.Namespace)
}
