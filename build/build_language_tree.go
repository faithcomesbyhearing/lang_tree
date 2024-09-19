package build

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"lang_tree/db"
	"os"
	"strconv"
)

func BuildLanguageTree() {
	var numLanguages = 26879
	languages := loadGlottoLanguoid()
	languages = loadIso6393(languages)
	languages = loadAIToolCompatibility(languages, "../data/espeak.tab", db.ESpeak, 1, 3)
	languages = loadAIToolCompatibility(languages, "../data/mms_asr.tab", db.MMSASR, 0, 1)
	languages = loadAIToolCompatibility(languages, "../data/mms_lid.tab", db.MMSLID, 0, 1)
	languages = loadAIToolCompatibility(languages, "../data/mms_tts.tab", db.MMSTTS, 0, 1)
	languages = loadAIToolCompatibility(languages, "../data/whisper.tab", db.Whisper, 1, 0)
	if len(languages) != numLanguages {
		fmt.Println("Load Iso6393: Expected ", numLanguages, " got ", len(languages))
		os.Exit(1)
	}
	outputJSON(languages)
}

func loadGlottoLanguoid() []db.Language {
	var languages []db.Language
	file, err := os.Open("../data/languoid.tab")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	first := true
	var record []string
	var count6393 = 0
	for {
		record, err = reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		if first {
			first = false
			continue
		}
		var lang db.Language
		lang.GlottoId = record[0]
		lang.FamilyId = record[1]
		lang.ParentId = record[2]
		lang.Name = record[3]
		lang.Bookkeeping, err = strconv.ParseBool(record[4])
		if err != nil {
			panic(err)
		}
		lang.Level = record[5]
		lang.Iso6393 = record[8]
		if lang.Iso6393 != "" {
			count6393++
		}
		lang.CountryIds = record[14]
		languages = append(languages, lang)
	}
	fmt.Println("Num iso639-3", count6393)
	return languages
}

func loadIso6393(languages []db.Language) []db.Language {
	var isoMap = make(map[string]string)
	var inGlotto = make(map[string]bool)
	var notInIso = make(map[string]bool)
	var notInGlotto = make(map[string]bool)
	fmt.Println("Num Glotto Records", len(languages))
	file, err := os.Open("../data/iso-639-3.tab")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = '\t'
	first := true
	var record []string
	for {
		record, err = reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		if first {
			first = false
			continue
		}
		iso6393 := record[0]
		iso6391 := record[3]
		isoMap[iso6393] = iso6391
	}
	fmt.Println("Num iso639-3 records", len(isoMap))
	for i := range languages {
		iso6393 := languages[i].Iso6393
		if iso6393 != "" {
			inGlotto[iso6393] = languages[i].Bookkeeping
			iso6391, ok := isoMap[iso6393]
			if ok {
				if iso6391 != "" {
					languages[i].Iso6391 = iso6391
				}
			} else {
				notInIso[iso6393] = true
				fmt.Println("Glotto ISO639-3 Not In ISO List", iso6393)
			}
		}
	}
	fmt.Println("Num iso639-3 values in glotto", len(inGlotto))
	fmt.Println("Num Glotto Not In ISO", len(notInIso))
	for iso6393 := range isoMap {
		bookkeeping, ok := inGlotto[iso6393]
		if !ok && !bookkeeping {
			notInGlotto[iso6393] = true
			fmt.Println("ISO639-3 Not In Glotto List", iso6393)
		}
	}
	fmt.Println("Num iso639-3 records not in glotto", len(notInGlotto))
	return languages
}

func loadAIToolCompatibility(languages []db.Language, filePath string, toolName string, iso3Col int, nameCol int) []db.Language {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	var toolMap = make(map[string]string)
	var record []string
	for {
		record, err = reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		toolMap[record[iso3Col]] = record[nameCol]
	}
	var usedMap = make(map[string]string)
	for i := range languages {
		name, ok := toolMap[languages[i].Iso6393]
		if ok {
			languages[i] = setLanguage(languages[i], toolName)
			usedMap[languages[i].Iso6393] = name
		} else {
			name, ok = toolMap[languages[i].Iso6391]
			if ok {
				languages[i] = setLanguage(languages[i], toolName)
				usedMap[languages[i].Iso6391] = name
			}
		}
	}
	var missingCount = 0
	for iso, name := range toolMap {
		_, ok := usedMap[iso]
		if !ok {
			fmt.Println("AI Tool", toolName, "Has iso code", iso, name, "but it has no match in table")
			missingCount++
		}
	}
	fmt.Println("Num ai-tool records not matching:", missingCount, "out of", len(toolMap))
	return languages
}

func setLanguage(language db.Language, toolName string) db.Language {
	switch toolName {
	case db.ESpeak:
		language.ESpeak = true
	case db.MMSASR:
		language.MMSASR = true
	case db.MMSLID:
		language.MMSLID = true
	case db.MMSTTS:
		language.MMSTTS = true
	case db.Whisper:
		language.Whisper = true
	default:
		panic("Unknown tool name: " + toolName)
	}
	return language
}

func outputJSON(languages []db.Language) {
	bytes, err := json.MarshalIndent(languages, "", "    ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("../db/language_tree.json", bytes, 0644)
	if err != nil {
		panic(err)
	}
}
