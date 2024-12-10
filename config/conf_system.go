package config

type System struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Env  string `yaml:"env"`
}

func (s System) Addr() string {
	return s.Host + s.Port
}
