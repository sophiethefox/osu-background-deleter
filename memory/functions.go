package memory

import (
	"log"
	"path/filepath"
	"strings"
	"time"

	"osu-background-deleter/mem"
)

// UpdateTime Intervall between value updates
var UpdateTime int

// UnderWine?
var UnderWine bool

// MemCycle test
var MemCycle bool

// SongsFolderPath is full path to osu! Songs. Gets set automatically on Windows (through memory)
var SongsFolderPath string

var allProcs []mem.Process
var process mem.Process
var procerr error

// Init the whole thing and get osu! memory values to start working with it.
func Init() {

	allProcs, procerr = mem.FindProcess(osuProcessRegex, "osu!lazer", "osu!framework")
	for {
		start := time.Now()
		if procerr != nil {
			DynamicAddresses.IsReady = false
			for procerr != nil {
				allProcs, procerr = mem.FindProcess(osuProcessRegex, "osu!lazer", "osu!framework")
				log.Println("It seems that we lost the process, retrying! ERROR:", procerr)
				time.Sleep(1 * time.Second)
			}
			err := initBase()
			for err != nil {
				log.Println("Failure mid getting offsets, retrying! ERROR:", err)
				err = initBase()
				time.Sleep(1 * time.Second)
			}
		}
		if !DynamicAddresses.IsReady {
			err := initBase()
			for err != nil {
				log.Println("Failure mid getting offsets, retrying! ERROR:", err)
				err = initBase()
				time.Sleep(1 * time.Second)
			}
		} else {
			err := mem.Read(process,
				&patterns.PreSongSelectAddresses,
				&menuData.PreSongSelectData)
			if err != nil {
				DynamicAddresses.IsReady = false
				log.Println("It appears that we lost the process, retrying! ERROR:", err)
				continue
			}
			MenuData.OsuStatus = menuData.Status

			bmUpdateData()
		}

		elapsed := time.Since(start)
		if MemCycle {
			log.Printf("Cycle took %s", elapsed)
		}
		time.Sleep(time.Duration(UpdateTime-int(elapsed.Milliseconds())) * time.Millisecond)

	}

}

var tempBeatmapString string
var tempGameMode int32 = 5

func bmUpdateData() error {
	mem.Read(process, &patterns, &menuData)

	bmString := menuData.Path
	if (strings.HasSuffix(bmString, ".osu") && tempBeatmapString != bmString) || (strings.HasSuffix(bmString, ".osu") && tempGameMode != menuData.MenuGameMode) { //On map/mode change
		for i := 0; i < 50; i++ {
			if menuData.BackgroundFilename != "" {
				break
			}
			time.Sleep(25 * time.Millisecond)
			mem.Read(process, &patterns, &menuData)
		}
		tempGameMode = menuData.MenuGameMode
		tempBeatmapString = bmString
		MenuData.Bm.BeatmapID = menuData.MapID
		MenuData.Bm.BeatmapSetID = menuData.SetID

		MenuData.GameMode = menuData.MenuGameMode
		MenuData.Bm.RandkedStatus = menuData.RankedStatus
		MenuData.Bm.BeatmapMD5 = menuData.MD5
		MenuData.Bm.Path = path{
			AudioPath:            menuData.AudioFilename,
			BGPath:               menuData.BackgroundFilename,
			BeatmapOsuFileString: menuData.Path,
			BeatmapFolderString:  menuData.Folder,
			FullMP3Path:          filepath.Join(SongsFolderPath, menuData.Folder, menuData.AudioFilename),
			FullDotOsu:           filepath.Join(SongsFolderPath, menuData.Folder, bmString),
			InnerBGPath:          filepath.Join(menuData.Folder, menuData.BackgroundFilename),
		}
	}

	return nil
}
