package tests

import (
	"AccelByteTakeHome/utils"
	"fmt"
	"testing"
	"time"
)

func TestFileNameToTime(t *testing.T) {
	months := []string{"jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sep", "oct", "nov", "dec"}
	for _, v := range months {
		for i := 1; i <= 31; i++ {
			fileName := ""
			if i <= 9 {
				fileName = fmt.Sprintf("0%d-%s.json", i, v)
			} else {
				fileName = fmt.Sprintf("%d-%s.json", i, v)
			}
			parsedTime := utils.FileNameToTime(fileName)
			t.Log(parsedTime)
		}
	}
}

func TestJsonFile(t *testing.T) {
	utils.JsonFile("/Users/haoranzhou/Documents/AccelByteTakeHome",
		time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2022, time.January, 1, 23, 59, 59, 0, time.UTC),
		make(chan map[string]int))
}
