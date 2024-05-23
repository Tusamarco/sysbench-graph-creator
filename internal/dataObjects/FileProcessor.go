package dataObjects

import (
	"bufio"
	"fmt"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	global "sysbench-graph-creator/internal/global"
	"time"
)

type FileProcessor struct {
	arDataFile       []DataFile
	arPathFiles      []string
	testCollectionAr []TestCollection
	MyScanner        *bufio.Scanner
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

/*
Function to pupulate the array in the fileparser with top objects TestCollection. Test collaction is the abstraction of the whole set of data existing in the file
*/
func (fileProc *FileProcessor) GetTestCollectionArray() ([]TestCollection, error) {

	filesLength := len(fileProc.arPathFiles)

	for i, path := range fileProc.arPathFiles {
		log.Infof("Processing [%d] files. Analyzing file [%d/%d] path: %s", filesLength, i+1, filesLength, path)

		testCollection, OK := fileProc.getTestCollectionData(path)
		if !OK {
			log.Errorf("Parsing Test Collection failed %s", path)
		}

		fileProc.testCollectionAr = append(fileProc.testCollectionAr, testCollection)
	}

	//we identify the last test executed and its end time
	fileProc.identifyEndTime()
	// we identify and set the post writes collections
	fileProc.identifyPerPostWrite()

	return fileProc.testCollectionAr, nil

}

func (fileProc *FileProcessor) identifyEndTime() {
	for i := 0; i < len(fileProc.testCollectionAr); i++ {

		mylen := len(fileProc.testCollectionAr[i].Tests)
		if mylen > 0 {
			mylen--
			//icounter := 0
			var maxDate time.Time
			for _, item := range fileProc.testCollectionAr[i].Tests {

				if !maxDate.After(item.DateEnd) {
					maxDate = item.DateEnd
				}
			}
			//log.Debugf("Processing %s %s", myrange, item.Name)
			fileProc.testCollectionAr[i].DateEnd = maxDate
			difference := fileProc.testCollectionAr[i].DateEnd.Sub(fileProc.testCollectionAr[i].DateStart)
			fileProc.testCollectionAr[i].ExecutionTime = int64(difference.Minutes())

		}

		//if mylen > 0 {
		//	myTest := fileProc.testCollectionAr[i].Tests[mylen-1]
		//}
	}
}

/*
Browse all the collections and identify who is the post write one comparing them by:
- name
- run
- date
- action type
- dimension
*/
func (fileProc *FileProcessor) identifyPerPostWrite() {

	for i := 0; i < len(fileProc.testCollectionAr); i++ {

		myTestCollection := fileProc.testCollectionAr[i]
		for y := 0; y < len(fileProc.testCollectionAr); y++ {
			if myTestCollection.Name == fileProc.testCollectionAr[y].Name &&
				myTestCollection.Dimension == fileProc.testCollectionAr[y].Dimension &&
				myTestCollection.ActionType == fileProc.testCollectionAr[y].ActionType &&
				myTestCollection.RunNumber == fileProc.testCollectionAr[y].RunNumber &&
				myTestCollection.MySQLProducer == fileProc.testCollectionAr[y].MySQLProducer &&
				myTestCollection.MySQLVersion == fileProc.testCollectionAr[y].MySQLVersion {
				if myTestCollection.DateStart.After(fileProc.testCollectionAr[y].DateStart) {
					myTestCollection.SelectPostWrites = POSTWRITE
					fileProc.testCollectionAr[i] = myTestCollection
					log.Debugf("Assign Post write to true to collection %s dimension %s run %d ", myTestCollection.Name, myTestCollection.Dimension, myTestCollection.RunNumber)
					break
				}
			}

		}
	}

}

/*
 */
func (fileProc *FileProcessor) getTestCollectionData(path string) (TestCollection, bool) {
	testCollection := new(TestCollection)
	testCollection.Tests = make(map[string]Test)
	//var err error
	//Open file and loop in to lines for meta
	file, err := os.Open(path)

	nameFiltered := global.ReplaceString(filepath.Base(file.Name()), "_runNumber[0-9]_", "_")
	nameFiltered = global.ReplaceString(nameFiltered, "_\\d{4}-\\d{2}-\\d{1,2}_\\d{2}_\\d{2}", "")

	testCollection.Name = nameFiltered
	testCollection.FileName = filepath.Base(file.Name())

	if err != nil {
		log.Error(err)
	}

	defer file.Close()
	fileProc.MyScanner = bufio.NewScanner(file)

	numberOfLines, err := global.LineCount(path)
	if err != nil {
		log.Error(err)
	}
	testCollection.PBarr = progressbar.Default(int64(numberOfLines))
	defer testCollection.PBarr.Set(numberOfLines)

	metaTop := true
	// first we retrive meta information about the tests
	for fileProc.MyScanner.Scan() {
		line := fileProc.MyScanner.Text()
		testCollection.PBarr.Add(1)

		if len(line) > 1 {
			// load the mata for the collection (whole file run)
			if strings.Contains(line, "META") && metaTop {
				if !testCollection.getTestCollectionMeta(line, path) {
					log.Error(fmt.Errorf("Cannot load Meta information for collection"))
				}
				metaTop = false
			}
			//load meta and data for each specific test and add the tests
			if strings.Contains(line, "SUBTEST:") && !metaTop {
				newTest, OK := testCollection.getTestMeta(line, path, fileProc)
				if !OK {
					log.Error(fmt.Errorf("Cannot load Meta information for test"))
				} else {
					testCollection.Tests[newTest.Name] = newTest
				}

			}

		}
	}
	return *testCollection, true
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
		//log.Debugf("Meta argument parsing %s", values)
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
			case "mysqlproducer":
				tescImpl.MySQLProducer = tescImpl.getProducer(values[1])
			case "mysqlversion":
				tescImpl.MySQLVersion = values[1]

			}
			if err != nil {
				log.Error(err)
				return false
			}
		}

	}
	return true
}

