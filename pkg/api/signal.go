// +gocover:ignore:file signal handling and process termination

package api

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/br-lemes/golem/pkg/console"
)

var actionInterrupt = newInterruptGuard()

type interruptGuard struct {
	mu        sync.Mutex
	inAction  bool
	pending   bool
	interrupt chan os.Signal
}

func newInterruptGuard() *interruptGuard {
	guard := &interruptGuard{interrupt: make(chan os.Signal, 1)}
	signal.Notify(guard.interrupt, syscall.SIGINT)
	go guard.handleSignals()
	return guard
}

func (guard *interruptGuard) handleSignals() {
	for range guard.interrupt {
		guard.mu.Lock()
		if guard.inAction {
			guard.pending = true
			guard.mu.Unlock()
			console.Errorf("signal: interrupt ")
			continue
		}
		guard.mu.Unlock()
		signal.Stop(guard.interrupt)
		os.Exit(130)
	}
}

func beginCriticalAction() func() {
	actionInterrupt.mu.Lock()
	actionInterrupt.inAction = true
	actionInterrupt.mu.Unlock()

	var once sync.Once
	return func() {
		once.Do(func() {
			actionInterrupt.mu.Lock()
			actionInterrupt.inAction = false
			pending := actionInterrupt.pending
			actionInterrupt.pending = false
			actionInterrupt.mu.Unlock()
			if pending {
				signal.Stop(actionInterrupt.interrupt)
				os.Exit(130)
			}
		})
	}
}
