package loggers

import (
	"log/syslog"
)

// https://stackoverflow.com/questions/26120698/how-to-change-the-date-time-format-of-gos-log-package
// Maybe add more  customization to the logs

type LogWriter struct {
	Writer *syslog.Writer
}

func NewWriter(host string) *LogWriter {
	w, err := syslog.Dial("udp", host, syslog.LOG_EMERG | syslog.LOG_KERN , "skele-bot")
	if err != nil {
		panic(err)
	}	

	return &LogWriter{
		Writer: w,
	}
}

func (l *LogWriter) LogInfo (s string) {
	l.Writer.Info(s)
}

func (l *LogWriter) LogError (s string) {
	l.Writer.Err(s)
}

func (l *LogWriter) LogNotice (s string) {
	l.Writer.Notice(s)
}
