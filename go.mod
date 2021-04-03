module github.com/slok/k8s-webhook-example

go 1.15

require (
	github.com/coreos/prometheus-operator v0.39.0
	github.com/oklog/run v1.1.0
	github.com/prometheus/client_golang v1.10.0
	github.com/sirupsen/logrus v1.8.1
	github.com/slok/go-http-metrics v0.6.1
	github.com/slok/kubewebhook/v2 v2.0.0
	github.com/stretchr/testify v1.7.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	k8s.io/api v0.20.5
	k8s.io/apimachinery v0.20.5
)

replace k8s.io/client-go => k8s.io/client-go v0.20.5
