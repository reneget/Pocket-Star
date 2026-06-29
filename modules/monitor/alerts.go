package monitor

import (
	"fmt"
	"sync"
	"time"

	"github.com/anomalyco/my-pretty-star/pkg/log"
)

type AlertHandler func(Alert)

type AlertEngine struct {
	mu       sync.RWMutex
	rules    []AlertRule
	handlers []AlertHandler
	store    *MetricsStore
}

func NewAlertEngine(store *MetricsStore) *AlertEngine {
	return &AlertEngine{
		store: store,
		rules: defaultRules(),
	}
}

func defaultRules() []AlertRule {
	return []AlertRule{
		{Check: CheckPeerAlive, Warning: 120, Critical: 300, Enabled: true},
		{Check: CheckLatency, Warning: 100, Critical: 500, Enabled: true},
		{Check: CheckCPU, Warning: 70, Critical: 90, Enabled: true},
		{Check: CheckRAM, Warning: 80, Critical: 95, Enabled: true},
		{Check: CheckDisk, Warning: 80, Critical: 95, Enabled: true},
		{Check: CheckContainerCPU, Warning: 80, Critical: 95, Enabled: true},
	}
}

func (e *AlertEngine) SetRules(rules []AlertRule) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.rules = rules
}

func (e *AlertEngine) AddHandler(h AlertHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.handlers = append(e.handlers, h)
}

func (e *AlertEngine) Evaluate(r CheckResult) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	for _, rule := range e.rules {
		if !rule.Enabled || rule.Check != r.Type {
			continue
		}
		if rule.Label != "" && rule.Label != r.Label {
			continue
		}

		var level AlertLevel
		var triggered bool

		switch r.Type {
		case CheckPeerAlive:
			if r.Value > rule.Critical {
				level = AlertCritical
				triggered = true
			} else if r.Value > rule.Warning {
				level = AlertWarning
				triggered = true
			}
		case CheckLatency:
			if r.Value > rule.Critical {
				level = AlertCritical
				triggered = true
			} else if r.Value > rule.Warning {
				level = AlertWarning
				triggered = true
			}
		case CheckCPU, CheckRAM, CheckDisk, CheckContainerCPU:
			if r.Value > rule.Critical {
				level = AlertCritical
				triggered = true
			} else if r.Value > rule.Warning {
				level = AlertWarning
				triggered = true
			}
		case CheckContainer:
			if !r.OK {
				level = AlertCritical
				triggered = true
			}
		}

		if triggered {
			a := Alert{
				Level:     level,
				Check:     r.Type,
				Label:     r.Label,
				Value:     r.Value,
				Threshold: rule.Warning,
				Message:   fmt.Sprintf("%s: %.1f (threshold: %.1f)", r.Type, r.Value, rule.Warning),
				Time:      time.Now().Unix(),
			}
			if !r.OK && r.Message != "" {
				a.Message = r.Message
			}
			if level == AlertCritical {
				a.Message = "CRITICAL: " + a.Message
				log.Default.Push(fmt.Sprintf("🚨 %s %s: %.1f", r.Type, r.Label, r.Value))
			} else {
				a.Message = "WARNING: " + a.Message
				log.Default.Push(fmt.Sprintf("⚠️ %s %s: %.1f", r.Type, r.Label, r.Value))
			}

			for _, h := range e.handlers {
				h(a)
			}
		}
	}
}

func (e *AlertEngine) EvaluateContainer(ct ContainerInfo) {
	r := CheckResult{
		Type: CheckContainer, OK: ct.Running, Value: float64(ct.Restarts), Label: ct.Name,
		Timestamp: time.Now().Unix(), Message: ct.Status,
	}
	if !ct.Running {
		r.Message = fmt.Sprintf("container %s is %s", ct.Name, ct.Status)
	}
	e.Evaluate(r)
}
