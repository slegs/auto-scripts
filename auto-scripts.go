package main

import (
	"fmt"
	"os"
	"os/exec"
	"log"
	"strings"
	"sort"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

var configpath string=os.Getenv("HOME") + "/.auto-scripts"
var currentpath string=os.Getenv("PWD") + "/"
var configname string="config"
var configtype string="yaml"

func main() {
	err:=initViper()

	if err != nil{
		log.Fatal("Error initialising \n %s ",err)
	}

	err = menu()

	if err != nil{
		log.Fatal("Error Menu \n %s ",err)
	}
}

func menu() error {
	app := cli.NewApp()
	app.Name = "auto-scripts"
	app.Description = "Auto Scripts - automating remote server management"
  app.Version = "1.0.0"

  app.Flags = []cli.Flag {
    cli.StringFlag{
      Name: "templates, t",
			Value: viper.GetString("templatefolder"),
      Usage: "Location of command templates",
    },
  }

	var tl=getTemplateList()
	var appCommands []cli.Command

	for _,template := range tl {
		var cmdArgs []cli.Flag
		var aL []string=strings.Split(getTemplateAttribute(template,"arguments"),",")

		cmdArgs=append(cmdArgs, cli.StringFlag{Name: "template",
											Value: getTemplateFilename(template),})
		for _,a := range aL {
			var f cli.Flag
			if getArgumentAttribute(a,"type") == "param" {
				f=cli.StringFlag{Name: a + "," + getArgumentAttribute(a,"flag"),
													Value: getArgumentAttribute(a,"default"),}
			} else if getArgumentAttribute(a,"type") == "flag"{
				f=cli.BoolFlag{Name: a + "," + getArgumentAttribute(a,"flag"),}
			}

			cmdArgs=append(cmdArgs, f)
		}

		appCommands=append(appCommands,cli.Command{
						Name: template,
						Usage: getTemplateFilename(template) + " - " +
			            	getTemplateAttribute(template,"description"),
						Flags: cmdArgs,
						Action: func(c *cli.Context) error {
											var (
												cmdOut []byte
												err    error
												args	[]string
											)
							        c.Command.VisibleFlags()


 											//fmt.Printf("Flags %#v\n", c.FlagNames())
											for _,a := range c.FlagNames() {
												//fmt.Printf("Argument %#v\n", a)
												if a == "template" {
														args = append(args,c.String(a))
												}
											}

											for _,a := range c.FlagNames() {
												//fmt.Printf("Argument %#v\n", a)
												if a != "template" {
													if c.String(a) != "" {
														args = append(args,"-" +  getArgumentAttribute(a,"flag"))
														args = append(args,c.String(a))
													}
												}
											}

											fmt.Printf("STARTING Execution of command %#v with args %#v",c.String("template"),args)

											cmdOut, err = exec.Command("/bin/bash",args...).Output()
											if err != nil {
												fmt.Printf("\nFAILED Execution of command %#v with args %#v\n",c.String("template"),args)
												return err
											}

											fmt.Printf("\nOUTPUT:\n%#s\n", cmdOut)


											fmt.Printf("SUCCESSFUL Execution of command %#v with args %#v\n",c.String("template"),args)
											return nil
										},
					})
		}

  app.Commands = appCommands;

  sort.Sort(cli.FlagsByName(app.Flags))
  sort.Sort(cli.CommandsByName(app.Commands))

  return app.Run(os.Args)

}


func initViper() error {
	var configfullpath string=configpath + "/" + configname + "." + configtype

	//Set viper config values
	viper.SetConfigType(configtype)
	viper.SetConfigName(configname)
	viper.AddConfigPath(configpath)

	//Create the config file if doesnt exist
	//Load in the default values
	err:=setDefaults(configfullpath)
	if err != nil{
			return fmt.Errorf("Error setting defaults: %s", err)
	}

	// Find and read the config file
	err = viper.ReadInConfig()
	if err != nil{
		return fmt.Errorf("Error reading config: %s", err)
	}

	// Create the templates folder if doesnt exist
	err=dirCreate(viper.GetString("templatefolder"))
	if err != nil{
		return fmt.Errorf("Error reading or creating template folder: %s", err)
	}

	// Check if Templates Exist
	return checkTemplates()

}

func checkTemplates() error{
	var tl=getTemplateList()
	for _,t := range tl {
		var fn=getTemplateFilename(t)
		if !fileExists(fn){
			return fmt.Errorf("Template: " + t + " Filename: " + fn + " does not exist. Please copy it into location.")
		}
	}
	return nil

}

func getArgumentList() []string{
	var tl=viper.GetStringMap("arguments")
	var argumentList []string
	for t := range tl {
		argumentList=append(argumentList,t)
	}
	return argumentList
}

func getArgumentAttribute(argument string, attribute string) string{
	 return viper.GetString("arguments." + argument + "." + attribute)
}

func getTemplateList() []string{
	var tl=viper.GetStringMap("templates")
	var templateList []string
	for t := range tl {
		templateList=append(templateList,t)
	}
	return templateList
}

func getTemplateAttribute(template string, attribute string) string{
	 return viper.GetString("templates." + template + "." + attribute)
}

func getTemplateFilename(template string) string{
	 return viper.GetString("templatefolder") + "/" + getTemplateAttribute(template,"filename")
}

func setDefaults(file string) error{
	err:=dirCreate(configpath) //if not exists
	if err != nil{
		return fmt.Errorf("Error reading or creating config path folder: %s", err)
	}

	if(!fileExists(file)){

		err=fileCreate(file) //if not exists

		if err != nil {
			return fmt.Errorf("Error reading or creating config file: %s", err)
		}

		viper.SetDefault("templatefolder", configpath + "/" + "templates")
		viper.SetDefault("templates",
			map[string]interface{}{"get": map[string]string{"filename": "get.sh",
				"description": "Pull files from server using rsync",
				"arguments": "backup,port,timestamp,directory,files,filter"},
				"local": map[string]string{"filename": "local.sh",
					"description": "Execute generic bash command on local machine",
					"arguments": "command"},
				"remote": map[string]string{"filename": "remote.sh",
					"description": "Execute generic bash command on local machine",
					"arguments": "command,port"},
				"manage": map[string]string{"filename": "manage.sh",
					"description": "Manage remote service using systemctl",
					"arguments": "service,port,servicecommand,lines"},
				"push": map[string]string{"filename": "push.sh",
					"description": "Push files from server using rsync",
					"arguments": "backup,port,timestamp,directory,remotebackup,files,filter"}})
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
						"default": ""},
					"files": map[string]string{"flag": "f",
						"type": "param",
						"mandatory": "false",
						"description": "Directory to store files (relative to local path)",
						"default": "files"},
					"directory": map[string]string{"flag": "d",
						"type": "param",
						"mandatory": "true",
						"description": "Directory for remote syncing in absolute then used relatively locally",
						"default": ""},
					"filter": map[string]string{"flag": "i",
						"type": "param",
						"mandatory": "false",
						"description": "For retieving a specific file or filter of files",
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
					"lines": map[string]string{"flag": "l",
						"type": "param",
						"mandatory": "false",
						"description": "Log lines to return",
						"default": "100"},
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
		return viper.WriteConfig() // writes current config to predefined path set by 'viper.AddConfigPath()' and 'viper.SetConfigName'
	}
	return nil
}

func fileExists(filename string) bool {

    _, err := os.Stat(filename)
    if os.IsNotExist(err) {
			return false
     }
		return true
}

func fileCreate(filename string) error {

    if !fileExists(filename) {
			var file, err = os.Create(filename)
			if err != nil{
				return(fmt.Errorf("Error create file %s",filename))
			}
			defer file.Close()
    }

		return nil
}

func dirCreate(dirname string) error {
    _, err := os.Stat(dirname)
		if os.IsNotExist(err) {
			return os.Mkdir(dirname, 0770)
		}
		return nil
}
