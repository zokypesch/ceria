package helper

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// ReadConfig for get config yaml file
type ReadConfig struct {
	prop string
}

// singleton helper
var cfg *ReadConfig

// NewReadConfigService for new service config
func NewReadConfigService() *ReadConfig {
	if cfg == nil {
		cfg = &ReadConfig{}
	}
	return cfg
}

// GetByName Find and read the config file
func (cfg *ReadConfig) GetByName(params string) string {

	name := params //strings.ToUpper(params)
	if name == "" {
		return ""
	}

	s := strings.Split(name, ".")

	if len(s) > 1 {
		mapCfg := viper.GetStringMapString(s[0])
		return mapCfg[s[1]]
	}
	return viper.GetString(s[0])
}

// Init initial function
func (cfg *ReadConfig) Init() error {
	var err error

	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath("../")
	viper.AddConfigPath(".")   // path to look for the config file in
	err = viper.ReadInConfig() // Find and read the config file

	if err != nil {
		return err
	}
	newCfg := fmt.Sprintf("config.%s", viper.Get("MODE").(string))
	viper.SetConfigName(newCfg)
	err = viper.ReadInConfig()

	return err
}
