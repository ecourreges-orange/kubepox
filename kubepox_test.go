package kubepox

import (
	"fmt"
	"testing"

	api "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// np2 is egress only for target pods with role=frontend
var defaultdenyingress = networking.NetworkPolicy{
	ObjectMeta: metav1.ObjectMeta{
		Name: "defaultdenyingress",
	},
	Spec: networking.NetworkPolicySpec{
		PolicyTypes: []networking.PolicyType{
			networking.PolicyTypeIngress,
		},
	},
}

// np2 is egress only for target pods with role=frontend
var defaultallowingress = networking.NetworkPolicy{
	ObjectMeta: metav1.ObjectMeta{
		Name: "defaultallowingress",
	},
	Spec: networking.NetworkPolicySpec{
		Ingress: []networking.NetworkPolicyIngressRule{
			networking.NetworkPolicyIngressRule{},
		},
	},
}

// np2 is egress only for target pods with role=frontend
var defaultdenyegress = networking.NetworkPolicy{
	ObjectMeta: metav1.ObjectMeta{
		Name: "defaultdenyegress",
	},
	Spec: networking.NetworkPolicySpec{
		PolicyTypes: []networking.PolicyType{
			networking.PolicyTypeEgress,
		},
	},
}

// np2 is egress only for target pods with role=frontend
var defaultallowegress = networking.NetworkPolicy{
	ObjectMeta: metav1.ObjectMeta{
		Name: "defaultallowegress",
	},
	Spec: networking.NetworkPolicySpec{
		PolicyTypes: []networking.PolicyType{
			networking.PolicyTypeEgress,
		},
		Egress: []networking.NetworkPolicyEgressRule{
			networking.NetworkPolicyEgressRule{},
		},
	},
}

// np2 is egress only for target pods with role=frontend
var defaultdenyall = networking.NetworkPolicy{
	ObjectMeta: metav1.ObjectMeta{
		Name: "defaultdenyall",
	},
	Spec: networking.NetworkPolicySpec{
		PolicyTypes: []networking.PolicyType{
			networking.PolicyTypeEgress,
			networking.PolicyTypeIngress,
		},
	},
}

// np1 is ingress only for target pods with role=frontend
var np1 = networking.NetworkPolicy{
	ObjectMeta: metav1.ObjectMeta{
		Name: "np1",
	},
	Spec: networking.NetworkPolicySpec{
		PodSelector: metav1.LabelSelector{
			MatchLabels: map[string]string{
				"role": "frontend",
			},
		},
		Ingress: []networking.NetworkPolicyIngressRule{
			networking.NetworkPolicyIngressRule{
				From: []networking.NetworkPolicyPeer{
					networking.NetworkPolicyPeer{
						PodSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"role": "backend",
							},
						},
					},
				},
			},
		},
	},
}

// np1 is ingress only for target pods with role=frontend
var np1namespacex = networking.NetworkPolicy{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "np1",
		Namespace: "x",
	},
	Spec: networking.NetworkPolicySpec{
		PodSelector: metav1.LabelSelector{
			MatchLabels: map[string]string{
				"role": "frontend",
			},
		},
		Ingress: []networking.NetworkPolicyIngressRule{
			networking.NetworkPolicyIngressRule{
				From: []networking.NetworkPolicyPeer{
					networking.NetworkPolicyPeer{
						PodSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"role": "backend",
							},
						},
					},
				},
			},
		},
	},
}

// np2 is egress only for target pods with role=frontend
var np2 = networking.NetworkPolicy{
	ObjectMeta: metav1.ObjectMeta{
		Name: "np2",
	},
	Spec: networking.NetworkPolicySpec{
		PodSelector: metav1.LabelSelector{
			MatchLabels: map[string]string{
				"role": "frontend",
			},
		},
		Egress: []networking.NetworkPolicyEgressRule{
			networking.NetworkPolicyEgressRule{
				To: []networking.NetworkPolicyPeer{
					networking.NetworkPolicyPeer{
						PodSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"role": "backend",
							},
						},
					},
				},
			},
		},
		PolicyTypes: []networking.PolicyType{
			networking.PolicyTypeEgress,
		},
	},
}

