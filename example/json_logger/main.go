package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tryfix/log"
	traceable_context "github.com/tryfix/traceable-context"
)

func main() {
	constructor := log.NewLog(log.WithOutput(log.OutJson))
	lg := constructor.Log()

	// usage of log
	lg.Print("message", "param1", "param2")
	lg.Println("message", "param1", "param2")
	lg.Printf("%s %s", "param1", "param2")
	lg.Info("message", "param1", "param2")
	lg.Trace("message")
	lg.Error("message")
	lg.Error("message", "param1", "param2")

	// log complex structures
	type ABC struct {
		CDE string
		def string
	}
	cpx := ABC{CDE: "some string", def: "some unexported string"}

	lg.Trace("log complex", "param1", cpx)
	lg.TraceContext(context.Background(), "log complex", "param1", cpx)

	// log error values
	err := fmt.Errorf("error: %s", "plain error")
	lg.Error("error", err)

	// log embedded errors
	type WrappedError struct {
		error
	}

	errWrapped := WrappedError{
		fmt.Errorf("inner error"),
	}

	lg.Error("wrapped error", errWrapped)
	lg.ErrorContext(context.Background(), "contexed-error", errWrapped)

	// log with a traceable context
	tCtx := traceable_context.WithUUID(uuid.New())
	ctx, fn := context.WithCancel(tCtx)
	defer fn()
	logger := constructor.Log(
		log.WithLevel(log.TRACE),
		log.WithFilePath(false),
		log.Prefixed("level-1"),
		log.WithCtxTraceExtractor(func(ctx context.Context) string {
			if trace := traceable_context.FromContext(ctx); trace != uuid.Nil {
				return trace.String()
			}

			return ""
		}))
	logger.ErrorContext(ctx, "message", "param1", "param2")
	logger.ErrorContext(ctx, "message")
	logger.WarnContext(ctx, "message", "param1", "param2")
	logger.ErrorContext(context.Background(), "message", "with empty trace")

	// sub logger with traceable context
	subLogger := logger.NewLog(log.Prefixed("sub-logger"))
	subLogger.ErrorContext(ctx, "message", "with trace")
	subLogger.ErrorContext(context.Background(), "message", "with empty trace")

	// prefixed log
	prefixedLogger := constructor.PrefixedLog(log.WithLevel(log.ERROR), log.WithFilePath(true))
	prefixedLogger.Info("module.sub-module", "message")
	prefixedLogger.Trace("module.sub-module", "message")
	prefixedLogger.Error("module.sub-module", "message")
	prefixedLogger.Error("module.sub-module", "message", "param1", "param2")

	// enable context reading
	// keys
	const k1 string = "key1"
	const k2 string = "key2"

	// get details from context
	lCtx := context.Background()
	lCtx = context.WithValue(lCtx, k1, "context_val_1")
	lCtx = context.WithValue(lCtx, k2, "context_val_2")

	ctxLogger := constructor.Log(
		log.WithLevel(log.TRACE),
		log.WithFilePath(false),
		log.Prefixed("context_logger"),
		log.WithCtxExtractor(func(ctx context.Context) []interface{} {
			return []interface{}{
				fmt.Sprintf("%s: %+v", k1, ctx.Value(k1)),
			}
		}),
	)

	ctxLogger.ErrorContext(lCtx, "message", "param1", "param2")
	ctxLogger.ErrorContext(lCtx, "message")
	ctxLogger.WarnContext(lCtx, "message", "param1", "param2")

	// child logger with additional context extraction functionality
	ctxChildLogger := ctxLogger.NewLog(log.Prefixed("context_child_logger"),
		log.WithCtxExtractor(func(ctx context.Context) []interface{} {
			return []interface{}{
				fmt.Sprintf("%s: %+v", k2, ctx.Value(k2)),
			}
		}),
		log.WithCtxMapExtractor(func(ctx context.Context) map[string]string {
			return map[string]string{
				k2: ctx.Value(k2).(string),
			}
		}),
	)

	ctxChildLogger.ErrorContext(lCtx, "message", "param1", "param2")
	ctxChildLogger.ErrorContext(lCtx, "message")
	ctxChildLogger.WarnContext(lCtx, "message", "param1", "param2")
}
