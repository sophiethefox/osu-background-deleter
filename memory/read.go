package memory

type PreSongSelectAddresses struct {
	Status        int64 `sig:"48 83 F8 04 73 1E"`
	SettingsClass int64 `sig:"83 E0 20 85 C0 7E 2F"`
}

type songsFolderD struct {
	SongsFolder string `mem:"[[Settings + 0xB8] + 0x4]"`
}

type PreSongSelectData struct {
	Status uint32 `mem:"[Status - 0x4]"`
}

type staticAddresses struct {
	PreSongSelectAddresses
	Base     int64 `sig:"F8 01 74 04 83 65"`
	PlayTime int64 `sig:"5E 5F 5D C3 A1 ?? ?? ?? ?? 89 ?? 04"`
	SkinData int64 `sig:"75 21 8B 1D"`
	Rulesets int64 `sig:"7D 15 A1 ?? ?? ?? ?? 85 C0"`
}

func (staticAddresses) Ruleset() string {
	return "[[Rulesets - 0xB] + 0x4]"
}

func (staticAddresses) Beatmap() string {
	return "[Base - 0xC]"
}

func (PreSongSelectAddresses) Settings() string {
	return "[SettingsClass + 0x8]"
}

type menuD struct {
	PreSongSelectData
	MenuGameMode       int32   `mem:"[Base - 0x33]"`
	Plays              int32   `mem:"[Base - 0x33] + 0xC"`
	Artist             string  `mem:"[[Beatmap] + 0x18]"`
	ArtistOriginal     string  `mem:"[[Beatmap] + 0x1C]"`
	Title              string  `mem:"[[Beatmap] + 0x24]"`
	TitleOriginal      string  `mem:"[[Beatmap] + 0x28]"`
	AR                 float32 `mem:"[Beatmap] + 0x2C"`
	CS                 float32 `mem:"[Beatmap] + 0x30"`
	HP                 float32 `mem:"[Beatmap] + 0x34"`
	OD                 float32 `mem:"[Beatmap] + 0x38"`
	StarRatingStruct   uint32  `mem:"[Beatmap] + 0x8C"`
	AudioFilename      string  `mem:"[[Beatmap] + 0x64]"`
	BackgroundFilename string  `mem:"[[Beatmap] + 0x68]"`
	Folder             string  `mem:"[[Beatmap] + 0x78]"`
	Creator            string  `mem:"[[Beatmap] + 0x7C]"`
	Name               string  `mem:"[[Beatmap] + 0x80]"`
	Path               string  `mem:"[[Beatmap] + 0x90]"`
	Difficulty         string  `mem:"[[Beatmap] + 0xAC]"`
	MapID              int32   `mem:"[Beatmap] + 0xC8"`
	SetID              int32   `mem:"[Beatmap] + 0xCC"`
	RankedStatus       int32   `mem:"[Beatmap] + 0x12C"` // unknown, unsubmitted, pending/wip/graveyard, unused, ranked, approved, qualified
	MD5                string  `mem:"[[Beatmap] + 0x6C]"`
	ObjectCount        int32   `mem:"[Beatmap] + 0xFC"`
}
