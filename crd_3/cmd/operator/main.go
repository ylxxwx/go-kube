package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	//kubeinformers "k8s.io/client-go/informers"
	"github.com/ylxxwx/go-kube/crd_3/pkg/controller"
	"github.com/ylxxwx/go-kube/crd_3/pkg/generated/clientset/versioned"
	fooInformers "github.com/ylxxwx/go-kube/crd_3/pkg/generated/informers/externalversions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var kubeconfig string

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "/home/kmt/.kube/config", "kubeconfig file")
}

func main() {
	fmt.Println("Enter main")
	flag.Parse()

	stopCh := SetupSignalHandler()
	config, err := rest.InClusterConfig()
	//config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	fcs, err := versioned.NewForConfig(config)
	cs, err := kubernetes.NewForConfig(config)

	kubeInformerFactory := fooInformers.NewSharedInformerFactory(fcs, time.Second*30)

	if err != nil {
		panic(err)
	}

	ctl := controller.NewController(cs, fcs, kubeInformerFactory.Foogroup().V1alpha1().Foos())
	kubeInformerFactory.Start(stopCh)

	if err = ctl.Run(2, stopCh); err != nil {
		//klog.Fatalf("Error running controller: %s", err.Error())
	}

}

var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}
var onlyOneSignalHandler = make(chan struct{})

func SetupSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler) // panics when called twice

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}
