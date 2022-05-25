package main

//
// start the coordinator process, which is implemented
// in ../mr/coordinator.go
//
// go run mrcoordinator.go pg*.txt
//
// Please do not change this file.
//

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmattson4/map-reduce/mr"
)

func main() {
	logFlag := flag.String("log", "coordinator.log", "Sets output location of logs")
	reduceCountFlag := flag.Int("n", 10, "Set count of intermediary reduce files")
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: mrcoordinator inputfiles...\n")
		os.Exit(1)
	}

	lf, err := os.OpenFile(*logFlag, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		log.Fatalf("Cannot open log file at location %v", *logFlag)
	}

	log := log.New(lf, "", log.LstdFlags)
	log.Print(os.Args)
	m := mr.MakeCoordinator(os.Args[1:], *reduceCountFlag, log, "grpc", 8089)

	for m.Done() == false {
		time.Sleep(time.Second)
	}
	time.Sleep(time.Second)
}
