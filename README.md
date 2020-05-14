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



[k8s-admission-webhooks]: https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/
[Kubewebhook]: https://github.com/slok/kubewebhook
