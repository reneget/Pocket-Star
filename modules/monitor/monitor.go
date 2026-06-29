package monitor

import (
	"fmt"
	"time"

	"github.com/anomalyco/my-pretty-star/pkg/log"
	"github.com/anomalyco/my-pretty-star/pkg/module"
)

var DefaultStore *MetricsStore
var DefaultChecker *Checker
var DefaultAlerts *AlertEngine

type MonitorModule struct {
	ctx     *module.Context
	status  module.Status
	store   *MetricsStore
	checker *Checker
	alerts  *AlertEngine
	stopCh  chan struct{}
}

func (m *MonitorModule) Name() string          { return "monitor" }
func (m *MonitorModule) Priority() int          { return 5 }
func (m *MonitorModule) Dependencies() []string { return []string{"core"} }
func (m *MonitorModule) Status() module.Status  { return m.status }

func (m *MonitorModule) Init(ctx *module.Context) error {
	m.ctx = ctx
	m.stopCh = make(chan struct{})

	store, err := NewMetricsStore(ctx.DataDir)
	if err != nil {
		return fmt.Errorf("monitor init: %w", err)
	}
	m.store = store
	DefaultStore = store

	m.alerts = NewAlertEngine(store)
	DefaultAlerts = m.alerts

	m.checker = NewChecker(store, m.alerts)
	DefaultChecker = m.checker

	m.status = module.StatusInactive
	log.Default.Push("monitor: initialized")
	return nil
}

func (m *MonitorModule) Start() error {
	if m.status == module.StatusActive {
		return nil
	}

	m.checker.RunAll()

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				m.checker.RunAll()
			case <-m.stopCh:
				return
			}
		}
	}()

	m.status = module.StatusActive
	log.Default.Push("monitor: started")
	return nil
}

func (m *MonitorModule) Stop() error {
	close(m.stopCh)
	if m.store != nil {
		m.store.Close()
	}
	m.status = module.StatusInactive
	log.Default.Push("monitor: stopped")
	return nil
}

func (m *MonitorModule) Store() *MetricsStore   { return m.store }
func (m *MonitorModule) Checker() *Checker       { return m.checker }
func (m *MonitorModule) Alerts() *AlertEngine    { return m.alerts }

func init() {
	module.Register(&MonitorModule{})
}
