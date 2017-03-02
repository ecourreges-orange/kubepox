package kubepox

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	api "k8s.io/client-go/pkg/api/v1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

// pod1 has one single label: "role": "WebFrontend"
const pod1s = `{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "creationTimestamp": "2017-02-28T19:18:52Z",
        "labels": {
            "role": "WebFrontend"
        },
        "name": "frontend",
        "namespace": "demo",
        "resourceVersion": "543485",
        "selfLink": "/api/v1/namespaces/demo/pods/frontend",
        "uid": "bd9c01fb-fdea-11e6-91a2-42010a8001b8"
    },
    "spec": {
        "containers": [
            {
                "image": "redis",
                "imagePullPolicy": "Always",
                "name": "redismaster",
                "ports": [
                    {
                        "containerPort": 6379,
                        "protocol": "TCP"
                    }
                ],
                "resources": {
                    "requests": {
                        "cpu": "100m",
                        "memory": "100Mi"
                    }
                },
                "terminationMessagePath": "/dev/termination-log",
                "volumeMounts": [
                    {
                        "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                        "name": "default-token-kb5rn",
                        "readOnly": true
                    }
                ]
            }
        ],
        "dnsPolicy": "ClusterFirst",
        "nodeName": "gke-test48-default-pool-1e87d955-mwh5",
        "restartPolicy": "Always",
        "securityContext": {},
        "serviceAccount": "default",
        "serviceAccountName": "default",
        "terminationGracePeriodSeconds": 30,
        "volumes": [
            {
                "name": "default-token-kb5rn",
                "secret": {
                    "defaultMode": 420,
                    "secretName": "default-token-kb5rn"
                }
            }
        ]
    },
    "status": {
        "conditions": [
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2017-02-28T19:18:52Z",
                "status": "True",
                "type": "Initialized"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2017-02-28T19:18:54Z",
                "status": "True",
                "type": "Ready"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2017-02-28T19:18:52Z",
                "status": "True",
                "type": "PodScheduled"
            }
        ],
        "containerStatuses": [
            {
                "containerID": "docker://48d8e102c086b5686abde297d983c3278af548a9551dc84facb8e77f5944bc75",
                "image": "redis",
                "imageID": "docker://sha256:1a8a9ee54eb755a427e00484a64fa4edbe9c2abe59ca468a4e452b343a2b57c2",
                "lastState": {},
                "name": "redismaster",
                "ready": true,
                "restartCount": 0,
                "state": {
                    "running": {
                        "startedAt": "2017-02-28T19:18:53Z"
                    }
                }
            }
        ],
        "hostIP": "10.128.0.8",
        "phase": "Running",
        "podIP": "10.4.2.7",
        "startTime": "2017-02-28T19:18:52Z"
    }
}
        `

// pod2 has one single label: "role": "External"
const pod2s = `{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "creationTimestamp": "2017-02-28T19:18:52Z",
        "labels": {
            "role": "External"
        },
        "name": "external",
        "namespace": "demo",
        "resourceVersion": "543486",
        "selfLink": "/api/v1/namespaces/demo/pods/external",
        "uid": "bd8de6cd-fdea-11e6-91a2-42010a8001b8"
    },
    "spec": {
        "containers": [
            {
                "image": "redis",
                "imagePullPolicy": "Always",
                "name": "redismaster",
                "ports": [
                    {
                        "containerPort": 80,
                        "protocol": "TCP"
                    }
                ],
                "resources": {
                    "requests": {
                        "cpu": "100m",
                        "memory": "100Mi"
                    }
                },
                "terminationMessagePath": "/dev/termination-log",
                "volumeMounts": [
                    {
                        "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                        "name": "default-token-kb5rn",
                        "readOnly": true
                    }
                ]
            }
        ],
        "dnsPolicy": "ClusterFirst",
        "nodeName": "gke-test48-default-pool-1e87d955-30r9",
        "restartPolicy": "Always",
        "securityContext": {},
        "serviceAccount": "default",
        "serviceAccountName": "default",
        "terminationGracePeriodSeconds": 30,
        "volumes": [
            {
                "name": "default-token-kb5rn",
                "secret": {
                    "defaultMode": 420,
                    "secretName": "default-token-kb5rn"
                }
            }
        ]
    },
    "status": {
        "conditions": [
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2017-02-28T19:18:52Z",
                "status": "True",
                "type": "Initialized"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2017-02-28T19:18:54Z",
                "status": "True",
                "type": "Ready"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2017-02-28T19:18:52Z",
                "status": "True",
                "type": "PodScheduled"
            }
        ],
        "containerStatuses": [
            {
                "containerID": "docker://a82690934217d85872a901c877be5343f32bb8d39a8a0c7deb98bbebf5c74922",
                "image": "redis",
                "imageID": "docker://sha256:1a8a9ee54eb755a427e00484a64fa4edbe9c2abe59ca468a4e452b343a2b57c2",
                "lastState": {},
                "name": "redismaster",
                "ready": true,
                "restartCount": 0,
                "state": {
                    "running": {
                        "startedAt": "2017-02-28T19:18:53Z"
                    }
                }
            }
        ],
        "hostIP": "10.128.0.9",
        "phase": "Running",
        "podIP": "10.4.1.8",
        "startTime": "2017-02-28T19:18:52Z"
    }
}
`

