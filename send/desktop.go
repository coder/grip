package send

import (
	"github.com/cdr/grip/level"
	"github.com/cdr/grip/message"
	"github.com/gen2brain/beeep"
	"github.com/pkg/errors"
)

type desktopNotify struct {
	*Base
}

// NewDesktopNotify constructs a sender that pushes messages
// to local system notification process.
func NewDesktopNotify(name string, l LevelInfo) (Sender, error) {
	s := &desktopNotify{
		Base: NewBase(name),
	}

	if err := s.SetLevel(l); err != nil {
		return nil, errors.Wrap(err, "problem seeting level on new sender")
	}

	return s, nil
}

// MakeDesktopNotify constructs a default sender that pushes messages
// to local system notification.
func MakeDesktopNotify(name string) (Sender, error) {
	s, err := NewDesktopNotify(name, LevelInfo{Threshold: level.Trace, Default: level.Debug})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return s, nil
}

func (s *desktopNotify) Send(m message.Composer) {
	if s.Level().ShouldLog(m) {
		if m.Priority() >= level.Critical {
			if err := beeep.Alert(s.Name(), m.String(), ""); err != nil {
				s.ErrorHandler()(err, m)
			}
		} else {
			if err := beeep.Notify(s.Name(), m.String(), ""); err != nil {
				s.ErrorHandler()(err, m)
			}
		}
	}
}
