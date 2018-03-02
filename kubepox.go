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
		// Validation of namespace
		if policy.Namespace != pod.Namespace {
			continue
		}

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
// returns nil if the policies in parameters are not applicable to Ingress
func ListIngressRulesPerPod(pod *api.Pod, allPolicies *networking.NetworkPolicyList) (*[]networking.NetworkPolicyIngressRule, error) {
	matchedPolicies, err := ListPoliciesPerPod(pod, allPolicies)
	if err != nil {
		return nil, err
	}
	return ingressSetGenerator(matchedPolicies)
}

// ListEgressRulesPerPod Generate a set of EgressRules that apply to the pod given in parameter.
// returns nil if the policies in parameters are not applicable to Egress
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
// returns nil if the policies in parameters are not applicable to Ingress
func ingressSetGenerator(policies *networking.NetworkPolicyList) (*[]networking.NetworkPolicyIngressRule, error) {
	ingressRules := []networking.NetworkPolicyIngressRule{}
	applicable := false

	for _, policy := range policies.Items {
		if IsPolicyApplicableToIngress(&policy) {
			applicable = true
			for _, singleRule := range policy.Spec.Ingress {
				ingressRules = append(ingressRules, singleRule)
			}
		}
	}
	if applicable {
		return &ingressRules, nil
	}
	return nil, nil
}

// generate a new table of IngressRules which are the Logical OR of all the existing IngressRules from all the policies given in parameter
// returns nil if the policies in parameters are not applicable to Egress
func egressSetGenerator(policies *networking.NetworkPolicyList) (*[]networking.NetworkPolicyEgressRule, error) {
	egressRules := []networking.NetworkPolicyEgressRule{}
	applicable := false

	for _, policy := range policies.Items {
		if IsPolicyApplicableToEgress(&policy) {
			applicable = true
			for _, singleRule := range policy.Spec.Egress {
				egressRules = append(egressRules, singleRule)
			}
		}
	}
	if applicable {
		return &egressRules, nil
	}
	return nil, nil
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

// IsPodSelected returns the selection status of the pod given as parameter over all the NetworkPolicies given as parameter.
// return status for ingress and egress
func IsPodSelected(pod *api.Pod, policies *networking.NetworkPolicyList) (bool, bool, error) {
	isApplicableToIngress := false
	isApplicableToEgress := false

	applicablePolicies, err := ListPoliciesPerPod(pod, policies)
	if err != nil {
		return false, false, nil
	}
	for _, policy := range applicablePolicies.Items {
		if IsPolicyApplicableToIngress(&policy) {
			isApplicableToIngress = true
		}
		if IsPolicyApplicableToEgress(&policy) {
			isApplicableToEgress = true
		}
		if isApplicableToIngress && isApplicableToEgress {
			return true, true, nil
		}
	}

	return isApplicableToIngress, isApplicableToEgress, nil
}

// IsPodSelectedIngress returns the selection status of the pod given as parameter over all the NetworkPolicies given as parameter.
// return status for ingress
func IsPodSelectedIngress(pod *api.Pod, policies *networking.NetworkPolicyList) (bool, error) {
	applicablePolicies, err := ListPoliciesPerPod(pod, policies)
	if err != nil {
		return false, nil
	}
	for _, policy := range applicablePolicies.Items {
		if IsPolicyApplicableToIngress(&policy) {
			return true, nil
		}
	}
	return false, nil
}

// IsPodSelectedEgress returns the selection status of the pod given as parameter over all the NetworkPolicies given as parameter.
// return status for egress
func IsPodSelectedEgress(pod *api.Pod, policies *networking.NetworkPolicyList) (bool, error) {
	applicablePolicies, err := ListPoliciesPerPod(pod, policies)
	if err != nil {
		return false, nil
	}
	for _, policy := range applicablePolicies.Items {
		if IsPolicyApplicableToEgress(&policy) {
			return true, nil
		}
	}
	return false, nil
}
