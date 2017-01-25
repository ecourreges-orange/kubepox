package kubepox

import (
	"encoding/json"
	"reflect"
	"testing"

	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/apis/extensions"
)

// pod1 has one single label: "role": "WebFrontend"
const pod1s = `{
   "Name": "frontend",
   "GenerateName": "",
   "Namespace": "demo",
   "SelfLink": "/api/v1/namespaces/demo/pods/frontend",
   "UID": "3d7d45f4-c609-11e6-a4fd-06e1269025c9",
   "ResourceVersion": "4964585",
   "Generation": 0,
   "CreationTimestamp": "2016-12-19T16:36:07Z",
   "DeletionTimestamp": null,
   "DeletionGracePeriodSeconds": null,
   "Labels": {
      "role": "WebFrontend"
   },
   "Annotations": null,
   "OwnerReferences": null,
   "Finalizers": null,
   "ClusterName": "",
   "Spec": {
      "Volumes": [
         {
            "Name": "default-token-79ah4",
            "HostPath": null,
            "EmptyDir": null,
            "GCEPersistentDisk": null,
            "AWSElasticBlockStore": null,
            "GitRepo": null,
            "Secret": {
               "SecretName": "default-token-79ah4",
               "Items": null,
               "DefaultMode": 420
            },
            "NFS": null,
            "ISCSI": null,
            "Glusterfs": null,
            "PersistentVolumeClaim": null,
            "RBD": null,
            "Quobyte": null,
            "FlexVolume": null,
            "Cinder": null,
            "CephFS": null,
            "Flocker": null,
            "DownwardAPI": null,
            "FC": null,
            "AzureFile": null,
            "ConfigMap": null,
            "VsphereVolume": null,
            "AzureDisk": null,
            "PhotonPersistentDisk": null
         }
      ],
      "InitContainers": null,
      "Containers": [
         {
            "Name": "redismaster",
            "Image": "redis",
            "Command": null,
            "Args": null,
            "WorkingDir": "",
            "Ports": [
               {
                  "Name": "",
                  "HostPort": 0,
                  "ContainerPort": 6379,
                  "Protocol": "TCP",
                  "HostIP": ""
               }
            ],
            "Env": null,
            "Resources": {
               "Limits": null,
               "Requests": {
                  "cpu": "100m",
                  "memory": "100Mi"
               }
            },
            "VolumeMounts": [
               {
                  "Name": "default-token-79ah4",
                  "ReadOnly": true,
                  "MountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                  "SubPath": ""
               }
            ],
            "LivenessProbe": null,
            "ReadinessProbe": null,
            "Lifecycle": null,
            "TerminationMessagePath": "/dev/termination-log",
            "ImagePullPolicy": "Always",
            "SecurityContext": null,
            "Stdin": false,
            "StdinOnce": false,
            "TTY": false
         }
      ],
      "RestartPolicy": "Always",
      "TerminationGracePeriodSeconds": 30,
      "ActiveDeadlineSeconds": null,
      "DNSPolicy": "ClusterFirst",
      "NodeSelector": null,
      "ServiceAccountName": "default",
      "NodeName": "ip-10-0-0-4.us-west-2.compute.internal",
      "SecurityContext": {
         "HostNetwork": false,
         "HostPID": false,
         "HostIPC": false,
         "SELinuxOptions": null,
         "RunAsUser": null,
         "RunAsNonRoot": null,
         "SupplementalGroups": null,
         "FSGroup": null
      },
      "ImagePullSecrets": null,
      "Hostname": "",
      "Subdomain": "",
      "Affinity": {
         "NodeAffinity": null,
         "PodAffinity": null,
         "PodAntiAffinity": null
      }
   },
   "Status": {
      "Phase": "Running",
      "Conditions": [
         {
            "Type": "Initialized",
            "Status": "True",
            "LastProbeTime": null,
            "LastTransitionTime": "2016-12-19T16:36:07Z",
            "Reason": "",
            "Message": ""
         },
         {
            "Type": "Ready",
            "Status": "True",
            "LastProbeTime": null,
            "LastTransitionTime": "2016-12-19T16:36:11Z",
            "Reason": "",
            "Message": ""
         },
         {
            "Type": "PodScheduled",
            "Status": "True",
            "LastProbeTime": null,
            "LastTransitionTime": "2016-12-19T16:36:07Z",
            "Reason": "",
            "Message": ""
         }
      ],
      "Message": "",
      "Reason": "",
      "HostIP": "10.0.0.4",
      "PodIP": "10.2.66.72",
      "StartTime": "2016-12-19T16:36:07Z",
      "InitContainerStatuses": null,
      "ContainerStatuses": [
         {
            "Name": "redismaster",
            "State": {
               "Waiting": null,
               "Running": {
                  "StartedAt": "2016-12-19T16:36:10Z"
               },
               "Terminated": null
            },
            "LastTerminationState": {
               "Waiting": null,
               "Running": null,
               "Terminated": null
            },
            "Ready": true,
            "RestartCount": 0,
            "Image": "redis",
            "ImageID": "docker://sha256:1c2ac2024e4b6d621cea1458923bdbd1806f2c7c50c8a7292e0e6551b8d768e3",
            "ContainerID": "docker://e12e3b2b3be33b6436e9d2df23559b3829ff6c64989b6b4bb23d5d660808de63"
         }
      ]
   }
}

        `

