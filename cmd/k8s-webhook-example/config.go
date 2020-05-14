package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

// CmdConfig represents the configuration of the command.
type CmdConfig struct {
	Debug           bool
	Development     bool
	ListenAddr      string
	MetricsPath     string
	TLSCertFilePath string
	TLSKeyFilePath  string
	LabelMarks      map[string]string
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
	app.Flag("listen-address", "the address where the HTTP server will be listening.").Default(":8080").StringVar(&c.ListenAddr)
	app.Flag("metrics-path", "the path where Prometheus metrics will be served.").Default("/metrics").StringVar(&c.MetricsPath)
	app.Flag("tls-cert-file-path", "the path for the webhook HTTP server TLS cert file.").StringVar(&c.TLSCertFilePath)
	app.Flag("tls-key-file-path", "the path for the webhook HTTP server TLS key file.").StringVar(&c.TLSKeyFilePath)
	app.Flag("webhook-label-marks", "the marks the webhook will set to all resources, if no marks, the label marker webhook will be disabled .").Short('l').StringMapVar(&c.LabelMarks)

	_, err := app.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}

	return c, nil
}
