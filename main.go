package main

import (
	"csvreader/teval"
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

func main() {
	log.SetFlags(0)

	sep := flag.String("s", ",", "change separator")
	shouldEvaluate := flag.Bool("e", false, "evaluate table")
	cpuProfile := flag.String("cpuprofile", "", "write cpu profile to file")
	memProfile := flag.String("memprofile", "", "write memory profile to file")
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatal("usage: csvreader *filename*")
	}

	// Start profiling
	if *cpuProfile != "" {
		file, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		if err := pprof.StartCPUProfile(file); err != nil {
			log.Fatal(err)
		}
		defer pprof.StopCPUProfile()
	}

	filename := flag.Args()[0]
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	table, err := teval.Read(file, *sep)
	if err != nil {
		log.Fatal(err)
	}

	if *shouldEvaluate {
		if err := table.EvalAll(); err != nil {
			log.Fatal(err)
		}
	}

	if *memProfile != "" {
		file, err := os.Create(*memProfile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		runtime.GC()
		if err := pprof.WriteHeapProfile(file); err != nil {
			log.Fatal(err)
		}
	}

	if err := teval.Write(os.Stdout, table, *sep); err != nil {
		log.Fatal(err)
	}
}
