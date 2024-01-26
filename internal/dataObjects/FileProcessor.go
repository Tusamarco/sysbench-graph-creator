package dataObjects

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func GetFileList(path string) (error, []DataFile) {
	var arDataFile []DataFile

	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() &&
				strings.Contains(info.Name(), ".csv") &&
				!strings.Contains(path, "/data/") {
				limiter := ""
				var myDataFile DataFile
				myDataFile.FullPath = path
				if strings.Contains(info.Name(), "_large_") {
					limiter = "_large_"
				} else {
					limiter = "_small_"
				}
				if strings.Contains(path, "sysbench") {
					myDataFile.Producer = "sysbench"
				}
				if strings.Contains(path, "tpcc") {
					myDataFile.Producer = "tpcc"
				}
				if strings.Contains(path, "dbt3") {
					myDataFile.Producer = "tpcc"
				}

				re := regexp.MustCompile(`(\d{4}-\d{2}-\d{1,2}_\d{2}_\d{2})`)
				match := re.FindStringSubmatch(path)

				if match[0] != "" {
					//strDate := match[1]
					myDataFile.RunDate, err = time.Parse("2006-01-02_04_05", match[0])
				}
				//Global.ReturnDateFromString(match[1], "0000-12-23_00_00")
				if err != nil {
					log.Warnf("Parsing error ", err)
					//return err
				}
				myDataFile.TestName = info.Name()[0:strings.Index(info.Name(), limiter)]
				//myDataFile
				fmt.Println(path, info.Size())
				arDataFile = append(arDataFile, myDataFile)
			}

			return nil
		})
	if err != nil {
		log.Println(err)
	}

	return err, arDataFile
}
