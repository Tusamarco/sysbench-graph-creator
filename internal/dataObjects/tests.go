package dataObjects

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"time"
)

const (
	READ  = 0
	WRITE = 1

	PREWRITE  = 0
	POSTWRITE = 1
)

type TestKey struct {
	Producer        string
	SelectPreWrites int
	TestName        string
	Dimension       string //Dimension is Large/Small
}
type ResultValueKey struct {
	Name   string
	Thread int
}

type ResultValue struct {
	ValueKey ResultValueKey
	Runs     map[int]float64
	STD      float64
	Lerror   float64
}

type ResultTest struct {
	Key             TestKey
	Producer        string
	ActionType      int //
	SelectPreWrites int
	TestName        string
	Dimension       string //Dimension is Large/Small
	Tests           map[TestKey]Test
	Results         ResultValue
	STD             float64
	Gerror          float64
}

type TestCollection struct {
	DateStart       time.Time //Date is coming from when was run the test
	DateEnd         time.Time //Date is coming from when was run the test
	Dimension       string    //Dimension is Large/Small
	ExecutionTime   int64
	TestName        string
	Producer        string //sysbench/tpcc/dbt3
	Tests           map[string]Test
	ActionType      int //
	SelectPreWrites int
	HostDB          string
	RunNumber       int
	Engine          string
	Name            string
	MySQLVersion    string
	MySQLProducer   string
	PBarr           *progressbar.ProgressBar
	FileName        string
}

type Test struct {
	Name          string
	DateStart     time.Time
	DateEnd       time.Time
	Dimension     string //Dimension is Large/Small
	ExecutionTime int64
	Labels        []string
	Threads       []int
	ThreadExec    map[int]Execution
	ActionType    int //
	Filter        string
	RunNumber     int
}

type Execution struct {
	Thread    int
	Command   string
	Result    map[string]float64
	DateStart time.Time
	DateEnd   time.Time
	Processed bool `default:false`
}

func getStringAction(code int) (string, error) {
	switch code {
	case 0:
		return "READ", nil
	case 1:
		return "WRITE", nil

	}
	return "", fmt.Errorf("Invalid code passed %d", code)
}
func getCodeAction(action string) (int, error) {
	switch action {
	case "select":
		return 0, nil
	case "write":
		return 1, nil

	}
	return 10, fmt.Errorf("Invalid action passed %s", action)
}
