package log

import (
	"context"

	"github.com/rs/zerolog"
)

// logParser contains parsing logic for a logger.
type jsonLogParser struct {
	*logOptions
	log zerolog.Logger
}

// newJsonLogParser creates a new instance of the json parser.
func newJsonLogParser(o *logOptions) jsonLogParser {
	return jsonLogParser{
		logOptions: o,
		log:        newZerolog(o),
	}
}

// print attaches concatenated v to the message field of the json as a single string.
func (l *jsonLogParser) print(v ...interface{}) {
	l.log.Print(v...)
}

// printf attaches the format parsed string to the message field of the json.
func (l *jsonLogParser) printf(format string, v ...interface{}) {
	l.log.Printf(format, v...)
}

// parse parses all additional data.
func (l *jsonLogParser) parse(ctx context.Context, event *zerolog.Event, prefix string, params ...interface{}) *zerolog.Event {
	event = l.withPrefix(event, prefix)
	event = l.withExtractedTrace(ctx, event)
	event = l.withExtractedCtx(ctx, event)
	event = l.withParams(event, params...)

	return event
}

// withLoggerPrefix attaches the logger prefix to the event.
func (l *jsonLogParser) withPrefix(event *zerolog.Event, prefix string) *zerolog.Event {
	const key string = "prefix"

	if l.prefix != "" {
		if prefix != "" {
			return event.Str(key, l.prefix+"."+prefix)
		}

		return event.Str(key, l.prefix)
	}

	if prefix != "" {
		return event.Str(key, prefix)
	}

	return event
}

// withExtractedTrace adds the extacted trace value to the event.
func (l *jsonLogParser) withExtractedTrace(ctx context.Context, event *zerolog.Event) *zerolog.Event {
	if l.ctxTraceExt != nil {
		if trace := l.ctxTraceExt(ctx); trace != "" {
			return event.Str("trace", trace)
		}
	}

	return event
}

// withExtractedCtx adds the extacted context values to the event.
func (l *jsonLogParser) withExtractedCtx(ctx context.Context, event *zerolog.Event) *zerolog.Event {
	if l.ctxExt != nil {
		if ctxData := l.ctxExt(ctx); len(ctxData) > 0 {
			return event.Interface("context", ctxData)
		}
	}

	return event
}

// withParams adds additional details to the event.
func (l *jsonLogParser) withParams(event *zerolog.Event, params ...interface{}) *zerolog.Event {
	if len(params) == 0 {
		return event
	}

	return event.Interface("details", params)
}
