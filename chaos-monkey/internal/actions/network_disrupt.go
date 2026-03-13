package actions

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/chaos-monkey/chaos-monkey/internal/k8s"
	"go.uber.org/zap"
)

const (
	pumbaNamespace     = "victim"
	pumbaLabelSelector = "app=pumba"
	pumbaContainer     = "pumba"
)

type NetworkDisrupt struct {
	client *k8s.Client
	delay  string
	log    *zap.Logger
}

func NewNetworkDisrupt(client *k8s.Client, delay string, log *zap.Logger) *NetworkDisrupt {
	return &NetworkDisrupt{client: client, delay: delay, log: log}
}

func (a *NetworkDisrupt) Name() string { return "network_disrupt" }

func (a *NetworkDisrupt) Execute(ctx context.Context) error {
	pods, err := a.client.ListDaemonSetPods(ctx, pumbaNamespace, pumbaLabelSelector)
	if err != nil {
		return fmt.Errorf("list pumba pods: %w", err)
	}

	if len(pods) == 0 {
		return fmt.Errorf("no Pumba DaemonSet pods found")
	}

	// Pick a random Pumba pod (one per node)
	target := pods[rand.Intn(len(pods))]
	a.log.Info("executing network disruption via Pumba",
		zap.String("pod", target.Name),
		zap.String("node", target.Spec.NodeName),
		zap.String("delay", a.delay),
	)

	// Parse delay value (e.g. "200ms" -> 200)
	delayMs := 200
	if a.delay != "" {
		var ms int
		fmt.Sscanf(a.delay, "%dms", &ms)
		if ms > 0 {
			delayMs = ms
		}
	}

	// Use tc netem to add delay on the node's default interface for 30s
	cmd := []string{
		"sh", "-c",
		fmt.Sprintf(
			"IF=$(ip route 2>/dev/null | awk '/default/ {print $5; exit}'); IF=${IF:-eth0}; "+
				"tc qdisc replace dev $IF root netem delay %dms 2>/dev/null; "+
				"sleep 30; "+
				"tc qdisc del dev $IF root 2>/dev/null || true",
			delayMs,
		),
	}

	stdout, stderr, err := a.client.ExecInPod(ctx, pumbaNamespace, target.Name, pumbaContainer, cmd)
	if err != nil {
		a.log.Warn("pumba exec output",
			zap.String("stdout", stdout),
			zap.String("stderr", stderr),
		)
		return fmt.Errorf("exec pumba: %w", err)
	}

	a.log.Info("network disruption applied",
		zap.String("pod", target.Name),
		zap.String("stdout", stdout),
	)
	return nil
}
