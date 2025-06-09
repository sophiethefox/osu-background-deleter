package gui

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"osu-background-deleter/config"
	"osu-background-deleter/memory"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/spf13/cast"

	"time"
)

var JSONByte []byte

type JSONMessage struct {
	Settings struct {
		Folders struct {
			Songs string `json:"songs"`
		} `json:"folders"`
	} `json:"settings"`
	Menu struct {
		BM struct {
			ID   json.Number `json:"id"`
			Path struct {
				Folder string `json:"folder"`
				BG     string `json:"bg"`
			} `json:"path"`
		} `json:"bm"`
	} `json:"menu"`
}

func CopyFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %v", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("Couldn't open dest file: %v", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return fmt.Errorf("Couldn't copy to dest from source: %v", err)
	}

	inputFile.Close() // for Windows, close before trying to remove: https://stackoverflow.com/a/64943554/246801

	// err = os.Remove(sourcePath)
	// if err != nil {
	// 	return fmt.Errorf("Couldn't remove source file: %v", err)
	// }
	return nil
}

func ReplaceWithBlackImage(path string) error {
	// Open original image to get dimensions
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode image (supports PNG and JPEG)
	img, format, err := image.DecodeConfig(file)

	if err != nil {
		return err
	}

	println(format)
	width, height := img.Width, img.Height

	// Create a new black RGBA image
	black := color.RGBA{0, 0, 0, 255}
	blackImage := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := range height {
		for x := range width {
			blackImage.Set(x, y, black)
		}
	}

	// Overwrite the original file
	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()

	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png":
		err = png.Encode(outFile, blackImage)
	case ".jpg", ".jpeg":
		err = jpeg.Encode(outFile, blackImage, &jpeg.Options{Quality: 90})
	default:
		err = jpeg.Encode(outFile, blackImage, &jpeg.Options{Quality: 90}) // fallback
	}

	return err
}

func loadImage(path string) *canvas.Image {
	if path == "\\\\" {
		return nil
	}
	file, err := os.Open(path)
	if err != nil {
		// log.Println("Failed to open image:", err)
		return nil
	}
	defer file.Close()

	// Automatically detect and decode image format
	img, _, err := image.Decode(file)
	if err != nil {
		// log.Println("Failed to decode image:", err)
		return nil
	}

	image := canvas.NewImageFromImage(img)
	image.FillMode = canvas.ImageFillContain
	return image
}

var (
	deletedLastOriginalFilename string
	deletedLastOutputFilename   string
	updatePending               bool
)

// todo: add deleting from thumbnail cache, /osu/Data/bt/[id]l?.jpg

func ClearAndMoveBackground(imagePath, id, filename string) {
	_lc := strings.ToLower(imagePath)
	// sanity check. ensure we're only deleting an image.
	if strings.HasSuffix(_lc, ".png") || strings.HasSuffix(_lc, ".jpg") || strings.HasSuffix(_lc, ".jpeg") {
		outDir := "./deleted_backgrounds" + "/" + id
		outFile := outDir + "/" + filename
		// create folder for beatmap id deleted files
		os.MkdirAll(outDir, os.ModePerm)

		// copy original background to deleted bg directory
		CopyFile(imagePath, outFile)

		// overwrite original bg w/ black img
		ReplaceWithBlackImage(imagePath)

		deletedLastOriginalFilename = imagePath
		deletedLastOutputFilename = outFile
		updatePending = true
	}
}

func UndoDeletion() {
	// move file from deleted backgrounds directory to the osu map directory
	// even more sanity checks yipppeeeeee!!! i dont like shipping code that involves deleting files
	lc_ := strings.ToLower(deletedLastOriginalFilename)
	lc__ := strings.ToLower(deletedLastOutputFilename)
	if strings.HasSuffix(lc_, ".png") || strings.HasSuffix(lc_, ".jpg") || strings.HasSuffix(lc_, ".jpeg") {
		if strings.HasSuffix(lc__, ".png") || strings.HasSuffix(lc__, ".jpg") || strings.HasSuffix(lc__, ".jpeg") {
			CopyFile(deletedLastOutputFilename, deletedLastOriginalFilename)
			os.Remove(deletedLastOutputFilename)
			updatePending = true
		}
	}

}

