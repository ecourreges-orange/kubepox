package main

import (
	"encoding/json"
	"flag"
	"fmt"

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

	listPolicies(myClient)

}

func listPolicies(myClient *client.Client) error {
	policies, err := myClient.Extensions().NetworkPolicies("").List(api.ListOptions{})
	if err != nil {
		fmt.Printf("test")
	}

	for _, policy := range policies.Items {
		fmt.Println("Existing policies:")
		pp, _ := json.MarshalIndent(&policy, "", "   ")
		fmt.Println(string(pp))
	}

	return nil
}

func podPerLabel(myClient *client.Client) error {
	return nil
}

func policiesPerPod(myClient *client.Client) error {
	return nil
}
