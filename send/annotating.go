package send

import (
	"errors"
	"strings"

	"cdr.dev/grip/message"
)

type annotatingSender struct {
	Sender
	annotations map[string]interface{}
}

// NewAnnotatingSender adds the annotations defined in the annotations
// map to every argument.
//
// Calling code should assume that the sender owns the annotations map
// and it should not attempt to modify that data after calling the
// sender constructor. Furthermore, since it owns the sender, callin Close on
// this sender will close the underlying sender.
//
// While you can wrap an existing sender with the annotator, changes
// to the annotating sender (e.g. level, formater, error handler) will
// propagate to the embedded sender.
func NewAnnotatingSender(s Sender, annotations map[string]interface{}) Sender {
	return &annotatingSender{
		Sender:      s,
		annotations: annotations,
	}
}

func (s *annotatingSender) Send(m message.Composer) {
	if !s.Sender.Level().ShouldLog(m) {
		return
	}

	errs := []string{}
	for k, v := range s.annotations {
		err := m.Annotate(k, v)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		s.ErrorHandler()(errors.New(strings.Join(errs, ";\n")), m)
	}

	s.Sender.Send(m)
}
