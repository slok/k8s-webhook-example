# k8s-webhook-example

A production ready [Kubernetes admission webhook][k8s-admission-webhooks] example using [Kubewebhook].

The example tries showing these:

- How to set up a production ready Kubernetes admission webhook.
  - Clean and decouple structure.
  - Metrics.
  - Gracefull shutdown.
  - Testing webhooks.
- Serve multiple webhooks on the same application.
- Mutating and validating webhooks with different use cases (check Webhooks section).

## Structure

The application is mainly structured in 3 parts:

- `main`: This is where everything is created, wired, configured and set up, [cmd/k8s-webhook-example](cmd/k8s-webhook-example/main.go).
- `http`: This is the package that configures the HTTP server, wires the routes and the webhook handlers. [internal/http/webhook](internal/http/webhook).
- Application services: These services have the domain logic of the validators and mutators:
  - [`mutation/mark`](internal/mutation/mark): Logic for `all-mark-webhook.slok.dev` webhook.
  - [`validation/ingress`](internal/validation/ingress): Logic for `ingress-validation-webhook.slok.dev` webhook.
  - [`mutation/prometheus`](internal/mutation/prometheus): Logic for `service-monitor-safer.slok.dev` webhook.

Apart from the webhook refering stuff we have other parts like:

- [Decoupled metrics](internal/metrics)
- [Decoupled logger](internal/log)
- [Application command line flags](cmd/k8s-webhook-example/config.go)

And finally there is an example of how we could deploy our webhooks on a production server:

- [Deploy](deploy)

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

### `service-monitor-safer.slok.dev`

- Webhook type: Mutating.
- Resources affected: [ServiceMonitors] (`monitoring.coreos.com/v1`) CRD.

This webhook show two things.

- Working with CRDs, in this case mutating them.
- Working with Static webhooks (specific type).

This webhook takes Prometheus `monitoring.coreos.com/v1/servicemonitors` CRs and sets safe scraping intervals, it checks the interval and in case is missing or is less that the minimum configured it will mutate the CR to set the minimum scrape interval.

This will show us how to deal with CRDs in webhooks, and also how we can make static webhooks to only work safely in a specific resource type.

The static webhooks are specially important on resources that are not known, these are:

- CRDs.
- Core resources that are on the cluster but not on the webhook libraries, because of Kubernetes different versions (new types and deprecations from version to version).

If we use dynamic webhook on unknown types by our webhook app, we will deal with `runtime.Unstructured`, this is not bad and is safe, it would add complexity to mutate/validate these objects, although for mutating/validating metadata fields (e.g `labels`), is easy and simple.

That said, most webhooks can/should use dynamic type webhooks because are common resources, like `ingress-validation-webhook.slok.dev`, `all-mark-webhook.slok.dev`, that use dynamic webhooks correctly.

[k8s-admission-webhooks]: https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/
[kubewebhook]: https://github.com/slok/kubewebhook
[servicemonitors]: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#servicemonitor
