# kubepox

Kubernetes network Policy eXploration tool


```
Usage:
kubepox [--config <config>][--namespace <namespace>] get (np|pod)
kubepox [--config <config>][--namespace <namespace>] affect (np|pod) (<name>)

Options:
--namespace=NAMESPACE Namespace to run the query in
--config=FILE path to the KubeConfig file.
```
## How does it work ?

* kubepox get  allows you to get all the NetworkPolicies and Pods. (JSON output, but same API objects as with Kubectl)
* kubepox affect allows you to retrieve the networkpolicyList and podList of affected objects based on a specific pod or policy.
