package main

import (
	"fmt"
	"os"
	"github.com/spf13/viper"
)


func main() {
	var configpath string=os.Getenv("HOME")+"/.auto-scripts"
	var configname string="config"
	var configtype string="yaml"
	var configfullpath string=configpath + "/" +configname+"."+configtype

	viper.SetConfigType(configtype)
	viper.SetConfigName(configname)
	viper.AddConfigPath(configpath)

	dirCreate(configpath) //if not exists
	fileCreate(configfullpath) //if not exists

	viper.SetDefault("TemplateFolder", "templates")
	viper.SetDefault("Templates", map[string]string{"Name": "Push Files", "FileName": "push-files.sh"})
	err := viper.WriteConfig() // writes current config to predefined path set by 'viper.AddConfigPath()' and 'viper.SetConfigName'
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func fileCreate(filename string) {
    _, err := os.Stat(filename)
    if os.IsNotExist(err) {
			var file, err = os.Create(filename)
			if err != nil {
				panic(fmt.Errorf("Fatal error create file: %s \n", err))
				return
			}
			defer file.Close()
    }
}

func dirCreate(dirname string) {
    _, err := os.Stat(dirname)
		if os.IsNotExist(err) {
			err = os.Mkdir(dirname, 0770)
			if err != nil {
				panic(fmt.Errorf("Fatal error create directory: %s \n", err))
			}
		}
}
