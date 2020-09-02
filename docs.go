/*
Package grip provides a flexible logging package for basic Go programs.
Drawing inspiration from Go and Python's standard library
logging, as well as systemd's journal service, and other logging
systems, Grip provides a number of very powerful logging
abstractions in one high-level package.

Logging Instances

The central type of the grip package is the Journaler type,
instances of which provide distinct log capturing system. For ease,
following from the Go standard library, the grip package provides
parallel public methods that use an internal "standard" Jouernaler
instance in the grip package, which has some defaults configured
and may be sufficient for many use cases.

Output

The send.Sender interface provides a way of changing the logging
backend, and the send package provides a number of alternate
implementations of logging systems, including: systemd's journal,
logging to standard output, logging to a file, and generic syslog
support.

Messages

The message.Composer interface is the representation of all
messages. They are implemented to provide a raw structured form as
well as a string representation for more conentional logging
output. Furthermore they are intended to be easy to produce, and defer
more expensive processing until they're being logged, to prevent
expensive operations producing messages that are below threshold.

Basic Logging

Loging helpers exist for the following levels:

   Emergency + (fatal/panic)
   Alert
   Critical
   Error
   Warning
   Notice
   Info
   Debug

These methods accept both strings (message content,) or types that
implement the message.MessageComposer interface. Composer types make
it possible to delay generating a message unless the logger is over
the logging threshold. Use this to avoid expensive serialization
operations for suppressed logging operations.

All levels also have additional methods with `ln` and `f` appended to
the end of the method name which allow Println() and Printf() style
functionality. You must pass printf/println-style arguments to these methods.

Conditional Logging

The Conditional logging methods take two arguments, a Boolean, and a
message argument. Messages can be strings, objects that implement the
MessageComposer interface, or errors. If condition boolean is true,
the threshold level is met, and the message to log is not an empty
string, then it logs the resolved message.

Use conditional logging methods to potentially suppress log messages
based on situations orthogonal to log level, with "log sometimes" or
"log rarely" semantics. Combine with MessageComposers to to avoid
expensive message building operations.
*/
package grip

// This file is intentionally documentation only.
