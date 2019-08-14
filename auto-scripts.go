package main

import (
	"fmt"
	"os"
	"github.com/spf13/viper"
)

func main() {
	initViper()
}

func initViper(){
	var configpath string=os.Getenv("HOME") + "/.auto-scripts"
	var configname string="config"
	var configtype string="yaml"
	var configfullpath string=configpath + "/" + configname + "." + configtype

	//Set viper config values
	viper.SetConfigType(configtype)
	viper.SetConfigName(configname)
	viper.AddConfigPath(configpath)

	//Create the config file if doesnt exist
	//Load in the default values
	setDefaults(configpath,configfullpath)

	// Find and read the config file
	err := viper.ReadInConfig()
	check(err,"panic")

	// Create the templates folder if doesnt exist
	dirCreate(viper.GetString("templatefolder"))

	// Check if Templates Exist
	if checkTemplates(viper.GetString("templatefolder")) {
		fmt.Println("Templates Check OK - Continuing")
	} else {
		fmt.Println("Templates Check FAILED - Exiting")
	}
}

func checkTemplates(path string) bool{
	var result bool=true
	var tl=viper.GetStringMap("templates")
	for t := range tl {
		var fn=path+"/"+viper.GetString("templates."+t+".filename")
		if !fileExists(fn){
			fmt.Println("Template: " + t + " Filename: " + fn + " does not exist. Please copy it into location.")
			result=false
		}
	}
	return result

}
func setDefaults(path string, file string){
	dirCreate(path) //if not exists
	if(fileCreate(file)){ //if not exists
		viper.SetDefault("templatefolder", path + "/" + "templates")
		viper.SetDefault("filesfolder", "files")
		viper.SetDefault("templates",
			map[string]interface{}{"get": map[string]string{"filename": "get-files.sh",
				"description": "Pull files from server using rsync",
				"arguments": "backup,port,timestamp,directory",
				"backup": "true"},
				"local": map[string]string{"filename": "local.sh",
					"description": "Execute generic bash command on local machine",
					"arguments": "command"},
				"remote": map[string]string{"filename": "remote.sh",
					"description": "Execute generic bash command on local machine",
					"arguments": "command,port"},
				"manage": map[string]string{"filename": "manage.sh",
					"description": "Manage remote service using systemctl",
					"arguments": "service,port,servicecommand"},
				"push": map[string]string{"filename": "push-files.sh",
					"description": "Push files from server using rsync",
					"arguments": "backup,port,timestamp,directory",
					"backup": "true"}})
		viper.SetDefault("Arguments",
			map[string]interface{}{"port": map[string]string{"flag": "p",
					"type": "param",
					"mandatory": "false",
					"description": "Remote Port",
					"default": "2121"},
					"timestamp": map[string]string{"flag": "t",
						"type": "param",
						"mandatory": "false",
						"description": "Timestamp to be used in backups",
						"default": "$(date +%Y%m%d%H%M%S)"},
					"directory": map[string]string{"flag": "d",
						"type": "param",
						"mandatory": "true",
						"description": "Directory for remote syncing in absolute then used relatively locally",
						"default": ""},
					"command": map[string]string{"flag": "c",
						"type": "param",
						"mandatory": "true",
						"description": "Command to be passed to script",
						"default": ""},
					"servicecommand": map[string]string{"flag": "v",
						"type": "param",
						"mandatory": "true",
						"description": "systemctl comands Start, Stop, Restart or Status",
						"default": ""},
					"service": map[string]string{"flag": "s",
						"type": "param",
						"mandatory": "true",
						"description": "Service name",
						"default": ""},
					"backup": map[string]string{"flag": "b",
						"type": "param",
						"mandatory": "false",
						"description": "Backup Directory (relative for local and absolute for remote)",
						"default": "backup"},
					"remotebackup": map[string]string{"flag": "r",
						"type": "flag",
						"mandatory": "false",
						"description": "Indicate if a full remote backup required",
						"default": ""}})
		err := viper.WriteConfig() // writes current config to predefined path set by 'viper.AddConfigPath()' and 'viper.SetConfigName'
		check(err,"panic")
	}
}

func fileExists(filename string) bool {

    _, err := os.Stat(filename)
    if os.IsNotExist(err) {
			return false
     }
		return true
}

func fileCreate(filename string) bool {

    if !fileExists(filename) {
			var file, err = os.Create(filename)
			check(err,"panic");
			return true
			defer file.Close()
    }
		return false
}

func dirCreate(dirname string) {
    _, err := os.Stat(dirname)
		if os.IsNotExist(err) {
			err = os.Mkdir(dirname, 0770)
			check(err,"panic");
		}
}

func check(e error, errtype string) {
			if e != nil {
				switch errtype {
						case "panic":
								panic(fmt.Errorf("%s \n", e))
				    default:
				        fmt.Errorf("%s \n", e)
				    }
			}
}
