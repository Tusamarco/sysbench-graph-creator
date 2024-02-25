package dataObjects

import (
	"github.com/montanaflynn/stats"
	log "github.com/sirupsen/logrus"
	"math"
	"sort"
	global "sysbench-graph-creator/internal/global"
)

//gonum.org/v1/gonum/stat
//github.com/montanaflynn/stats

type Calculator struct {
	LocalCollection map[int]TestCollection
	TestResults     map[TestKey]ResultTest //TestKey is the key for the map

}

func (calcIMpl *Calculator) Init() {
	calcIMpl.LocalCollection = make(map[int]TestCollection)
	calcIMpl.TestResults = make(map[TestKey]ResultTest)
}

func (calcIMpl *Calculator) BuildResults(testCollections []TestCollection) map[TestKey]ResultTest {
	//emptyArray := []ResultTest{}
	calcIMpl.LocalCollection = calcIMpl.getCollectionMap(testCollections)
	log.Debugf("Imported %d collections of %d", len(calcIMpl.LocalCollection), len(testCollections))
	calcIMpl.loopCollections()

	return calcIMpl.TestResults
}

func (calcIMpl *Calculator) getCollectionMap(collections []TestCollection) map[int]TestCollection {
	collectionMap := make(map[int]TestCollection)

	if len(collections) > 0 {
		for x := 0; x < len(collections); x++ {
			if collections[x].TestName != "" {
				collectionMap[x] = collections[x]
			} else {
				log.Debugf("Why is this collection empty?")
			}
		}
	}

	return collectionMap

}

func (calcIMpl *Calculator) loopCollections() {
	myTempCollectionMap := make(map[int]TestCollection)

	for id, myCollection := range calcIMpl.LocalCollection {
		//add the collections to local object for processing
		myTempCollectionMap[myCollection.RunNumber] = myCollection
		delete(calcIMpl.LocalCollection, id)

		for id2, myCollection2 := range calcIMpl.LocalCollection {

			if myCollection.TestName == myCollection2.TestName &&
				myCollection.Dimension == myCollection2.Dimension &&
				myCollection.MySQLProducer == myCollection2.MySQLProducer &&
				myCollection.MySQLVersion == myCollection2.MySQLVersion &&
				myCollection.ActionType == myCollection2.ActionType &&
				myCollection.SelectPostWrites == myCollection2.SelectPostWrites {

				log.Debugf("id %d Collection: %s Run %d", id, myCollection.Name, myCollection.RunNumber)
				log.Debugf("if %d Collection: %s Run %d", id2, myCollection2.Name, myCollection2.RunNumber)
				myTempCollectionMap[myCollection2.RunNumber] = myCollection2
				delete(calcIMpl.LocalCollection, id2)

			}
		}
		log.Debugf(" Identified %d collections mathing same execution. Name %s dimension %s producer %s %s actiontype %d PostWrite %d ",
			len(myTempCollectionMap),
			myTempCollectionMap[0].Name,
			myTempCollectionMap[0].Dimension,
			myTempCollectionMap[0].MySQLProducer,
			myTempCollectionMap[0].MySQLVersion,
			myTempCollectionMap[0].ActionType,
			myTempCollectionMap[0].SelectPostWrites)

		calcIMpl.loopTests(myTempCollectionMap)

	}

}

/*
We loop all the tests create an array of them group by tests
calculate the results
*/
func (calcIMpl *Calculator) loopTests(colectionMap map[int]TestCollection) bool {
	// first we identify which collection has most tests (in case we have some of them failed and not processed.
	myHigherTestLen := 0
	var leadCollection TestCollection
	var OK bool

	// We identify the leading collection, remove it by the map and then start to use it to extract the tests from all
	for _, col := range colectionMap {
		if len(col.Tests) > myHigherTestLen {
			myHigherTestLen = len(col.Tests)
			leadCollection = col
		}
	}
	// last check to verify we are working on a valid collection
	if leadCollection.Name != "" && len(leadCollection.Tests) > 0 {
		delete(colectionMap, leadCollection.RunNumber)
		for id, myTest := range leadCollection.Tests {
			testAr := []Test{myTest}

			log.Debugf("Processing results for Name %s dimension %s producer %s %s actiontype %d PostWrite %d ",
				myTest.Name,
				myTest.Dimension,
				leadCollection.MySQLProducer,
				leadCollection.MySQLVersion,
				myTest.ActionType,
				leadCollection.SelectPostWrites)

			for _, col := range colectionMap {
				testAr = append(testAr, col.Tests[id])
			}

			OK = calcIMpl.calculateTestResultTest(testAr, leadCollection)

		}
		return OK
	}

	return false

}

