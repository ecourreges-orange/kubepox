package kubepox

import (
	"encoding/json"
	"fmt"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/labels"

	client "k8s.io/kubernetes/pkg/client/unversioned"
)

func PrintPolicies(c *client.Client) error {
	policies, err := c.Extensions().NetworkPolicies("").List(api.ListOptions{})
	if err != nil {
		return err
	}

	for _, policy := range policies.Items {
		fmt.Println("Existing policies:")
		pp, _ := json.MarshalIndent(&policy, "", "   ")
		fmt.Println(string(pp))
	}

	return nil
}

func PrintPods(c *client.Client) error {
	pods, err := c.Pods("").List(api.ListOptions{})
	if err != nil {
		return err
	}

	for _, pod := range pods.Items {
		fmt.Println("Existing pods:")
		pp, _ := json.MarshalIndent(&pod, "", "   ")
		fmt.Println(string(pp))
	}

	return nil
}

// ListPoliciesPerPod returns all the NetworkPolicies that are associated with a pod.
func ListPoliciesPerPod(c *client.Client, pod *api.Pod) (*extensions.NetworkPolicyList, error) {
	return nil, nil
}

// ListPodsPerPolicy returns all the Pods that are affected by a policy.
func ListPodsPerPolicy(c *client.Client, np *extensions.NetworkPolicy) (*api.PodList, error) {
	labels := labels.Set(np.Spec.PodSelector.MatchLabels)
	selector := np.Spec.PodSelector.MatchExpressions
	fmt.Printf("%+v\n", labels)
	fmt.Printf("%+v\n", selector)
	return nil, nil

}
