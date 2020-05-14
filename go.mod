module github.com/slok/k8s-webhook-example

go 1.14

require (
	github.com/oklog/run v1.1.0
	github.com/prometheus/client_golang v1.6.0
	github.com/sirupsen/logrus v1.6.0
	github.com/slok/go-http-metrics v0.6.1
	github.com/slok/kubewebhook v0.9.2-0.20200512083742-39c518541235
	github.com/stretchr/testify v1.5.1
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
)
