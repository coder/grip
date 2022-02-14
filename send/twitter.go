package send

import (
	"context"
	"fmt"
	"log"
	"os"

	"cdr.dev/grip/level"
	"cdr.dev/grip/message"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/pkg/errors"
)

type twitterLogger struct {
	twitter twitterClient
	*Base
}

// TwitterOptions describes the credentials required to connect to the
// twitter API. While the name is used for internal reporting, the
// other values should be populated with credentials obtained from the
// Twitter API.
type TwitterOptions struct {
	Name           string
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

func (opts *TwitterOptions) resolve(ctx context.Context) *twitter.Client {
	return twitter.NewClient(oauth1.NewConfig(opts.ConsumerKey, opts.ConsumerSecret).
		Client(ctx, oauth1.NewToken(opts.AccessToken, opts.AccessSecret)))
}

// MakeTwitterLogger constructs a default sender implementation that
// posts messages to a twitter account. The implementation does not
// rate limit outgoing messages, which should be the responsibility of
// the caller.
func MakeTwitterLogger(ctx context.Context, opts *TwitterOptions) (Sender, error) {
	return NewTwitterLogger(ctx, opts, LevelInfo{level.Trace, level.Trace})
}

// NewTwitterLogger constructs a sender implementation that posts
// messages to a twitter account, with configurable level
// information. The implementation does not rate limit outgoing
// messages, which should be the responsibility of the caller.
func NewTwitterLogger(ctx context.Context, opts *TwitterOptions, l LevelInfo) (Sender, error) {
	s := &twitterLogger{
		twitter: newTwitterClient(ctx, opts),
		Base:    NewBase(opts.Name),
	}

	if err := s.SetLevel(l); err != nil {
		return nil, errors.Wrap(err, "invalid level specification")
	}

	fallback := log.New(os.Stdout, "", log.LstdFlags)
	if err := s.SetErrorHandler(ErrorHandlerFromLogger(fallback)); err != nil {
		return nil, err
	}

	s.reset = func() {
		fallback.SetPrefix(fmt.Sprintf("[%s] ", s.Name()))
	}

	s.SetName(opts.Name)

	if err := s.twitter.Verify(); err != nil {
		return nil, errors.Wrap(err, "problem connecting to twitter")
	}

	return s, nil
}

func (s *twitterLogger) Send(m message.Composer) {
	if s.Level().ShouldLog(m) {
		if err := s.twitter.Send(m.String()); err != nil {
			s.ErrorHandler()(err, m)
		}
	}
}

type twitterClient interface {
	Verify() error
	Send(string) error
}

type twitterClientImpl struct {
	twitter *twitter.Client
}

func newTwitterClient(ctx context.Context, opts *TwitterOptions) twitterClient {
	return &twitterClientImpl{twitter: opts.resolve(ctx)}
}

func (tc *twitterClientImpl) Verify() error {
	_, _, err := tc.twitter.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{})
	return errors.Wrap(err, "could not verify account")
}

func (tc *twitterClientImpl) Send(in string) error {
	_, _, err := tc.twitter.Statuses.Update(in, nil)
	return errors.WithStack(err)
}
