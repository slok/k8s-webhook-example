# k8s-webhook-example

A production ready [Kubernetes admission webhook][k8s-admission-webhooks] example using [Kubewebhook].

The example tries showing these:

- How to set up a production ready Kubernetes admission webhook.
  - Clean and decouple structure.
  - Metrics.
  - Gracefull shutdown.
  - Testing webhooks.
- Serve multiple webhooks on the same application.
- Mutating and validating webhooks with different combination of resource types, e.g:
  - Metadata based mutation/validation on all resource types.
  - Single specific types.
  - Multiple sppecific types.
  - CRD types.
- Mutating/validating webhook chains (sequential mutations/validations).

## Webhooks

### Label marker

- Webhook type: Mutating.
- Effect: Adds a label.
- Resources affected: All resources.
- Shows: How to mutate any resource in a generic way.

### Ingress single host validation

- Webhook type: Validating.
- Effect: Validates that an ingress only has one host.
- Resources affected: Ingresses.
- Shows: How to validate an specific type.

### Ingress host regex validation

- Webhook type: Validating.
- Effect: Validates that an ingress host matches a regex.
- Resources affected: Ingresses.
- Shows: How to validate an specific type chained with the ingress single host validation.

### Single replica pods are marked

- Webhook type: Mutating.
- Effect: Pods with <2 replicas will be marked as dangerous.
- Resources affected: pods.
- Shows: Mutating specific single type with logic conditions.

### Deployments and statefulsets policy check

- Webhook type: Validating.
- Effect: pods with <2 replicas will be marked as dangerous.
- Resources affected: pods.
- Shows: Validating specific multiple types.

[k8s-admission-webhooks]: https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/
[Kubewebhook]: https://github.com/slok/kubewebhook
