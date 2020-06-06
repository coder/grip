====================================
``grip`` -- A Golang Logging Library
====================================

Grip is a high level 

#. Provide a common logging interface with support for multiple
   output/messaging backends.

#. Provides some simple methods for errors, particularly when
   you want to accumulate and then return errors.

#. Provides tools for collecting structured logging information.

*You just get a grip, folks.*

Use
---

Download: ::

   go get -u github.com/deciduosity/grip

Import: ::

   import "github.com/deduosity/grip"

Design
------

Interface
~~~~~~~~~

Grip provides three main interfaces: 

- The ``send.Sender`` interfaces which implements sending messages to various
  output sources. Provides sending as well as the ability to support error
  handling, and message formating support.

- The ``message.Composer`` which wraps messages providing both "string"
  formating as well as a "raw" serialized approach. 

- The ``grip.Journaler`` interface provides a high level logging interface,
  and is mirrored in the package's public interface as a defult logger. 

Goals
~~~~~

- Provide robust high-level abstractions for applications to manage messaging,
  logging, and metrics collection.

- Integrate with other logging systems (e.g. standard library logging,
  standard output of subprocesses, other libraries, etc.)

- Minimize operational complexity and dependencies for having robust logging
  (e.g. make it possible to log effectively from within a program without
  requiring log relays or collection agents.)

Development
-----------

Grip is relatively stable, though there are additional features and areas of
development: 

- structured metrics collection. This involves adding a new interface as a
  superset of the Composer interface, and providing ways of filtering these
  messages out to provide better tools for collecting diagnostic data from
  applications.

- additional Sender implementations to support additional output formats and
  needs.

- better integration with recent development in error wrapping in the go
  standard library.
  
If you encounter a problem please feel free to create a github issue or open a
pull request.

Features
--------

Output Formats
~~~~~~~~~~~~~~

Grip supports a number of different logging output backends:

- systemd's journal (linux-only)
- syslog (unix-only)
- writing messages to standard output.
- writing messages to a file.
- sending messages to a slack's channel
- sending messages to a user via XMPP (jabber.)
- creating or commeting on a jira ticket
- creating or commenting on a github issue
- sending a desktop notification
- sending an email. 
- sending log output.
- create a tweet.  

See the documentation of the `Sender interface
<https://godoc.org/github.com/tychoish/grip/send#Sender>`_ for more
information on building new senders. The `base sender implementation
<https://godoc.org/github.com/tychoish/grip/send#Base>`_ implements most of
the interface, except for the Send method.

In addition to a collection of useful output implementations, grip also
provides tools for managing output including: 

- the `multi sender
  <https://godoc.org/github.com/deciduosity/grip/send#NewConfiguredMultiSender>`_
  for combining multiple senders to "tee" the output together,
  
- the `buffering sender
  <https://godoc.org/github.com/deciduosity/grip/send#NewBufferedSender>`_ for
  wrapping a sender with a buffer that will batch messages after reciving a
  specified number of messages, or on a specific interval.

- the `io.Writer
  <https://godoc.org/github.com/deciduosity/grip/send#WriterSender>`_ to convert a
  sender implementation to an io.Writer, to be able to use grip fundamentals
  in situations that call for ``io.Writers`` (e.g. the output of
  subprocesses,.

- the `WrapWriter
  <https://godoc.org/github.com/deciduosity/grip/send#WrapWriter>`_ to use an
  arbitrary ``io.Writer`` interface as a sender.

Logging
~~~~~~~

Provides a fully featured level-based logging system with multiple
backends (e.g. send.Sender). By default logging messages are printed
to standard output, but backends exists for many possible targets. The
interface for logging is provided by the Journaler interface.

By default ``grip.std`` defines a standard global  instances
that you can use with a set of ``grip.<Level>`` functions, or you can
create your own ``Journaler`` instance and embed it in your own
structures and packages.

Defined helpers exist for the following levels/actions:

- ``Debug``
- ``Info``
- ``Notice``
- ``Warning``
- ``Error``
- ``Critical``
- ``Alert``
- ``Emergency``
- ``EmergencyPanic``
- ``EmergencyFatal``

Helpers ending with ``Panic`` call ``panic()`` after logging the message
message, and helpers ending with ``Fatal`` call ``os.Exit(1)`` after logging
the message. These are primarily for handling errors in your main() function
and should be used sparingly, if at all, elsewhere.

Sender instances have a notion of "default" log levels and thresholds, which
provide the basis for verbosity control and sane default behavior. The default
level defines the priority/level of any message with an invalid priority
specified. The threshold level, defines the minimum priority or level that
``grip`` sends to the logging system. It's not possible to suppress the
highest log level, ``Emergency`` messages will always log.

``Journaler`` objects have additional methods (also
available as functions in the ``grip`` package to manage and configure the
instance.

Error Collector for "Continue on Error" Semantics
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

If you want to do something other than ignore or simply log errors, but don't
want to abort after an error, the `Catcher Interface
<https://godoc.org/github.com/deciduosity/grip#Catcher>`_ provides a threadsafe
way of aggregating errors. Consider: ::

   func doStuff(dirname string) (error) {
           files, err := ioutil.ReadDir(dirname)
           if err != nil {
                   // should abort here because we shouldn't continue.
                   return err
           }

           catcher := grip.NewCatcher()t
           for _, f := range files {
               err = doStuffToFile(f.Name())
               catcher.Add(err)
           }

           return catcher.Resolve()
   }

Grip provides several error catchers (which are independent of the logging
infrastructure.) They are Basic, Simple, and Extended. These variants differ
on how the collected errors are represented in the final error object. Basic
uses the ``Error()`` method of component errors, Simple users
``fmt.Sprintf("%s", err)`` and Extended users ``fmt.Sprintf("%+v",
err)``. There are also Timestamp methods that annotate all errors with a
timestamp of when the error was collected to improve debugability in longer
running asynchronous contexts: these collectors rely on ``WrapErrorTime`` to
annotate the timestamp, which may be useful in other contexts.

Conditional Logging
~~~~~~~~~~~~~~~~~~~

``grip`` incldues support for conditional logging, so that you can
only log a message in certain situations, by adding a Boolean argument
to the logging call. Use this to implement "log sometimes" messages to
minimize verbosity without complicating the calling code around the
logging, or simplify logging call sites. These methods have a ``<Level>When```
format.

This is syntactic sugar around the `message.When
<https://godoc.org/github.com/deciduosity/grip/message#When>`_ message type, but
can reduce a lot of nesting and call-site complexity.

Composed Logging
~~~~~~~~~~~~~~~~

If the production of the log message is resource intensive or
complicated, you may wish to use a "composed logging," which delays
the generation of the log message from the logging call site to the
message propagation, to avoid generating the log message unless
necessary. Rather than passing the log message as a string, pass the
logging function an instance of a type that implements the
``Composer`` interface.

Grip uses composers internally, but you can pass composers directly to
any of the basic logging method (e.g. ``Info()``, ``Debug()``) for
composed logging.

Grip includes a number of message types, including those that collect
system information, process information, stacktraces, or simple
user-specified structured information.
