package scheduler

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/YxTiBlya/ci-core/logger"
)

type Scheduler struct {
	log        *logger.Logger
	components []Component
}

func NewScheduler(components ...Component) *Scheduler {
	sch := &Scheduler{
		log: logger.New("scheduler"),
	}

	cmps := make([]Component, len(components))
	copy(cmps, components)
	sch.components = cmps

	return sch
}

func (sch *Scheduler) Start() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // TODO: make configurable?
	defer cancel()

	for _, component := range sch.components {
		sch.log.Info().Str("component", component.name).Msg("starting component")
		if err := component.svc.Start(ctx); err != nil {
			sch.log.Fatal().Err(err).Str("component", component.name).Msg("failed to start component")
		}
	}

	select {
	case <-ctx.Done():
		sch.log.Fatal().Err(ErrContextTimeout).Msg("failed to start components")
	default:
	}

	sch.log.Info().Msg("all components is started")

	sch.gracefulShutDown()
}

func (sch *Scheduler) gracefulShutDown() {
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quitCh)

	sig := <-quitCh
	sch.log.Info().Str("signal", sig.String()).Msg("received signal")
	sch.stop()
	sch.log.Info().Msg("graceful shutdown done")
}

func (sch *Scheduler) stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // TODO: make configurable?
	defer cancel()

	for _, component := range sch.components {
		if err := component.svc.Stop(ctx); err != nil {
			sch.log.Error().Err(err).Str("component", component.name).Msg("failed to stop component")
		}
	}

	select {
	case <-ctx.Done():
		sch.log.Fatal().Err(ErrContextTimeout).Msg("failed to stop components")
	default:
	}
}
