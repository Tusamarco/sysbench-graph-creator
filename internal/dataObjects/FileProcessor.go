package dataObjects

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
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
				strings.Contains(info.Name(), ".csv") &&
				!strings.Contains(path, "/data/") {

				fileProc.arPathFiles = append(fileProc.arPathFiles, path)

				//re := regexp.MustCompile(`(\d{4}-\d{2}-\d{1,2}_\d{2}_\d{2})`)
				//match := re.FindStringSubmatch(path)
				//
				//if match[0] != "" {
				//	//strDate := match[1]
				//	myDataFile.RunDate, err = time.Parse("2006-01-02_04_05", match[0])
				//}
				//Global.ReturnDateFromString(match[1], "0000-12-23_00_00")
			}

			return nil
		})
	if err != nil {
		log.Error(err)
	}

	return err
}

func (fileProc *FileProcessor) GetTestCollectionArray() ([]TestCollection, error) {

	for i, path := range fileProc.arPathFiles {
		log.Debugf("Processing file %d: %s", i+1, path)
		// create and fill test collection
		filename := path[strings.LastIndex(path, "/")+1:]
		testCollection := fileProc.getTestCollectionMeta(filename, path)

		file, err := os.Open(path)
		if err != nil {
			log.Error(err)
			return fileProc.testCollectionAr, nil
		}
		defer file.Close()

		//Open file and loop in to lines for meta
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) > 1 {
				if line[0:4] == "META" {
					log.Debugf("META :%s", line)
				}

			}

		}

		fileProc.testCollectionAr = append(fileProc.testCollectionAr, testCollection)
	}

	return fileProc.testCollectionAr, nil

}

func (fileProc *FileProcessor) getTestCollectionMeta(filename string, path string) TestCollection {
	testCollection := new(TestCollection)
	var limiter string

	if strings.Contains(filename, "_large_") {
		limiter = "_large_"
		testCollection.Dimension = "large"
	} else {
		limiter = "_small_"
		testCollection.Dimension = "small"
	}
	if strings.Contains(filename, "_write_") {
		testCollection.ActionType = WRITE
	} else {
		testCollection.ActionType = READ
	}
	if strings.Contains(path, "sysbench") {
		testCollection.Producer = "sysbench"
	}
	if strings.Contains(path, "tpcc") {
		testCollection.Producer = "tpcc"
	}
	if strings.Contains(path, "dbt3") {
		testCollection.Producer = "dbt3"
	}

	re := regexp.MustCompile(`(\d{4}-\d{2}-\d{1,2}_\d{2}_\d{2})`)
	match := re.FindStringSubmatch(filename)

	if match[0] != "" {
		//strDate := match[1]
		myTime, err := time.Parse("2006-01-02_04_05", match[0])
		if err != nil {
			log.Warnf("Parsing error ", err)
			//return err
		}

		testCollection.DateStart = myTime
	}
	//Global.ReturnDateFromString(match[1], "0000-12-23_00_00")
	testCollection.TestName = filename[0:strings.Index(filename, limiter)]
	return *testCollection
}
