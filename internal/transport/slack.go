package transport

import (
	"context"
	"os"

	"emperror.dev/errors"
	"github.com/rs/zerolog"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type (
	Bot interface {
		Run(ctx context.Context) error
	}

	SlackConfig struct {
		appToken  string
		authToken string
		channelID string
	}

	SlackServer struct {
		client *socketmode.Client
		logger zerolog.Logger
	}

	debugLogger struct {
		logger zerolog.Logger
	}
)

func NewSlackConfigFromEnv() (SlackConfig, error) {
	appToken, exists := os.LookupEnv("SLACK_APP_TOKEN")
	if !exists {
		return SlackConfig{}, errors.New("SLACK_APP_TOKEN environment variable not set")
	}

	authToken, exists := os.LookupEnv("SLACK_AUTH_TOKEN")
	if !exists {
		return SlackConfig{}, errors.New("SLACK_AUTH_TOKEN environment variable not set")
	}

	channelID, exists := os.LookupEnv("SLACK_CHANNEL_ID")
	if !exists {
		return SlackConfig{}, errors.New("SLACK_CHANNEL_ID environment variable not set")
	}

	return SlackConfig{
		appToken:  appToken,
		authToken: authToken,
		channelID: channelID,
	}, nil
}

func NewSlackClient(config SlackConfig, logger zerolog.Logger) *socketmode.Client {
	debugLog := newDebugLogger(logger.With().Str("bot", "slack_socket").Logger())

	client := slack.New(
		config.authToken,
		slack.OptionLog(debugLog),
		slack.OptionDebug(true),
		slack.OptionAppLevelToken(config.appToken),
	)

	return socketmode.New(
		client,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(debugLog),
	)
}

func NewSlackBot(client *socketmode.Client, logger zerolog.Logger) Bot {
	return &SlackServer{
		client: client,
		logger: logger,
	}
}

func newDebugLogger(logger zerolog.Logger) debugLogger {
	return debugLogger{logger: logger}
}

func (s SlackServer) Run(ctx context.Context) error {
	s.logger.Info().Msg("🚀 Starting Slack Server")

	return s.client.RunContext(ctx)
}

func (l debugLogger) Output(i int, msg string) error {
	l.logger.Debug().Msg(msg)

	return nil
}
