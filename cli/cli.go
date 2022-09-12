package cli

import (
	"AccelByteTakeHome/counter"
	"AccelByteTakeHome/utils"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type CommandLine struct {
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" generateSummary ")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 1 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) GenerateSummary(path, fileType string, startTime, endTime time.Time, outputFileType, outputFileName, outPath string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}
	wg := sync.WaitGroup{}
	resMap := make(map[string]int)
	for _, file := range files {
		tempPath := fmt.Sprintf("%s/%s", path, file.Name())
		extension := filepath.Ext(tempPath)
		cMap := make(chan map[string]int, 10)
		if extension == ".json" && fileType == "json" {
			if utils.ValidFile(file.Name(), startTime, endTime) {
				wg.Add(1)
				go func() {
					// utils.JsonFile(tempPath, startTime, endTime, cMap)
					jsonCounter := counter.NewCounter(fileType, tempPath, startTime, endTime, cMap)
					jsonCounter.Count(tempPath, startTime, endTime, cMap)
					close(cMap)
				}()
				go func() {
					defer wg.Done()
					for c := range cMap {
						for key, value := range c {
							resMap[key] += value
						}
					}
				}()
			}
		}
		if extension == ".csv" && fileType == "csv" {
			if utils.ValidFile(file.Name(), startTime, endTime) {
				wg.Add(1)
				go func() {
					// utils.CsvFile(tempPath, startTime, endTime, cMap)
					csvCounter := counter.NewCounter(fileType, tempPath, startTime, endTime, cMap)
					csvCounter.Count(tempPath, startTime, endTime, cMap)
					close(cMap)
				}()
				go func() {
					defer wg.Done()
					for c := range cMap {
						for key, value := range c {
							resMap[key] += value
						}
					}
				}()
			}
		}
	}
	wg.Wait()
	fmt.Println(fmt.Sprintf("start time: %s", startTime))
	fmt.Println(fmt.Sprintf("end time: %s", endTime))
	fmt.Println(resMap)
	utils.GenerateReport(resMap, outputFileType, outputFileName, outPath)
}

func (cli *CommandLine) Run(path string) {

	cli.validateArgs()

	generateSummaryCmd := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	dLong := generateSummaryCmd.String("directory", "", "the directory path containing the files")
	dShort := generateSummaryCmd.String("d", "", "the directory path containing the files")
	tLong := generateSummaryCmd.String("type", "", "the type of input files")
	tShort := generateSummaryCmd.String("t", "", "the type of input files")
	startTime := generateSummaryCmd.String("startTime", "", "this is the start time, inclusively")
	endTime := generateSummaryCmd.String("endTime", "", "this is the end time, exclusively")
	outputFileType := generateSummaryCmd.String("outputFileType", "json", "the output file type of the summary")
	outputFileName := generateSummaryCmd.String("outputFileName", "out", "the output file name")

	err := generateSummaryCmd.Parse(os.Args[2:])
	if err != nil {
		log.Panic(err)
	}

	if generateSummaryCmd.Parsed() {
		if *dLong == "" && *dShort == "" || (*tLong == "" && *tShort == "") || *startTime == "" || *endTime == "" {
			generateSummaryCmd.Usage()
			runtime.Goexit()
		}
		d := utils.GetFinalStringFromTwoStringWithSameMeaning(*dLong, *dShort)
		t := utils.GetFinalStringFromTwoStringWithSameMeaning(*tLong, *tShort)
		start, startErr := time.ParseInLocation(time.RFC3339, *startTime, time.UTC)
		end, endErr := time.ParseInLocation(time.RFC3339, *endTime, time.UTC)
		if startErr != nil || endErr != nil {
			fmt.Println("invalid time format")
			runtime.Goexit()
		}
		start = start.UTC()
		end = end.UTC()
		cli.GenerateSummary(d, t, start, end, *outputFileType, *outputFileName, path)
	}
}
