package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/aporeto-inc/kubepox"

	"github.com/docopt/docopt-go"
	"k8s.io/kubernetes/pkg/api"
	apiu "k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
	client "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
)

//Todo: Make it clean and a real executable with flags.
var (
	kubeconfig = flag.String("kubeconfig", "/Users/bvandewa/.kube/config", "absolute path to the kubeconfig file")
)

func main() {

	usage := `

	Usage:
  kubepox [--config <config>][--namespace <namespace>] get-all (policies|pods)
  kubepox [--config <config>][--namespace <namespace>] get-pods <policy>
  kubepox [--config <config>][--namespace <namespace>] get-policies <pod>
  kubepox [--config <config>][--namespace <namespace>] get-rules <pod> [human]

  Options:
	--namespace=NAMESPACE Namespace to run the query in
	--config=FILE path to the KubeConfig file.
	`

	arguments, _ := docopt.Parse(usage, nil, true, "Naval Fate 2.0", false)

	// Get location of the Kubeconfig file. By default in your home.
	var kubeconfig string
	if arguments["--config"] == nil {
		kubeconfig = os.Getenv("HOME") + "/.kube/config"
	} else {
		kubeconfig = arguments["--config"].(string)
	}

	// Get namespace, by default it will be "default"
	var namespace string
	if arguments["--namespace"] == nil {
		namespace = "default"
	} else {
		namespace = arguments["--namespace"].(string)
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Printf("Error opening Kubeconfig: %v\n", err)
		os.Exit(1)
	}

	myClient, err := client.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating REST Kube Client: %v\n", err)
		os.Exit(1)
	}

	// Display all policies. Similar to kubectl describe policies in json
	if arguments["get-all"].(bool) && arguments["policies"].(bool) {

		policies, err := myClient.Extensions().NetworkPolicies(namespace).List(api.ListOptions{})
		if err != nil {
			fmt.Printf("Couldn't get Network Policy: %v\n", err)
			os.Exit(1)
		}
		renderPolicies(policies)
		os.Exit(0)
	}
	// Display all pods. Similar to kubectl describe pods in json
	if arguments["get-all"].(bool) && arguments["pods"].(bool) {

		pods, err := myClient.Pods(namespace).List(api.ListOptions{})
		if err != nil {
			fmt.Printf("Couldn't get all the pods %v\n", err)
			os.Exit(1)
		}
		renderPods(pods)
		os.Exit(0)
	}

	// Get all the pods that get affected by the policy
	if arguments["get-pods"].(bool) {
		// Get the Policy in argument
		np, err := myClient.Extensions().NetworkPolicies(namespace).Get(arguments["<policy>"].(string))
		if err != nil {
			fmt.Printf("Couldn't get Network Policy: %v\n", err)
			os.Exit(1)
		}
		allPods, err := myClient.Pods(namespace).List(api.ListOptions{})
		if err != nil {
			fmt.Printf("Couldn't get all the pods %v\n", err)
			os.Exit(1)
		}
		matchedPods, err := kubepox.ListPodsPerPolicy(np, allPods)
		if err != nil {
			fmt.Printf("Error getting matching pods: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Matched pods for policy %s :\n", np.Name)
		renderPods(matchedPods)
		os.Exit(0)
	}

	// Get all the policies that get applied to a Pod.
	if arguments["get-policies"].(bool) {

		pod, err := myClient.Pods(namespace).Get(arguments["<pod>"].(string))
		if err != nil {
			fmt.Printf("Couldn't get target pod %v\n", err)
			os.Exit(1)
		}

		allPolicies, err := myClient.Extensions().NetworkPolicies(namespace).List(api.ListOptions{})
		if err != nil {
			fmt.Printf("Couldn't get all Network Policies: %v\n", err)
		}

		matchedPolicies, err := kubepox.ListPoliciesPerPod(pod, allPolicies)
		if err != nil {
			fmt.Printf("Error getting matching policies: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Applied policies for pod %s :\n", pod.Name)
		renderPolicies(matchedPolicies)
		os.Exit(0)
	}

	// Get all the IngressRules that get applied to a Pod.
	if arguments["get-rules"].(bool) {

		pod, err := myClient.Pods(namespace).Get(arguments["<pod>"].(string))
		if err != nil {
			fmt.Printf("Couldn't get target pod %v\n", err)
			os.Exit(1)
		}

		allPolicies, err := myClient.Extensions().NetworkPolicies(namespace).List(api.ListOptions{})
		if err != nil {
			fmt.Printf("Couldn't get all Network Policies: %v\n", err)
			os.Exit(1)
		}

		matchedRules, err := kubepox.ListIngressRulesPerPod(pod, allPolicies)
		if err != nil {
			fmt.Printf("Couldn't get all the rules: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("WhiteList for pod %s :\n\n", pod.Name)
		if arguments["human"].(bool) {
			renderIngressRulesHuman(matchedRules)
			os.Exit(0)
		}
		renderIngressRules(matchedRules)

	}

}

func renderPolicies(policies *extensions.NetworkPolicyList) {
	for count, policy := range policies.Items {
		fmt.Printf("POLICY %d\n", count+1)
		pp, _ := json.MarshalIndent(&policy, "", "   ")
		fmt.Println(string(pp))
	}
}

func renderPods(pods *api.PodList) {
	for count, pod := range pods.Items {
		fmt.Printf("POD %d\n", count+1)
		pp, _ := json.MarshalIndent(&pod, "", "   ")
		fmt.Println(string(pp))
	}
}

func renderIngressRules(ingressRules *[]extensions.NetworkPolicyIngressRule) {
	for count, rule := range *ingressRules {
		fmt.Printf("RULE %d\n", count)
		pp, _ := json.MarshalIndent(&rule, "", "   ")
		fmt.Println(string(pp))
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func portsRepresentation(rule *extensions.NetworkPolicyIngressRule) string {
	if len(rule.Ports) == 0 {
		return "ALL"
	}
	entryString := ""
	for count, port := range rule.Ports {
		entryString += string(*port.Protocol)
		entryString += ":"
		entryString += port.Port.String()
		if count == len(rule.Ports)-1 {
			break
		}
		entryString += ", "
	}
	return entryString
}

func entryFromRule(rule *extensions.NetworkPolicyIngressRule, ruleCount, entryCount int) (string, error) {
	entryString := ""
	entryString += strconv.Itoa(ruleCount+1) + "\t" + strconv.Itoa(entryCount+1) + "\t"

	selector, err := apiu.LabelSelectorAsSelector(rule.From[entryCount].PodSelector)
	if err != nil {
		return "", err
	}
	entryString += selector.String()
	entryString += "\t"
	entryString += portsRepresentation(rule)
	entryString += "\t\n"
	return entryString, nil
}

func renderIngressRulesHuman(ingressRules *[]extensions.NetworkPolicyIngressRule) {
	w := tabwriter.NewWriter(os.Stdout, 10, 0, 3, '-', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "RULE\tSELECTOR\tFROM PODS\tALLOWED TRAFFIC\t")
	for ruleCount, rule := range *ingressRules {
		for entryCount := 0; entryCount < len(rule.From); entryCount++ {
			entryString, err := entryFromRule(&rule, ruleCount, entryCount)
			if err != nil {
				fmt.Println("error while trying to render")
				os.Exit(1)
			}
			fmt.Fprint(w, entryString)
		}
	}
	w.Flush()
}
