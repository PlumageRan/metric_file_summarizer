package main

import (
	"AccelByteTakeHome/cli"
	"log"
	"os"
)

func main() {
	defer os.Exit(0)
	cmd := cli.CommandLine{}
	path, err := os.Getwd()
	if err != nil {
		log.Fatalf("can not get working directory, err is %+v", err)
	}
	cmd.Run(path)
}
