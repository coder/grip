// Functional Messages
//
// Grip can automatically convert three types of functions into
// messages:
//
//    func() Fields
//    func() Composer
//    func() error
//
// The benefit of these functions is that they're only called if
// the message is above the logging threshold. In the case of
// conditional logging (i.e. When), if the conditional is false, then
// the function is never called.
//
// in the case of all the buffered sending implementation, the
// function call can be deferred and run outside of the main thread,
// and so may be an easy way to defer message production outside in
// cases where messages may be complicated.
//
// Additionally, the message conversion in grip's logging method can
// take these function types and convert them to these messages, which
// can clean up some call-site operations, and makes it possible to
// use defer with io.Closer methods without wrapping the method in an
// additional function, as in:
//
//     defer grip.Error(file.Close)
//
// Although the WrapErrorFunc method, as in the following may permit
// useful annotation, as follows, which has the same "lazy" semantics.
//
//     defer grip.Error(message.WrapErrorFunc(file.Close, message.Fields{}))
//
package message

import (
	"fmt"
	"io"

	"github.com/cdr/grip/level"
	"github.com/pkg/errors"
)

// FieldsProducer is a function that returns a structured message body
// as a way of writing simple Composer implementations in the form
// anonymous functions, as in:
//
//    grip.Info(func() message.Fields {return message.Fields{"message": "hello world!"}})
//
// Grip can automatically convert these functions when passed to a
// logging function.
//
// If the Fields object is nil or empty then no message is logged.
type FieldsProducer func() Fields

type fieldsProducerMessage struct {
	fp     FieldsProducer
	cached Composer
	level  level.Priority
}

// NewFieldsProducerMessage constructs a lazy FieldsProducer wrapping
// message at the specified level.
//
// FieldsProducer functions are only called, before calling the
// Loggable, String, Raw, or Annotate methods. Changing the priority
// does not call the function. In practice, if the priority of the
// message is below the logging threshold, then the function will
// never be called.
func NewFieldsProducerMessage(p level.Priority, fp FieldsProducer) Composer {
	return &fieldsProducerMessage{level: p, fp: fp}
}

// MakeFieldsProducerMessage constructs a lazy FieldsProducer wrapping
// message at the specified level.
//
// FieldsProducer functions are only called, before calling the
// Loggable, String, Raw, or Annotate methods. Changing the priority
// does not call the function. In practice, if the priority of the
// message is below the logging threshold, then the function will
// never be called.
func MakeFieldsProducerMessage(fp FieldsProducer) Composer {
	return &fieldsProducerMessage{fp: fp}
}

// MakeConvertedFieldsProducer converts a generic map to a fields
// producer, as the message types are equivalent.
func MakeConvertedFieldsProducer(mp func() map[string]interface{}) Composer {
	return MakeFieldsProducerMessage(func() Fields {
		return mp()
	})
}

// NewConvertedFieldsProducer converts a generic map to a fields
// producer at the specified priority, as the message types are equivalent,
func NewConvertedFieldsProducer(p level.Priority, mp func() map[string]interface{}) Composer {
	return NewFieldsProducerMessage(p, func() Fields {
		return mp()
	})
}

func (fp *fieldsProducerMessage) resolve() {
	if fp.cached == nil {
		if fp.fp == nil {
			fp.cached = NewFields(fp.level, Fields{})
		} else {
			fp.cached = NewFields(fp.level, fp.fp())
		}
	}
}

func (fp *fieldsProducerMessage) Annotate(k string, v interface{}) error {
	fp.resolve()
	return fp.cached.Annotate(k, v)
}

func (fp *fieldsProducerMessage) SetPriority(p level.Priority) error {
	if !p.IsValid() {
		return errors.New("invalid level")
	}
	fp.level = p
	if fp.cached != nil {
		return fp.cached.SetPriority(fp.level)
	}

	return nil
}

func (fp *fieldsProducerMessage) Loggable() bool {
	if fp.fp == nil {
		return false
	}

	fp.resolve()
	return fp.cached.Loggable()
}

func (fp *fieldsProducerMessage) Priority() level.Priority { return fp.level }
func (fp *fieldsProducerMessage) String() string           { fp.resolve(); return fp.cached.String() }
func (fp *fieldsProducerMessage) Raw() interface{}         { fp.resolve(); return fp.cached.Raw() }

////////////////////////////////////////////////////////////////////////

// ComposerProducer constructs a lazy composer, and makes it easy to
// implement new Composers as functions returning an existing composer
// type. Consider the following:
//
//    grip.Info(func() message.Composer { return WrapError(validateRequest(req), message.Fields{"op": "name"})})
//
// Grip can automatically convert these functions when passed to a
// logging function.
//
// If the Fields object is nil or empty then no message is ever logged.
type ComposerProducer func() Composer

type composerProducerMessage struct {
	cp     ComposerProducer
	cached Composer
	level  level.Priority
}

// NewComposerMessage constructs a message, with the given priority,
// that will call the ComposerProducer function lazily during logging.
//
// ComposerProducer functions are only called, before calling the
// Loggable, String, Raw, or Annotate methods. Changing the priority
// does not call the function. In practice, if the priority of the
// message is below the logging threshold, then the function will
// never be called.
func NewComposerProducerMessage(p level.Priority, cp ComposerProducer) Composer {
	return &composerProducerMessage{level: p, cp: cp}
}

// MakeComposerMessage constructs a message that will call the
// ComposerProducer function lazily during logging.
//
// ComposerProducer functions are only called, before calling the
// Loggable, String, Raw, or Annotate methods. Changing the priority
// does not call the function. In practice, if the priority of the
// message is below the logging threshold, then the function will
// never be called.
func MakeComposerProducerMessage(cp ComposerProducer) Composer {
	return &composerProducerMessage{cp: cp}
}