var (
	lastText           string
	lastImage          string
	lastImageBeatmapId string
	lastImageFilename  string
	updateMutex        sync.Mutex
	text               *widget.Label
	imageWidget        *canvas.Image
)

func Clamp(f, low, high int) int {
	if f < low {
		return low
	}
	if f > high {
		return high
	}
	return f
}

func Init(processAlreadyRunning bool) {
	println(processAlreadyRunning)
	flag.Parse()
	log.SetFlags(0)

	a := app.New()
	Window := a.NewWindow("osu! Map Background Deleter")

	text = widget.NewLabelWithStyle("Waiting for data...", fyne.TextAlignLeading, fyne.TextStyle{Monospace: false})
	text.Wrapping = fyne.TextWrapWord
	textContainer := container.NewStack(text)
	textContainer.Resize(fyne.NewSize(480, text.MinSize().Height))

	imageWidget = canvas.NewImageFromResource(nil)
	imageWidget.FillMode = canvas.ImageFillContain
	imageWidget.SetMinSize(fyne.NewSize(600, 350))

	button1 := widget.NewButton("Delete", func() {
		log.Println("Deleting " + lastImage)
		ClearAndMoveBackground(lastImage, lastImageBeatmapId, lastImageFilename)
	})

	button2 := widget.NewButton("Undo", func() {
		log.Println("Undoing deletion of " + lastImage)
		UndoDeletion()
	})
	buttons := container.NewHBox(button1, button2)

	content := container.NewVBox(
		textContainer,
		imageWidget,
		buttons,
	)

	if cast.ToBool(config.Config["minimize_to_tray"]) {
		if desk, ok := a.(desktop.App); ok {
			m := fyne.NewMenu("osu! Background Deleter",
				fyne.NewMenuItem("Show", func() {
					Window.Show()
				}))
			desk.SetSystemTrayMenu(m)
		}

		Window.SetCloseIntercept(func() {
			Window.Hide()
		})
	} else {
		Window.SetCloseIntercept(func() {
			a.Quit()
			syscall.Exit(0)
		})
	}

	Window.SetContent(content)
	Window.Resize(fyne.NewSize(600, 400))
	Window.SetFixedSize(false)

	if processAlreadyRunning {
		d := dialog.NewInformation("osu! Map Background Deleter", "The app is already running! Please check your system tray.", Window)
		d.SetOnClosed(func() {
			a.Quit()
			os.Exit(0)
		})
		d.Show()
	}

	Window.Show()

	go func() {
		ticker := time.NewTicker(time.Millisecond * 100)
		for range ticker.C {
			fyne.Do(func() {
				var err error

				type wsStruct struct { //order sets here
					A memory.InSettingsValues `json:"settings"`
					B memory.InMenuValues     `json:"menu"`
				}
				// for {
				group := wsStruct{
					A: memory.SettingsData,
					B: memory.MenuData,
				}

				JSONByte, err = json.Marshal(group)
				if err != nil {
					println("")
				}

				// probs change this. unnecessary de/serialization idk
				var msg JSONMessage

				err_ := json.Unmarshal(JSONByte, &msg)
				if err_ != nil {
					log.Println("json unmarshal:", err_)
					return
				}

				songFolder := msg.Settings.Folders.Songs
				folder := msg.Menu.BM.Path.Folder
				lastImageBeatmapId = msg.Menu.BM.ID.String()
				lastImageFilename = msg.Menu.BM.Path.BG
				bg := songFolder + "\\" + folder + "\\" + lastImageFilename

				displayText := "Background file: " + bg

				updateMutex.Lock()
				sameImage := (bg == lastImage)
				updateMutex.Unlock()

				if sameImage && !updatePending {
					return
				}

				updateMutex.Lock()
				defer updateMutex.Unlock()

				if displayText != lastText {
					text.SetText(displayText)
					lastText = displayText
				}

				if (bg != lastImage || updatePending) && bg != "" {
					if newImg := loadImage(bg); newImg != nil {
						imageWidget.Image = newImg.Image
						imageWidget.SetMinSize(fyne.NewSize(600, 350))
						imageWidget.Refresh()
						lastImage = bg
					}
				}
				updatePending = false
			})
		}

	}()
	a.Run()
}
