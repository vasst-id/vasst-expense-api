package logs

import (
	"io"
	"os"
	"path"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Options struct {
	FileDirectory, FileName     string
	MaxSize, MaxAge, MaxBackups int
	ConsoleLog                  bool
	OutputLevel                 zerolog.Level
	DisableCompress             bool
}

type Logger struct {
	*zerolog.Logger
	outputLevel zerolog.Level
}

func New(opt Options) *Logger {
	var wr []io.Writer
	log.Logger = log.Output(zerolog.Logger{})

	if opt.ConsoleLog {
		wr = append(wr, zerolog.ConsoleWriter{Out: os.Stderr})
	}

	if opt.FileDirectory != "" {
		wr = append(wr, rotatedLogFile(opt))
	}

	mw := io.MultiWriter(wr...)

	logger := zerolog.New(mw).With().Timestamp().Logger()

	return &Logger{
		Logger:      &logger,
		outputLevel: opt.OutputLevel,
	}
}

func rotatedLogFile(opt Options) io.Writer {
	if err := os.MkdirAll(opt.FileDirectory, 0744); err != nil {
		log.Error().Err(err).Str("path", opt.FileDirectory).Msg("can't create log directory")
		return nil
	}

	return &lumberjack.Logger{
		Filename:   path.Join(opt.FileDirectory, opt.FileName),
		MaxBackups: opt.MaxBackups,       // files
		MaxSize:    opt.MaxSize,          // megabytes
		MaxAge:     opt.MaxAge,           // days
		Compress:   !opt.DisableCompress, // enable compression by default
	}
}

// Output implements nsq.logger / std log Output
func (l *Logger) Output(calldepth int, s string) error {
	l.WithLevel(l.outputLevel).Msg(s)
	return nil
}
