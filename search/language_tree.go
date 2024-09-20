package search

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"lang_tree/db"
	"os"
)

// https://github.com/espeak-ng/espeak-ng/blob/master/docs/languages.md?plain=1

type LanguageTree struct {
	ctx    context.Context
	Table  []db.Language
	Roots  []*db.Language
	idMap  map[string]*db.Language
	isoMap map[string]*db.Language
}

func NewLanguageTree(ctx context.Context) LanguageTree {
	var l LanguageTree
	l.ctx = ctx
	l.idMap = make(map[string]*db.Language)
	l.isoMap = make(map[string]*db.Language)
	return l
}

func (l *LanguageTree) Load() error {
	err := l.loadTable()
	if err != nil {
		return err
	}
	err = l.buildTree()
	return err
}

func (l *LanguageTree) Search(iso6393 string, search string) ([]*db.Language, int, error) {
	var results []*db.Language
	var distance int
	var lang *db.Language
	lang, ok := l.isoMap[iso6393]
	if !ok {
		err := errors.New("iso code " + iso6393 + " is not known.")
		return results, distance, err
	}
	if !l.validateSearch(search) {
		err := errors.New("Search parameter" + search + "is not known")
		return results, distance, err
	}
	results, distance = l.searchRelatives(lang, search)
	return results, distance, nil
}

func (l *LanguageTree) loadTable() error {
	// Read json file of languages
	filename := "../db/language_tree.json"
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	// Parse json into Language slice
	err = json.Unmarshal(content, &l.Table)
	return err
}

func (l *LanguageTree) buildTree() error {
	// Make Map of GlottoId
	for i := range l.Table {
		lang := l.Table[i]
		l.idMap[lang.GlottoId] = &lang
	}
	// Build Tree
	var idMapCount int
	var parentIdCount int
	for glottoId, lang := range l.idMap {
		idMapCount++
		if lang.ParentId != "" {
			parentIdCount++
			parent, ok := l.idMap[lang.ParentId]
			if !ok {
				return errors.New("Missing parent id " + lang.ParentId + " is not known.")
			}
			lang.Parent = parent
			l.idMap[glottoId] = lang
			lang.Parent.Children = append(lang.Parent.Children, lang)
			l.idMap[lang.ParentId] = lang.Parent
		}
		if lang.Iso6393 != "" {
			l.isoMap[lang.Iso6393] = lang
		}
	}
	fmt.Println("count: ", idMapCount)
	fmt.Println("parent count: ", parentIdCount)
	// Build root
	for _, lang := range l.idMap {
		if lang.Parent == nil {
			l.Roots = append(l.Roots, lang)
		}
	}
	return nil
}

func (l *LanguageTree) searchRelatives(start *db.Language, search string) ([]*db.Language, int) {
	var finalLang, results []*db.Language
	var hierDown int
	var hierUp = -1
	var limit = 1000
	for limit > 0 && start != nil {
		hierUp++
		fmt.Println("\nSearching", start.Name, search, limit)
		results, hierDown = l.descendantSearch(start, search, limit)
		fmt.Println("hierUp", hierUp, "hierDown", hierDown, "num", len(results))
		for _, result := range results {
			fmt.Println("descendentSearch lang.Name", result.Name)
		}
		if len(results) > 0 {
			finalLang = results
			limit = hierDown - 1
		}
		start = start.Parent
	}
	return finalLang, hierUp + hierDown
}

// DescendantSearch performs a breadth-first search of the LanguageTree
func (l *LanguageTree) descendantSearch(start *db.Language, search string, limit int) ([]*db.Language, int) {
	var results []*db.Language
	var depth int
	if start == nil {
		return results, depth
	}
	var queue LanguageQueue
	queue.Enqueue(start, 0)
	for !queue.IsEmpty() {
		item := queue.Dequeue()
		if (len(results) > 0 && item.Depth > depth) || item.Depth > limit {
			return results, depth
		}
		depth = item.Depth
		if l.isMatch(item.Lang, search) {
			results = append(results, item.Lang)
		}
		fmt.Printf("Depth: %d, Name: %s, GlottoId: %s\n", item.Depth, item.Lang.Name, item.Lang.GlottoId)

		for _, child := range item.Lang.Children {
			queue.Enqueue(child, item.Depth+1)
		}
	}
	return results, depth
}

func (l *LanguageTree) validateSearch(search string) bool {
	switch search {
	case db.ESpeak:
		return true
	case db.MMSASR:
		return true
	case db.MMSLID:
		return true
	case db.MMSTTS:
		return true
	case db.Whisper:
		return true
	default:
		return false
	}
}

func (l *LanguageTree) isMatch(lang *db.Language, search string) bool {
	switch search {
	case db.ESpeak:
		return lang.ESpeak
	case db.MMSASR:
		return lang.MMSASR
	case db.MMSLID:
		return lang.MMSLID
	case db.MMSTTS:
		return lang.MMSTTS
	case db.Whisper:
		return lang.Whisper
	default:
		return false
	}
}
