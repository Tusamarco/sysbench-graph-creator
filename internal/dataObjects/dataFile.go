package dataObjects

import "time"

type DataFile struct {
	FullPath string
	TestName string
	Producer string
	TestType string
	RunDate  time.Time
}
