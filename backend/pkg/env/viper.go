package env

import (
	"github.com/Novando/pintartek/pkg/logger"
	"github.com/spf13/viper"
	"strings"
)

// InitViper
// Initialize Viper to use the config file as env variable
func InitViper(path string, logger *logger.Logger) error {
	var configName string
	splitPaths := strings.Split(path, "/")
	if len(splitPaths) > 0 {
		for i := 0; i < len(splitPaths); i++ {
			configName = splitPaths[i]
		}
	}
	splitNames := strings.Split(configName, ".")
	if len(splitNames) < 2 {
		logger.Fatalf("Failed to parse config name")
	}
	formatName := splitNames[len(splitNames)-1]
	viper.SetConfigName(strings.TrimRight(configName, "."+formatName))
	viper.SetConfigType(formatName)
	viper.AddConfigPath(strings.TrimRight(path, configName))
	err := viper.ReadInConfig()
	if err != nil {
		logger.Infof("Configs file: %v", err)
	}
	return err
}