// pod2 has one single label: "role": "External"
const pod2s = `{
   "Name": "external",
   "GenerateName": "",
   "Namespace": "demo",
   "SelfLink": "/api/v1/namespaces/demo/pods/external",
   "UID": "3d732f0d-c609-11e6-a4fd-06e1269025c9",
   "ResourceVersion": "4964580",
   "Generation": 0,
   "CreationTimestamp": "2016-12-19T16:36:06Z",
   "DeletionTimestamp": null,
   "DeletionGracePeriodSeconds": null,
   "Labels": {
      "role": "External"
   },
   "Annotations": null,
   "OwnerReferences": null,
   "Finalizers": null,
   "ClusterName": "",
   "Spec": {
      "Volumes": [
         {
            "Name": "default-token-79ah4",
            "HostPath": null,
            "EmptyDir": null,
            "GCEPersistentDisk": null,
            "AWSElasticBlockStore": null,
            "GitRepo": null,
            "Secret": {
               "SecretName": "default-token-79ah4",
               "Items": null,
               "DefaultMode": 420
            },
            "NFS": null,
            "ISCSI": null,
            "Glusterfs": null,
            "PersistentVolumeClaim": null,
            "RBD": null,
            "Quobyte": null,
            "FlexVolume": null,
            "Cinder": null,
            "CephFS": null,
            "Flocker": null,
            "DownwardAPI": null,
            "FC": null,
            "AzureFile": null,
            "ConfigMap": null,
            "VsphereVolume": null,
            "AzureDisk": null,
            "PhotonPersistentDisk": null
         }
      ],
      "InitContainers": null,
      "Containers": [
         {
            "Name": "redismaster",
            "Image": "redis",
            "Command": null,
            "Args": null,
            "WorkingDir": "",
            "Ports": [
               {
                  "Name": "",
                  "HostPort": 0,
                  "ContainerPort": 80,
                  "Protocol": "TCP",
                  "HostIP": ""
               }
            ],
            "Env": null,
            "Resources": {
               "Limits": null,
               "Requests": {
                  "cpu": "100m",
                  "memory": "100Mi"
               }
            },
            "VolumeMounts": [
               {
                  "Name": "default-token-79ah4",
                  "ReadOnly": true,
                  "MountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                  "SubPath": ""
               }
            ],
            "LivenessProbe": null,
            "ReadinessProbe": null,
            "Lifecycle": null,
            "TerminationMessagePath": "/dev/termination-log",
            "ImagePullPolicy": "Always",
            "SecurityContext": null,
            "Stdin": false,
            "StdinOnce": false,
            "TTY": false
         }
      ],
      "RestartPolicy": "Always",
      "TerminationGracePeriodSeconds": 30,
      "ActiveDeadlineSeconds": null,
      "DNSPolicy": "ClusterFirst",
      "NodeSelector": null,
      "ServiceAccountName": "default",
      "NodeName": "ip-10-0-0-4.us-west-2.compute.internal",
      "SecurityContext": {
         "HostNetwork": false,
         "HostPID": false,
         "HostIPC": false,
         "SELinuxOptions": null,
         "RunAsUser": null,
         "RunAsNonRoot": null,
         "SupplementalGroups": null,
         "FSGroup": null
      },
      "ImagePullSecrets": null,
      "Hostname": "",
      "Subdomain": "",
      "Affinity": {
         "NodeAffinity": null,
         "PodAffinity": null,
         "PodAntiAffinity": null
      }
   },
   "Status": {
      "Phase": "Running",
      "Conditions": [
         {
            "Type": "Initialized",
            "Status": "True",
            "LastProbeTime": null,
            "LastTransitionTime": "2016-12-19T16:36:06Z",
            "Reason": "",
            "Message": ""
         },
         {
            "Type": "Ready",
            "Status": "True",
            "LastProbeTime": null,
            "LastTransitionTime": "2016-12-19T16:36:10Z",
            "Reason": "",
            "Message": ""
         },
         {
            "Type": "PodScheduled",
            "Status": "True",
            "LastProbeTime": null,
            "LastTransitionTime": "2016-12-19T16:36:06Z",
            "Reason": "",
            "Message": ""
         }
      ],
      "Message": "",
      "Reason": "",
      "HostIP": "10.0.0.4",
      "PodIP": "10.2.66.71",
      "StartTime": "2016-12-19T16:36:06Z",
      "InitContainerStatuses": null,
      "ContainerStatuses": [
         {
            "Name": "redismaster",
            "State": {
               "Waiting": null,
               "Running": {
                  "StartedAt": "2016-12-19T16:36:09Z"
               },
               "Terminated": null
            },
            "LastTerminationState": {
               "Waiting": null,
               "Running": null,
               "Terminated": null
            },
            "Ready": true,
            "RestartCount": 0,
            "Image": "redis",
            "ImageID": "docker://sha256:1c2ac2024e4b6d621cea1458923bdbd1806f2c7c50c8a7292e0e6551b8d768e3",
            "ContainerID": "docker://5c07497ca46686553f90572952214709ca5717286fc2460a81e814ce73a66b41"
         }
      ]
   }
}
`

