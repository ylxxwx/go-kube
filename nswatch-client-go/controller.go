package main

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	corev1informer "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
)

const controllerAgentName = "myctl"

type Controller struct {
	clientSet *kubernetes.Clientset
	nsLister  corelisters.NamespaceLister
	nsSynced  cache.InformerSynced
	workqueue workqueue.RateLimitingInterface
	recorder  record.EventRecorder
}

func NewController(cs *kubernetes.Clientset, informer corev1informer.NamespaceInformer) *Controller {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: cs.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(new interface{}) {
			ns, ok := new.(*corev1.Namespace)
			if ok {
				fmt.Printf("Create ns:%s, ver:%s\n", ns.Name, ns.ResourceVersion)
			}
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				fmt.Println("EventHandler AddFunc, key:", key)
			}
		},
		UpdateFunc: func(old, new interface{}) {
			newNs := new.(*corev1.Namespace)
			oldNs := old.(*corev1.Namespace)
			fmt.Printf("EventHandler UpdateFunc,%s new v:%s, old v:%s\n", newNs.Name, newNs.ResourceVersion, oldNs.ResourceVersion)
			new.(*corev1.Namespace).Spec.Finalizers = []corev1.FinalizerName{}
		},
		DeleteFunc: func(new interface{}) {
			newNs := new.(*corev1.Namespace)
			fmt.Println("EventHandler DelFunc:", newNs.Name)
		},
	})

	return &Controller{
		clientSet: cs,
		nsLister:  informer.Lister(),
		nsSynced:  informer.Informer().HasSynced,
		workqueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "NSs"),
		recorder:  recorder,
	}
}

func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	//defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	klog.Info("Starting Foo controller")

	// Wait for the caches to be synced before starting workers
	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.nsSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	// Launch two workers to process Foo resources
	//for i := 0; i < workers; i++ {
	//	go wait.Until(c.runWorker, time.Second, stopCh)
	//}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}
