# k8s-webhook-example

A production ready [Kubernetes admission webhook][k8s-admission-webhooks] example using [Kubewebhook].

The example tries showing these:

- How to set up a production ready Kubernetes admission webhook.
  - Clean and decouple structure.
  - Metrics.
  - Gracefull shutdown.
  - Testing webhooks.
- Serve multiple webhooks on the same application.
- Mutating and validating webhooks with different use cases (check Webhooks section)
- Mutating/validating webhook chains (sequential mutations/validations).

## Webhooks

### `all-mark-webhook.slok.dev`

- Webhook type: Mutating.
- Resources affected: `deployments`, `daemonsets`, `cronjobs`, `jobs`, `statefulsets`, `pods`

This webhooks shows how to add a label to all the specified types in a generic way. 

We use dynamic webhooks without the requirement to know what type of objects we are dealing with. This is becase all the types implement `metav1.Object` interface that accesses to the metadata of the object. In this case our domain logic doesn't need to know what type is.

### `ingress-validation-webhook.slok.dev`

- Webhook type: Validating.
- Resources affected: Ingresses.

This webhook has a chain of validation on ingress objects, it is composed of 2 validations:

- Check an ingress has a single host/rule.
- Check an ingress host matches specific regexes.

This webhook shows two things:

First, shows how to create a chain of validations for a single webhook handler.

Second, it shows how to deal with specific types of resources in different group/versions, for this it uses a dynamic webhook (like `all-mark-webhook.slok.dev`) but this instead, typecasts to the specific types, in this case, the webhook validates all available ingresses, specifically `extensions/v1beta1` and `networking.k8s.io/v1beta1`.


[k8s-admission-webhooks]: https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/
[Kubewebhook]: https://github.com/slok/kubewebhook
