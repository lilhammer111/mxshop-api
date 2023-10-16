package config

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	Goods     string `mapstructure:"order"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataId"`
	Group     string `mapstructure:"group"`
}

type ServerConfig struct {
	Name             string       `mapstructure:"name" json:"name"`
	Host             string       `mapstructure:"host" json:"host"`
	Port             int          `mapstructure:"port" json:"port"`
	Tags             []string     `mapstructure:"tags" json:"tags"`
	GoodsSrvInfo     SrvConfig    `mapstructure:"goods_srv" json:"goods_srv"`
	OrderSrvInfo     SrvConfig    `mapstructure:"order_srv" json:"order_srv"`
	InventorySrvInfo SrvConfig    `mapstructure:"inventory_srv" json:"inventory_srv"`
	JWTInfo          JWTConfig    `mapstructure:"jwt" json:"jwt"`
	ConsulInfo       ConsulConfig `mapstructure:"consul" json:"consul"`
	AliPayInfo       AliPayConfig `mapstructure:"alipay" json:"alipay"`
}

type SrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type AliPayConfig struct {
	AppID        string `mapstructure:"app_id" json:"app_id"`
	PrivateKey   string `mapstructure:"private_key" json:"private_key"`
	AliPublicKey string `mapstructure:"ali_public_key" json:"ali_public_key"`
	NotifyURL    string `mapstructure:"notify_url" json:"notify_url"`
	ReturnURL    string `mapstructure:"return_url" json:"return_url"`
}
