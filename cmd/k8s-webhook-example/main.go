package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/slok/k8s-webhook-example/internal/http/webhook"
	"github.com/slok/k8s-webhook-example/internal/log"
	internalprometheus "github.com/slok/k8s-webhook-example/internal/metrics/prometheus"
	"github.com/slok/k8s-webhook-example/internal/mutation/mark"
	"github.com/slok/k8s-webhook-example/internal/validation/ingress"
)

var (
	// Version is set at compile time.
	Version = "dev"
)

func runApp() error {
	cfg, err := NewCmdConfig()
	if err != nil {
		return fmt.Errorf("could not get commandline configuration: %w", err)
	}

	// Set up logger.
	logrusLog := logrus.New()
	logrusLogEntry := logrus.NewEntry(logrusLog).WithField("app", "k8s-webhook-example")
	if cfg.Debug {
		logrusLogEntry.Logger.SetLevel(logrus.DebugLevel)
	}
	if !cfg.Development {
		logrusLogEntry.Logger.SetFormatter(&logrus.JSONFormatter{})
	}
	logger := log.NewLogrus(logrusLogEntry).WithKV(log.KV{"version": Version})

	// Dependencies.
	metricsRec := internalprometheus.NewRecorder(prometheus.DefaultRegisterer)

	var marker mark.Marker
	if len(cfg.LabelMarks) > 0 {
		marker = mark.NewLabelMarker(cfg.LabelMarks)
		logger.Infof("label marker webhook enabled")
	} else {
		marker = mark.DummyMarker
		logger.Warningf("label marker webhook disabled")
	}

	var ingressHostValidator ingress.Validator
	if len(cfg.IngressHostRegexes) > 0 {
		ingressHostValidator, err = ingress.NewHostRegexValidator(cfg.IngressHostRegexes)
		if err != nil {
			return fmt.Errorf("could not create ingress regex host validator: %w", err)
		}
		logger.Infof("ingress host regex validation webhook enabled")
	} else {
		ingressHostValidator = ingress.DummyValidator
		logger.Warningf("ingress host regex validation webhook disabled")
	}

	var ingressSingleHostValidator ingress.Validator
	if cfg.EnableIngressSingleHost {
		ingressSingleHostValidator = ingress.SingleHostValidator
		logger.Infof("ingress single host validation webhook enabled")
	} else {
		ingressSingleHostValidator = ingress.DummyValidator
		logger.Warningf("ingress single host validation webhook disabled")
	}

	// Prepare run entrypoints.
	var g run.Group

	// OS signals.
	{
		sigC := make(chan os.Signal, 1)
		exitC := make(chan struct{})
		signal.Notify(sigC, syscall.SIGTERM, syscall.SIGINT)

		g.Add(
			func() error {
				select {
				case s := <-sigC:
					logger.Infof("signal %s received", s)
					return nil
				case <-exitC:
					return nil
				}
			},
			func(_ error) {
				close(exitC)
			},
		)
	}

	// Metrics HTTP server.
	{
		logger := logger.WithKV(log.KV{"addr": cfg.MetricsListenAddr, "http-server": "metrics"})
		mux := http.NewServeMux()

		// Metrics.
		mux.Handle(cfg.MetricsPath, promhttp.Handler())

		// Pprof.
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

		// Health checks.
		mux.HandleFunc("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

		server := http.Server{Addr: cfg.MetricsListenAddr, Handler: mux}

		g.Add(
			func() error {
				logger.Infof("http server listening...")
				return server.ListenAndServe()
			},
			func(_ error) {
				logger.Infof("start draining connections")
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				err := server.Shutdown(ctx)
				if err != nil {
					logger.Errorf("error while shutting down the server: %s", err)
				} else {
					logger.Infof("server stopped")
				}
			},
		)
	}

	// Webhook HTTP server.
	{
		logger := logger.WithKV(log.KV{"addr": cfg.WebhookListenAddr, "http-server": "webhooks"})

		// Webhook handler.
		wh, err := webhook.New(webhook.Config{
			Marker:                     marker,
			IngressRegexHostValidator:  ingressHostValidator,
			IngressSingleHostValidator: ingressSingleHostValidator,
			MetricsRecorder:            metricsRec,
			Logger:                     logger,
		})
		if err != nil {
			return fmt.Errorf("could not create webhooks handler: %w", err)
		}

		mux := http.NewServeMux()
		mux.Handle("/", wh)
		server := http.Server{Addr: cfg.WebhookListenAddr, Handler: mux}

		g.Add(
			func() error {
				if cfg.TLSCertFilePath == "" || cfg.TLSKeyFilePath == "" {
					logger.Warningf("webhook running without TLS")
					logger.Infof("http server listening...")
					return server.ListenAndServe()
				}

				logger.Infof("https server listening...")
				return server.ListenAndServeTLS(cfg.TLSCertFilePath, cfg.TLSKeyFilePath)
			},
			func(_ error) {
				logger.Infof("start draining connections")
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				err := server.Shutdown(ctx)
				if err != nil {
					logger.Errorf("error while shutting down the server: %s", err)
				} else {
					logger.Infof("server stopped")
				}
			},
		)
	}

	err = g.Run()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := runApp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running app: %s", err)
		os.Exit(1)
	}

	os.Exit(0)
}
