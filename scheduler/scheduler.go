package scheduler

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type Scheduler struct {
	log        *zap.SugaredLogger // TODO: i do interface?
	components []Component
}

func NewScheduler(logger *zap.SugaredLogger, components ...Component) *Scheduler {
	sch := &Scheduler{}
	sch.log = logger

	cmps := make([]Component, len(components))
	copy(cmps, components)
	sch.components = cmps

	return sch
}

func (sch *Scheduler) Start() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // TODO: make configurable?
	defer cancel()

	for _, component := range sch.components {
		sch.log.Infof("starting component %s", component.name)
		if err := component.svc.Start(ctx); err != nil {
			sch.log.Fatal("failed to start component ", component.name, zap.Error(err))
		}
	}

	select {
	case <-ctx.Done():
		sch.log.Fatal("failed to start components", zap.Error(ErrContextTimeout))
	default:
	}

	sch.log.Info("all components is started")

	sch.gracefulShutDown()
}

func (sch *Scheduler) gracefulShutDown() {
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quitCh)

	sig := <-quitCh
	sch.log.Info("received signal", zap.Stringer("signal", sig))
	sch.stop()
	sch.log.Info("graceful shutdown done")
}

func (sch *Scheduler) stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // TODO: make configurable?
	defer cancel()

	for _, component := range sch.components {
		if err := component.svc.Stop(ctx); err != nil {
			sch.log.Error("failed to stop component", zap.String("component", component.name), zap.Error(err))
		}
	}

	select {
	case <-ctx.Done():
		sch.log.Fatal("failed to stop components", zap.Error(ErrContextTimeout))
	default:
	}
}
