# kubepox

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
