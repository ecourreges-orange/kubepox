package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/aporeto-inc/kubepox"

	"github.com/docopt/docopt-go"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
)

//Todo: Make it clean and a real executable with flags.
var (
	kubeconfig = flag.String("kubeconfig", "/Users/bvandewa/.kube/config", "absolute path to the kubeconfig file")
)

func main() {

	usage := `

	Usage:
	kubepox [--config <config>][--namespace <namespace>] get (np|pod)
	kubepox [--config <config>][--namespace <namespace>] affect (np|pod) (<name>)

  Options:
	--namespace=NAMESPACE Namespace to run the query in
	--config=FILE path to the KubeConfig file.
	`

	arguments, _ := docopt.Parse(usage, nil, true, "Naval Fate 2.0", false)
	fmt.Println(arguments)

	var kubeconfig string
	if arguments["--config"] == nil {
		kubeconfig = os.Getenv("HOME") + "/.kube/config"
	} else {
		kubeconfig = arguments["--config"].(string)
	}

	var namespace string
	if arguments["--namespace"] == nil {
		namespace = ""
	} else {
		namespace = arguments["--namespace"].(string)
	}
	fmt.Println(namespace)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Printf("Error opening Kubeconfig: %v\n", err)
		os.Exit(1)
	}

	myClient, err := client.New(config)
	if err != nil {
		fmt.Printf("Error creating REST Kube Client: %v\n", err)
		os.Exit(1)
	}

	if arguments["get"].(bool) && arguments["np"].(bool) {
		kubepox.PrintPolicies(myClient)
	} else if arguments["get"].(bool) && arguments["pod"].(bool) {
		kubepox.PrintPods(myClient)
	} else if arguments["affect"].(bool) && arguments["np"].(bool) {
		// Get the Policy in argument
		np, err := myClient.Extensions().NetworkPolicies(namespace).Get(arguments["<name>"].(string))
		if err != nil {
			fmt.Printf("Couldn't get Network Policy: %v\n", err)
		}
		allPods, err := myClient.Pods(namespace).List(api.ListOptions{})
		if err != nil {
			fmt.Printf("Couldn't get all the pods %v\n", err)
		}
		matchedPods, err := kubepox.ListPodsPerPolicy(np, allPods)
		if err != nil {
			fmt.Printf("Error getting matching pods: %v\n", err)
		}
		fmt.Println("Resulting Pods:")

		renderPods(matchedPods)

	} else if arguments["affect"].(bool) && arguments["pod"].(bool) {
		pod, err := myClient.Pods(namespace).Get(arguments["<name>"].(string))
		if err != nil {
			fmt.Printf("Couldn't get target pod %v\n", err)
		}

		allPolicies, err := myClient.Extensions().NetworkPolicies(namespace).List(api.ListOptions{})
		if err != nil {
			fmt.Printf("Couldn't get all Network Policies: %v\n", err)
		}

		matchedPolicies, err := kubepox.ListPoliciesPerPod(pod, allPolicies)
		if err != nil {
			fmt.Printf("Error getting matching policies: %v\n", err)
		}
		fmt.Println("Resulting Policies:")
		renderPolicies(matchedPolicies)

	}

}

func renderPolicies(policies *extensions.NetworkPolicyList) {
	for _, policy := range policies.Items {
		fmt.Printf("\n\n\n\n")
		pp, _ := json.MarshalIndent(&policy, "", "   ")
		fmt.Println(string(pp))
	}
}

func renderPods(pods *api.PodList) {
	for _, pod := range pods.Items {
		fmt.Printf("\n\n\n\n")
		pp, _ := json.MarshalIndent(&pod, "", "   ")
		fmt.Println(string(pp))
	}
}
