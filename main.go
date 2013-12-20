package main

import (
	. "github.com/Zephyyrr/BiG"
	"flag"
	"log"
	"os"
	"io"
	"time"
	"bufio"
)

const (
	SPEED = time.Millisecond
)

var (
	verbose = flag.Bool("v", false, "Verbose output.")
	timelimit = flag.Duration("t", 0, "Timelimit for execution. 0 for no limit.")
	speed = flag.Duration("s", SPEED, "time between two ticks in the VM.")
	fastmode = flag.Bool("f", false, "Fast-mode. No delay.")
)

func main() {
	flag.Parse()
	if flag.Arg(0) == "" {
		log.Println("Needs a file on commandline. Exiting.")
		os.Exit(1)
	}
		file, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Println(flag.Arg(0))
		log.Fatal(err)
	}
	fs := loadFungeSpace(file)
	vm := NewVM(fs)
	
	timer := time.NewTimer(*timelimit)
	if *timelimit <= 0 {
		timer.Stop()
	}
	done := make(chan bool)
	if *verbose {
		log.Println("Program:\n", vm.FS)
		log.Println("Timelimit:", *timelimit)
		log.Println("---BEGIN---")
	}
	if *fastmode {
		go func() {
			for !vm.Done() {
				vm.Tick()
			}
		}()
	} else {
		ticker := time.NewTicker(*speed)
		go vm.Run(ticker)
	}
	go func(){
		for !vm.Done(){
			time.Sleep(10*time.Millisecond)
		}
		done<-true
	}()
	
	select {
		case <-timer.C:
			if *verbose {
				log.Println("---END---")
			}
			log.Println("Execution took to long.")
			os.Exit(1)
		case <-done: 
			if *verbose {
				log.Println("---END---")
			}
	}
}

func loadFungeSpace(in io.Reader) FungeSpace {
	buff := bufio.NewReaderSize(in, 80)
	fs := new(Befunge93Space)
	l, _, err := buff.ReadLine()
	for i := 0; err == nil && i < 25 ; i++ {
		for j, v := range l {
			fs[i][j] = int8(v)
		}
		l, _, err = buff.ReadLine()
	}
	return fs
}