// np3 is ingress and egress for target pods role=frontend
var np3 = networking.NetworkPolicy{
	ObjectMeta: metav1.ObjectMeta{
		Name: "np3",
	},
	Spec: networking.NetworkPolicySpec{
		PodSelector: metav1.LabelSelector{
			MatchLabels: map[string]string{
				"role": "frontend",
			},
		},
		Ingress: []networking.NetworkPolicyIngressRule{
			networking.NetworkPolicyIngressRule{
				From: []networking.NetworkPolicyPeer{
					networking.NetworkPolicyPeer{
						PodSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"role": "backend",
							},
						},
					},
				},
			},
		},
		Egress: []networking.NetworkPolicyEgressRule{
			networking.NetworkPolicyEgressRule{
				To: []networking.NetworkPolicyPeer{
					networking.NetworkPolicyPeer{
						PodSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"role": "backend",
							},
						},
					},
				},
			},
		},
	},
}

// pod1 is a target pod with role=frontend
var pod1 = api.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name: "pod1",
		Labels: map[string]string{
			"role": "frontend",
		},
	},
}

// pod1 is a target pod with role=frontend
var pod1namespacex = api.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pod1",
		Namespace: "x",
		Labels: map[string]string{
			"role": "frontend",
		},
	},
}

// pod2 is a target pod with role=backend
var pod2 = api.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name: "pod2",
		Labels: map[string]string{
			"role": "backend",
		},
	},
}

func TestListPoliciesPerPod(t *testing.T) {

	type testStruct struct {
		Policies networking.NetworkPolicyList
		Pod      api.Pod
		Result   networking.NetworkPolicyList
	}

	tests := []testStruct{
		testStruct{
			Policies: buildNetworkPolicyList(),
			Pod:      pod1,
			Result:   buildNetworkPolicyList(),
		},
		testStruct{
			Policies: buildNetworkPolicyList(defaultdenyall),
			Pod:      pod1,
			Result:   buildNetworkPolicyList(defaultdenyall),
		},
		testStruct{
			Policies: buildNetworkPolicyList(defaultdenyingress, defaultallowingress, defaultdenyegress, defaultdenyegress, defaultdenyall),
			Pod:      pod2,
			Result:   buildNetworkPolicyList(defaultdenyingress, defaultallowingress, defaultdenyegress, defaultdenyegress, defaultdenyall),
		},
		testStruct{
			Policies: buildNetworkPolicyList(np1),
			Pod:      pod1,
			Result:   buildNetworkPolicyList(np1),
		},
		testStruct{
			Policies: buildNetworkPolicyList(np1, np2),
			Pod:      pod1,
			Result:   buildNetworkPolicyList(np1, np2),
		},
		testStruct{
			Policies: buildNetworkPolicyList(np1, np2, np3),
			Pod:      pod1,
			Result:   buildNetworkPolicyList(np1, np2, np3),
		},
		testStruct{
			Policies: buildNetworkPolicyList(),
			Pod:      pod2,
			Result:   buildNetworkPolicyList(),
		},
		testStruct{
			Policies: buildNetworkPolicyList(np1),
			Pod:      pod2,
			Result:   buildNetworkPolicyList(),
		},
		testStruct{
			Policies: buildNetworkPolicyList(np1, np2),
			Pod:      pod2,
			Result:   buildNetworkPolicyList(),
		},
		testStruct{
			Policies: buildNetworkPolicyList(np1, np2, np3),
			Pod:      pod2,
			Result:   buildNetworkPolicyList(),
		},

		// different namespace tests
		testStruct{
			Policies: buildNetworkPolicyList(np1, np2, np3),
			Pod:      pod1namespacex,
			Result:   buildNetworkPolicyList(),
		},
		testStruct{
			Policies: buildNetworkPolicyList(np1namespacex),
			Pod:      pod1,
			Result:   buildNetworkPolicyList(),
		},
		testStruct{
			Policies: buildNetworkPolicyList(np1namespacex),
			Pod:      pod1namespacex,
			Result:   buildNetworkPolicyList(np1namespacex),
		},
	}

	for i, test := range tests {
		t.Log("Testing ListPolicyPerPod ", i)
		result, err := ListPoliciesPerPod(&test.Pod, &test.Policies)
		if err != nil {
			t.Errorf("Error on ListPolicyPerPod for test %d", i)
		}

		if err := testNPListEquality(*result, test.Result); err != nil {
			t.Errorf("Error on  ListPolicyPerPod test %ds : %s ", i, err)
		}
	}

}

