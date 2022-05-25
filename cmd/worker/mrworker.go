package main

//
// start a worker process, which is implemented
// in ../mr/worker.go. typically there will be
// multiple worker processes, talking to one coordinator.
//
// go run mrworker.go wc.so
//
// Please do not change this file.
//

import (
	"flag"
	"fmt"
	"log"
	"os"
	"plugin"

	"github.com/jmattson4/map-reduce/mr"
)

func main() {
	pid := os.Getpid()
	logFlag := flag.String("log", fmt.Sprintf("worker-%v.log", pid), "Sets output location of logs")
	flag.Parse()

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: mrworker xxx.so\n")
		os.Exit(1)
	}

	lf, err := os.OpenFile(*logFlag, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		log.Fatalf("Cannot open log file at location %v", *logFlag)
	}

	log := log.New(lf, "", log.LstdFlags)

	mapf, reducef := loadPlugin(os.Args[1])

	mr.Worker(mapf, reducef, log)
}

//
// load the application Map and Reduce functions
// from a plugin file, e.g. ../mrapps/wc.so
//
func loadPlugin(filename string) (func(string, string) []mr.KeyValue, func(string, []string) string) {
	p, err := plugin.Open(filename)
	if err != nil {
		log.Fatalf("cannot load plugin %v", filename)
	}
	xmapf, err := p.Lookup("Map")
	if err != nil {
		log.Fatalf("cannot find Map in %v", filename)
	}
	mapf := xmapf.(func(string, string) []mr.KeyValue)
	xreducef, err := p.Lookup("Reduce")
	if err != nil {
		log.Fatalf("cannot find Reduce in %v", filename)
	}
	reducef := xreducef.(func(string, []string) string)

	return mapf, reducef
}
