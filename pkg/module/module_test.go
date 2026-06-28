package module

import "testing"

func TestRegisterAndGet(t *testing.T) {
	before := len(GetRegistry())

	m := &testModule{name: "test"}
	Register(m)

	after := len(GetRegistry())
	if after != before+1 {
		t.Errorf("expected %d modules, got %d", before+1, after)
	}

	found := false
	for _, reg := range GetRegistry() {
		if reg.Name() == "test" {
			found = true
			break
		}
	}
	if !found {
		t.Error("registered module not found in registry")
	}
}

type testModule struct {
	name string
	ctx  *Context
	st   Status
}

func (m *testModule) Name() string          { return m.name }
func (m *testModule) Priority() int          { return 0 }
func (m *testModule) Dependencies() []string { return nil }
func (m *testModule) Status() Status        { return m.st }
func (m *testModule) Init(ctx *Context) error {
	m.ctx = ctx
	m.st = StatusInactive
	return nil
}
func (m *testModule) Start() error {
	m.st = StatusActive
	return nil
}
func (m *testModule) Stop() error {
	m.st = StatusInactive
	return nil
}
