package actions

import (
	"context"
	"fmt"
	"time"

	"github.com/chaos-monkey/chaos-monkey/internal/k8s"
	"go.uber.org/zap"
)

type ScaleZero struct {
	client     *k8s.Client
	namespace  string
	deployment string
	duration   time.Duration
	log        *zap.Logger
}

func NewScaleZero(client *k8s.Client, namespace, deployment string, duration time.Duration, log *zap.Logger) *ScaleZero {
	return &ScaleZero{
		client:     client,
		namespace:  namespace,
		deployment: deployment,
		duration:   duration,
		log:        log,
	}
}

func (a *ScaleZero) Name() string { return "scale_zero" }

func (a *ScaleZero) Execute(ctx context.Context) error {
	// Get current replica count before scaling down
	currentReplicas, err := a.client.GetDeploymentReplicas(ctx, a.namespace, a.deployment)
	if err != nil {
		return fmt.Errorf("get current replicas: %w", err)
	}

	if currentReplicas == 0 {
		return fmt.Errorf("deployment %s is already at 0 replicas", a.deployment)
	}

	a.log.Info("scaling deployment to 0",
		zap.String("deployment", a.deployment),
		zap.Int32("previous_replicas", currentReplicas),
		zap.Duration("down_duration", a.duration),
	)

	if err := a.client.ScaleDeployment(ctx, a.namespace, a.deployment, 0); err != nil {
		return fmt.Errorf("scale to 0: %w", err)
	}

	select {
	case <-ctx.Done():
		a.log.Warn("context cancelled during scale-zero, restoring replicas")
		a.client.ScaleDeployment(context.Background(), a.namespace, a.deployment, currentReplicas)
		return ctx.Err()
	case <-time.After(a.duration):
	}

	a.log.Info("restoring replicas", zap.String("deployment", a.deployment), zap.Int32("replicas", currentReplicas))
	if err := a.client.ScaleDeployment(ctx, a.namespace, a.deployment, currentReplicas); err != nil {
		return fmt.Errorf("restore replicas: %w", err)
	}

	a.log.Info("scale-zero complete", zap.String("deployment", a.deployment))
	return nil
}
