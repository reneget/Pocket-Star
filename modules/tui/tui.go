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
}

func (m *TUIModule) Name() string             { return "tui" }
func (m *TUIModule) Priority() int             { return 10 }
func (m *TUIModule) Dependencies() []string    { return []string{"core"} }
func (m *TUIModule) Status() module.Status    { return m.status }

func (m *TUIModule) Init(ctx *module.Context) error {
	m.ctx = ctx
	m.model = initialModel(ctx)
	m.status = module.StatusInactive
	return nil
}

func (m *TUIModule) Start() error {
	if m.status == module.StatusActive {
		return nil
	}

	m.prog = tea.NewProgram(
		m.model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	go func() {
		if _, err := m.prog.Run(); err != nil {
			log.Default.Push(fmt.Sprintf("tui exited: %v", err))
		}
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

func Run(ctx *module.Context) error {
	m := &TUIModule{}
	if err := m.Init(ctx); err != nil {
		return fmt.Errorf("tui init: %w", err)
	}
	if err := m.Start(); err != nil {
		return fmt.Errorf("tui start: %w", err)
	}

	prog := tea.NewProgram(
		m.model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := prog.Run(); err != nil {
		return fmt.Errorf("tui: %w", err)
	}

	return nil
}

func init() {
	module.Register(&TUIModule{})
}
