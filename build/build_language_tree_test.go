package build

import (
	"encoding/json"
	"fmt"
	"lang_tree/db"
	"os"
	"testing"
)

func TestBuildLanguageTree(t *testing.T) {
	var numLanguages = 26879
	languages := loadGlottoLanguoid()
	languages = loadIso6393(languages)
	languages = loadAIToolCompatibility(languages, "db/language/espeak.tab", db.ESpeak, 1, 3)
	languages = loadAIToolCompatibility(languages, "db/language/mms_asr.tab", db.MMSASR, 0, 1)
	languages = loadAIToolCompatibility(languages, "db/language/mms_lid.tab", db.MMSLID, 0, 1)
	languages = loadAIToolCompatibility(languages, "db/language/mms_tts.tab", db.MMSTTS, 0, 1)
	languages = loadAIToolCompatibility(languages, "db/language/whisper.tab", db.Whisper, 1, 0)
	if len(languages) != numLanguages {
		fmt.Println("Load Iso6393: Expected ", numLanguages, " got ", len(languages))
		os.Exit(1)
	}
	outputJSON(languages)
}

func outputJSON(languages []db.Language) {
	bytes, err := json.MarshalIndent(languages, "", "    ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("db/language/language_tree.jason", bytes, 0644)
	if err != nil {
		panic(err)
	}
}
