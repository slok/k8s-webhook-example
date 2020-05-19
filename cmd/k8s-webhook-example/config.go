package main

import (
	"os"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

// CmdConfig represents the configuration of the command.
type CmdConfig struct {
	Debug                   bool
	Development             bool
	WebhookListenAddr       string
	MetricsListenAddr       string
	MetricsPath             string
	TLSCertFilePath         string
	TLSKeyFilePath          string
	EnableIngressSingleHost bool
	IngressHostRegexes      []string
	MinSMScrapeInterval     time.Duration

	LabelMarks map[string]string
}

// NewCmdConfig returns a new command configuration.
func NewCmdConfig() (*CmdConfig, error) {
	c := &CmdConfig{
		LabelMarks: map[string]string{},
	}
	app := kingpin.New("k8s-webhook-example", "A Kubernetes production-ready admission webhook example.")
	app.Version(Version)

	app.Flag("debug", "Enable debug mode.").BoolVar(&c.Debug)
	app.Flag("development", "Enable development mode.").BoolVar(&c.Development)
	app.Flag("webhook-listen-address", "the address where the HTTPS server will be listening to serve the webhooks.").Default(":8080").StringVar(&c.WebhookListenAddr)
	app.Flag("metrics-listen-address", "the address where the HTTP server will be listening to serve metrics, healthchecks, profiling...").Default(":8081").StringVar(&c.MetricsListenAddr)
	app.Flag("metrics-path", "the path where Prometheus metrics will be served.").Default("/metrics").StringVar(&c.MetricsPath)
	app.Flag("tls-cert-file-path", "the path for the webhook HTTPS server TLS cert file.").StringVar(&c.TLSCertFilePath)
	app.Flag("tls-key-file-path", "the path for the webhook HTTPS server TLS key file.").StringVar(&c.TLSKeyFilePath)
	app.Flag("webhook-label-marks", "a map of labels the webhook will set to all resources, if no labels, the label marker webhook will be disabled. Can repeat flag").Short('l').StringMapVar(&c.LabelMarks)
	app.Flag("webhook-enable-ingress-single-host", "enables validation of ingress to have only a single host/rule.").Short('s').BoolVar(&c.EnableIngressSingleHost)
	app.Flag("webhook-ingress-host-regex", "a list of regexes that will validate ingress hosts matching against this regexes, no host disables validation webhook. Can repeat flag.").Short('h').StringsVar(&c.IngressHostRegexes)
	app.Flag("webhook-sm-min-scrape-interval", "the minimum screate interval service monitors can have.").DurationVar(&c.MinSMScrapeInterval)

	_, err := app.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}

	return c, nil
}
