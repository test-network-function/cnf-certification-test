package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// Custom log levels.
const CustomLevelFatal = slog.Level(12)

var CustomLevelNames = map[slog.Leveler]string{
	CustomLevelFatal: "FATAL",
}

type CustomHandler struct {
	opts  slog.HandlerOptions
	attrs []slog.Attr
	mu    *sync.Mutex
	out   io.Writer
}

func NewCustomHandler(out io.Writer, opts *slog.HandlerOptions) *CustomHandler {
	h := &CustomHandler{out: out, mu: &sync.Mutex{}}
	if opts != nil {
		h.opts = *opts
	}
	if h.opts.Level == nil {
		h.opts.Level = slog.LevelInfo
	}

	return h
}

func (h *CustomHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

// The Handle method will write a log line with the following format:
// LOG_LEVEL [TIME] [SOURCE_FILE] [CUSTOM_ATTRS] MSG
//
//nolint:gocritic // r param is heavy but defined in the slog.Handler interface
func (h *CustomHandler) Handle(_ context.Context, r slog.Record) error {
	var buf []byte
	// Level
	var levelAttr slog.Attr
	if h.opts.ReplaceAttr != nil {
		levelAttr = h.opts.ReplaceAttr(nil, slog.Any(slog.LevelKey, r.Level))
	} else {
		levelAttr = slog.Any(slog.LevelKey, r.Level)
	}
	buf = h.appendAttr(buf, levelAttr)
	// Time
	if !r.Time.IsZero() {
		buf = h.appendAttr(buf, slog.Time(slog.TimeKey, r.Time))
	}
	// Source
	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		buf = h.appendAttr(buf, slog.String(slog.SourceKey, fmt.Sprintf("%s: %d", filepath.Base(f.File), f.Line)))
	}
	// Attributes
	for _, attr := range h.attrs {
		buf = h.appendAttr(buf, attr)
	}
	// Message
	buf = h.appendAttr(buf, slog.String(slog.MessageKey, r.Message))
	buf = append(buf, "\n"...)
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.out.Write(buf)
	return err
}

// Not implemented. Returns the nil handler.
func (h *CustomHandler) WithGroup(_ string) slog.Handler {
	return nil
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}

	// Create a new handler with default attributes
	h2 := *h
	// A deep copy of the attributes is required
	h2.attrs = make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(h2.attrs, h.attrs)
	h2.attrs = append(h2.attrs, attrs...)

	return &h2
}

func (h *CustomHandler) appendAttr(buf []byte, a slog.Attr) []byte {
	// Resolve the Attr's value before doing anything else.
	a.Value = a.Value.Resolve()
	// Ignore empty Attrs.
	if a.Equal(slog.Attr{}) {
		return buf
	}
	switch a.Value.Kind() {
	case slog.KindString:
		if a.Key == slog.MessageKey {
			buf = fmt.Appendf(buf, "%s", a.Value.String())
		} else {
			buf = fmt.Appendf(buf, "[%s] ", a.Value.String())
		}
	case slog.KindTime:
		buf = fmt.Appendf(buf, "[%s] ", a.Value.Time().Format(time.StampMilli))
	default:
		if a.Key == slog.LevelKey {
			buf = fmt.Appendf(buf, "%-5s ", a.Value.String())
		} else {
			buf = fmt.Appendf(buf, "[%s: %s] ", a.Key, a.Value)
		}
	}
	return buf
}
