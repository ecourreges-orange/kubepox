# kubepox

[![Twitter URL](https://img.shields.io/badge/twitter-follow-blue.svg)](https://twitter.com/aporeto_trireme) [![Slack URL](https://img.shields.io/badge/slack-join-green.svg)](https://triremehq.slack.com/messages/general/) [![Documentation](https://img.shields.io/badge/docs-godoc-blue.svg)](https://godoc.org/github.com/aporeto-inc/kubepox)

Kubernetes network Policy eXploration tool

## Library

kubepox is a lightweight library that implements the selection logic used by Kubernetes NetworkPolicies as defined on those specs:
- https://kubernetes.io/docs/concepts/services-networking/network-policies/
- https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.9/#networkpolicy-v1-networking

kubepox takes Kubernetes Pods and NetworkPolicies as input. The implementation need to get those objects, typically from Kubernetes API.

Kubepox is used by the [Trireme-Kubernetes](https://github.com/aporeto-inc/trireme-kubernetes) project as well as the [Aporeto product](https://console.aporeto.com) to enforce pods based on Kubernetes Network-Policies

Kubepox implements the following logic:

- Return all the NetworkPolicies that apply to a pod out of a list:
```
func ListPoliciesPerPod(pod *api.Pod, allPolicies *networking.NetworkPolicyList)
```
- Return the list of Ingress or Egress Rules (from NetworkPolicies) that apply to a pod:
```
func ListIngressRulesPerPod(pod *api.Pod, allPolicies *networking.NetworkPolicyList)
func ListEgressRulesPerPod(pod *api.Pod, allPolicies *networking.NetworkPolicyList)
```
- List all the pods (out of a pod list) that get affected by a policy:
```
func ListPodsPerPolicy(np *networking.NetworkPolicy, allPods *api.PodList)
```
- Decide if a  policy applies to Ingress and//or Egress:
```
func IsPolicyApplicableToIngress(policy *networking.NetworkPolicy)
func IsPolicyApplicableToEgress(policy *networking.NetworkPolicy)
```

- Decide if a Pod gets affected on Ingress//Egress by at least one of the Policies out of a list:
```
func IsPodSelected(pod *api.Pod, policies *networking.NetworkPolicyList)
```

## CLI implementation

As an example, Kubepox can be used with a CLI tool that connects to Kubernetes API  in order to display the policy logic

```
Usage:
kubepox [--config <config>][--namespace <namespace>] get-all (policies|pods)
kubepox [--config <config>][--namespace <namespace>] get-pods <policy>
kubepox [--config <config>][--namespace <namespace>] get-policies <pod>
kubepox [--config <config>][--namespace <namespace>] get-rules <pod>

Options:
--namespace=NAMESPACE Namespace to run the query in (default is "default")
--config=FILE path to the kubeConfig file. (default is ~/.kube/kubeconfig)
```
## How does it work ?

* `kubepox get-all`  retrieves all the NetworkPolicies and Pods. (JSON output, but same API objects as with Kubectl)
* `kubepox get-pods`  retrieves the  podList of affected pods based on a specific policy.
* `kubepox get-policies` retrieves all the policies that apply to a specific pod
* `kubepox get-rules` retrieves all the rules that apply to a specific rule (union of policy rules)

## Example: Rules applied per pod

It is now very easy to see the agglomerate of all the rules that get applied to your Pods. For example:

```
sharma:kubepox bvandewa$ ./kubepox  get-rules redis-django human
Allowed traffic rules for pod redis-django :

------RULE|-----ENTRY|----------------------------------------------------POD SELECTOR|---AND PORT MATCH|
---------1|---------1|----------------------------------------here=frontend,there=ceci|---------tcp:8000|
---------1|---------2|-------------------------------------------------------test=this|-----------------|
---------2|---------1|---role=frontend,testads in (asda,asdd,asdr),tet=tatata,web=ceci|---------tcp:6379|
---------2|---------2|-------------------------------------------------------test=this|---------udp:5000|
```

This comes from the following policies that the pod `redis-django` matches.


Those policies:

```
apiVersion: extensions/v1beta1
kind: NetworkPolicy
metadata:
 name: test-network-policy
spec:
 podSelector:
  matchLabels:
    role: db
 ingress:
  - from:
     - podSelector:
        matchLabels:
         role: frontend
         web: ceci
         tet: tatata
        matchExpressions:
         - key: testads
           operator: In
           values: [asdr,asda,asdd]
     - podSelector:
        matchLabels:
          test: this
    ports:
     - protocol: tcp
       port: 6379
     - protocol: udp
       port: 5000

```

And

```
apiVersion: extensions/v1beta1
kind: NetworkPolicy
metadata:
 name: test-network-policy
spec:
 podSelector:
  matchLabels:
    role: db
 ingress:
  - from:
     - podSelector:
        matchLabels:
         here: frontend
         there: ceci
     - podSelector:
        matchLabels:
          test: this
    ports:
     - protocol: tcp
       port: 8000

```
