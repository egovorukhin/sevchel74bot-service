package server

import "github.com/egovorukhin/egolog"

type Logger struct {
	Format   string `yaml:"format"`
	Filename string `yaml:"filename"`
	Time     Time   `yaml:"time"`
}

type Time struct {
	Format   string `yaml:"format"`
	Zone     string `yaml:"zone"`
	Interval int    `yaml:"interval"`
}

func (l *Logger) Write(data []byte) (n int, err error) {
	egolog.Infofn(l.Filename, string(data))
	return len(data), nil
}