func TestListIngressRulesPerPod(t *testing.T) {

	type testStruct struct {
		Policies networking.NetworkPolicyList
		Pod      api.Pod
		Result   []networking.NetworkPolicyIngressRule
	}

	tests := []testStruct{
		testStruct{
			Policies: buildNetworkPolicyList(defaultdenyingress),
			Pod:      pod1,
			Result:   []networking.NetworkPolicyIngressRule{},
		},
		testStruct{
			Policies: buildNetworkPolicyList(defaultallowingress),
			Pod:      pod1,
			Result: []networking.NetworkPolicyIngressRule{
				networking.NetworkPolicyIngressRule{},
			},
		},
		testStruct{
			Policies: buildNetworkPolicyList(defaultdenyingress, np1),
			Pod:      pod1,
			Result: []networking.NetworkPolicyIngressRule{
				np1.Spec.Ingress[0],
			},
		},
		testStruct{
			Policies: buildNetworkPolicyList(defaultallowingress, np1),
			Pod:      pod1,
			Result: []networking.NetworkPolicyIngressRule{
				networking.NetworkPolicyIngressRule{},
				np1.Spec.Ingress[0],
			},
		},

		// One of the policy not apploed for this specific pod.
		testStruct{
			Policies: buildNetworkPolicyList(defaultdenyingress, np1),
			Pod:      pod2,
			Result:   []networking.NetworkPolicyIngressRule{},
		},
		testStruct{
			Policies: buildNetworkPolicyList(defaultallowingress, np1),
			Pod:      pod2,
			Result: []networking.NetworkPolicyIngressRule{
				networking.NetworkPolicyIngressRule{},
			},
		},

		// Policies not applicatble to Ingress
		testStruct{
			Policies: buildNetworkPolicyList(defaultdenyegress),
			Pod:      pod1,
			Result:   nil,
		},
		testStruct{
			Policies: buildNetworkPolicyList(defaultallowegress),
			Pod:      pod1,
			Result:   nil,
		},

		// Mix of applicable Ingress and non-applicable Egress Policies
		testStruct{
			Policies: buildNetworkPolicyList(defaultdenyegress, np1),
			Pod:      pod1,
			Result: []networking.NetworkPolicyIngressRule{
				np1.Spec.Ingress[0],
			},
		},
		testStruct{
			Policies: buildNetworkPolicyList(defaultallowegress, np1),
			Pod:      pod1,
			Result: []networking.NetworkPolicyIngressRule{
				np1.Spec.Ingress[0],
			},
		},

		// Mix of non-applicable Ingress and non-applicable Egress Policies
		testStruct{
			Policies: buildNetworkPolicyList(defaultdenyegress, np1),
			Pod:      pod2,
			Result:   nil,
		},
		testStruct{
			Policies: buildNetworkPolicyList(defaultallowegress, np1),
			Pod:      pod2,
			Result:   nil,
		},

		// Mix of non-applicable Ingress and non-applicable Egress Policies
		testStruct{
			Policies: buildNetworkPolicyList(defaultdenyall),
			Pod:      pod1,
			Result:   []networking.NetworkPolicyIngressRule{},
		},
		testStruct{
			Policies: buildNetworkPolicyList(defaultdenyall, np1),
			Pod:      pod1,
			Result: []networking.NetworkPolicyIngressRule{
				np1.Spec.Ingress[0],
			},
		},
		// Mix of non-applicable Ingress and non-applicable Egress Policies
		testStruct{
			Policies: buildNetworkPolicyList(defaultdenyall),
			Pod:      pod2,
			Result:   []networking.NetworkPolicyIngressRule{},
		},
		testStruct{
			Policies: buildNetworkPolicyList(defaultdenyall, np1),
			Pod:      pod2,
			Result:   []networking.NetworkPolicyIngressRule{},
		},
	}

	for i, test := range tests {
		t.Log("Testing ListPolicyPerPod Ingress ", i)
		result, err := ListIngressRulesPerPod(&test.Pod, &test.Policies)
		if err != nil {
			t.Errorf("Error on ListPolicyPerPod for test %d", i)
		}

		if result == nil && test.Result == nil {
			continue
		}

		if err := testNPIngressRuleListEquality(*result, test.Result); err != nil {
			t.Errorf("Error on  ListPolicyPerPod test %ds : %s ", i, err)
		}
	}

}