// pod3 has one single label: "role": "BusinessBackend"
const pod3s = `{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "creationTimestamp": "2017-02-28T19:18:52Z",
        "labels": {
            "role": "BusinessBackend"
        },
        "name": "backend",
        "namespace": "demo",
        "resourceVersion": "543487",
        "selfLink": "/api/v1/namespaces/demo/pods/backend",
        "uid": "bdabca07-fdea-11e6-91a2-42010a8001b8"
    },
    "spec": {
        "containers": [
            {
                "image": "nginx",
                "imagePullPolicy": "Always",
                "name": "nginx",
                "ports": [
                    {
                        "containerPort": 6379,
                        "protocol": "TCP"
                    }
                ],
                "resources": {},
                "terminationMessagePath": "/dev/termination-log",
                "volumeMounts": [
                    {
                        "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                        "name": "default-token-kb5rn",
                        "readOnly": true
                    }
                ]
            }
        ],
        "dnsPolicy": "ClusterFirst",
        "nodeName": "gke-test48-default-pool-1e87d955-30r9",
        "restartPolicy": "Always",
        "securityContext": {},
        "serviceAccount": "default",
        "serviceAccountName": "default",
        "terminationGracePeriodSeconds": 30,
        "volumes": [
            {
                "name": "default-token-kb5rn",
                "secret": {
                    "defaultMode": 420,
                    "secretName": "default-token-kb5rn"
                }
            }
        ]
    },
    "status": {
        "conditions": [
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2017-02-28T19:18:52Z",
                "status": "True",
                "type": "Initialized"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2017-02-28T19:18:54Z",
                "status": "True",
                "type": "Ready"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2017-02-28T19:18:52Z",
                "status": "True",
                "type": "PodScheduled"
            }
        ],
        "containerStatuses": [
            {
                "containerID": "docker://f66867c34c3f505b7296eb790fec350323d13f6f2a03f7030fa59b7171822a54",
                "image": "nginx",
                "imageID": "docker://sha256:db079554b4d2f7c65c4df3adae88cb72d051c8c3b8613eb44e86f60c945b1ca7",
                "lastState": {},
                "name": "nginx",
                "ready": true,
                "restartCount": 0,
                "state": {
                    "running": {
                        "startedAt": "2017-02-28T19:18:54Z"
                    }
                }
            }
        ],
        "hostIP": "10.128.0.9",
        "phase": "Running",
        "podIP": "10.4.1.9",
        "startTime": "2017-02-28T19:18:52Z"
    }
}`

