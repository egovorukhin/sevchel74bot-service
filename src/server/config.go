package server

type Server struct {
	Certificate *Certificate `yaml:"certificate"`
	Http        Http         `yaml:"http"`
	Webhook     Webhook      `yaml:"webhook"`
}

type Timeout struct {
	Read  int `yaml:"read"`
	Write int `yaml:"write"`
	Idle  int `yaml:"idle"`
}

type Buffer struct {
	Read  int `yaml:"read"`
	Write int `yaml:"write"`
}

type Certificate struct {
	Cert       string `yaml:"cert"`
	Key        string `yaml:"key"`
	ClientCert string `yaml:"clientCert"`
}
