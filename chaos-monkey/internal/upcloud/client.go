// Package upcloud is a stub — node-level chaos actions (reboot/stop) require
// UpCloud API basic auth (username+password) which is not available when using
// API tokens. Node actions are disabled; pod_kill, network_disrupt, and
// scale_zero run via the Kubernetes API instead.
package upcloud

import "go.uber.org/zap"

type Client struct {
	log *zap.Logger
}

func New(log *zap.Logger) *Client {
	return &Client{log: log}
}
