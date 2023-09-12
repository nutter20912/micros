package config

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var BasePath string

func init() {
	BasePath, _ = filepath.Abs(".")
}

// 載入配置
func Init(configName string) {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigFile(fmt.Sprintf("%s/config/%s.yaml", BasePath, configName))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}
