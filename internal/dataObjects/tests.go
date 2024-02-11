package dataObjects

import (
	"fmt"
	"time"
)

const (
	READ  = 0
	WRITE = 1

	PREWRITE  = 0
	POSTWRITE = 1
)

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
	Thread  int
	Command string
	Result  map[string]float64
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
