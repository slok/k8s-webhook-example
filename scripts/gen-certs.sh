#! /bin/bash

WEBHOOK_NS=k8s-webhook-example
WEBHOOK_NAME=k8s-webhook-example
WEBHOOK_SVC=${WEBHOOK_NAME}.${WEBHOOK_NS}.svc
OUT_CERT_FILE=./deploy/app-certs.yaml
OUT_WEBBOK_FILE=./deploy/webhooks.yaml

# Create certs for our webhook
openssl genrsa -out webhookCA.key 2048
openssl req -new -key ./webhookCA.key -subj "/CN=${WEBHOOK_SVC}" -out ./webhookCA.csr 
openssl x509 -req -days 365 -in webhookCA.csr -signkey webhookCA.key -out webhook.crt

# Create certs secrets for k8s
kubectl -n ${WEBHOOK_NS} create secret generic \
    ${WEBHOOK_NAME}-certs \
    --from-file=key.pem=./webhookCA.key \
    --from-file=cert.pem=./webhook.crt \
    --dry-run -o yaml > ${OUT_CERT_FILE}

# Set the CABundle on the webhook registration
CA_BUNDLE=$(cat ./webhook.crt | base64 -w0)
sed "s/CA_BUNDLE/${CA_BUNDLE}/" ./deploy/webhooks.yaml.tpl > ${OUT_WEBBOK_FILE}


# Clean
rm ./webhookCA* && rm ./webhook.crt

# Add note of autogenerated file.
sed -i '1i# File autogenerated by ./scripts/gen-certs.sh' ${OUT_CERT_FILE}
sed -i '1i# File autogenerated by ./scripts/gen-certs.sh' ${OUT_WEBBOK_FILE}