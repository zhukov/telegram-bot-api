package tgbotapi

import (
	"context"
	"errors"
	stdlog "log"
	"log/slog"
	"os"
)

// BotLogger is an interface that represents the required methods to log data.
//
// Instead of requiring the standard logger, we can just specify the methods we
// use and allow users to pass anything that implements these.
type BotLogger interface {
	Println(v ...any)
	Printf(format string, v ...any)
}

var errNilLogger = errors.New("logger is nil")

var log BotLogger = stdlog.New(os.Stderr, "", stdlog.LstdFlags)

// SetLogger specifies the logger that the package should use.
func SetLogger(logger BotLogger) error {
	if logger == nil {
		return errNilLogger
	}
	log = logger
	return nil
}

func (bot *BotAPI) debugLoggingEnabled() bool {
	return bot.Debug && !bot.loggingDisabled
}

func (bot *BotAPI) logDebug(ctx context.Context, msg string, args ...any) {
	if !bot.debugLoggingEnabled() {
		return
	}
	bot.logMessage(ctx, slog.LevelDebug, msg, args...)
}

func (bot *BotAPI) logRequestDebug(ctx context.Context, endpoint string, debugInfo requestDebug) {
	if !bot.debugLoggingEnabled() {
		return
	}
	if bot.logger != nil {
		args := []any{"endpoint", endpoint, "params", debugInfo.params}
		if debugInfo.fileCount > 0 {
			args = append(args, "file_count", debugInfo.fileCount)
		}
		bot.logger.DebugContext(ctx, "telegram request", args...)
		return
	}
	if debugInfo.fileCount > 0 {
		log.Printf("Endpoint: %s, params: %v, with %d files\n", endpoint, debugInfo.params, debugInfo.fileCount)
		return
	}
	log.Printf("Endpoint: %s, params: %v\n", endpoint, debugInfo.params)
}

func (bot *BotAPI) logResponseDebug(ctx context.Context, endpoint string, response string) {
	if !bot.debugLoggingEnabled() {
		return
	}
	if bot.logger != nil {
		bot.logger.DebugContext(ctx, "telegram response", "endpoint", endpoint, "response", response)
		return
	}
	log.Printf("Endpoint: %s, response: %s\n", endpoint, response)
}

func (bot *BotAPI) logUpdateError(ctx context.Context, err error) {
	if bot.loggingDisabled {
		return
	}
	if bot.logger != nil {
		bot.logger.ErrorContext(ctx, "telegram get updates failed", "error", err)
		bot.logger.InfoContext(ctx, "telegram get updates retry scheduled", "delay", "3s")
		return
	}
	log.Println(err)
	log.Println("Failed to get updates, retrying in 3 seconds...")
}

func (bot *BotAPI) logMessage(ctx context.Context, level slog.Level, msg string, args ...any) {
	if bot.logger != nil {
		bot.logger.Log(ctx, level, msg, args...)
		return
	}

	switch len(args) {
	case 0:
		log.Println(msg)
	default:
		log.Println(append([]any{msg}, args...)...)
	}
}
