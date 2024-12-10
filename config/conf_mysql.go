package config

type Mysql struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Config   string `yaml:"config"` // 高级配置， 例如charset
	DB       string `yaml:"db"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	LogLevel string `yaml:"log_level"` // 日志等级，debug是输出全部sql， dev， release
}

func (m Mysql) Dsn() string {
	return m.User + ":" + m.Password + "@tcp(" + m.Host + m.Port + ")/" + m.DB + "?" + m.Config
}