// for each test we build an object and if multiple run we calculate the median, std and gerror
func (calcIMpl *Calculator) calculateTestResultTest(tests []Test, leadCollection TestCollection) bool {
	var newTestResult ResultTest
	var newTestResultKey TestKey

	newTestResultKey.MySQLProducer = leadCollection.MySQLProducer
	newTestResultKey.MySQLVersion = leadCollection.MySQLVersion
	newTestResultKey.TestName = tests[0].Name
	newTestResultKey.TestCollectionName = leadCollection.TestName
	newTestResultKey.Dimension = tests[0].Dimension
	newTestResultKey.ActionType = tests[0].ActionType
	newTestResultKey.SelectPreWrites = leadCollection.SelectPostWrites
	newTestResult.Key = newTestResultKey

	if len(tests) > 1 {
		newTestResult.Executions = len(tests)
		_, newTestResult.Labels = calcIMpl.transformLablesForMultipleExecutions(tests)
		newTestResult.STD, newTestResult.Gerror = calcIMpl.getLabelSTDGerror(newTestResult.Labels)
	} else {
		newTestResult.STD = 0
		newTestResult.Gerror = 0
		newTestResult.Executions = 1
		_, newTestResult.Labels = calcIMpl.transformLablesForSingleExecution(tests[0])
	}

	calcIMpl.TestResults[newTestResultKey] = newTestResult

	return true
}

// Before processing we transform the dataset from rows into column to be able to calculate the median, std and gerror [Multi run case]
func (calcIMpl *Calculator) transformLablesForMultipleExecutions(test []Test) (bool, map[string][]ResultValue) {
	labels := make(map[string][]ResultValue)
	for _, label := range test[0].Labels {
		resultValueAr := []ResultValue{}
		log.Debugf("processing label %s", label)

		//we need to loop all the threads and get the values for the label
		tempValuesAr := []float64{}

		for _, thread := range test[0].ThreadExec {
			threadI := thread.Thread
			for _, th := range test {
				for thLabel, thResult := range th.ThreadExec[threadI].Result {
					if thLabel == label {
						log.Debugf("Processing main: %s current: %s  Execution: %d Thread: %d result %.4f", label, thLabel, th.RunNumber, threadI, thResult)
						tempValuesAr = append(tempValuesAr, thResult)
					}

				}
			}

			//calculate final value, std and gerror
			resultValueAr = append(resultValueAr, evaluateMultipleExecutionsValues(tempValuesAr, label, threadI))
			tempValuesAr = []float64{}
			log.Debugf("")
		}

		sort.SliceStable(resultValueAr, func(i, j int) bool {
			return resultValueAr[i].ThreadNumber < resultValueAr[j].ThreadNumber
		})

		labels[label] = resultValueAr

	}
	return true, labels
}

// here we do the sdt calculation and gerror
func evaluateMultipleExecutionsValues(arValues []float64, label string, threadId int) ResultValue {

	medianValue, _ := stats.Median(arValues)
	medianValue, _ = stats.Round(medianValue, 2)
	stdValue, _ := stats.StandardDeviationSample(arValues)
	stdValue, _ = stats.Round(stdValue, 2)
	errorV := stdValue / medianValue * 100
	log.Debugf("Thread %d Label %s Median %.4f Std %.4f  Distance(error) %.4f", threadId, label, medianValue, stdValue, errorV)

	return ResultValue{threadId, label, medianValue, stdValue, errorV}
}

