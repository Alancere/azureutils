package logger

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
)

func LogDemo() {
	// slog.SetDefault(slog.Default().WithGroup("Alancere"))

	// slog.Info("my info 1")

	// logger := slog.Default().With("id", "!23131231")
	// parserLogger := logger.WithGroup("parser")
	// parserLogger.Info("my parse logger")

	// mylog := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// mylog.Info("xxx", "a", 1)

	// slog.Info("###release-v2", "msg", "my error msg")
	// slog.Error("###release-v2", "err", "my error msg")

	// mylog := slog.Default().WithGroup("Alancere")
	// mylog.Info("mylog info")
	// mylog.Debug("mylog debug")
	// mylog.Error("mylog error")
	// mylog.WithGroup("msg").Info("group info")

	// x := slog.With("a", "b", "c", "d")
	// y := x.WithGroup("a")
	// y.Info("y1", "y2", "v1")
	// y.Error("y2", "y2", "v2")

	// x := slog.New(new(XHandler))
	// x.Info("hahhaha")

	// slog.NewTextHandler()
	// l := slog.With("Command", "generator release-v2")
	// l.Info("my msg")
	// l.Debug("my debug")

	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		// AddSource: true,
		Level: slog.LevelInfo,
	})).With("Command", "release-v2")
	slog.SetDefault(l)

	slog.Info("my info")
	slog.Debug("my debug")
	slog.Warn("my warn")
	slog.Error("my error")

	log.Println("log print")
	fmt.Println("fmt print")
}

type XHandler struct{}

// _  = XHandler.(slog.Handler)

func (h *XHandler) Enabled(_ context.Context, level slog.Level) bool {
	return true
}

func (h *XHandler) Handle(_ context.Context, record slog.Record) error {
	fmt.Println(record)
	return nil
}

func (h *XHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *XHandler) WithGroup(name string) slog.Handler {
	return h
}
