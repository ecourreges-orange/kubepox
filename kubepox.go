package kubepox

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/labels"

	apiu "k8s.io/kubernetes/pkg/api/unversioned"
)

// ListPoliciesPerPod returns all the NetworkPolicies that are associated with a pod.
func ListPoliciesPerPod(pod *api.Pod, allPolicies *extensions.NetworkPolicyList) (*extensions.NetworkPolicyList, error) {

	matchedPolicies := extensions.NetworkPolicyList{
		Items: []extensions.NetworkPolicy{},
	}
	podLabels := labels.Set(pod.GetLabels())

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

	selector, err := apiu.LabelSelectorAsSelector(&np.Spec.PodSelector)
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
