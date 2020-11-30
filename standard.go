package grip

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cdr/grip/level"
	"github.com/cdr/grip/logging"
	"github.com/cdr/grip/send"
)

var std = NewJournaler("grip")

func init() {
	if !strings.Contains(os.Args[0], "go-build") {
		std.SetName(filepath.Base(os.Args[0]))
	}

	sender, err := send.NewNativeLogger(std.Name(), std.GetSender().Level())
	std.Alert(std.SetSender(sender))
	std.Alert(err)
}

// SetDefaultStandardLogger set's the standard library's global
// logging instance to use grip's global logger at the specified
// level.
func SetDefaultStandardLogger(p level.Priority) {
	log.SetFlags(0)
	log.SetOutput(send.MakeWriterSender(std.GetSender(), p))
}

// MakeStandardLogger constructs a standard library logging instance
// that logs all messages to the global grip logging instance.
func MakeStandardLogger(p level.Priority) *log.Logger {
	return send.MakeStandardLogger(std.GetSender(), p)
}

// NewJournaler creates a new Journaler instance. The Sender method is a
// non-operational bootstrap method that stores default and threshold
// types, as needed. You must use SetSender() or the
// UseSystemdLogger(), UseNativeLogger(), or UseFileLogger() methods
// to configure the backend.
func NewJournaler(name string) Journaler {
	return logging.NewGrip(name)
}

// GetSender returns the current Journaler's sender instance. Use this in
// combination with SetSender to have multiple Journaler instances
// backed by the same send.Sender instance.
func GetSender() send.Sender {
	return std.GetSender()
}

// GetDefaultJournaler returns the default journal instance used by
// this library. This call is not thread safe relative to other
// logging calls, or SetDefaultJournaler call, although all journaling
// methods are safe.
func GetDefaultJournaler() Journaler {
	return std
}

// SetDefaultJournaler allows you to override the standard logger,
// that is used by calls in the grip package. This call is not thread
// safe relative to other logging calls, or the GetDefaultJournaler
// call, although all journaling methods are safe: as a result be sure
// to only call this method during package and process initialization.
func SetDefaultJournaler(l Journaler) {
	std = l
}

// Name of the logger instance
func Name() string {
	return std.Name()
}

// SetName declare a name string for the logger, including in the logging
// message. Typically this is included on the output of the command.
func SetName(name string) {
	std.SetName(name)
}

// SetLevel sets the default and threshold level in the underlying sender.
func SetLevel(info send.LevelInfo) error {
	return std.SetLevel(info)
}

// SetSender swaps send.Sender() implementations in a logging
// instance. Calls the Close() method on the existing instance before
// changing the implementation for the current instance.
func SetSender(s send.Sender) error {
	return std.SetSender(s)
}

func Log(l level.Priority, msg interface{}) {
	std.Log(l, msg)
}
func Logf(l level.Priority, msg string, a ...interface{}) {
	std.Logf(l, msg, a...)
}
func Logln(l level.Priority, a ...interface{}) {
	std.Logln(l, a...)
}
func LogWhen(conditional bool, l level.Priority, m interface{}) {
	std.LogWhen(conditional, l, m)
}

// Leveled Logging Methods
// Emergency-level logging methods

func EmergencyFatal(msg interface{}) {
	std.EmergencyFatal(msg)
}
func Emergency(msg interface{}) {
	std.Emergency(msg)
}
func Emergencyf(msg string, a ...interface{}) {
	std.Emergencyf(msg, a...)
}
func Emergencyln(a ...interface{}) {
	std.Emergencyln(a...)
}
func EmergencyPanic(msg interface{}) {
	std.EmergencyPanic(msg)
}
func EmergencyWhen(conditional bool, m interface{}) {
	std.EmergencyWhen(conditional, m)
}

// Alert-level logging methods

func Alert(msg interface{}) {
	std.Alert(msg)
}
func Alertf(msg string, a ...interface{}) {
	std.Alertf(msg, a...)
}
func Alertln(a ...interface{}) {
	std.Alertln(a...)
}
func AlertWhen(conditional bool, m interface{}) {
	std.AlertWhen(conditional, m)
}

// Critical-level logging methods

func Critical(msg interface{}) {
	std.Critical(msg)
}
func Criticalf(msg string, a ...interface{}) {
	std.Criticalf(msg, a...)
}
func Criticalln(a ...interface{}) {
	std.Criticalln(a...)
}
func CriticalWhen(conditional bool, m interface{}) {
	std.CriticalWhen(conditional, m)
}

// Error-level logging methods

func Error(msg interface{}) {
	std.Error(msg)
}
func Errorf(msg string, a ...interface{}) {
	std.Errorf(msg, a...)
}
func Errorln(a ...interface{}) {
	std.Errorln(a...)
}
func ErrorWhen(conditional bool, m interface{}) {
	std.ErrorWhen(conditional, m)
}

// Warning-level logging methods

func Warning(msg interface{}) {
	std.Warning(msg)
}
func Warningf(msg string, a ...interface{}) {
	std.Warningf(msg, a...)
}
func Warningln(a ...interface{}) {
	std.Warningln(a...)
}
func WarningWhen(conditional bool, m interface{}) {
	std.WarningWhen(conditional, m)
}

// Notice-level logging methods

func Notice(msg interface{}) {
	std.Notice(msg)
}
func Noticef(msg string, a ...interface{}) {
	std.Noticef(msg, a...)
}
func Noticeln(a ...interface{}) {
	std.Noticeln(a...)
}
func NoticeWhen(conditional bool, m interface{}) {
	std.NoticeWhen(conditional, m)
}

// Info-level logging methods

func Info(msg interface{}) {
	std.Info(msg)
}
func Infof(msg string, a ...interface{}) {
	std.Infof(msg, a...)
}
func Infoln(a ...interface{}) {
	std.Infoln(a...)
}
func InfoWhen(conditional bool, message interface{}) {
	std.InfoWhen(conditional, message)
}

// Debug-level logging methods

func Debug(msg interface{}) {
	std.Debug(msg)
}
func Debugf(msg string, a ...interface{}) {
	std.Debugf(msg, a...)
}
func Debugln(a ...interface{}) {
	std.Debugln(a...)
}
func DebugWhen(conditional bool, m interface{}) {
	std.DebugWhen(conditional, m)
}