func TestListEgressRulesPerPod(t *testing.T) {

	type testStruct struct {
		Policies networking.NetworkPolicyList
		Pod      api.Pod
		Result   []networking.NetworkPolicyEgressRule
	}

	tests := []testStruct{
		testStruct{
			Policies: buildNetworkPolicyList(defaultdenyegress),
			Pod:      pod1,
			Result:   []networking.NetworkPolicyEgressRule{},
		},
		testStruct{
			Policies: buildNetworkPolicyList(defaultallowegress),
			Pod:      pod1,
			Result: []networking.NetworkPolicyEgressRule{
				networking.NetworkPolicyEgressRule{},
			},
		},
		testStruct{
			Policies: buildNetworkPolicyList(defaultdenyegress, np3),
			Pod:      pod1,
			Result: []networking.NetworkPolicyEgressRule{
				np3.Spec.Egress[0],
			},
		},
		testStruct{
			Policies: buildNetworkPolicyList(defaultallowegress, np3),
			Pod:      pod1,
			Result: []networking.NetworkPolicyEgressRule{
				networking.NetworkPolicyEgressRule{},
				np3.Spec.Egress[0],
			},
		},

		// One of the policy not apploed for this specific pod.
		testStruct{
			Policies: buildNetworkPolicyList(defaultdenyegress, np3),
			Pod:      pod2,
			Result:   []networking.NetworkPolicyEgressRule{},
		},
		testStruct{
			Policies: buildNetworkPolicyList(defaultallowegress, np3),
			Pod:      pod2,
			Result: []networking.NetworkPolicyEgressRule{
				networking.NetworkPolicyEgressRule{},
			},
		},

		// Policies not applicatble to Ingress
		testStruct{
			Policies: buildNetworkPolicyList(defaultdenyingress),
			Pod:      pod1,
			Result:   nil,
		},
		testStruct{
			Policies: buildNetworkPolicyList(defaultallowingress),
			Pod:      pod1,
			Result:   nil,
		},

		// Mix of applicable Ingress and non-applicable Egress Policies
		testStruct{
			Policies: buildNetworkPolicyList(defaultdenyingress, np3),
			Pod:      pod1,
			Result: []networking.NetworkPolicyEgressRule{
				np3.Spec.Egress[0],
			},
		},
		testStruct{
			Policies: buildNetworkPolicyList(defaultallowingress, np3),
			Pod:      pod1,
			Result: []networking.NetworkPolicyEgressRule{
				np3.Spec.Egress[0],
			},
		},

		// Mix of non-applicable Ingress and non-applicable Egress Policies
		testStruct{
			Policies: buildNetworkPolicyList(defaultdenyingress, np1),
			Pod:      pod2,
			Result:   nil,
		},
		testStruct{
			Policies: buildNetworkPolicyList(defaultallowingress, np1),
			Pod:      pod2,
			Result:   nil,
		},

		// Mix of non-applicable Ingress and non-applicable Egress Policies
		testStruct{
			Policies: buildNetworkPolicyList(defaultdenyall),
			Pod:      pod1,
			Result:   []networking.NetworkPolicyEgressRule{},
		},
		testStruct{
			Policies: buildNetworkPolicyList(defaultdenyall, np3),
			Pod:      pod1,
			Result: []networking.NetworkPolicyEgressRule{
				np3.Spec.Egress[0],
			},
		},
		// Mix of non-applicable Ingress and non-applicable Egress Policies
		testStruct{
			Policies: buildNetworkPolicyList(defaultdenyall),
			Pod:      pod2,
			Result:   []networking.NetworkPolicyEgressRule{},
		},
		testStruct{
			Policies: buildNetworkPolicyList(defaultdenyall, np3),
			Pod:      pod2,
			Result:   []networking.NetworkPolicyEgressRule{},
		},
	}

	for i, test := range tests {
		t.Log("Testing ListPolicyPerPod Egress ", i)
		result, err := ListEgressRulesPerPod(&test.Pod, &test.Policies)
		if err != nil {
			t.Errorf("Error on ListPolicyPerPod Egress for test %d", i)
		}

		if result == nil && test.Result == nil {
			continue
		}
		if result != nil && test.Result == nil {
			t.Errorf("Issue 1 %d", i)
		}
		if result == nil && test.Result != nil {
			t.Errorf("Issue 2 %d", i)
		}

		if err := testNPEgressRuleListEquality(*result, test.Result); err != nil {
			t.Errorf("Error on  ListPolicyPerPod test %ds : %s ", i, err)
		}
	}

}