// pod4 got multiple labels: "role": "WebFrontend", "job": "worker",
const pod4s = `{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "creationTimestamp": "2017-02-28T19:18:52Z",
        "labels": {
            "role": "WebFrontend",
            "job": "worker"
        },
        "name": "frontend",
        "namespace": "demo",
        "resourceVersion": "543485",
        "selfLink": "/api/v1/namespaces/demo/pods/frontend",
        "uid": "bd9c01fb-fdea-11e6-91a2-42010a8001b8"
    },
    "spec": {
        "containers": [
            {
                "image": "redis",
                "imagePullPolicy": "Always",
                "name": "redismaster",
                "ports": [
                    {
                        "containerPort": 6379,
                        "protocol": "TCP"
                    }
                ],
                "resources": {
                    "requests": {
                        "cpu": "100m",
                        "memory": "100Mi"
                    }
                },
                "terminationMessagePath": "/dev/termination-log",
                "volumeMounts": [
                    {
                        "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                        "name": "default-token-kb5rn",
                        "readOnly": true
                    }
                ]
            }
        ],
        "dnsPolicy": "ClusterFirst",
        "nodeName": "gke-test48-default-pool-1e87d955-mwh5",
        "restartPolicy": "Always",
        "securityContext": {},
        "serviceAccount": "default",
        "serviceAccountName": "default",
        "terminationGracePeriodSeconds": 30,
        "volumes": [
            {
                "name": "default-token-kb5rn",
                "secret": {
                    "defaultMode": 420,
                    "secretName": "default-token-kb5rn"
                }
            }
        ]
    },
    "status": {
        "conditions": [
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2017-02-28T19:18:52Z",
                "status": "True",
                "type": "Initialized"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2017-02-28T19:18:54Z",
                "status": "True",
                "type": "Ready"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2017-02-28T19:18:52Z",
                "status": "True",
                "type": "PodScheduled"
            }
        ],
        "containerStatuses": [
            {
                "containerID": "docker://48d8e102c086b5686abde297d983c3278af548a9551dc84facb8e77f5944bc75",
                "image": "redis",
                "imageID": "docker://sha256:1a8a9ee54eb755a427e00484a64fa4edbe9c2abe59ca468a4e452b343a2b57c2",
                "lastState": {},
                "name": "redismaster",
                "ready": true,
                "restartCount": 0,
                "state": {
                    "running": {
                        "startedAt": "2017-02-28T19:18:53Z"
                    }
                }
            }
        ],
        "hostIP": "10.128.0.8",
        "phase": "Running",
        "podIP": "10.4.2.7",
        "startTime": "2017-02-28T19:18:52Z"
    }
}`

// Policy matches label "role": "WebFrontend".
const policy1s = `{
   "metadata": {
      "name": "frontend-policy",
      "namespace": "default",
      "selfLink": "/apis/extensions/v1beta1/namespaces/default/networkpolicies/frontend-policy",
      "uid": "22c9ad41-89b1-11e6-a046-0800277021d9",
      "resourceVersion": "2482",
      "generation": 1,
      "creationTimestamp": "2016-10-03T21:34:16Z"
   },
   "spec": {
      "podSelector": {
         "matchLabels": {
            "role": "WebFrontend"
         }
      },
      "ingress": [
         {
            "ports": [
               {
                  "protocol": "tcp",
                  "port": 80
               },
               {
                  "protocol": "tcp",
                  "port": 8080
               }
            ],
            "from": [
               {
                  "podSelector": {
                     "matchLabels": {
                        "role": "WebFrontend"
                     }
                  }
               },
               {
                  "podSelector": {
                     "matchLabels": {
                        "role": "BusinessBackend"
                     }
                  }
               }
            ]
         }
      ]
   }
}`

// Policy matches label "role": "BusinessBackend".
const policy2s = `{
   "metadata": {
      "name": "backend-policy",
      "namespace": "default",
      "selfLink": "/apis/extensions/v1beta1/namespaces/default/networkpolicies/backend-policy",
      "uid": "22cde507-89b1-11e6-a046-0800277021d9",
      "resourceVersion": "2483",
      "generation": 1,
      "creationTimestamp": "2016-10-03T21:34:16Z"
   },
   "spec": {
      "podSelector": {
         "matchLabels": {
            "role": "BusinessBackend"
         }
      },
      "ingress": [
         {
            "from": [
               {
                  "podSelector": {
                     "matchLabels": {
                        "role": "WebFrontend"
                     }
                  }
               },
               {
                  "podSelector": {
                     "matchLabels": {
                        "role": "BusinessBackend"
                     }
                  }
               },
               {
                  "podSelector": {
                     "matchLabels": {
                        "role": "Database"
                     }
                  }
               }
            ]
         }
      ]
   }
}`

