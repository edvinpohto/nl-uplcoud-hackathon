package actions

import (
	"context"
	"fmt"
)

type NodeReboot struct{}

func NewNodeReboot() *NodeReboot { return &NodeReboot{} }

func (a *NodeReboot) Name() string { return "node_reboot" }

func (a *NodeReboot) Execute(_ context.Context) error {
	return fmt.Errorf("node_reboot disabled: UpCloud Go SDK requires username/password auth")
}
