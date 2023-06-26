package util

// Config store all configuration of the application
// The values are read by viper from config file or environment variables
import (
	"github.com/spf13/viper"
)

type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

// LoadConfig read configuration from file environment variable
func LoadConfig(path string) (config Config, err error) {
	// absPath, err := filepath.Abs(path)
	// if err != nil {
	// 	return
	// }
	// println(absPath)
	// viper.AddConfigPath(absPath)
	viper.SetConfigFile("D:\\Studying\\go\\Simple Bank Project\\.env")
	// viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		return
	}
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
