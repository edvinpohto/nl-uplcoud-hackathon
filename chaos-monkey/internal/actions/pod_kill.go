package actions

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/chaos-monkey/chaos-monkey/internal/k8s"
	"go.uber.org/zap"
)

type PodKill struct {
	client    *k8s.Client
	namespace string
	log       *zap.Logger
}

func NewPodKill(client *k8s.Client, namespace string, log *zap.Logger) *PodKill {
	return &PodKill{client: client, namespace: namespace, log: log}
}

func (a *PodKill) Name() string { return "pod_kill" }

func (a *PodKill) Execute(ctx context.Context) error {
	pods, err := a.client.ListRunningPods(ctx, a.namespace)
	if err != nil {
		return fmt.Errorf("list pods: %w", err)
	}

	if len(pods) == 0 {
		return fmt.Errorf("no running pods in namespace %s", a.namespace)
	}

	target := pods[rand.Intn(len(pods))]
	a.log.Info("killing pod", zap.String("pod", target.Name), zap.String("namespace", a.namespace))

	if err := a.client.DeletePod(ctx, a.namespace, target.Name); err != nil {
		return fmt.Errorf("delete pod %s: %w", target.Name, err)
	}

	a.log.Info("pod killed", zap.String("pod", target.Name))
	return nil
}
