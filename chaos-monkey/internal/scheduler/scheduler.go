package scheduler

import (
	"context"
	"math/rand"
	"time"

	"github.com/chaos-monkey/chaos-monkey/internal/metrics"
	"go.uber.org/zap"
)

// Action is a chaos action that can be executed.
type Action interface {
	Name() string
	Execute(ctx context.Context) error
}

// WeightedAction pairs an action with its selection weight.
type WeightedAction struct {
	Action Action
	Weight int
}

// Scheduler runs chaos actions on a weighted random schedule.
type Scheduler struct {
	actions     []WeightedAction
	minInterval time.Duration
	maxInterval time.Duration
	log         *zap.Logger
}

func New(actions []WeightedAction, minInterval, maxInterval time.Duration, log *zap.Logger) *Scheduler {
	return &Scheduler{
		actions:     actions,
		minInterval: minInterval,
		maxInterval: maxInterval,
		log:         log,
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	s.log.Info("chaos scheduler started",
		zap.Duration("min_interval", s.minInterval),
		zap.Duration("max_interval", s.maxInterval),
		zap.Int("num_actions", len(s.actions)),
	)

	for {
		interval := s.randomInterval()
		s.log.Info("next chaos action in", zap.Duration("interval", interval))

		select {
		case <-ctx.Done():
			s.log.Info("scheduler shutting down")
			return
		case <-time.After(interval):
		}

		action := s.selectAction()
		if action == nil {
			s.log.Warn("no action selected, skipping")
			continue
		}

		s.runAction(ctx, action)
	}
}

func (s *Scheduler) runAction(ctx context.Context, action Action) {
	name := action.Name()
	s.log.Info("executing chaos action", zap.String("action", name))

	start := time.Now()
	err := action.Execute(ctx)
	elapsed := time.Since(start)

	metrics.ActionsTotal.WithLabelValues(name).Inc()
	metrics.ActionDuration.WithLabelValues(name).Observe(elapsed.Seconds())
	metrics.LastActionTime.WithLabelValues(name).SetToCurrentTime()

	if err != nil {
		metrics.ActionErrors.WithLabelValues(name).Inc()
		s.log.Error("chaos action failed",
			zap.String("action", name),
			zap.Duration("elapsed", elapsed),
			zap.Error(err),
		)
	} else {
		s.log.Info("chaos action complete",
			zap.String("action", name),
			zap.Duration("elapsed", elapsed),
		)
	}
}

// selectAction picks a random action based on weights.
func (s *Scheduler) selectAction() Action {
	totalWeight := 0
	for _, wa := range s.actions {
		totalWeight += wa.Weight
	}

	if totalWeight == 0 {
		return nil
	}

	r := rand.Intn(totalWeight)
	cumulative := 0
	for _, wa := range s.actions {
		cumulative += wa.Weight
		if r < cumulative {
			return wa.Action
		}
	}
	return s.actions[len(s.actions)-1].Action
}

func (s *Scheduler) randomInterval() time.Duration {
	delta := s.maxInterval - s.minInterval
	if delta <= 0 {
		return s.minInterval
	}
	jitter := time.Duration(rand.Int63n(int64(delta)))
	return s.minInterval + jitter
}
