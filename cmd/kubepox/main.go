package main

import (
	"flag"
	"fmt"

	"github.com/aporeto-inc/kubepox"

	"k8s.io/kubernetes/pkg/api"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
)

//Todo: Make it clean and a real executable with flags.
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
	//policy, _ := myClient.Extensions().NetworkPolicies("").List(api.ListOptions{})
	//matchedPods, _ := kubepox.ListPodsPerPolicy(myClient, &policy.Items[0])
	//fmt.Printf("%+v\n", matchedPods)
	pods, _ := myClient.Pods("default").List(api.ListOptions{})
	//fmt.Printf("%+v\n", pods)
	podToTest := pods.Items[1]
	fmt.Println("Testing policies for Pod: " + podToTest.GetName())
	listOfPolicies, _ := kubepox.ListPoliciesPerPod(myClient, &pods.Items[1])
	fmt.Printf("\n\n%+v\n", listOfPolicies)
}
