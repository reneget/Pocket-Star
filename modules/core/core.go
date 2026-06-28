package core

import (
	"fmt"

	"github.com/anomalyco/my-pretty-star/pkg/module"
)

type CoreModule struct {
	ctx    *module.Context
	status module.Status
}

func (m *CoreModule) Name() string             { return "core" }
func (m *CoreModule) Priority() int             { return 0 }
func (m *CoreModule) Dependencies() []string    { return nil }
func (m *CoreModule) Status() module.Status    { return m.status }

func (m *CoreModule) Init(ctx *module.Context) error {
	m.ctx = ctx
	m.status = module.StatusInactive
	return nil
}

func (m *CoreModule) Start() error {
	if m.ctx == nil {
		return fmt.Errorf("core module not initialized")
	}
	fmt.Printf("core: starting VPN (network=%s, hub=%s)\n", m.ctx.Network, m.ctx.HubAddr)
	m.status = module.StatusActive
	return nil
}

func (m *CoreModule) Stop() error {
	m.status = module.StatusInactive
	return nil
}

func init() {
	module.Register(&CoreModule{})
}
