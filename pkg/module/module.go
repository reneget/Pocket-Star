package module

type Status string

const (
	StatusInactive Status = "inactive"
	StatusActive   Status = "active"
	StatusError    Status = "error"
)

type Context struct {
	DataDir string
	Network string
	HubAddr string
}

type Module interface {
	Name() string
	Priority() int
	Dependencies() []string
	Init(ctx *Context) error
	Start() error
	Stop() error
	Status() Status
}

var registry []Module

func Register(m Module) {
	registry = append(registry, m)
}

func GetRegistry() []Module {
	return registry
}
