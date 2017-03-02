package kubepox

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	api "k8s.io/client-go/pkg/api/v1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

// ListPoliciesPerPod returns all the NetworkPolicies that are associated with a pod.
func ListPoliciesPerPod(pod *api.Pod, allPolicies *extensions.NetworkPolicyList) (*extensions.NetworkPolicyList, error) {

	matchedPolicies := extensions.NetworkPolicyList{
		Items: []extensions.NetworkPolicy{},
	}
	podLabels := labels.Set(pod.GetLabels())
	fmt.Printf("Labels: %+v \n\n\n", pod.GetLabels())

	// Iterate over all policies and find the one that apply to the pod.
	for _, policy := range allPolicies.Items {
		fmt.Printf("Policy: %+v \n\n\n", policy)
		selector, err := metav1.LabelSelectorAsSelector(&policy.Spec.PodSelector)
		if err != nil {
			return nil, err
		}
		if selector.Matches(podLabels) {
			matchedPolicies.Items = append(matchedPolicies.Items, policy)
		}
	}

	return &matchedPolicies, nil
}

// ListIngressRulesPerPod Generate a set of IngressRules that apply to the pod given in parameter.
func ListIngressRulesPerPod(pod *api.Pod, allPolicies *extensions.NetworkPolicyList) (*[]extensions.NetworkPolicyIngressRule, error) {
	matchedPolicies, err := ListPoliciesPerPod(pod, allPolicies)
	if err != nil {
		return nil, err
	}
	return ingressSetGenerator(matchedPolicies)
}

// ListPodsPerPolicy returns all the Pods that are affected by a policy out of the list.
func ListPodsPerPolicy(np *extensions.NetworkPolicy, allPods *api.PodList) (*api.PodList, error) {

	selector, err := metav1.LabelSelectorAsSelector(&np.Spec.PodSelector)
	if err != nil {
		return nil, err
	}

	matchedPods := api.PodList{
		Items: []api.Pod{},
	}

	// Match pods based on the Label Selector that came with the policy
	for _, pod := range allPods.Items {
		if selector.Matches(labels.Set(pod.GetLabels())) {
			matchedPods.Items = append(matchedPods.Items, pod)
		}
	}

	return &matchedPods, nil
}

// generate a new table of IngressRules which are the Logical OR of all the existing IngressRules from all the policies given in parameter
func ingressSetGenerator(policies *extensions.NetworkPolicyList) (*[]extensions.NetworkPolicyIngressRule, error) {
	ingressRules := []extensions.NetworkPolicyIngressRule{}
	for _, policy := range policies.Items {
		for _, singleRule := range policy.Spec.Ingress {
			ingressRules = append(ingressRules, singleRule)
		}
	}
	return &ingressRules, nil
}
