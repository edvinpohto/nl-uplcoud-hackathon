package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/chaos-monkey/chaos-monkey/internal/actions"
	"github.com/chaos-monkey/chaos-monkey/internal/config"
	k8sclient "github.com/chaos-monkey/chaos-monkey/internal/k8s"
	"github.com/chaos-monkey/chaos-monkey/internal/scheduler"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	log, _ := zap.NewProduction()
	defer log.Sync()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config", zap.Error(err))
	}

	// Build Kubernetes client for Cluster A (victim)
	k8s, err := k8sclient.NewFromBase64(cfg.KubeconfigData)
	if err != nil {
		log.Fatal("failed to create k8s client", zap.Error(err))
	}
	log.Info("connected to victim cluster")

	// Node-level actions (reboot/stop) are disabled — UpCloud Go SDK requires
	// username/password auth which is unavailable when using API tokens.
	weightedActions := []scheduler.WeightedAction{
		{
			Action: actions.NewPodKill(k8s, cfg.VictimNamespace, log),
			Weight: cfg.WeightPodKill,
		},
		{
			Action: actions.NewNetworkDisrupt(k8s, cfg.NetworkDisruptDelay, log),
			Weight: cfg.WeightNetworkDisrupt,
		},
		{
			Action: actions.NewScaleZero(k8s, cfg.VictimNamespace, "victim-app", cfg.ScaleZeroDuration, log),
			Weight: cfg.WeightScaleZero,
		},
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	sched := scheduler.New(weightedActions, cfg.MinInterval, cfg.MaxInterval, log)
	go sched.Run(ctx)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	srv := &http.Server{
		Addr:    ":" + cfg.MetricsPort,
		Handler: mux,
	}

	go func() {
		log.Info("metrics server listening", zap.String("port", cfg.MetricsPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("metrics server error", zap.Error(err))
		}
	}()

	<-ctx.Done()
	log.Info("shutting down chaos monkey")
}
