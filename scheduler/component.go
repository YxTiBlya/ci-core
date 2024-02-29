package scheduler

import "context"

type Component struct {
	name string
	svc  Service
}

type Service interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

func NewComponent(name string, svc Service) Component {
	return Component{name: name, svc: svc}
}
