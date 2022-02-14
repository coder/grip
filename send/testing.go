package send

import (
	"testing"

	"cdr.dev/grip/level"
	"cdr.dev/grip/message"
)

type testLogger struct {
	t *testing.T
	*Base
}

// NewTesting constructs a fully configured Sender implementation that
// logs using the testing.T's logging facility for better integration
// with unit tests. Construct and register such a sender for
// grip.Journaler instances that you use inside of tests to have
// logging that correctly respects go test's verbosity.
//
// By default, this constructor creates a sender with a level threshold
// of "debug" and a default log level of "info."
func NewTesting(t *testing.T) Sender {
	return MakeTesting(t, t.Name(), LevelInfo{Threshold: level.Debug, Default: level.Info})
}

// MakeTesting produces a sender implementation that logs using the
// testing.T's logging facility for better integration with unit
// tests. Construct and register such a sender for grip.Journaler
// instances that you use inside of tests to have logging that
// correctly respects go test's verbosity.
func MakeTesting(t *testing.T, name string, l LevelInfo) Sender {
	s, err := setup(&testLogger{t: t, Base: NewBase(name)}, name, l)
	if err != nil {
		t.Fatalf("problem setting up logger %s: %v", name, err)
	}
	return s
}

func (s *testLogger) Send(m message.Composer) {
	if s.Level().ShouldLog(m) {
		out, err := s.formatter(m)
		if err != nil {
			s.t.Logf("formating message [type=%T]: %v", m, err)
			return
		}

		s.t.Log(out)
	}
}
