package counter

import (
	"AccelByteTakeHome/utils"
	"time"
)

type Counter interface {
	Count(path string, startTime, endTime time.Time, cMap chan map[string]int)
}

type JsonCounter struct {
}

func (J JsonCounter) Count(path string, startTime, endTime time.Time, cMap chan map[string]int) {
	utils.JsonFile(path, startTime, endTime, cMap)
}

type CsvCounter struct {
}

func (C CsvCounter) Count(path string, startTime, endTime time.Time, cMap chan map[string]int) {
	utils.CsvFile(path, startTime, endTime, cMap)
}

func NewCounter(fileType string, path string, startTime, endTime time.Time, cMap chan map[string]int) Counter {
	switch fileType {
	case "json":
		return JsonCounter{}
	case "csv":
		return CsvCounter{}
	}
	return nil
}