func (cp *composerProducerMessage) resolve() {
	if cp.cached == nil {
		cp.cached = cp.cp()
		if cp.cached == nil {
			cp.cached = NewSimpleFields(cp.level, Fields{})
		} else {
			_ = cp.cached.SetPriority(cp.level)
		}
	}
}

func (cp *composerProducerMessage) Annotate(k string, v interface{}) error {
	cp.resolve()
	return cp.cached.Annotate(k, v)
}

func (cp *composerProducerMessage) SetPriority(p level.Priority) error {
	if !p.IsValid() {
		return errors.New("invalid level")
	}

	cp.level = p
	if cp.cached != nil {
		return cp.cached.SetPriority(cp.level)
	}

	return nil
}

func (cp *composerProducerMessage) Loggable() bool {
	if cp.cp == nil {
		return false
	}

	cp.resolve()
	return cp.cached.Loggable()
}

func (cp *composerProducerMessage) Priority() level.Priority { return cp.level }
func (cp *composerProducerMessage) String() string           { cp.resolve(); return cp.cached.String() }
func (cp *composerProducerMessage) Raw() interface{}         { cp.resolve(); return cp.cached.Raw() }

////////////////////////////////////////////////////////////////////////

// ErrorProducer is a function that returns an error, and is used for
// constructing message that lazily wraps the resulting function which
// is called when the message is dispatched.
//
// If you pass one of these functions to a logging method, the
// ConvertToComposer operation will construct a lazy Composer based on
// this function, as in:
//
//    grip.Error(func() error { return errors.New("error message") })
//
// It may be useful also to pass a "closer" function in this form, as
// in:
//
//    grip.Error(file.Close)
//
// As a special case the WrapErrorFunc method has the same semantics
// as other ErrorProducer methods, but makes it possible to annotate
// an error.
type ErrorProducer func() error

type errorProducerMessage struct {
	ep     ErrorProducer
	cached Composer
	level  level.Priority
}

// NewErrorProducerMessage returns a mesage that wrapps an error
// producing function, at the specified level. If the function returns
// then there is never a message logged.
//
// ErrorProducer functions are only called, before calling the
// Loggable, String, Raw, or Annotate methods. Changing the priority
// does not call the function. In practice, if the priority of the
// message is below the logging threshold, then the function will
// never be called.
func NewErrorProducerMessage(p level.Priority, ep ErrorProducer) Composer {
	return &errorProducerMessage{level: p, ep: ep}
}

// MakeErrorProducerMessage returns a mesage that wrapps an error
// producing function. If the function returns then there is never a
// message logged.
//
// ErrorProducer functions are only called, before calling the
// Loggable, String, Raw, or Annotate methods. Changing the priority
// does not call the function. In practice, if the priority of the
// message is below the logging threshold, then the function will
// never be called.
func MakeErrorProducerMessage(ep ErrorProducer) Composer {
	return &errorProducerMessage{ep: ep}
}

func (ep *errorProducerMessage) resolve() {
	if ep.cached == nil {
		ep.cached = NewErrorMessage(ep.level, ep.ep())
	}
}

func (ep *errorProducerMessage) Annotate(k string, v interface{}) error {
	ep.resolve()
	return ep.cached.Annotate(k, v)
}

func (ep *errorProducerMessage) SetPriority(p level.Priority) error {
	if !p.IsValid() {
		return errors.New("invalid level")
	}

	ep.level = p
	if ep.cached != nil {
		return ep.cached.SetPriority(ep.level)
	}

	return nil
}

func (ep *errorProducerMessage) Loggable() bool {
	if ep.ep == nil {
		return false
	}

	ep.resolve()
	return ep.cached.Loggable()
}

func (ep *errorProducerMessage) Priority() level.Priority { return ep.level }
func (ep *errorProducerMessage) String() string           { ep.resolve(); return ep.cached.String() }
func (ep *errorProducerMessage) Raw() interface{}         { ep.resolve(); return ep.cached.Raw() }
func (ep *errorProducerMessage) Error() string            { return ep.String() }
func (ep *errorProducerMessage) Unwrap() error            { return ep.Cause() }

func (ep *errorProducerMessage) Cause() error {
	ep.resolve()

	switch err := ep.cached.(type) {
	case *errorComposerWrap:
		return err.err
	case *errorMessage:
		return err.err
	case error:
		return err
	default:
		return nil
	}
}

func (ep *errorProducerMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", errors.Cause(ep))
			_, _ = io.WriteString(s, ep.String())
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, ep.Error())
	}
}

// WrapErrorFunc produces a lazily-composed wrapped error message. The
// function is only called is
//
// The resulting method itself implements the "error" interface
// (supporing unwrapping,) as well as the composer type, so you can
// return the result of this function as an error to avoid needing to
// manage multiple error annotations.
func WrapErrorFunc(ep ErrorProducer, m interface{}) Composer {
	return errorComposerShim{MakeComposerProducerMessage(func() Composer { return WrapError(ep(), m) })}
}

type errorComposerShim struct {
	Composer
}

func (ecs errorComposerShim) Error() string { return ecs.Composer.String() }
func (ecs errorComposerShim) Unwrap() error { return ecs.Cause() }

func (ecs errorComposerShim) Cause() error {
	switch err := ecs.Composer.(type) {
	case *errorComposerWrap:
		return err.err
	case *errorMessage:
		return err.err
	case error:
		return err
	default:
		return nil
	}
}

func (ecs errorComposerShim) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", errors.Cause(ecs))
			_, _ = io.WriteString(s, ecs.String())
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, ecs.Error())
	}
}
