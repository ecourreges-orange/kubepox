package kubepox

import (
	"encoding/json"
	"fmt"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/labels"

	apiu "k8s.io/kubernetes/pkg/api/unversioned"
	client "k8s.io/kubernetes/pkg/client/unversioned"
)

// PrintPolicies pretty print all the policies
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

// PrintPods pretty print all the pods
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

	matchedPolicies := extensions.NetworkPolicyList{
		Items: []extensions.NetworkPolicy{},
	}
	podLabels := labels.Set(pod.GetLabels())

	allPolicies, err := c.Extensions().NetworkPolicies("").List(api.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Iterate over all policies and find the one that apply to the pod.
	for _, policy := range allPolicies.Items {
		selector, err := apiu.LabelSelectorAsSelector(&policy.Spec.PodSelector)
		if err != nil {
			return nil, err
		}
		if selector.Matches(podLabels) {
			matchedPolicies.Items = append(matchedPolicies.Items, policy)
		}
	}

	return &matchedPolicies, nil
}

// ListPodsPerPolicy returns all the Pods that are affected by a policy.
func ListPodsPerPolicy(c *client.Client, np *extensions.NetworkPolicy) (*api.PodList, error) {

	selector, err := apiu.LabelSelectorAsSelector(&np.Spec.PodSelector)
	if err != nil {
		return nil, err
	}

	// Match pods based on the Label Selector that came with the policy
	matchedPods, err := c.Pods("").List(api.ListOptions{LabelSelector: selector})
	if err != nil {
		return nil, err
	}
	return matchedPods, nil
}