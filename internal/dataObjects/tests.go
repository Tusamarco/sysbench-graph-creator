package dataObjects

import "time"

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
	Producer        string //MySQL/Percona Server/PXC /Maria
	Version         string //Producer version
	Tests           map[string]Test
	ActionType      int //
	SelectPreWrites int
	HostDB          string
	RunNumber       int
}

type Test struct {
	Date          string
	Dimension     string //Dimension is Large/Small
	ExecutionTime int64
	Labels        []string
	TestType      string //sysbench/tpcc/dbt3
	Threads       []int
	runs          map[int]run
	ActionType    int //
}

type run struct {
	thread int
	result map[string]float64
}
