package kubepox

import (
	"fmt"
	"testing"

	api "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

// pod2 is a target pod with role=backend
var pod2 = api.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name: "pod2",
		Labels: map[string]string{
			"role": "backend",
		},
	},
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
	}

	for i, test := range tests {
		t.Log("Testing ListPolicyPerPod Ingress ", i)
		result, err := ListPoliciesPerPod(&test.Pod, &test.Policies)
		if err != nil {
			t.Errorf("Error on ListPolicyPerPod for test %d", i)
		}

		if err := testNPListEquality(*result, test.Result); err != nil {
			t.Errorf("Error on  ListPolicyPerPod test %ds : %s ", i, err)
		}
	}

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
