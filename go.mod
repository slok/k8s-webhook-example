module github.com/slok/k8s-webhook-example

go 1.15

require (
	github.com/oklog/run v1.1.0
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.61.0
	github.com/prometheus/client_golang v1.12.2
	github.com/sirupsen/logrus v1.8.1
	github.com/slok/go-http-metrics v0.6.1
	github.com/slok/kubewebhook/v2 v2.1.1-0.20210813062814-0d6b91199b6d
	github.com/stretchr/testify v1.8.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	k8s.io/api v0.25.4
	k8s.io/apimachinery v0.25.4
)
