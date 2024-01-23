package dataObjects

const (
	READ  = 0
	WRITE = 1
)

type TestCollection struct {
	DateStart     string //Date is coming from when was run the test
	DateEnd       string //Date is coming from when was run the test
	Dimension     string //Dimension is Large/Small
	ExecutionTime int64
	TestName      string
	Producer      string //MySQL/Percona Server/PXC /Maria
	Version       string //Producer version
	Tests         map[string]Test
}

type Test struct {
	Date          string
	Dimension     string //Dimension is Large/Small
	ExecutionTime int64
	Labels        []string
	TestType      string //sysbench/tpcc/dbt3
	Threads       []int
	ActionType    int //
	runs          map[int]run
}

type run struct {
	thread int
	result map[string]float64
}