// Before processing we transform the dataset from rows into column to be able to calculate the median, std and gerror [Single run case]
func (calcIMpl *Calculator) transformLablesForSingleExecution(test Test) (bool, map[string][]ResultValue) {
	labels := make(map[string][]ResultValue)

	for _, label := range test.Labels {
		resultValueAr := []ResultValue{}
		log.Debugf("processing label %s", label)

		for thID, th := range test.ThreadExec {

			for thLabel, thResult := range th.Result {
				if thLabel == label {
					resultValueAr = append(resultValueAr, ResultValue{thID, label, thResult, 0, 0})
				}

			}

		}
		sort.SliceStable(resultValueAr, func(i, j int) bool {
			return resultValueAr[i].ThreadNumber < resultValueAr[j].ThreadNumber
		})
		labels[label] = resultValueAr

	}

	return true, labels
}

func (calcIMpl *Calculator) GroupByProducers() []Producer {
	producersAr := []Producer{}
	for key, _ := range calcIMpl.TestResults {
		present := false
		for _, producer := range producersAr {
			if producer.MySQLProducer == key.MySQLProducer && producer.MySQLVersion == key.MySQLVersion {
				present = true
			}
		}
		if !present {
			newProducer := Producer{key.MySQLProducer, key.MySQLVersion, make(map[TestKey]ResultTest), []TestType{}, 0.0, 0.0}
			log.Debugf("Adding producer %v", newProducer)
			producersAr = append(producersAr, newProducer)
		}

	}
	producersAr = calcIMpl.assignTestsResultsToProducers(producersAr)
	producersAr = calcIMpl.calculateProducerSTDGerror(producersAr)
	return producersAr

}
func (calcIMpl *Calculator) assignTestsResultsToProducers(producersAr []Producer) []Producer {

	for idx, producer := range producersAr {

		for key, value := range calcIMpl.TestResults {
			if producer.MySQLProducer == key.MySQLProducer && producer.MySQLVersion == key.MySQLVersion {
				producer.TestsResults[key] = value
				present := false
				for _, testType := range producer.TestsTypes {
					if testType.Name == key.TestName && testType.Dimension == key.Dimension && testType.SelectPreWrites == key.SelectPreWrites {
						present = true
					}
				}
				if !present {
					newTestType := TestType{key.TestName, key.Dimension, key.SelectPreWrites}
					producer.TestsTypes = append(producer.TestsTypes, newTestType)
				}
			}
		}
		producersAr[idx] = producer
	}
	return producersAr
}

func (calcIMpl *Calculator) getLabelSTDGerror(labels map[string][]ResultValue) (float64, float64) {
	valuesSTDAr := []float64{0}
	valuesGerrAr := []float64{0}
	for _, resultValueAr := range labels {

		for _, resultValue := range resultValueAr {
			if resultValue.Value > 0 && !math.IsNaN(resultValue.STD) {
				valuesSTDAr = append(valuesSTDAr, resultValue.STD)
				valuesGerrAr = append(valuesGerrAr, resultValue.Lerror)
			}
		}
		stdValue := 0.00
		if len(valuesSTDAr) > 1 {
			stdValue, _ = stats.StandardDeviationSample(valuesSTDAr)
		}

		stdValue, _ = stats.Round(stdValue, 2)
		gerrValue := global.Average(valuesGerrAr)
		return stdValue, gerrValue
	}
	return 0, 0
}

func (calcIMpl *Calculator) calculateProducerSTDGerror(ar []Producer) []Producer {
	valuesSTDAr := []float64{0}
	valuesGerrAr := []float64{0}
	for idx, producer := range ar {
		for _, testResult := range producer.TestsResults {
			valuesSTDAr = append(valuesSTDAr, testResult.STD)
			valuesGerrAr = append(valuesGerrAr, testResult.Gerror)
		}
		stdValue := 0.0
		if len(valuesSTDAr) > 1 {
			stdValue, _ = stats.StandardDeviationSample(valuesSTDAr)
			stdValue, _ = stats.Round(stdValue, 2)
		}
		gerrValue := global.Average(valuesGerrAr)
		producer.STD = stdValue
		producer.Gerror = gerrValue
		ar[idx] = producer
	}
	return ar

}