// pod3 has one single label: "role": "BusinessBackend"
const pod3s = `{
   "Name": "backend",
   "GenerateName": "",
   "Namespace": "demo",
   "SelfLink": "/api/v1/namespaces/demo/pods/backend",
   "UID": "3d86e2c0-c609-11e6-a4fd-06e1269025c9",
   "ResourceVersion": "4964577",
   "Generation": 0,
   "CreationTimestamp": "2016-12-19T16:36:07Z",
   "DeletionTimestamp": null,
   "DeletionGracePeriodSeconds": null,
   "Labels": {
      "role": "BusinessBackend"
   },
   "Annotations": null,
   "OwnerReferences": null,
   "Finalizers": null,
   "ClusterName": "",
   "Spec": {
      "Volumes": [
         {
            "Name": "default-token-79ah4",
            "HostPath": null,
            "EmptyDir": null,
            "GCEPersistentDisk": null,
            "AWSElasticBlockStore": null,
            "GitRepo": null,
            "Secret": {
               "SecretName": "default-token-79ah4",
               "Items": null,
               "DefaultMode": 420
            },
            "NFS": null,
            "ISCSI": null,
            "Glusterfs": null,
            "PersistentVolumeClaim": null,
            "RBD": null,
            "Quobyte": null,
            "FlexVolume": null,
            "Cinder": null,
            "CephFS": null,
            "Flocker": null,
            "DownwardAPI": null,
            "FC": null,
            "AzureFile": null,
            "ConfigMap": null,
            "VsphereVolume": null,
            "AzureDisk": null,
            "PhotonPersistentDisk": null
         }
      ],
      "InitContainers": null,
      "Containers": [
         {
            "Name": "nginx",
            "Image": "nginx",
            "Command": null,
            "Args": null,
            "WorkingDir": "",
            "Ports": [
               {
                  "Name": "",
                  "HostPort": 0,
                  "ContainerPort": 6379,
                  "Protocol": "TCP",
                  "HostIP": ""
               }
            ],
            "Env": null,
            "Resources": {
               "Limits": null,
               "Requests": null
            },
            "VolumeMounts": [
               {
                  "Name": "default-token-79ah4",
                  "ReadOnly": true,
                  "MountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                  "SubPath": ""
               }
            ],
            "LivenessProbe": null,
            "ReadinessProbe": null,
            "Lifecycle": null,
            "TerminationMessagePath": "/dev/termination-log",
            "ImagePullPolicy": "Always",
            "SecurityContext": null,
            "Stdin": false,
            "StdinOnce": false,
            "TTY": false
         }
      ],
      "RestartPolicy": "Always",
      "TerminationGracePeriodSeconds": 30,
      "ActiveDeadlineSeconds": null,
      "DNSPolicy": "ClusterFirst",
      "NodeSelector": null,
      "ServiceAccountName": "default",
      "NodeName": "ip-10-0-0-5.us-west-2.compute.internal",
      "SecurityContext": {
         "HostNetwork": false,
         "HostPID": false,
         "HostIPC": false,
         "SELinuxOptions": null,
         "RunAsUser": null,
         "RunAsNonRoot": null,
         "SupplementalGroups": null,
         "FSGroup": null
      },
      "ImagePullSecrets": null,
      "Hostname": "",
      "Subdomain": "",
      "Affinity": {
         "NodeAffinity": null,
         "PodAffinity": null,
         "PodAntiAffinity": null
      }
   },
   "Status": {
      "Phase": "Running",
      "Conditions": [
         {
            "Type": "Initialized",
            "Status": "True",
            "LastProbeTime": null,
            "LastTransitionTime": "2016-12-19T16:36:07Z",
            "Reason": "",
            "Message": ""
         },
         {
            "Type": "Ready",
            "Status": "True",
            "LastProbeTime": null,
            "LastTransitionTime": "2016-12-19T16:36:09Z",
            "Reason": "",
            "Message": ""
         },
         {
            "Type": "PodScheduled",
            "Status": "True",
            "LastProbeTime": null,
            "LastTransitionTime": "2016-12-19T16:36:07Z",
            "Reason": "",
            "Message": ""
         }
      ],
      "Message": "",
      "Reason": "",
      "HostIP": "10.0.0.5",
      "PodIP": "10.2.69.69",
      "StartTime": "2016-12-19T16:36:07Z",
      "InitContainerStatuses": null,
      "ContainerStatuses": [
         {
            "Name": "nginx",
            "State": {
               "Waiting": null,
               "Running": {
                  "StartedAt": "2016-12-19T16:36:09Z"
               },
               "Terminated": null
            },
            "LastTerminationState": {
               "Waiting": null,
               "Running": null,
               "Terminated": null
            },
            "Ready": true,
            "RestartCount": 0,
            "Image": "nginx",
            "ImageID": "docker://sha256:abf312888d132e461c61484457ee9fd0125d666672e22f972f3b8c9a0ed3f0a1",
            "ContainerID": "docker://1eb816fda2428572ef2f12521ea45363a785466a0dfca3567ed2bf661a232d15"
         }
      ]
   }
}`

