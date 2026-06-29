package module

import (
	"fmt"
	"sort"
	"strings"
)

type ModuleInfo struct {
	Name         string
	Priority     int
	Dependencies []string
	Status       Status
	Enabled      bool
}

type Manager struct {
	ctx     *Context
	modules map[string]Module
	order   []string
	enabled map[string]bool
	active  map[string]Module
}

func NewManager(ctx *Context) *Manager {
	m := &Manager{
		ctx:     ctx,
		modules: make(map[string]Module),
		enabled: make(map[string]bool),
		active:  make(map[string]Module),
	}

	reg := GetRegistry()
	sort.Slice(reg, func(i, j int) bool {
		return reg[i].Priority() < reg[j].Priority()
	})

	for _, mod := range reg {
		m.modules[mod.Name()] = mod
		m.order = append(m.order, mod.Name())
	}

	for _, name := range m.order {
		m.enabled[name] = true
	}

	return m
}

func (m *Manager) InitAll() error {
	for _, name := range m.order {
		mod := m.modules[name]
		if err := mod.Init(m.ctx); err != nil {
			return fmt.Errorf("init %s: %w", name, err)
		}
	}
	return nil
}

func (m *Manager) StartAll() error {
	for _, name := range m.order {
		if !m.enabled[name] {
			continue
		}
		if err := m.checkDeps(name); err != nil {
			return err
		}
		mod := m.modules[name]
		if err := mod.Start(); err != nil {
			return fmt.Errorf("start %s: %w", name, err)
		}
		m.active[name] = mod
	}
	return nil
}

func (m *Manager) StopAll() {
	for i := len(m.order) - 1; i >= 0; i-- {
		name := m.order[i]
		if mod, ok := m.active[name]; ok {
			mod.Stop()
			delete(m.active, name)
		}
	}
}

func (m *Manager) Enable(name string) error {
	if _, ok := m.modules[name]; !ok {
		return fmt.Errorf("module %q not found", name)
	}
	m.enabled[name] = true

	if mod, ok := m.active[name]; ok && mod.Status() == StatusActive {
		return nil
	}

	if err := m.checkDeps(name); err != nil {
		m.enabled[name] = false
		return err
	}

	mod := m.modules[name]
	if err := mod.Start(); err != nil {
		m.enabled[name] = false
		return fmt.Errorf("start %s: %w", name, err)
	}
	m.active[name] = mod
	return nil
}

func (m *Manager) Disable(name string) error {
	if _, ok := m.modules[name]; !ok {
		return fmt.Errorf("module %q not found", name)
	}
	m.enabled[name] = false

	if mod, ok := m.active[name]; ok {
		mod.Stop()
		delete(m.active, name)
	}

	for _, depName := range m.order {
		mod := m.modules[depName]
		for _, dep := range mod.Dependencies() {
			if dep == name {
				m.Disable(depName)
			}
		}
	}

	return nil
}

func (m *Manager) Get(name string) Module {
	return m.modules[name]
}

func (m *Manager) IsEnabled(name string) bool {
	return m.enabled[name]
}

func (m *Manager) IsActive(name string) bool {
	_, ok := m.active[name]
	return ok
}

func (m *Manager) List() []ModuleInfo {
	var out []ModuleInfo
	for _, name := range m.order {
		mod := m.modules[name]
		info := ModuleInfo{
			Name:         name,
			Priority:     mod.Priority(),
			Dependencies: mod.Dependencies(),
			Status:       mod.Status(),
			Enabled:      m.enabled[name],
		}
		out = append(out, info)
	}
	return out
}

func (m *Manager) checkDeps(name string) error {
	mod := m.modules[name]
	for _, dep := range mod.Dependencies() {
		if !m.enabled[dep] {
			return fmt.Errorf("module %q depends on %q which is disabled", name, dep)
		}
		if _, ok := m.active[dep]; !ok {
			if err := m.Enable(dep); err != nil {
				return fmt.Errorf("enable dependency %s for %s: %w", dep, name, err)
			}
		}
	}
	return nil
}

func ModuleListString(infos []ModuleInfo) string {
	var b strings.Builder
	for _, info := range infos {
		status := "● active"
		if !info.Enabled {
			status = "○ disabled"
		} else if info.Status != StatusActive {
			status = "◌ inactive"
		}
		deps := ""
		if len(info.Dependencies) > 0 {
			deps = fmt.Sprintf(" (depends: %s)", strings.Join(info.Dependencies, ", "))
		}
		fmt.Fprintf(&b, "  %-10s  %s%s\n", info.Name, status, deps)
	}
	return b.String()
}
