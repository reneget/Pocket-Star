package tui

import (
	"fmt"

	"github.com/anomalyco/my-pretty-star/pkg/log"
	"github.com/anomalyco/my-pretty-star/pkg/module"
	tea "github.com/charmbracelet/bubbletea"
)

type TUIModule struct {
	ctx    *module.Context
	status module.Status
	model  model
	prog   *tea.Program
	doneCh chan struct{}
	mgr    *module.Manager
}

func (m *TUIModule) Name() string             { return "tui" }
func (m *TUIModule) Priority() int             { return 10 }
func (m *TUIModule) Dependencies() []string    { return nil }
func (m *TUIModule) Status() module.Status    { return m.status }
func (m *TUIModule) SetModuleManager(mgr *module.Manager) { m.mgr = mgr }

func (m *TUIModule) Init(ctx *module.Context) error {
	m.ctx = ctx
	m.doneCh = make(chan struct{})
	m.status = module.StatusInactive
	return nil
}

func (m *TUIModule) Start() error {
	if m.status == module.StatusActive {
		return nil
	}

	// model создаётся здесь, после того как mgr уже установлен
	m.model = initialModel(m.ctx, m.mgr)

	m.prog = tea.NewProgram(
		m.model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	go func() {
		if _, err := m.prog.Run(); err != nil {
			log.Default.Push(fmt.Sprintf("tui exited: %v", err))
		}
		close(m.doneCh)
	}()

	m.status = module.StatusActive
	log.Default.Push("tui: started")
	return nil
}

func (m *TUIModule) Stop() error {
	if m.prog != nil {
		m.prog.Quit()
	}
	m.status = module.StatusInactive
	log.Default.Push("tui: stopped")
	return nil
}

func (m *TUIModule) Wait() {
	<-m.doneCh
}

func init() {
	module.Register(&TUIModule{})
}
