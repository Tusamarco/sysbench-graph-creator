package dataObjects

import (
	log "github.com/sirupsen/logrus"
	"sort"
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

func (calcIMpl *Calculator) BuildResults(testCollections []TestCollection) []ResultTest {
	emptyArray := []ResultTest{}
	calcIMpl.LocalCollection = calcIMpl.getCollectionMap(testCollections)
	log.Debugf("Imported %d collections of %d", len(calcIMpl.LocalCollection), len(testCollections))
	calcIMpl.loopCollections()

	return emptyArray
}

func (calcIMpl *Calculator) getCollectionMap(collections []TestCollection) map[int]TestCollection {
	collectionMap := make(map[int]TestCollection)

	if len(collections) > 0 {
		for x := 0; x < len(collections); x++ {
			if collections[x].TestName != "" {
				collectionMap[x] = collections[x]
			} else {
				log.Debugf("Whay is this collection empty?")
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
			if len(testAr) < 1 {
				log.Warnf("Yeah that is weird no test has been found %s", id)
				return false
			}
			OK = calcIMpl.calculateTestResultForSingleTest(testAr, leadCollection)

		}
		return OK
	}

	return false

}

func (calcIMpl *Calculator) calculateTestResultForSingleTest(tests []Test, leadCollection TestCollection) bool {
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

	if len(tests) > 10 {

	} else {
		newTestResult.STD = 0
		newTestResult.Gerror = 0
		_, newTestResult.Labels = calcIMpl.transformLablesForSingleTest(tests[0])
	}

	calcIMpl.TestResults[newTestResultKey] = newTestResult

	//i := 0
	//
	////we use the first entry as the leading test, but we need to check if there are more than 1 otherwise it will got an error
	//if len(tests) > 1 {
	//	i = 1
	//}
	//
	//for i = i; i < len(tests); i++ {
	//
	//}

	return true
}

func (calcIMpl *Calculator) transformLablesForSingleTest(test Test) (bool, map[string][]ResultValue) {
	labels := make(map[string][]ResultValue)

	for _, label := range test.Labels {
		resultValueAr := []ResultValue{}
		log.Debugf("processing label %s", label)

		for thID, th := range test.ThreadExec {

			for thLabel, thResult := range th.Result {
				if thLabel == label {
					resultValueAr = append(resultValueAr, ResultValue{thID, thResult, 0, 0})
				}

			}

		}
		sort.SliceStable(resultValueAr, func(i, j int) bool {
			return resultValueAr[i].threadNumber < resultValueAr[j].threadNumber
		})
		labels[label] = resultValueAr

	}

	return false, labels
}