// Policy matches label "role": "Database".
const policy3s = `{
   "metadata": {
      "name": "database-policy",
      "namespace": "default",
      "selfLink": "/apis/extensions/v1beta1/namespaces/default/networkpolicies/database-policy",
      "uid": "22d227d2-89b1-11e6-a046-0800277021d9",
      "resourceVersion": "2484",
      "generation": 1,
      "creationTimestamp": "2016-10-03T21:34:16Z"
   },
   "spec": {
      "podSelector": {
         "matchLabels": {
            "role": "Database"
         }
      },
      "ingress": [
         {
            "from": [
               {
                  "podSelector": {
                     "matchLabels": {
                        "role": "BusinessBackend"
                     }
                  }
               },
               {
                  "podSelector": {
                     "matchLabels": {
                        "role": "Database"
                     }
                  }
               }
            ]
         }
      ]
   }
}`

// Policy matches label "role": "Database".
const policy4s = `{
   "metadata": {
      "name": "database-policy",
      "namespace": "default",
      "selfLink": "/apis/extensions/v1beta1/namespaces/default/networkpolicies/database-policy",
      "uid": "22d227d2-89b1-11e6-a046-0800277021d9",
      "resourceVersion": "2484",
      "generation": 1,
      "creationTimestamp": "2016-10-03T21:34:16Z"
   },
   "spec": {
      "podSelector": {
         "matchLabels": {
            "job": "worker"
         }
      },
      "ingress": [
         {
            "from": [
               {
                  "podSelector": {
                     "matchLabels": {
                        "xcvxcv": "asdasd"
                     }
                  }
               },
               {
                  "podSelector": {
                     "matchLabels": {
                        "ZXZX": "qwewqeqwe"
                     }
                  }
               }
            ]
         }
      ]
   }
}`

