package db

type ToolName string

const (
	ESpeak  ToolName = "espeak"
	MMSASR  ToolName = "mms_asr"
	MMSLID  ToolName = "mms_lid"
	MMSTTS  ToolName = "mms_tts"
	Whisper ToolName = "whisper"
)

type Language struct {
	GlottoId    string      `json:"id"`
	FamilyId    string      `json:"family_id"`
	ParentId    string      `json:"parent_id"`
	Name        string      `json:"name"`
	Bookkeeping bool        `json:"bookkeeping"`
	Level       string      `json:"level"` //(language, dialect, family)
	Iso6393     string      `json:"iso639_3"`
	CountryIds  string      `json:"country_ids"`
	Iso6391     string      `json:"iso639_1"`
	ESpeak      bool        `json:"espeak"`
	MMSASR      bool        `json:"mms_asr"`
	MMSLID      bool        `json:"mms_lid"`
	MMSTTS      bool        `json:"mms_tts"`
	Whisper     bool        `json:"whisper"`
	Parent      *Language   `json:"-"`
	Children    []*Language `json:"-"`
}
