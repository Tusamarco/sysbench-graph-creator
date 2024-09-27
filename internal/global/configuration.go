package Global

import (
	"github.com/Tusamarco/toml"
	log "github.com/sirupsen/logrus"
	"os"

	"syscall"
)

// commandline params
type Params struct {
	ConfigFile         string
	ConfigPath         string
	SourceDataPath     string
	DestinationPath    string
	CsvDestinationPath string
	FilterByProducer   string
	FilterByVersion    string
	FilterByDimension  string
	FilterByTitle      string
	Labels             string
	ConvertChartsToCsv bool
	PrintCharts        bool
	PrintData          bool
	FilterByPrePost    string
	TestName           string
}

// Global scheduler conf
type GlobalDef struct {
	TestName    string //global test name ie: Percona Server VS MySQL
	LogLevel    string
	LogTarget   string // #stdout | file
	LogFile     string //"/tmp/pscheduler"
	Performance bool
}

type Parser struct {
	SourceDataPath  string `toml:sourceDataPath`
	FilterOutliners bool   `toml:filterOutliners`
	DistanceLabel   string `toml:distanceLabel`
}

type Color struct {
	Colors []string `toml:"color"`
}

type Render struct {
	GraphType           string `toml:graphType`
	DestinationPath     string `toml:destinationPath`
	PrintStats          bool   `toml:printStats`
	PrintData           bool   `toml:printData`
	HttpServerPort      int    `toml:httpServerPort`
	HttpServerIp        string `toml:httpServerIp`
	Labels              string `toml:labels`
	StatsLabels         string `toml:statslabels`
	ReadSummaryLabel    string `toml:readSummaryLabel`
	WriteSummaryLabel   string `toml:writeSummaryLabel`
	ChartHeight         int    `toml:chartHeight`
	ChartWidth          int    `toml:chartWidth`
	PrintCharts         bool   `toml:printCharts`
	PrintChartsFormat   string `toml:printChartsFormat`
	ConvertChartsToCsv  bool   `toml:convertChartsToCsv`
	CsvDestinationPath  string `toml:csvDestinationPath`
	HtmlDestinationPath string
	FilterTestsByTitle  string `toml:filterTestsByTitle`
	FilterByDimension   string `toml:filterByDimension`
	FilterByVersion     string `toml:filterByVersion`
	FilterByProducer    string `toml:filterByProducer`
	FilterByPrePost     string `toml:filterByPrePost`
}

// Main structure working as container for the configuration sections
// So far only 2 but this may increase like logs for instance
type Configuration struct {
	Parser Parser    `toml:"parser"`
	Render Render    `toml:"render"`
	Global GlobalDef `toml:"global"`
	Colors Color     `toml:"colors"`
}

// Methods to return the config as map
func GetConfig(path string) Configuration {
	var config Configuration
	config.fillDefaults()
	if _, err := toml.DecodeFile(path, &config); err != nil {
		log.Error(err)
		syscall.Exit(2)
	}
	return config
}

func GetParams() Params {
	var params Params
	return params
}

func (conf *Configuration) fillDefaults() {
	//conf.Parser.sourceDataPath=""
}

// We assign the value coming from command line to config
func (conf *Configuration) ParseCommandLine(params Params) {

	if params.SourceDataPath != "" {
		conf.Parser.SourceDataPath = params.SourceDataPath
	}
	if params.DestinationPath != "" {
		conf.Render.DestinationPath = params.DestinationPath
	}
	if params.TestName != "" {
		conf.Global.TestName = params.TestName
	}

	if params.CsvDestinationPath != "" {
		conf.Render.CsvDestinationPath = params.CsvDestinationPath
	}
	if params.Labels != "" && conf.Render.Labels == "" {
		conf.Render.Labels = params.Labels
	}
	if params.PrintData && !conf.Render.PrintData {
		conf.Render.PrintData = params.PrintData
	}

	if params.PrintCharts && !conf.Render.PrintCharts {
		conf.Render.PrintCharts = params.PrintCharts
	}

	if params.ConvertChartsToCsv && !conf.Render.ConvertChartsToCsv {
		conf.Render.ConvertChartsToCsv = params.ConvertChartsToCsv
	}

	if params.FilterByVersion != "" {
		conf.Render.FilterByVersion = params.FilterByVersion
	}
	if params.FilterByProducer != "" {
		conf.Render.FilterByProducer = params.FilterByProducer
	}
	if params.FilterByTitle != "" {
		conf.Render.FilterTestsByTitle = params.FilterByTitle
	}
	if params.FilterByDimension != "" {
		conf.Render.FilterByDimension = params.FilterByDimension
	}
	if params.FilterByPrePost != "" {
		conf.Render.FilterByPrePost = params.FilterByPrePost
	}

}

// We check paths and if not existing we will create them
func (conf *Configuration) CheckPaths() {
	paths := []string{conf.Render.DestinationPath, conf.Render.HtmlDestinationPath, conf.Render.CsvDestinationPath}
	for _, path := range paths {
		if path != "" {
			if !CheckIfPathExists(path) {
				CreatePath(path)
				log.Infof("Path %s not exists. Will create", path)
			}
		}
	}
}

// we perform sanity check on the incoming paramters, wew ill fix when possible exit otherwise
func (conf *Configuration) SanityChecks() {

	if !CheckIfPathExists(conf.Parser.SourceDataPath) {
		log.Errorf("Source Path %s  does not exist cannot proceed.", conf.Parser.SourceDataPath)
		os.Exit(1)
	}

}
