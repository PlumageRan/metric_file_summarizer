package utils

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func CsvFile(path string, startTime, endTime time.Time, cMap chan map[string]int) {

	resMap := make(map[string]int)

	fs, err := os.Open(path)
	if err != nil {
		log.Fatalf("can not open the file, err is %+v", err)
	}
	defer fs.Close()
	r := csv.NewReader(fs) //read line by line
	for {
		row, err := r.Read()
		if err != nil && err != io.EOF {
			log.Fatalf("can not read, err is %+v", err)
		}
		if err == io.EOF {
			break
		}
		checkTime, checkErr := time.Parse(time.RFC3339, row[0])
		value, convertErr := strconv.Atoi(row[2])
		if checkErr == nil && convertErr == nil && InTimeSpan(startTime, endTime, checkTime) {
			resMap[row[1]] += value
		}
	}
	cMap <- resMap
}
