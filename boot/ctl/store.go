package ctl

type I interface {
	Create() error
	Destroy() error
	Start() error
	IsAsync() bool
}

type serviced struct {
	id        uint
	name      string
	implement I
}

type Created struct {
	Name string
	Func func() I
}

type Destroyed struct {
	Name string
	Func func() error
}

var (
	serviceMap      = make(map[string]struct{})
	createServices  = make([]Created, 0)
	destroyServices = make([]Destroyed, 0)
)
