package main

import (
	"fmt"

	"github.com/spf13/viper"
)


func main() {
	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.auto-scripts")
	viper.SetDefault("TemplateFolder", "$HOME/.auto-scripts/templates")
	viper.SetDefault("Templates", map[string]string{"Name": "Push Files", "FileName": "push-files.sh"})
	err := viper.WriteConfig() // writes current config to predefined path set by 'viper.AddConfigPath()' and 'viper.SetConfigName'
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

}
