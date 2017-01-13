# kubepox

[![Twitter URL](https://img.shields.io/badge/twitter-follow-blue.svg)](https://twitter.com/aporeto_trireme) [![Slack URL](https://img.shields.io/badge/slack-join-green.svg)](https://triremehq.slack.com/messages/general/) [![Documentation](https://img.shields.io/badge/docs-godoc-blue.svg)](https://godoc.org/github.com/aporeto-inc/trireme)
[![Analytics](https://ga-beacon.appspot.com/UA-90298211-1/welcome-page)](https://github.com/igrigorik/ga-beacon)

Kubernetes network Policy eXploration tool

## Library

kubepox is a super lightweight library that implements a simple set of logic rules in order to get exactly the Policies/Rules that apply to a specific Pod.
kubepox doesn't connect to Kubernetes API. It takes into parameter the result of API Calls. Those API Calls have to be handled outside of kubepox (look at the CLI code for examples).


## CLI Tool

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
