package dataObjects

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"time"
)

const (
	READ           = 0
	WRITE          = 1
	READ_AND_WRITE = 10

	PREWRITE  = 0
	POSTWRITE = 1

	SMALL = "small"
	LARGE = "large"

	TPCC = "tpcc"
	SB   = "sysbench"
	DBT3 = "dbt3"
)

func DIMENSIONS() []string {
	return []string{SMALL, LARGE}
}

func ACTIONTYPES() []int {
	return []int{READ, WRITE, READ_AND_WRITE}
}
func PREPOSTWRITE() []int {
	return []int{PREWRITE, POSTWRITE}
}

type Producer struct {
	MySQLProducer       string
	MySQLVersion        string
	TestsResults        []ResultTest
	TestsTypes          []TestType
	TestCollectionsName string
	STDReadPre          float64
	GerrorReadPre       float64
	STDReadPost         float64
	GerrorReadPost      float64
	STDRWrite           float64
	GerrorWrite         float64
	Color               string
}

type TestType struct {
	Name               string
	Dimension          string
	SelectPreWrites    int
	ActionType         int
	TestCollectionName string
}
type TestKey struct {
	ActionType         int
	TestCollectionName string
	MySQLProducer      string
	MySQLVersion       string
	SelectPreWrites    int
	TestName           string
	Dimension          string //Dimension is Large/Small
}
type ResultValueKey struct {
	Name   string
	Thread int
}

type ResultValue struct {
	ThreadNumber int
	Label        string
	Value        float64
	STD          float64
	Lerror       float64
}

type ResultTest struct {
	Key        TestKey
	Labels     map[string][]ResultValue
	Executions int
	STD        float64
	Gerror     float64
}

type TestCollection struct {
	DateStart        time.Time //Date is coming from when was run the test
	DateEnd          time.Time //Date is coming from when was run the test
	Dimension        string    //Dimension is Large/Small
	ExecutionTime    int64
	TestName         string
	Producer         string //sysbench/tpcc/dbt3
	Tests            map[string]Test
	ActionType       int //
	SelectPostWrites int
	HostDB           string
	RunNumber        int
	Engine           string
	Name             string
	MySQLVersion     string
	MySQLProducer    string
	PBarr            *progressbar.ProgressBar
	FileName         string
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
	case 10:
		return "read/write", nil

	}
	return "", fmt.Errorf("Invalid code passed %d", code)
}
func getCodeAction(action string) (int, error) {
	switch action {
	case "select":
		return 0, nil
	case "write":
		return 1, nil
	case "read/write":
		return 10, nil

	}
	return 10, fmt.Errorf("Invalid action passed %s", action)
}
