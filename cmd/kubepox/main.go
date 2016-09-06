package main

import (
	"flag"
	"fmt"

	"github.com/aporeto-inc/kubepox"

	"k8s.io/kubernetes/pkg/api"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
)

var (
	kubeconfig = flag.String("kubeconfig", "/Users/bvandewa/.kube/config", "absolute path to the kubeconfig file")
)

func main() {
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Printf("test")
	}

	myClient, err := client.New(config)
	if err != nil {
		fmt.Printf("test")
	}

	kubepox.PrintPolicies(myClient)
	//kubepox.PrintPods(myClient)
	policy, _ := myClient.Extensions().NetworkPolicies("").List(api.ListOptions{})
	kubepox.ListPodsPerPolicy(myClient, &policy.Items[0])
}
