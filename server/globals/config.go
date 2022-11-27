package globals

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fatih/structs"
	"github.com/jeremywohl/flatten"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Dashboard struct {
		Subdomain string
		Admin     struct {
			Username string
			Password string
		}
	}

	Server struct {
		Hostname string
		Port     uint16
	}
}

func ConfigInit(file string) *AppConfig {
	var appConf AppConfig

	viper.SetConfigType("yaml")

	viper.SetConfigName(filepath.Base(file))
	viper.AddConfigPath(filepath.Dir(file))
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("LS")
	BindStruct(viper.GetViper(), &appConf)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(fmt.Errorf("fatal error: couldn't load config file: %W | %T", err, err))
		}
	}

	if err := viper.Unmarshal(&appConf); err != nil {
		panic(fmt.Errorf("fatal error: couldn't Unmarshal config: %W", err))
	}

	return &appConf
}

// Viper in it'scurent version isn't able to read environment variables
// for attributes that don't have a default nor have been declared previously
// e.g. in a yaml. This becomes apparent when using `viper.Unmarshal`
//
// This method hints viper the possible environment variables given
// a config struct preventing the above mentioned problem.
func BindStruct(v *viper.Viper, input interface{}) error {
	// Transform config struct to map
	confMap := structs.Map(input)

	// Flatten nested conf map
	structKeys, err := flatten.Flatten(confMap, "", flatten.DotStyle)
	if err != nil {
		return err
	}

	for key, _ := range structKeys {
		if err := v.BindEnv(key); err != nil {
			return err
		}
	}

	return nil
}
