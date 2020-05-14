apiVersion: admissionregistration.k8s.io/v1beta1
kind: MutatingWebhookConfiguration
metadata:
  name: all-mark-webhook
  labels:
    app: all-mark-webhook
    kind: mutator
webhooks:
  - name: all-mark-webhook.slok.dev
    clientConfig:
      service:
        name: k8s-webhook-example
        namespace: k8s-webhook-example
        path: /wh/mutating/allmark
      caBundle: CA_BUNDLE
    rules:
      - operations: ["CREATE", "UPDATE"]
        apiGroups: ["*"]
        apiVersions: ["*"]
        resources: ["*"]