package main

import (
	"flag"
	"log"
	"os"
	"osu-background-deleter/config"
	"osu-background-deleter/gui"
	"osu-background-deleter/mem"
	"osu-background-deleter/memory"
	"osu-background-deleter/proc"
	"runtime"

	"github.com/spf13/cast"
)

func main() {
	if proc.CheckIfProcessRunning("osu! Background Deleter") {
		gui.Init(true)
		return
	}

	config.Init()
	updateTimeFlag := flag.Int("update", cast.ToInt(config.Config["update"]), "How fast should we update the values? (in milliseconds)")

	isRunningInWINE := flag.Bool("wine", cast.ToBool(config.Config["wine"]), "Running under WINE?")
	songsFolderFlag := flag.String("path", config.Config["path"], `Path to osu! Songs directory ex: /mnt/ps3drive/osu\!/Songs`)
	memDebugFlag := flag.Bool("memdebug", cast.ToBool(config.Config["memdebug"]), `Enable verbose memory debugging?`)
	memCycleTestFlag := flag.Bool("memcycletest", cast.ToBool(config.Config["memcycletest"]), `Enable memory cycle time measure?`)
	// disablecgo := flag.Bool("cgodisable", cast.ToBool(config.Config["cgodisable"]), `Disable everything non memory-reader related? (pp counters)`)
	flag.Parse()
	// cgo := *disablecgo
	mem.Debug = *memDebugFlag
	memory.MemCycle = *memCycleTestFlag
	memory.UpdateTime = *updateTimeFlag
	memory.SongsFolderPath = *songsFolderFlag
	memory.UnderWine = *isRunningInWINE
	if runtime.GOOS != "windows" && memory.SongsFolderPath == "auto" {
		log.Fatalln("Please specify path to osu!Songs (see --help)")
	}
	if memory.SongsFolderPath != "auto" {
		if _, err := os.Stat(memory.SongsFolderPath); os.IsNotExist(err) {
			log.Fatalln(`Specified Songs directory does not exist on the system! (try setting to "auto" if you are on Windows or make sure that the path is correct)`)
		}
	}

	go memory.Init()
	gui.Init(false)
}
