package logging

import (
	"errors"
	"os"
	"sync"

	"cdr.dev/grip/level"
	"cdr.dev/grip/message"
	"cdr.dev/grip/send"
)

// Grip provides the core implementation of the Logging interface. The
// interface is mirrored in the "grip" package's public interface, to
// provide a single, global logging interface that requires minimal
// configuration.
type Grip struct {
	impl         send.Sender
	defaultLevel level.Priority
	mu           sync.RWMutex
}

// MakeGrip builds a new logging interface from a sender implmementation
func MakeGrip(s send.Sender) *Grip {
	return &Grip{
		impl:         s,
		defaultLevel: level.Info,
	}
}

// NewGrip takes the name for a logging instance and creates a new
// Grip instance with configured with a local, standard output logging.
// The default level is "Notice" and the threshold level is "info."
func NewGrip(name string) *Grip {
	sender, _ := send.NewNativeLogger(name,
		send.LevelInfo{
			Threshold: level.Trace,
			Default:   level.Trace,
		})

	return &Grip{impl: sender}
}

func (g *Grip) Name() string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.impl.Name()
}

func (g *Grip) SetName(n string) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	g.impl.SetName(n)
}

func (g *Grip) SetLevel(info send.LevelInfo) error {
	g.mu.RLock()
	defer g.mu.RUnlock()
	sl := g.impl.Level()

	if !info.Default.IsValid() {
		info.Default = sl.Default
	}

	if !info.Threshold.IsValid() {
		info.Threshold = sl.Threshold
	}

	return g.impl.SetLevel(info)
}

func (g *Grip) Send(m interface{}) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	g.impl.Send(message.ConvertToComposer(g.defaultLevel, m))
}

// SetSender swaps send.Sender() implementations in a logging
// instance. Calls the Close() method on the existing instance before
// changing the implementation for the current instance. SetSender
// will configure the incoming sender to have the same name as well as
// default and threshold level as the outgoing sender.
func (g *Grip) SetSender(s send.Sender) error {
	if s == nil {
		return errors.New("cannot set the sender to nil")
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	if err := s.SetLevel(g.impl.Level()); err != nil {
		return err
	}

	if err := g.impl.Close(); err != nil {
		return err
	}

	s.SetName(g.impl.Name())
	g.impl = s

	return nil
}

// GetSender returns the current Journaler's sender instance. Use this in
// combination with SetSender() to have multiple Journaler instances
// backed by the same send.Sender instance.
func (g *Grip) GetSender() send.Sender {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.impl
}

// Internal

// For sending logging messages, in most cases, use the
// Journaler.sender.Send() method, but we have a couple of methods to
// use for the Panic/Fatal helpers.
func (g *Grip) sendPanic(m message.Composer) {
	// the Send method in the Sender interface will perform this
	// check but to add fatal methods we need to do this here.
	g.mu.RLock()
	defer g.mu.RUnlock()

	if g.impl.Level().ShouldLog(m) {
		g.impl.Send(m)
		panic(m.String())
	}
}

func (g *Grip) sendFatal(m message.Composer) {
	// the Send method in the Sender interface will perform this
	// check but to add fatal methods we need to do this here.
	g.mu.RLock()
	defer g.mu.RUnlock()

	if g.impl.Level().ShouldLog(m) {
		g.impl.Send(m)
		os.Exit(1)
	}
}

func (g *Grip) send(m message.Composer) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	g.impl.Send(m)
}
