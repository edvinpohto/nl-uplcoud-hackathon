package actions

import (
	"context"
	"fmt"
)

type NodeStop struct{}

func NewNodeStop() *NodeStop { return &NodeStop{} }

func (a *NodeStop) Name() string { return "node_stop" }

func (a *NodeStop) Execute(_ context.Context) error {
	return fmt.Errorf("node_stop disabled: UpCloud Go SDK requires username/password auth")
}
