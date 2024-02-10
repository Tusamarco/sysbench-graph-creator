package dataObjects

import (
	"bufio"
	"fmt"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	global "sysbench-graph-creator/internal/global"
)

type FileProcessor struct {
	arDataFile       []DataFile
	arPathFiles      []string
	testCollectionAr []TestCollection
}

// This function will recursively look for summary files and collect them into an array of strings
func (fileProc *FileProcessor) GetFileList(path string) error {
	//var arDataFile []DataFile

	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() &&
				(strings.Contains(info.Name(), ".csv") || strings.Contains(path, ".txt")) &&
				!strings.Contains(info.Name(), "_warmup_") {

				fileProc.arPathFiles = append(fileProc.arPathFiles, path)
			}

			return nil
		})
	if err != nil {
		log.Error(err)
	}

	return err
}

func (fileProc *FileProcessor) GetTestCollectionArray() ([]TestCollection, error) {

	filesLength := len(fileProc.arPathFiles)

	for i, path := range fileProc.arPathFiles {
		log.Infof("Processing [%d] files. Analyzing file [%d/%d] path: %s", filesLength, i+1, filesLength, path)

		file, err := os.Open(path)
		if err != nil {
			log.Error(err)
			return fileProc.testCollectionAr, nil
		}
		defer file.Close()

		//Open file and loop in to lines for meta
		scanner := bufio.NewScanner(file)
		testCollection, OK := fileProc.getTestCollectionData(*scanner, path)
		if !OK {
			log.Errorf("Parsing Test Collection failed %s", path)
		}

		fileProc.testCollectionAr = append(fileProc.testCollectionAr, testCollection)
	}

	return fileProc.testCollectionAr, nil

}

func (fileProc *FileProcessor) getTestCollectionData(scanner bufio.Scanner, path string) (TestCollection, bool) {
	testCollection := new(TestCollection)
	//var err error

	numberOfLines, err := global.LineCount(path)
	if err != nil {
		log.Error(err)
	}
	barLine := progressbar.Default(int64(numberOfLines))
	metaTop := true
	// first we retrive meta information about the tests
	for scanner.Scan() {
		line := scanner.Text()
		barLine.Add(1)

		if len(line) > 1 {
			// load the mata for the collection (whole file run)
			if strings.Contains(line, "META") && metaTop {
				if !testCollection.getTestCollectionMeta(line, path) {
					log.Error(fmt.Errorf("Cannot load Meta information for collection"))
				}
				metaTop = true
			}
			//load meta and data for each specific test and add the tests
			if strings.Contains(line, "META") && !metaTop {
				if !testCollection.getTestMeta(line, path, scanner, barLine) {
					log.Error(fmt.Errorf("Cannot load Meta information for test"))
				}

			}

		}
	}
	return *testCollection, false
}

func (tescImpl *TestCollection) getTestCollectionMeta(meta string, path string) bool {

	meta = strings.TrimSpace(meta)
	metaTag := strings.Split(meta[5:], ";")
	length := len(metaTag)
	var err error

	/*
		Parse the meta information for the top test collection
		META: testIdentifyer=PS8042_iron_ssd2;dimension=large;actionType=select;runNumber=1;host=10.30.12.4;producer=sysbench;execDate=;engine=innodb
	*/
	for i := 0; i < length; i++ {
		values := strings.Split(metaTag[i], "=")
		log.Debugf("Meta argument parsing %s", values)
		if len(values) > 0 {
			trimmed := strings.Trim(values[0], " ")
			switch trimmed {
			case "testIdentifyer":
				tescImpl.TestName = values[1]
			case "dimension":
				tescImpl.Dimension = values[1]
			case "actionType":
				tescImpl.ActionType, err = getCodeAction(values[1])
			case "runNumber":
				tescImpl.RunNumber, _ = strconv.Atoi(values[1])
			case "host":
				tescImpl.HostDB = values[1]
			case "producer":
				tescImpl.Producer = values[1]
			case "execDate":
				tescImpl.DateStart, err = global.ParsetimeLocal(values[1], path)
			case "engine":
				tescImpl.Engine = values[1]

			}
			if err != nil {
				log.Error(err)
				return false
			}
		}

	}
	return true
}

/*
here we start to read the data for each test and return only when we have collect all the information in the summary
Here we also associate the single run thread run

	so we have 2 objects
	   testRunsCollection
	        |- runs

ie: for the test select inlist. We have the top object containing all the information related to the specific test:
META: testIdentifyer=PS8042_iron_ssd2;dimension=large;actionType=select;runNumber=1;execCommand=run;subtest=select_run_inlist;execDate=2024-02-02_12_12_27;engine=innodb

then it has an array of runs and each run is related to a number of threads bound to that run.
Each run will report information as :
TEST SUMMARY:
TotalTime,RunningThreads,totalEvents,Events/s,Tot Operations,operations/s,tot reads,reads/s,Tot writes,writes/s,oterOps/s,latencyPct95(μs) ,Tot errors,errors/s,Tot reconnects,reconnects/s,Latency(ms) min, Latency(ms) max, Latency(ms) avg, Latency(ms) sum
200,1,2642.00,13.21,2642.00,13.21,2642.00,13.21,0.00,0.00,0.00,137.35,0.00,0.00,0.00,0.00,0.04,0.22,0.08,200.00
======================================
RUNNING Test PS8042_iron_ssd2 sysbench select_run_inlist (filter: select) Thread=1 [END] 2024-02-02_12_15_47
======================================
*/
func (tescImpl *TestCollection) getTestMeta(meta string, path string, scanner bufio.Scanner, barrLine *progressbar.ProgressBar) bool {

	return true
}