// pod4 got multiple labels: "role": "WebFrontend", "job": "worker",
const pod4s = `{
   "Name": "backend",
   "GenerateName": "",
   "Namespace": "demo",
   "SelfLink": "/api/v1/namespaces/demo/pods/backend",
   "UID": "3d86e2c0-c609-11e6-a4fd-06e1269025c9",
   "ResourceVersion": "4964577",
   "Generation": 0,
   "CreationTimestamp": "2016-12-19T16:36:07Z",
   "DeletionTimestamp": null,
   "DeletionGracePeriodSeconds": null,
   "Labels": {
      "role": "WebFrontend",
			"job": "worker"
   },
   "Annotations": null,
   "OwnerReferences": null,
   "Finalizers": null,
   "ClusterName": "",
   "Spec": {
      "Volumes": [
         {
            "Name": "default-token-79ah4",
            "HostPath": null,
            "EmptyDir": null,
            "GCEPersistentDisk": null,
            "AWSElasticBlockStore": null,
            "GitRepo": null,
            "Secret": {
               "SecretName": "default-token-79ah4",
               "Items": null,
               "DefaultMode": 420
            },
            "NFS": null,
            "ISCSI": null,
            "Glusterfs": null,
            "PersistentVolumeClaim": null,
            "RBD": null,
            "Quobyte": null,
            "FlexVolume": null,
            "Cinder": null,
            "CephFS": null,
            "Flocker": null,
            "DownwardAPI": null,
            "FC": null,
            "AzureFile": null,
            "ConfigMap": null,
            "VsphereVolume": null,
            "AzureDisk": null,
            "PhotonPersistentDisk": null
         }
      ],
      "InitContainers": null,
      "Containers": [
         {
            "Name": "nginx",
            "Image": "nginx",
            "Command": null,
            "Args": null,
            "WorkingDir": "",
            "Ports": [
               {
                  "Name": "",
                  "HostPort": 0,
                  "ContainerPort": 6379,
                  "Protocol": "TCP",
                  "HostIP": ""
               }
            ],
            "Env": null,
            "Resources": {
               "Limits": null,
               "Requests": null
            },
            "VolumeMounts": [
               {
                  "Name": "default-token-79ah4",
                  "ReadOnly": true,
                  "MountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                  "SubPath": ""
               }
            ],
            "LivenessProbe": null,
            "ReadinessProbe": null,
            "Lifecycle": null,
            "TerminationMessagePath": "/dev/termination-log",
            "ImagePullPolicy": "Always",
            "SecurityContext": null,
            "Stdin": false,
            "StdinOnce": false,
            "TTY": false
         }
      ],
      "RestartPolicy": "Always",
      "TerminationGracePeriodSeconds": 30,
      "ActiveDeadlineSeconds": null,
      "DNSPolicy": "ClusterFirst",
      "NodeSelector": null,
      "ServiceAccountName": "default",
      "NodeName": "ip-10-0-0-5.us-west-2.compute.internal",
      "SecurityContext": {
         "HostNetwork": false,
         "HostPID": false,
         "HostIPC": false,
         "SELinuxOptions": null,
         "RunAsUser": null,
         "RunAsNonRoot": null,
         "SupplementalGroups": null,
         "FSGroup": null
      },
      "ImagePullSecrets": null,
      "Hostname": "",
      "Subdomain": "",
      "Affinity": {
         "NodeAffinity": null,
         "PodAffinity": null,
         "PodAntiAffinity": null
      }
   },
   "Status": {
      "Phase": "Running",
      "Conditions": [
         {
            "Type": "Initialized",
            "Status": "True",
            "LastProbeTime": null,
            "LastTransitionTime": "2016-12-19T16:36:07Z",
            "Reason": "",
            "Message": ""
         },
         {
            "Type": "Ready",
            "Status": "True",
            "LastProbeTime": null,
            "LastTransitionTime": "2016-12-19T16:36:09Z",
            "Reason": "",
            "Message": ""
         },
         {
            "Type": "PodScheduled",
            "Status": "True",
            "LastProbeTime": null,
            "LastTransitionTime": "2016-12-19T16:36:07Z",
            "Reason": "",
            "Message": ""
         }
      ],
      "Message": "",
      "Reason": "",
      "HostIP": "10.0.0.5",
      "PodIP": "10.2.69.69",
      "StartTime": "2016-12-19T16:36:07Z",
      "InitContainerStatuses": null,
      "ContainerStatuses": [
         {
            "Name": "nginx",
            "State": {
               "Waiting": null,
               "Running": {
                  "StartedAt": "2016-12-19T16:36:09Z"
               },
               "Terminated": null
            },
            "LastTerminationState": {
               "Waiting": null,
               "Running": null,
               "Terminated": null
            },
            "Ready": true,
            "RestartCount": 0,
            "Image": "nginx",
            "ImageID": "docker://sha256:abf312888d132e461c61484457ee9fd0125d666672e22f972f3b8c9a0ed3f0a1",
            "ContainerID": "docker://1eb816fda2428572ef2f12521ea45363a785466a0dfca3567ed2bf661a232d15"
         }
      ]
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
