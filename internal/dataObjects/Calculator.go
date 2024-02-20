package dataObjects

import log "github.com/sirupsen/logrus"

//gonum.org/v1/gonum/stat
//github.com/montanaflynn/stats

type Calculator struct {
	LocalCollection map[int]TestCollection
	TestResults     map[TestKey]ResultTest //TestKey is the key for the map

}

func (calcIMpl *Calculator) BuildResults(testCollections []TestCollection) []ResultTest {
	emptyArray := []ResultTest{}
	calcIMpl.LocalCollection = calcIMpl.getCollectionMap(testCollections)
	log.Debugf("Imported %d collections", len(calcIMpl.LocalCollection))
	calcIMpl.loopCollections()

	return emptyArray
}

func (calcIMpl *Calculator) getCollectionMap(collections []TestCollection) map[int]TestCollection {
	collectionMap := make(map[int]TestCollection)

	if len(collections) > 0 {
		for x := 0; x < len(collections); x++ {
			if collections[x].TestName != "" {
				collectionMap[x] = collections[x]
			}
		}
	}

	return collectionMap

}

func (calcIMpl *Calculator) loopCollections() {
	for id, myCollection := range calcIMpl.LocalCollection {

		delete(calcIMpl.LocalCollection, id)
		for id2, myCollection2 := range calcIMpl.LocalCollection {

			if myCollection.TestName == myCollection2.TestName &&
				myCollection.Dimension == myCollection2.Dimension &&
				myCollection.Producer == myCollection2.Producer &&
				myCollection.ActionType == myCollection2.ActionType &&
				myCollection.SelectPreWrites == myCollection2.SelectPreWrites {

				log.Debugf("id %d Collection: %v", myCollection, id)
				log.Debugf("if %d Collection: %v", myCollection2, id2)
			}
		}

	}

}
