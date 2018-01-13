package kubepox

import (
	api "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// ListPoliciesPerPod returns all the NetworkPolicies that are associated with a pod.
func ListPoliciesPerPod(pod *api.Pod, allPolicies *networking.NetworkPolicyList) (*networking.NetworkPolicyList, error) {
	matchedPolicies := networking.NetworkPolicyList{
		Items: []networking.NetworkPolicy{},
	}
	podLabels := labels.Set(pod.GetLabels())

	// Iterate over all policies and find the one that apply to the pod.
	for _, policy := range allPolicies.Items {
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
func ListIngressRulesPerPod(pod *api.Pod, allPolicies *networking.NetworkPolicyList) (*[]networking.NetworkPolicyIngressRule, error) {
	matchedPolicies, err := ListPoliciesPerPod(pod, allPolicies)
	if err != nil {
		return nil, err
	}
	return ingressSetGenerator(matchedPolicies)
}

// ListEgressRulesPerPod Generate a set of EgressRules that apply to the pod given in parameter.
func ListEgressRulesPerPod(pod *api.Pod, allPolicies *networking.NetworkPolicyList) (*[]networking.NetworkPolicyEgressRule, error) {
	matchedPolicies, err := ListPoliciesPerPod(pod, allPolicies)
	if err != nil {
		return nil, err
	}
	return egressSetGenerator(matchedPolicies)
}

// ListPodsPerPolicy returns all the Pods that are affected by a policy out of the list.
func ListPodsPerPolicy(np *networking.NetworkPolicy, allPods *api.PodList) (*api.PodList, error) {

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
func ingressSetGenerator(policies *networking.NetworkPolicyList) (*[]networking.NetworkPolicyIngressRule, error) {
	ingressRules := []networking.NetworkPolicyIngressRule{}

	for _, policy := range policies.Items {
		if IsPolicyApplicableToIngress(&policy) {
			for _, singleRule := range policy.Spec.Ingress {
				ingressRules = append(ingressRules, singleRule)
			}
		}
	}
	return &ingressRules, nil
}

// generate a new table of IngressRules which are the Logical OR of all the existing IngressRules from all the policies given in parameter
func egressSetGenerator(policies *networking.NetworkPolicyList) (*[]networking.NetworkPolicyEgressRule, error) {
	egressRules := []networking.NetworkPolicyEgressRule{}

	for _, policy := range policies.Items {
		if IsPolicyApplicableToEgress(&policy) {
			for _, singleRule := range policy.Spec.Egress {
				egressRules = append(egressRules, singleRule)
			}
		}
	}
	return &egressRules, nil
}

// IsPolicyApplicableToIngress returns true if the policy is applicable for Ingress traffic
func IsPolicyApplicableToIngress(policy *networking.NetworkPolicy) bool {

	// Logic: Policy applies to ingress only IF:
	// - flag is not set
	// - flag is set with an entry to type Ingress (even if no section Egress exists)

	if policy.Spec.PolicyTypes == nil {
		return true
	}

	for _, ptype := range policy.Spec.PolicyTypes {
		if ptype == networking.PolicyTypeIngress {
			return true
		}
	}

	return false
}

// IsPolicyApplicableToEgress returns true if the policy is applicable for Egress traffic
func IsPolicyApplicableToEgress(policy *networking.NetworkPolicy) bool {

	// Logic: Policy applies to egress only IF:
	// - flag is not set but egress section is present
	// - flag is set with an entry to type Egress (even if no section Egress exists)

	if policy.Spec.PolicyTypes == nil {
		if policy.Spec.Egress != nil {
			return true
		}
		return false
	}

	for _, ptype := range policy.Spec.PolicyTypes {
		if ptype == networking.PolicyTypeEgress {
			return true
		}
	}

	return false
}
