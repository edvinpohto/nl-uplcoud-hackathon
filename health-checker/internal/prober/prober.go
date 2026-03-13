package prober

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/chaos-monkey/health-checker/internal/metrics"
	"go.uber.org/zap"
)

type Prober struct {
	targetURL string
	interval  time.Duration
	timeout   time.Duration
	client    *http.Client
	log       *zap.Logger
}

func New(targetURL string, interval, timeout time.Duration, log *zap.Logger) *Prober {
	return &Prober{
		targetURL: targetURL,
		interval:  interval,
		timeout:   timeout,
		client: &http.Client{
			Timeout: timeout,
		},
		log: log,
	}
}

func (p *Prober) Run(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	consecutiveFailures := 0

	p.log.Info("starting prober", zap.String("target", p.targetURL), zap.Duration("interval", p.interval))

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.probe(&consecutiveFailures)
		}
	}
}

func (p *Prober) probe(consecutiveFailures *int) {
	start := time.Now()
	metrics.ProbeTotal.Inc()

	req, err := http.NewRequest(http.MethodGet, p.targetURL, nil)
	if err != nil {
		p.recordFailure(consecutiveFailures, fmt.Sprintf("create request: %v", err))
		return
	}

	resp, err := p.client.Do(req)
	elapsed := time.Since(start)

	if err != nil {
		p.recordFailure(consecutiveFailures, fmt.Sprintf("http get: %v", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		metrics.ResponseTime.Observe(elapsed.Seconds())
		metrics.ProbeSuccess.Set(1)
		metrics.ConsecutiveFailures.Set(0)
		*consecutiveFailures = 0
		p.log.Debug("probe ok", zap.Int("status", resp.StatusCode), zap.Duration("latency", elapsed))
	} else {
		p.recordFailure(consecutiveFailures, fmt.Sprintf("unexpected status %d", resp.StatusCode))
	}
}

func (p *Prober) recordFailure(consecutiveFailures *int, reason string) {
	*consecutiveFailures++
	metrics.ProbeSuccess.Set(0)
	metrics.ProbeErrors.Inc()
	metrics.ConsecutiveFailures.Set(float64(*consecutiveFailures))
	p.log.Warn("probe failed", zap.String("reason", reason), zap.Int("consecutive_failures", *consecutiveFailures))
}
