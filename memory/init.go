package memory

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"osu-background-deleter/mem"
)

var osuProcessRegex = regexp.MustCompile(`.*osu!\.exe.*`)
var patterns staticAddresses

var menuData menuD
var songsFolderData songsFolderD

func resolveSongsFolder() (string, error) {
	var err error
	osuExecutablePath, err := process.ExecutablePath()
	if err != nil {
		return "", err
	}
	if !strings.Contains(osuExecutablePath, `:\`) {
		log.Println("Automatic executable path finder has failed. Please try again or manually specify it. (see --help) GOT: ", osuExecutablePath)
		time.Sleep(5 * time.Second)
		return "", errors.New("osu! executable was not found")
	}
	rootFolder := strings.TrimSuffix(osuExecutablePath, "osu!.exe")
	songsFolder := filepath.Join(rootFolder, "Songs")
	if songsFolderData.SongsFolder == "Songs" || songsFolderData.SongsFolder == "CompatibilityContext" { //dirty hack to fix old stable offset
		return songsFolder, nil
	}
	return songsFolderData.SongsFolder, nil
}

func initBase() error {
	var err error

	allProcs, err = mem.FindProcess(osuProcessRegex, "osu!lazer", "osu!framework")
	if err != nil {
		return err
	}
	process = allProcs[0]

	err = mem.ResolvePatterns(process, &patterns.PreSongSelectAddresses)
	if err != nil {
		return err
	}

	err = mem.Read(process,
		&patterns.PreSongSelectAddresses,
		&menuData.PreSongSelectData)
	if err != nil {
		return err
	}
	fmt.Println("[MEMORY] Got osu!status addr...")

	if runtime.GOOS == "windows" && SongsFolderPath == "auto" {
		err = mem.Read(process,
			&patterns.PreSongSelectAddresses,
			&songsFolderData)
		if err != nil {
			return err
		}
		SongsFolderPath, err = resolveSongsFolder()
		if err != nil {
			log.Fatalln(err)
		}
	}
	fmt.Println("[MEMORY] Songs folder:", SongsFolderPath)
	pepath, err := process.ExecutablePath()
	if err != nil {
		panic(err)
	}
	SettingsData.Folders.Game = filepath.Dir(pepath)

	fmt.Println("[MEMORY] Resolving patterns...")
	err = mem.ResolvePatterns(process, &patterns)
	if err != nil {
		return err
	}

	SettingsData.Folders.Songs = SongsFolderPath

	DynamicAddresses.IsReady = true

	return nil
}
