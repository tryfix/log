package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tryfix/log"
	traceableContext "github.com/tryfix/traceable-context"
)

func main() {
	// usage of log
	log.Info(`message`, `param1`, `param2`)
	log.Trace(`message`)
	log.Error(`message`)
	log.Error(log.WithPrefix(`prefix`, `message`), `param1`, `param2`)

	// log with a traceable context
	tCtx := traceableContext.WithUUID(uuid.New())
	ctx, fn := context.WithCancel(tCtx)
	defer fn()
	logger := log.Constructor.Log(
		log.WithColors(true),
		log.WithLevel(log.TRACE),
		log.WithFilePath(false),
		log.WithFuncPath(true),
		log.Prefixed(`level-1`),
		log.WithCtxTraceExtractor(func(ctx context.Context) string {
			if trace := traceableContext.FromContext(ctx); trace != uuid.Nil {
				return trace.String()
			}

			return ""
		}))
	logger.ErrorContext(ctx, `message`, `param1`, `param2`)
	logger.ErrorContext(ctx, `message`)
	logger.ErrorContext(ctx, `message`)
	logger.ErrorContext(ctx, log.WithPrefix(`prefix`, `message`))
	logger.WarnContext(ctx, log.WithPrefix(`prefix`, `message`), `param1`, `param2`)

	type WrappedError struct {
		error
	}

	errWrapped := WrappedError{
		fmt.Errorf(`wrapped error: %s`, `inner error`),
	}

	logger.Error("error", errWrapped)
	logger.ErrorContext(context.Background(), "contexed-error", errWrapped)

	// sub logger with traceable context
	subLogger := logger.NewLog(log.Prefixed("sub-logger"))
	subLogger.ErrorContext(ctx, "message", "with trace")
	subLogger.ErrorContext(context.Background(), "message", "with empty trace")

	// prefixed log
	prefixedLogger := log.Constructor.PrefixedLog(
		log.WithLevel(log.ERROR),
		log.WithFuncPath(true),
		log.WithFilePath(true))
	prefixedLogger.Info(`module.sub-module`, `message`)
	prefixedLogger.Trace(`module.sub-module`, `message`)
	prefixedLogger.Error(`module.sub-module`, `message`)
	prefixedLogger.Error(`module.sub-module`, `message`, `param1`, `param2`)

	// enable context reading
	// keys
	const k1 = "key1"
	const k2 = "key2"

	// get details from context
	lCtx := context.Background()
	lCtx = context.WithValue(lCtx, k1, "context_val_1")
	lCtx = context.WithValue(lCtx, k2, "context_val_2")

	ctxLogger := log.Constructor.Log(log.WithColors(true),
		log.WithLevel(log.TRACE),
		log.WithFilePath(false),
		log.Prefixed(`context_logger`),
		log.WithCtxExtractor(func(ctx context.Context) []interface{} {
			return []interface{}{
				fmt.Sprintf("%s: %+v", k1, ctx.Value(k1)),
			}
		}),
		log.WithCtxMapExtractor(func(ctx context.Context) map[string]string {
			return map[string]string{
				k1 + `mp`: ctx.Value(k1).(string),
				k2 + `mp`: ctx.Value(k2).(string),
			}
		}),
	)

	ctxLogger.ErrorContext(lCtx, `message`)
	ctxLogger.ErrorContext(lCtx, `message`, `param1`, `param2`)
	ctxLogger.ErrorContext(lCtx, log.WithPrefix(`prefix`, `message`))
	ctxLogger.WarnContext(lCtx, log.WithPrefix(`prefix`, `message`), `param1`, `param2`)

	// child logger with additional context extraction functionality
	ctxChildLogger := ctxLogger.NewLog(log.Prefixed(`context_child_logger`),
		log.WithCtxExtractor(func(ctx context.Context) []interface{} {
			return []interface{}{
				fmt.Sprintf("%s: %+v", k2, ctx.Value(k2)),
			}
		}),
	)

	ctxChildLogger.ErrorContext(lCtx, `message`, `param1`, `param2`)
	ctxChildLogger.ErrorContext(lCtx, `message`)
	ctxChildLogger.ErrorContext(lCtx, `message`)
	ctxChildLogger.ErrorContext(lCtx, log.WithPrefix(`prefix`, `message`))
	ctxChildLogger.WarnContext(lCtx, log.WithPrefix(`prefix`, `message`), `param1`, `param2`)
	ctxChildLogger.Println(`param1`, `param2`)
	ctxChildLogger.Print(`param1`, `param2`)
}
