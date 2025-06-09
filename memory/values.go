package memory

//InMenuValues inside osu!memory
type InMenuValues struct {
	MainMenuValues MainMenuValues `json:"mainMenu"`
	OsuStatus      uint32         `json:"state"`
	GameMode       int32          `json:"gameMode"`
	Bm             bm             `json:"bm"`
}

type folders struct {
	Game  string `json:"game"`
	Songs string `json:"songs"`
}

type MainMenuValues struct {
	BassDensity float64 `json:"bassDensity"`
}

//InSettingsValues are values represented inside settings class, could be dynamic
type InSettingsValues struct {
	ShowInterface bool    `json:"showInterface"` //dynamic in gameplay
	Folders       folders `json:"folders"`
}

type bm struct {
	BeatmapID      int32  `json:"id"`
	BeatmapSetID   int32  `json:"set"`
	BeatmapMD5     string `json:"md5"`
	RandkedStatus  int32  `json:"rankedStatus"` //unknown, unsubmitted, pending/wip/graveyard, unused, ranked, approved, qualified
	Path           path   `json:"path"`
	HitObjectStats string `json:"-"`
	BeatmapString  string `json:"-"`
}

type path struct {
	InnerBGPath          string `json:"full"`
	BeatmapFolderString  string `json:"folder"`
	BeatmapOsuFileString string `json:"file"`
	BGPath               string `json:"bg"`
	AudioPath            string `json:"audio"`
	FullMP3Path          string `json:"-"`
	FullDotOsu           string `json:"-"`
}

type dynamicAddresses struct {
	IsReady bool
}

//MenuData contains raw values taken from osu! memory
var MenuData = InMenuValues{}

//SettingsData contains raw values taken from osu! memory
var SettingsData = InSettingsValues{}

//DynamicAddresses are in-between pointers that lead to values
var DynamicAddresses = dynamicAddresses{}
