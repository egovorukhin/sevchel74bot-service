package logger

import "github.com/egovorukhin/egolog"

type Config struct {
	DirPath  string    `yaml:"dirPath"`
	Info     string    `yaml:"info"`
	Error    string    `yaml:"error"`
	Debug    string    `yaml:"debug"`
	Rotation *Rotation `yaml:"rotation,omitempty"`
}

type Rotation struct {
	Size   int    `yaml:"size"`
	Format string `yaml:"format"`
	Path   string `yaml:"path"`
}

func Init(config Config) error {
	cfg := egolog.Config{
		DirPath:  config.DirPath,
		FileName: "app",
		Info:     egolog.Flags(config.Info),
		Error:    egolog.Flags(config.Error),
		Debug:    egolog.Flags(config.Debug),
		Rotation: &egolog.Rotation{
			Size:   config.Rotation.Size,
			Format: config.Rotation.Format,
			Path:   config.Rotation.Path,
		},
	}
	return egolog.InitLogger(cfg, callback)
}

func callback(egolog.InfoLog) {}
