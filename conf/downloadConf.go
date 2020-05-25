package conf

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var serviceConfig ServiceConf
var viperLock sync.RWMutex

type BangumiConf struct {
	Title     string `mapstructure:"Title"`
	SearchExp string `mapstructure:"SearchExp"`
}

type DownLoadConf struct {
	BangumiConfigs []BangumiConf `mapstructure:"BangumiConfs"`
}

type ServiceConf struct {
	DownLoadConf `mapstructure:"DownLoadConf"`
}

func init() {
	viper.SetConfigType("yaml")
	viper.SetConfigFile("./resource/dwConf.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	opt := viper.DecodeHook(
		yamlStringToStruct(serviceConfig),
	)
	if err := viper.Unmarshal(&serviceConfig, opt); err != nil {
		panic(fmt.Sprintf("unable to decode into struct, %v", err))
	}
}

func GetServiceConfig() ServiceConf {
	viperLock.RLock()
	defer viperLock.RUnlock()
	return serviceConfig
}

func yamlStringToStruct(m interface{}) func(rf reflect.Kind, rt reflect.Kind, data interface{}) (interface{}, error) {
	return func(rf reflect.Kind, rt reflect.Kind, data interface{}) (interface{}, error) {
		if rf != reflect.String || rt != reflect.Struct {
			return data, nil
		}

		raw := data.(string)
		if raw == "" {
			return m, nil
		}

		return m, yaml.UnmarshalStrict([]byte(raw), &m)
	}
}
