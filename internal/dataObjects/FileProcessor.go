package dataObjects

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

type FileProcessor struct {
	arDataFile     []DataFile
	arPathFiles    []string
	testCollection []TestCollection
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
		file, err := os.Open(path)
		if err != nil {
			log.Error(err)
			return fileProc.testCollection, nil
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

	}

	return fileProc.testCollection, nil

}
