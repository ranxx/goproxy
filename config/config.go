package config

import (
	"encoding/json"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config 配置
type Config struct {
	App     string   `json:"app" yaml:"app"`
	Server  Server   `json:"server" yaml:"server"`
	Clients []string `json:"clients" yaml:"clients"`
}

// Server 服务配置
type Server struct {
	IP      string `json:"ip" yaml:"ip"`
	Port    int    `json:"port" yaml:"port"`
	APIPort int    `json:"apiPort" yaml:"apiPort"`
}

// var
var (
	// Cfg 配置文件
	Cfg *Config

	// 解析函数
	parser map[string]func(file string) *Config = map[string]func(file string) *Config{
		"yaml": ParseYamlFile,
		"json": ParseJSONFile,
		"":     ParseYamlFile,
	}
)

// ParseYamlFile 解析 yaml 配置文件
func ParseYamlFile(file string) *Config {
	body, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	Cfg = new(Config)
	if err = yaml.Unmarshal(body, Cfg); err != nil {
		panic(err)
	}
	return Cfg
}

// ParseJSONFile 解析 json 配置文件
func ParseJSONFile(file string) *Config {
	body, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	Cfg = new(Config)
	if err = json.Unmarshal(body, Cfg); err != nil {
		panic(err)
	}
	return Cfg
}

// ParseFile 解析配置文件，会根据 file 的后缀名匹配不同的解析函数
// 目前支持 yaml，json；其余默认使用 yaml 方式进行解析
func ParseFile(file string) *Config {
	// 截取后缀名
	lindex := strings.LastIndex(file, ".")
	if lindex <= 0 || lindex == len(file)-1 {
		return ParseYamlFile(file)
	}

	suffix := file[lindex+1:]
	v, ok := parser[suffix]
	if !ok {
		return ParseYamlFile(file)
	}
	return v(file)
}
