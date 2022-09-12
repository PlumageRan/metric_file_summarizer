package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v3"
)

const (
	THOUSAND = 1000
)

func GetFinalStringFromTwoStringWithSameMeaning(a, b string) string {
	if a == b || b == "" {
		return a
	}
	return b
}

func FileNameToTime(fileName string) time.Time {
	fileNameWithoutFormat := strings.Split(fileName, ".")[0]
	// pattern := "(0[1-9]|[12][0-9]|3[01])-[?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec].(?:json|csv)"
	// ok, err := regexp.Match(pattern, []byte(fileNameWithoutFormat))
	// fmt.Println(ok)
	// if !ok {
	// 	fmt.Println(err)
	// 	return time.Unix(0, 0)
	// }
	dayAndMonthSlice := strings.Split(fileNameWithoutFormat, "-")
	if len(dayAndMonthSlice) < 2 {
		return time.Unix(0, 0)
	}
	day, month := dayAndMonthSlice[0], dayAndMonthSlice[1]
	toParse := fmt.Sprintf("%s %s %d 00:00 UTC", day, month, time.Now().Year()%THOUSAND)
	parsedTime, err := time.Parse(time.RFC822, toParse)
	if err != nil {
		fmt.Println(err)
		return time.Unix(0, 0)
	}
	return parsedTime
}

func InTimeSpan(start, end, check time.Time) bool {
	return check.Equal(start) || (check.After(start) && check.Before(end))
}

func SameDay(time1, time2 time.Time) bool {
	return time1.Year() == time2.Year() && time1.YearDay() == time2.YearDay()
}

func ValidFile(fileName string, startTime, endTime time.Time) bool {
	fileTime := FileNameToTime(fileName)
	return InTimeSpan(startTime, endTime, fileTime) || SameDay(startTime, fileTime)
}

type OutMetric struct {
	LevelName string `json:"level_name"`
	Value     int    `json:"value"`
}

func GenerateReport(resMap map[string]int, outType, outName, outPath string) {
	filePath := fmt.Sprintf("%s/%s.%s", outPath, outName, outType)
	var data []OutMetric
	for key, value := range resMap {
		data = append(data, OutMetric{
			LevelName: key,
			Value:     value,
		})
	}
	var out []byte
	var err error
	if outType == "json" {
		out, _ = json.MarshalIndent(data, "", " ")
	}
	if outType == "yaml" {
		out, err = yaml.Marshal(data)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = ioutil.WriteFile(filePath, out, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Create report: SUCCESS.")
}
