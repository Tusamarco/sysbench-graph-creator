package Global

import (
	"github.com/Tusamarco/toml"
	log "github.com/sirupsen/logrus"
	"syscall"
)

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
}

// Main structure working as container for the configuration sections
// So far only 2 but this may increase like logs for instance
type Configuration struct {
	Parser Parser    `toml:"parser"`
	Render Render    `toml:"render"`
	Global GlobalDef `toml:"global"`
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

func (conf *Configuration) fillDefaults() {
	//conf.Parser.sourceDataPath=""
}
