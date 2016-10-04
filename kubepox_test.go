package kubepox

import (
	"encoding/json"
	"reflect"
	"testing"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

// pod1 has one single label: "role": "WebFrontend"
const pod1s = `{
   "metadata": {
      "name": "frontend",
      "namespace": "default",
      "selfLink": "/api/v1/namespaces/default/pods/frontend",
      "uid": "e4913a63-89c2-11e6-a046-0800277021d9",
      "resourceVersion": "10956",
      "creationTimestamp": "2016-10-03T23:41:23Z",
      "labels": {
         "role": "WebFrontend"
      }
   },
   "spec": {
      "volumes": [
         {
            "name": "default-token-nxt3i",
            "secret": {
               "secretName": "default-token-nxt3i",
               "defaultMode": 420
            }
         }
      ],
      "containers": [
         {
            "name": "redismaster",
            "image": "redis",
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
            "volumeMounts": [
               {
                  "name": "default-token-nxt3i",
                  "readOnly": true,
                  "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
               }
            ],
            "terminationMessagePath": "/dev/termination-log",
            "imagePullPolicy": "Always"
         }
      ],
      "restartPolicy": "Always",
      "terminationGracePeriodSeconds": 30,
      "dnsPolicy": "ClusterFirst",
      "serviceAccountName": "default",
      "nodeName": "127.0.0.1",
      "securityContext": {}
   },
   "status": {
      "phase": "Running",
      "conditions": [
         {
            "type": "Initialized",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2016-10-03T23:41:23Z"
         },
         {
            "type": "Ready",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2016-10-03T23:41:28Z"
         },
         {
            "type": "PodScheduled",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2016-10-03T23:41:23Z"
         }
      ],
      "hostIP": "127.0.0.1",
      "podIP": "172.17.0.3",
      "startTime": "2016-10-03T23:41:23Z",
      "containerStatuses": [
         {
            "name": "redismaster",
            "state": {
               "running": {
                  "startedAt": "2016-10-03T23:41:27Z"
               }
            },
            "lastState": {},
            "ready": true,
            "restartCount": 0,
            "image": "redis",
            "imageID": "docker://sha256:1aa84b1b434e43bbd5a5577e334050e9bc5984aec570c981c7357e6bb6a6df1f",
            "containerID": "docker://8e00084abde70221b123263e0fcff055a88ce8141965d8143875974a5d71d078"
         }
      ]
   }
}`

// pod2 has one single label: "role": "External"
const pod2s = `{
   "metadata": {
      "name": "external",
      "namespace": "default",
      "selfLink": "/api/v1/namespaces/default/pods/external",
      "uid": "e489731b-89c2-11e6-a046-0800277021d9",
      "resourceVersion": "10949",
      "creationTimestamp": "2016-10-03T23:41:23Z",
      "labels": {
         "role": "External"
      }
   },
   "spec": {
      "volumes": [
         {
            "name": "default-token-nxt3i",
            "secret": {
               "secretName": "default-token-nxt3i",
               "defaultMode": 420
            }
         }
      ],
      "containers": [
         {
            "name": "redismaster",
            "image": "redis",
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
            "volumeMounts": [
               {
                  "name": "default-token-nxt3i",
                  "readOnly": true,
                  "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
               }
            ],
            "terminationMessagePath": "/dev/termination-log",
            "imagePullPolicy": "Always"
         }
      ],
      "restartPolicy": "Always",
      "terminationGracePeriodSeconds": 30,
      "dnsPolicy": "ClusterFirst",
      "serviceAccountName": "default",
      "nodeName": "127.0.0.1",
      "securityContext": {}
   },
   "status": {
      "phase": "Running",
      "conditions": [
         {
            "type": "Initialized",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2016-10-03T23:41:23Z"
         },
         {
            "type": "Ready",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2016-10-03T23:41:26Z"
         },
         {
            "type": "PodScheduled",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2016-10-03T23:41:23Z"
         }
      ],
      "hostIP": "127.0.0.1",
      "podIP": "172.17.0.2",
      "startTime": "2016-10-03T23:41:23Z",
      "containerStatuses": [
         {
            "name": "redismaster",
            "state": {
               "running": {
                  "startedAt": "2016-10-03T23:41:25Z"
               }
            },
            "lastState": {},
            "ready": true,
            "restartCount": 0,
            "image": "redis",
            "imageID": "docker://sha256:1aa84b1b434e43bbd5a5577e334050e9bc5984aec570c981c7357e6bb6a6df1f",
            "containerID": "docker://94d9a2a1d318715f6ccd7d6d557084134913fecb66c16f649f5a0608e4d2d505"
         }
      ]
   }
}`

// pod3 has one single label: "role": "BusinessBackend"
const pod3s = `{
   "metadata": {
      "name": "backend",
      "namespace": "default",
      "selfLink": "/api/v1/namespaces/default/pods/backend",
      "uid": "e497cea7-89c2-11e6-a046-0800277021d9",
      "resourceVersion": "10962",
      "creationTimestamp": "2016-10-03T23:41:23Z",
      "labels": {
         "role": "BusinessBackend"
      }
   },
   "spec": {
      "volumes": [
         {
            "name": "default-token-nxt3i",
            "secret": {
               "secretName": "default-token-nxt3i",
               "defaultMode": 420
            }
         }
      ],
      "containers": [
         {
            "name": "nginx",
            "image": "nginx",
            "ports": [
               {
                  "containerPort": 6379,
                  "protocol": "TCP"
               }
            ],
            "resources": {},
            "volumeMounts": [
               {
                  "name": "default-token-nxt3i",
                  "readOnly": true,
                  "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
               }
            ],
            "terminationMessagePath": "/dev/termination-log",
            "imagePullPolicy": "Always"
         }
      ],
      "restartPolicy": "Always",
      "terminationGracePeriodSeconds": 30,
      "dnsPolicy": "ClusterFirst",
      "serviceAccountName": "default",
      "nodeName": "127.0.0.1",
      "securityContext": {}
   },
   "status": {
      "phase": "Running",
      "conditions": [
         {
            "type": "Initialized",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2016-10-03T23:41:23Z"
         },
         {
            "type": "Ready",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2016-10-03T23:41:29Z"
         },
         {
            "type": "PodScheduled",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2016-10-03T23:41:23Z"
         }
      ],
      "hostIP": "127.0.0.1",
      "podIP": "172.17.0.4",
      "startTime": "2016-10-03T23:41:23Z",
      "containerStatuses": [
         {
            "name": "nginx",
            "state": {
               "running": {
                  "startedAt": "2016-10-03T23:41:29Z"
               }
            },
            "lastState": {},
            "ready": true,
            "restartCount": 0,
            "image": "nginx",
            "imageID": "docker://sha256:ba6bed934df2e644fdd34e9d324c80f3c615544ee9a93e4ce3cfddfcf84bdbc2",
            "containerID": "docker://a5435af8e358fb413826b156b9244e503f8ede3f44e04e433f1258efe4074963"
         }
      ]
   }
}
`

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
	result, _ := ListPoliciesPerPod(&pod1, &policyList)
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
