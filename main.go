package main

import (
	"csvreader/teval"
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
)

func main() {
	sep := flag.String("s", ",", "change separator")
	shouldEvaluate := flag.Bool("e", false, "evaluate table")
	cpuProfile := flag.String("cpuprofile", "", "write cpu profile to file")
	memProfile := flag.String("memprofile", "", "write memory profile to file")
	flag.Parse()
	if flag.NArg() < 1 {
		exit(fmt.Errorf("usage: csvreader *filename*"))
	}

	// Start profiling
	if *cpuProfile != "" {
		file, err := os.Create(*cpuProfile)
		if err != nil {
			exit(err)
		}
		defer file.Close()

		if err := pprof.StartCPUProfile(file); err != nil {
			exit(err)
		}
		defer pprof.StopCPUProfile()
	}

	filename := flag.Args()[0]
	file, err := os.Open(filename)
	if err != nil {
		exit(err)
	}
	defer file.Close()

	table, err := teval.Read(file, *sep)
	if err != nil {
		exit(err)
	}

	if *shouldEvaluate {
		if err := table.EvalAll(); err != nil {
			exit(err)
		}
	}

	if *memProfile != "" {
		file, err := os.Create(*memProfile)
		if err != nil {
			exit(err)
		}
		defer file.Close()

		runtime.GC()
		if err := pprof.WriteHeapProfile(file); err != nil {
			exit(err)
		}
	}

	if err := teval.Write(os.Stdout, table, *sep); err != nil {
		exit(err)
	}
}

func exit(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}