func (tescImpl *TestCollection) getProducer(name string) string {
	if len(name) < 1 {
		return name
	}

	index := len(name) - 1
	index1 := strings.Index(name, " -")
	index2 := strings.Index(name, " (")

	if index2 < index1 && index2 > 1 || index1 < 0 {
		index = index2
	} else if index1 < 0 && index2 < 0 {
		index = index
	} else {
		index = index1
	}

	name = name[:index]
	return name
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
func (tescImpl *TestCollection) getTestMeta(line string, path string, fileProc *FileProcessor) (Test, bool) {
	var err error
	var newTest Test
	newTest.init()

	value := strings.Split(strings.ReplaceAll(line, " ", ""), ":")

	newTest.Name = value[1]
	fileProc.MyScanner.Scan()
	line = fileProc.MyScanner.Text()
	tescImpl.PBarr.Add(1)

	if strings.Contains(line, "BLOCK: [START]") {
		newTest.DateStart, err = global.ParsetimeLocal(line, "")
		if err != nil {
			log.Error(err)
		}

		re := regexp.MustCompile(`^.*(\(filter.*\))`)
		match := re.FindStringSubmatch(line)
		if len(match) > 0 {
			if match[1] != "" {
				value = strings.Split(strings.ReplaceAll(strings.ReplaceAll(match[1], "(", ""), ")", ""), ":")
				if len(value) > 0 {
					newTest.Filter = strings.TrimSpace(value[1])
				}
			}
		}

	}

	fileProc.MyScanner.Scan()
	line = fileProc.MyScanner.Text()
	tescImpl.PBarr.Add(1)

	if strings.Contains(line, "META:") {
		line = strings.ReplaceAll(line, " ", "")
		metaTag := strings.Split(line[5:], ";")
		length := len(metaTag)
		var err error

		//META: testIdentifyer=PS8042_iron_ssd2;dimension=large;actionType=select;runNumber=1;execCommand=run;subtest=select_run_inlist;execDate=2024-02-02_12_12_27;engine=innodb
		for i := 0; i < length; i++ {
			values := strings.Split(metaTag[i], "=")
			//log.Debugf("Meta argument parsing %s", values)
			if len(values) > 0 {
				trimmed := strings.Trim(values[0], " ")
				switch trimmed {
				case "dimension":
					newTest.Dimension = values[1]
				case "actionType":
					newTest.ActionType, err = getCodeAction(values[1])
				case "runNumber":
					newTest.RunNumber, err = strconv.Atoi(values[1])
				}
			}
			if err != nil {
				log.Error(err)
				var errTest Test
				return errTest, false
			}
		}

	}

	for fileProc.MyScanner.Scan() {
		line = fileProc.MyScanner.Text()
		tescImpl.PBarr.Add(1)
		runExecuteInFull := false
		//lastRunningThreadNumber := 0
		if strings.Contains(line, "THREADS=") {
			//iThtread, _ := strconv.Atoi(line[8:])
			newRun, OK := newTest.getAllRuns(fileProc, *tescImpl)
			line = fileProc.MyScanner.Text()

			if !OK {
				log.Error("Error while processing runs ")
			}
			if !reflect.ValueOf(newRun).IsZero() {
				newTest.ThreadExec[newRun.Thread] = newRun
				newTest.Threads = append(newTest.Threads, newRun.Thread)
				//lastRunningThreadNumber = newRun.Thread
			}
		}
		if strings.Contains(line, "BLOCK: [END]") {
			newTest.DateEnd, err = global.ParsetimeLocal(line, "")
			if err != nil {
				log.Error(err)
			}
			return newTest, true
		}
		if strings.Contains(line, "SUBTEST:") {
			// todo if we reach here then the test had some issue and we need to manage the return in some way
			if !runExecuteInFull {

			}

		}
	}

	return newTest, true
}

func (test *Test) init() {
	test.ThreadExec = make(map[int]Execution)
	test.Labels = []string{}
	test.Threads = []int{}

}

func (test *Test) getAllRuns(fileProc *FileProcessor, tescImpl TestCollection) (Execution, bool) {
	var newRun Execution
	var errExecution Execution
	newRun.Result = make(map[string]float64)
	var err error

	line := fileProc.MyScanner.Text()
	if strings.Contains(line, "THREADS=") {
		iThtread, _ := strconv.Atoi(line[8:])
		newRun.Thread = iThtread
	}

	for fileProc.MyScanner.Scan() {
		tescImpl.PBarr.Add(1)
		//time.Sleep(1 * time.Microsecond / 10)

		line := fileProc.MyScanner.Text()
		//log.Debugf(line)

		if strings.Contains(line, "RUNNING ") && strings.Contains(line, "[START]") {
			newRun.DateStart, err = global.ParsetimeLocal(line, "")
			if err != nil {
				log.Error(err)
				return errExecution, false
			}
		}
		if strings.Contains(line, "RUNNING ") && strings.Contains(line, "[END]") {
			newRun.DateEnd, err = global.ParsetimeLocal(line, "")
			if err != nil {
				log.Error(err)
				return errExecution, false
			}

			return newRun, true
		}

		/*
			TODO We should never reach this condition, if we do something during the tests failed and we may have a corrupted log file
		*/
		if strings.Contains(line, "THREADS=") || strings.Contains(line, "SUBTEST:") {
			log.Error("It seems a test failed while executing. Results for this run are not correct")
			log.Errorf("Test Name: %s, Threads %d ", tescImpl.Name, newRun.Thread)

			return newRun, false
		}

		if strings.Contains(line, "Executing:") {
			newRun.Command = line[11:]
		}

		if strings.Contains(line, "TEST SUMMARY:") {
			fileProc.MyScanner.Scan()
			line = fileProc.MyScanner.Text()
			test.Labels = strings.Split(line, ",")
			fileProc.MyScanner.Scan()
			line = fileProc.MyScanner.Text()
			arResults := strings.Split(line, ",")
			ilen := len(arResults)
			if ilen == len(test.Labels) {
				for i := 0; i < ilen; i++ {
					newRun.Result[test.Labels[i]], err = strconv.ParseFloat(arResults[i], 64)
					if err != nil {
						log.Error(err)
						return errExecution, false
					}
				}
				newRun.Processed = true
			} else {
				log.Errorf("Error in assign results. Lenght of Labels and data doesn't match. Labels %d; Results %d. ", len(test.Labels), ilen)
				return errExecution, false
			}
		}

	}

	return newRun, true
}