func TestListPodsPerPolicy(t *testing.T) {
	type testStruct struct {
		Policy networking.NetworkPolicy
		Pods   api.PodList
		Result api.PodList
	}

	tests := []testStruct{
		testStruct{
			Policy: np1,
			Pods:   buildPodList(pod1),
			Result: buildPodList(pod1),
		},
		testStruct{
			Policy: np1,
			Pods:   buildPodList(pod2),
			Result: buildPodList(),
		},
		testStruct{
			Policy: np1,
			Pods:   buildPodList(pod1, pod2),
			Result: buildPodList(pod1),
		},
		testStruct{
			Policy: defaultdenyall,
			Pods:   buildPodList(pod1, pod2),
			Result: buildPodList(pod1, pod2),
		},
		testStruct{
			Policy: defaultdenyall,
			Pods:   buildPodList(),
			Result: buildPodList(),
		},
	}

	for i, test := range tests {
		t.Log("Testing ListPodsPerPolicy ", i)
		result, err := ListPodsPerPolicy(&test.Policy, &test.Pods)
		if err != nil {
			t.Errorf("Error on ListPodsPerPolicy for test %d", i)
		}

		if err := testPodListEquality(*result, test.Result); err != nil {
			t.Errorf("Error on ListPodsPerPolicy test %ds : %s ", i, err)
		}
	}

}

func buildNetworkPolicyList(nps ...networking.NetworkPolicy) networking.NetworkPolicyList {
	return networking.NetworkPolicyList{
		Items: nps,
	}
}

func buildPodList(pods ...api.Pod) api.PodList {
	return api.PodList{
		Items: pods,
	}
}

func testNPIngressRuleListEquality(resultList, expectedList []networking.NetworkPolicyIngressRule) error {
	//fmt.Printf("RESULT: %+v, \n EXPECTED: %+v \n", resultList, expectedList)
	if len(resultList) != len(expectedList) {
		return fmt.Errorf("Got %d element, expected %d element", len(resultList), len(expectedList))
	}

	for i, expect := range expectedList {
		result := resultList[i]
		results := result.String()
		expects := expect.String()
		if results != expects {
			return fmt.Errorf("Rule %d Got %s , expected %s", i, results, expects)
		}
	}

	return nil
}

func testNPEgressRuleListEquality(resultList, expectedList []networking.NetworkPolicyEgressRule) error {
	if len(resultList) != len(expectedList) {
		return fmt.Errorf("Got %d element, expected %d element", len(resultList), len(expectedList))
	}

	for i, expect := range expectedList {
		result := resultList[i]
		results := result.String()
		expects := expect.String()
		if results != expects {
			return fmt.Errorf("Rule %d Got %s , expected %s", i, results, expects)
		}
	}

	return nil
}

func testNPListEquality(result, expected networking.NetworkPolicyList) error {
	if len(result.Items) != len(expected.Items) {
		return fmt.Errorf("Got %d element, expected %d element", len(result.Items), len(expected.Items))
	}

MainLoop1:
	for _, expectedPolicy := range expected.Items {
		for _, resultPolicy := range result.Items {
			if expectedPolicy.Name == resultPolicy.Name {
				continue MainLoop1
			}
		}
		return fmt.Errorf("Couldn't find expected np %s element in result", expectedPolicy.Name)
	}

MainLoop2:
	for _, resultPolicy := range result.Items {
		for _, expectedPolicy := range expected.Items {
			if expectedPolicy.Name == resultPolicy.Name {
				continue MainLoop2
			}
		}
		return fmt.Errorf("Couldn't find result np %s element in result", resultPolicy.Name)
	}

	return nil
}

func testPodListEquality(result, expected api.PodList) error {
	if len(result.Items) != len(expected.Items) {
		return fmt.Errorf("Got %d element, expected %d element", len(result.Items), len(expected.Items))
	}

MainLoop1:
	for _, expectedPod := range expected.Items {
		for _, resultPod := range result.Items {
			if expectedPod.Name == resultPod.Name {
				continue MainLoop1
			}
		}
		return fmt.Errorf("Couldn't find expected np %s element in result", expectedPod.Name)
	}

MainLoop2:
	for _, resultPod := range result.Items {
		for _, expectedPod := range expected.Items {
			if expectedPod.Name == resultPod.Name {
				continue MainLoop2
			}
		}
		return fmt.Errorf("Couldn't find result np %s element in result", resultPod.Name)
	}

	return nil
}
