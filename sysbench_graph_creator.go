package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	DO "sysbench-graph-creator/internal/dataObjects"
	global "sysbench-graph-creator/internal/global"
)

var version = "0.1.0"

var configFile string
var configPath string
var sourcePath string
var destinationPath string
var testName string

func main() {
	const (
		Separator = string(os.PathSeparator)
	)

	//initialize help
	help := new(global.HelpText)
	help.Init()

	//return version adn exit
	if len(os.Args) <= 1 &&
		os.Args[1] == "--version" {
		fmt.Println("Sysbench graph Creator version: ", version)
		exitWithCode(0)
	}

	//Manage config and parameters from conf file [start]
	flag.StringVar(&configFile, "configfile", "", "Config file name for the script")
	flag.StringVar(&configPath, "configpath", "", "Config file path")
	flag.StringVar(&sourcePath, "sourcepath", "", "source path")
	flag.StringVar(&destinationPath, "destinationpath", "", "destination path")
	flag.StringVar(&testName, "testname", "", "name of the test")
	//flag.StringVar(nil, "version", pxc_scheduler_handler_version, "version: ")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n%s\n", help.GetHelpText())
	}
	flag.Parse()

	//check for current params
	if len(os.Args) < 2 || configFile == "" {
		fmt.Println("You must at least pass the --configfile=xxx parameter ")
		exitWithCode(1)
	}
	var currPath, err = os.Getwd()

	if configPath != "" {
		if configPath[len(configPath)-1:] == Separator {
			currPath = configPath
		} else {
			currPath = configPath + Separator
		}
	} else {
		currPath = currPath + Separator + "config" + Separator
	}

	if err != nil {
		fmt.Print("Problem loading the config")
		exitWithCode(1)
	}

	//Return our full configuration from file
	var config = global.GetConfig(currPath + configFile)
	//initialize the log system
	if !global.InitLog(config) {
		fmt.Println("Not able to initialize log system exiting")
		exitWithCode(1)
	}

	//commandline override config file
	portingCommandOption(config)

	//now the show begins
	myArFiles, err1 := DO.GetFileList(config.Parser.SourceDataPath)
	if err1 != nil {
		log.Error(err1)
		exitWithCode(1)
	}

	log.Debug(len(myArFiles))
}
func portingCommandOption(config global.Configuration) {
	if sourcePath != "" {
		config.Parser.SourceDataPath = sourcePath
	}
	if destinationPath != "" {
		config.Render.DestinationPath = destinationPath
	}
	if testName != "" {
		config.Global.TestName = testName
	}
}

func exitWithCode(errorCode int) {
	log.Debug("Exiting execution with code ", errorCode)
	os.Exit(errorCode)
}
