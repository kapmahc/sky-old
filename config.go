package sky

import "github.com/spf13/viper"

// IsProduction production mode ?
func IsProduction() bool {
	return viper.GetString("env") == "production"
}
