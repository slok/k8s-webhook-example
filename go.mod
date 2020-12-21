module github.com/slok/k8s-webhook-example

go 1.15

require (
	github.com/coreos/prometheus-operator v0.39.0
	github.com/oklog/run v1.1.0
	github.com/prometheus/client_golang v1.8.0
	github.com/sirupsen/logrus v1.6.0
	github.com/slok/go-http-metrics v0.6.1
	github.com/slok/kubewebhook/v2 v2.0.0-20201221081759-77a78fd03dd6
	github.com/stretchr/testify v1.6.1
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	k8s.io/api v0.19.6
	k8s.io/apimachinery v0.19.6
)

replace k8s.io/client-go => k8s.io/client-go v0.19.6
