package dataObjects

import (
	log "github.com/sirupsen/logrus"
	"os"
	Global "sysbench-graph-creator/internal/global"
)

func GetFileList(config Global.Configuration) (bool, error) {

	f, err := os.Open(config.Parser.SourceDataPath)
	if err != nil {
		log.Error(err)
		return false, err
	}
	files, err := f.Readdir(0)
	if err != nil {
		log.Error(err)
		return false, err
	}

	for _, file := range files {
		log.Info(file.Name(), file.IsDir())
	}
	return true, nil
}
