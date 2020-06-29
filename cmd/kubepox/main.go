package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/aporeto-inc/kubepox"

	"github.com/docopt/docopt-go"

	api "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
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

	arguments, _ := docopt.Parse(usage, nil, true, "KubePox", false)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5000)*time.Millisecond)

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

	myClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating REST Kube Client: %v\n", err)
		os.Exit(1)
	}

	// Display all policies. Similar to kubectl describe policies in json
	if arguments["get-all"].(bool) && arguments["policies"].(bool) {

		policies, err := myClient.NetworkingV1().NetworkPolicies(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Couldn't get Network Policy: %v\n", err)
			os.Exit(1)
		}
		renderPolicies(policies)
		os.Exit(0)
	}
	// Display all pods. Similar to kubectl describe pods in json
	if arguments["get-all"].(bool) && arguments["pods"].(bool) {

		pods, err := myClient.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
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
		np, err := myClient.NetworkingV1().NetworkPolicies(namespace).Get(ctx, arguments["<policy>"].(string), metav1.GetOptions{})
		if err != nil {
			fmt.Printf("Couldn't get Network Policy: %v\n", err)
			os.Exit(1)
		}
		allPods, err := myClient.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
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

		pod, err := myClient.CoreV1().Pods(namespace).Get(ctx, arguments["<pod>"].(string), metav1.GetOptions{})
		if err != nil {
			fmt.Printf("Couldn't get target pod %v\n", err)
			os.Exit(1)
		}

		allPolicies, err := myClient.NetworkingV1().NetworkPolicies(namespace).List(ctx, metav1.ListOptions{})
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

		pod, err := myClient.CoreV1().Pods(namespace).Get(ctx, arguments["<pod>"].(string), metav1.GetOptions{})
		if err != nil {
			fmt.Printf("Couldn't get target pod %v\n", err)
			os.Exit(1)
		}

		allPolicies, err := myClient.NetworkingV1().NetworkPolicies(namespace).List(ctx, metav1.ListOptions{})
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
	defer cancel()
}

func renderPolicies(policies *networking.NetworkPolicyList) {
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

func renderIngressRules(ingressRules *[]networking.NetworkPolicyIngressRule) {
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

func portsRepresentation(rule *networking.NetworkPolicyIngressRule) string {
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

func entryFromRule(rule *networking.NetworkPolicyIngressRule, ruleCount, entryCount int) (string, error) {
	entryString := ""
	entryString += strconv.Itoa(ruleCount+1) + "\t" + strconv.Itoa(entryCount+1) + "\t"

	selector, err := metav1.LabelSelectorAsSelector(rule.From[entryCount].PodSelector)
	if err != nil {
		return "", err
	}
	entryString += selector.String()
	entryString += "\t"
	entryString += portsRepresentation(rule)
	entryString += "\t\n"
	return entryString, nil
}

func renderIngressRulesHuman(ingressRules *[]networking.NetworkPolicyIngressRule) {
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