func TestSingleLabelMatch(t *testing.T) {
	pod1 := api.Pod{}
	pod2 := api.Pod{}
	pod3 := api.Pod{}

	policy1 := extensions.NetworkPolicy{}
	policy2 := extensions.NetworkPolicy{}
	policy3 := extensions.NetworkPolicy{}

	json.Unmarshal([]byte(pod1s), &pod1)
	json.Unmarshal([]byte(pod2s), &pod2)
	json.Unmarshal([]byte(pod3s), &pod3)

	json.Unmarshal([]byte(policy1s), &policy1)
	json.Unmarshal([]byte(policy2s), &policy2)
	json.Unmarshal([]byte(policy3s), &policy3)

	policyList := extensions.NetworkPolicyList{
		Items: []extensions.NetworkPolicy{policy1,
			policy2,
			policy3,
		},
	}

	podList := api.PodList{
		Items: []api.Pod{pod1,
			pod2,
			pod3,
		},
	}

	// Testing frontend
	t.Log("Testing Pod frontend single policy match")
	fmt.Printf("POD1: %+v \n\n\n", pod1)

	result, _ := ListPoliciesPerPod(&pod1, &policyList)
	fmt.Printf("RESULTS %+v", result)
	if len(result.Items) != 1 {
		t.Errorf("Expected 1 policy match for frontend, got %d", len(result.Items))
	}
	if !reflect.DeepEqual(result.Items[0], policy1) {
		t.Errorf("Expected PolicyMatch for WebFrontend")
	}

	resultRules, _ := ListIngressRulesPerPod(&pod1, &policyList)
	if !reflect.DeepEqual(*resultRules, policy1.Spec.Ingress) {
		t.Errorf("Expected RuleMatch for WebFrontend")
	}

	// Testing External
	t.Log("Testing Pod External no policy match")
	result, _ = ListPoliciesPerPod(&pod2, &policyList)
	if len(result.Items) != 0 {
		t.Errorf("Expected 0 policy match for external, got %d", len(result.Items))
	}

	// Testing BusinessBackend
	t.Log("Testing Pod Backend single policy match")
	result, _ = ListPoliciesPerPod(&pod3, &policyList)
	if len(result.Items) != 1 {
		t.Errorf("Expected 1 policy match for BusinessBackend, got %d", len(result.Items))
	}
	if !reflect.DeepEqual(result.Items[0], policy2) {
		t.Errorf("Expected PolicyMatch for BusinessBackend")
	}

	resultRules, _ = ListIngressRulesPerPod(&pod3, &policyList)
	if !reflect.DeepEqual(*resultRules, policy2.Spec.Ingress) {
		t.Errorf("Expected RuleMatch for BusinessBackend")
	}

	// Testing Policy1:
	t.Log("Testing PolicyFrontend Policy")
	resultPods, _ := ListPodsPerPolicy(&policy1, &podList)
	if len(resultPods.Items) != 1 {
		t.Errorf("Expected 1 Pod match for policy frontend, got %d", len(resultPods.Items))
	}
	if !reflect.DeepEqual(resultPods.Items[0], pod1) {
		t.Errorf("Failed pod match for Frontend ListPodsPerPolicy")
	}

	// Testing Policy2:
	t.Log("Testing BusinessBackend Policy")
	resultPods, _ = ListPodsPerPolicy(&policy2, &podList)
	if len(resultPods.Items) != 1 {
		t.Errorf("Expected 1 Pod match for policy frontend, got %d", len(resultPods.Items))
	}
	if !reflect.DeepEqual(resultPods.Items[0], pod3) {
		t.Errorf("Failed pod match for Frontend ListPodsPerPolicy")
	}

	// Testing Policy 3
	t.Log("Testing Database Policy")
	resultPods, _ = ListPodsPerPolicy(&policy3, &podList)
	if len(resultPods.Items) != 0 {
		t.Errorf("Expected 1 Pod match for policy frontend, got %d", len(resultPods.Items))
	}
}

func TestMultipleLabelMatch(t *testing.T) {
	pod4 := api.Pod{}

	policy1 := extensions.NetworkPolicy{}
	policy2 := extensions.NetworkPolicy{}
	policy3 := extensions.NetworkPolicy{}
	policy4 := extensions.NetworkPolicy{}

	json.Unmarshal([]byte(pod4s), &pod4)

	json.Unmarshal([]byte(policy1s), &policy1)
	json.Unmarshal([]byte(policy2s), &policy2)
	json.Unmarshal([]byte(policy3s), &policy3)
	json.Unmarshal([]byte(policy4s), &policy4)

	policyList := extensions.NetworkPolicyList{
		Items: []extensions.NetworkPolicy{policy1,
			policy2,
			policy3,
			policy4,
		},
	}

	expectedPolicies := extensions.NetworkPolicyList{Items: []extensions.NetworkPolicy{policy1, policy4}}
	expectedRules := append(policy1.Spec.Ingress, policy4.Spec.Ingress...)
	// Testing pod4
	t.Log("Testing Pod 4 Multible policy match")
	resultPolicies, _ := ListPoliciesPerPod(&pod4, &policyList)
	if len(resultPolicies.Items) != 2 {
		t.Errorf("Expected 2 policy match for pod4, got %d", len(resultPolicies.Items))
	}
	if !reflect.DeepEqual(*resultPolicies, expectedPolicies) {
		t.Errorf("Expected Policy Match for pod4: \n got %+v ,\n expected %+v", resultPolicies, expectedPolicies)
	}

	resultRules, _ := ListIngressRulesPerPod(&pod4, &policyList)
	if len(*resultRules) != 2 {
		t.Errorf("Expected 2 rule match for pod4, got %d", len(*resultRules))
	}
	if !reflect.DeepEqual(*resultRules, expectedRules) {
		t.Errorf("Expected rule match for pod4: \n got %+v ,\n  expected %+v", resultRules, expectedRules)
	}

}
