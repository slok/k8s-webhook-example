package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

// CmdConfig represents the configuration of the command.
type CmdConfig struct {
	Debug             bool
	Development       bool
	WebhookListenAddr string
	MetricsListenAddr string
	MetricsPath       string
	TLSCertFilePath   string
	TLSKeyFilePath    string
	LabelMarks        map[string]string
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
	app.Flag("webhook-label-marks", "the marks the webhook will set to all resources, if no marks, the label marker webhook will be disabled .").Short('l').StringMapVar(&c.LabelMarks)

	_, err := app.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}

	return c, nil
}
