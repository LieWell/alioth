package core

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var GlobalConfig YamlConfig

type YamlConfig struct {
	Server YamlServer `yaml:"server"`
	Http   YamlHttp   `yaml:"http"`
	Mysql  YamlMysql  `yaml:"mysql"`
	Zap    YamlZap    `yaml:"zap"`
	JWT    YamlJWT    `yaml:"jwt"`
}

type YamlServer struct {
	Register bool `yaml:"register"` // 是否开启注册服务
}

type YamlHttp struct {
	Listen    string `yaml:"listen"`
	ListenTLS string `yaml:"listenTLS"`
	CertFile  string `yaml:"certFile"`
	KeyFile   string `yaml:"keyFile"`
}

type YamlMysql struct {
	Host               string `yaml:"host"`
	Port               string `yaml:"port"`
	Username           string `yaml:"username"`
	Password           string `yaml:"password"`
	Database           string `yaml:"database"`
	MaxOpenConnections int    `yaml:"maxOpenConnections"`
	MaxIdleConnections int    `yaml:"maxIdleConnections"`
}

type YamlZap struct {
	Level   string `yaml:"level"`
	File    string `yaml:"file"`
	MaxSize int    `yaml:"maxSize"` // 单位 Mi
	MaxAge  int    `yaml:"maxAge"`  // 单位 天
}

type YamlJWT struct {
	Secret   string   `yaml:"secret"`   // 密钥,签名字符串
	Expire   int      `yaml:"expire"`   // 过期时间,单位秒
	Issuer   string   `yaml:"issuer"`   // 签发人
	Audience []string `yaml:"audience"` // 受众
}

func LoadYamlConfig(filepath string) {

	raw, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("read config file error: %v", err)
	}

	var cfg YamlConfig
	err = yaml.Unmarshal(raw, &cfg)
	if err != nil {
		log.Fatalf("unmarshal config file error: %v", err)
	}

	// 记录所有配置的值
	GlobalConfig = cfg
}
