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

//var configFile string
//var configPath string
//var sourcePath string
//var destinationPath string
//var testName string

func main() {
	const (
		Separator = string(os.PathSeparator)
	)
	params := global.GetParams()
	//initialize help
	help := new(global.HelpText)
	help.Init()

	//return version adn exit
	if len(os.Args) == 1 {
		if os.Args[1] == "--version" {
			fmt.Println("Sysbench graph Creator version: ", version)
			exitWithCode(0)
		}
	}

	//Manage config and parameters from conf file [start]
	//Mandatory from command line
	flag.StringVar(&params.ConfigFile, "configfile", "", "Config file name for the script")
	flag.StringVar(&params.ConfigPath, "configpath", "", "Config file path")

	//mandatory (at least in config)
	flag.StringVar(&params.SourceDataPath, "sourceDataPath", "", "source path")
	flag.StringVar(&params.DestinationPath, "destinationPath", "", "destination path")
	flag.StringVar(&params.CsvDestinationPath, "csvDestinationPath", "", "csv destination path")
	flag.StringVar(&params.TestName, "testName", "", "Name of the test. IE 'Comparison MySQL 5.7 VS 8.0.38")

	//optional
	flag.StringVar(&params.FilterByProducer, "filterByProducer", "", "filter by producer(s) name, comma separated list")
	flag.StringVar(&params.FilterByVersion, "filterByVersion", "", "filter by version(s) name, comma separated list")
	flag.StringVar(&params.FilterByDimension, "filterByDimension", "", "filter by dimension(s) name, comma separated list")
	flag.StringVar(&params.FilterByTitle, "filterByTitle", "", "filter by test name(s) name, comma separated list")
	flag.StringVar(&params.FilterByPrePost, "filterByPrePost", "", "filter by pre or post write action , comma separated list [pre|post]. Default: pre,post")

	flag.StringVar(&params.Labels, "labels", "TotalTime,Events/s,operations/s,writes/s,reads/s,latencyPct95(μs)", "list of labels to use (comma separated) default: TotalTime,Events/s,operations/s,writes/s,reads/s,latencyPct95(μs)")
	flag.BoolVar(&params.ConvertChartsToCsv, "convertCsv", false, "if to convert to csv [false|true]")
	flag.BoolVar(&params.PrintCharts, "printCharts", false, "if to create jpeg images of the charts [false|true]")
	flag.BoolVar(&params.PrintData, "printData", true, "if to create html file version")

	//flag.StringVar(nil, "version", pxc_scheduler_handler_version, "version: ")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n%s\n", help.GetHelpText())
		flag.PrintDefaults()
	}
	flag.Parse()

	//check for current params
	if len(os.Args) < 2 || params.ConfigFile == "" {
		fmt.Println("You must at least pass the --configfile=xxx and --configpath parameters ")
		exitWithCode(1)
	}
	var currPath, err = os.Getwd()

	if params.ConfigPath != "" {
		if params.ConfigPath[len(params.ConfigPath)-1:] == Separator {
			currPath = params.ConfigPath
		} else {
			currPath = params.ConfigPath + Separator
		}
	} else {
		currPath = currPath + Separator + "config" + Separator
	}

	if err != nil {
		fmt.Print("Problem loading the config")
		exitWithCode(1)
	}

	//Return our full configuration from file
	var config = global.GetConfig(currPath + params.ConfigFile)

	//parameters from command line have higher priority so we parse and replace
	config.ParseCommandLine(params)

	//Sanity Checks before taking any action
	config.SanityChecks()

	//check for path and create the ones not existing
	config.CheckPaths()

	//initialize the log system
	if !global.InitLog(config) {
		fmt.Println("Not able to initialize log system exiting")
		exitWithCode(1)
	}

	//now the show begins
	fileProc := new(DO.FileProcessor)
	err1 := fileProc.GetFileList(config.Parser.SourceDataPath)
	testCollection, err1 := fileProc.GetTestCollectionArray()
	calculator := new(DO.Calculator)
	calculator.Init(config)
	testsResults := calculator.BuildResults(testCollection)
	producersAr := calculator.GroupByProducers()

	log.Infof("Test Results %d", len(testsResults))
	log.Infof("Test collection %d", len(testCollection))
	log.Infof("# of producers %d", len(producersAr))
	log.Infof("Producers STD and Distance")
	for _, producer := range producersAr {
		log.Infof("Producer: %s: %s: test name: %s", producer.MySQLProducer, producer.MySQLVersion, producer.TestCollectionsName)
		log.Infof("		READS PRE WRITES  STD: %.4f Dist(pct): %.4f", producer.STDReadPre, producer.GerrorReadPre)
		log.Infof("		READS POST WRITES  STD: %.4f Dist(pct): %.4f", producer.STDReadPost, producer.GerrorReadPost)
		log.Infof("		WRITES  STD: %.4f Dist(pct): %.4f", producer.STDRWrite, producer.GerrorWrite)
	}
	if err1 != nil {
		log.Error(err1)
		exitWithCode(1)
	}
	graph := new(DO.GraphGenerator)
	graph.Init(config, producersAr, testCollection)

	if graph.RenderReults() {
		if graph.BuildPage() {
			graph.ActivateHTTPServer()
		}

	}
	//graph.Test3()

	exitWithCode(0)
	//log.Debug(len(myArFiles))
}

func exitWithCode(errorCode int) {
	log.Debug("Exiting execution with code ", errorCode)
	os.Exit(errorCode)
}
