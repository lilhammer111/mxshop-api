package config

type ServerConfig struct {
	Name         string         `mapstructure:"name" json:"name"`
	Host         string         `mapstructure:"host" json:"host"`
	Port         int            `mapstructure:"port" json:"port"`
	Tags         []string       `mapstructure:"tags" json:"tags"`
	GoodsSrvInfo GoodsSrvConfig `mapstructure:"goods_srv" json:"goods_srv"`
	JWTInfo      JWTConfig      `mapstructure:"jwt" json:"jwt"`
	ConsulInfo   ConsulConfig   `mapstructure:"consul" json:"consul"`
}

type GoodsSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	Goods     string `mapstructure:"goods"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataId"`
	Group     string `mapstructure:"group"`
}
