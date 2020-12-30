package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/dundee/gdu/cli"
	"github.com/gdamore/tcell/v2"
)

// AppVersion stores the current version of the app
var AppVersion = "development"

func main() {
	logFile := flag.String("log-file", "/dev/null", "Path to a logfile")
	ignorePath := flag.String("ignore", "/proc,/dev,/sys,/run", "Absolute paths to ignore (separated by comma)")
	showVersion := flag.Bool("v", false, "Prints version")
	flag.Parse()

	if *showVersion {
		fmt.Println("Version:\t", AppVersion)
		return
	}

	f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	screen.Init()

	ui := cli.CreateUI(screen)
	ui.SetIgnorePaths(strings.Split(*ignorePath, ","))

	args := flag.Args()

	if len(args) == 1 {
		ui.AnalyzePath(args[0])
	} else {
		if runtime.GOOS == "linux" {
			ui.ListDevices()
		} else {
			ui.AnalyzePath(".")
		}
	}

	ui.StartUILoop()
}
