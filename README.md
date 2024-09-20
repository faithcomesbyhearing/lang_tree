# lang_tree

Currently, AI language tools, such as speech-to-text, text-to-speech, 
and language-ident have been developed for only a small percentage of the 
known languages.  As a result when doing AI work on a language that 
none of the tools process, it is essential to find a language closely 
related, which does have AI tool support.

The lang_tree program uses a glotto log hierarch of languages to find 
related languages.  This tool is able to find related languages for eSpeak, 
mms-language-ident, mms-speech-to-text, mms-text-to-speech, 
whisper-speech-to-text, whisper-translation.

Searches are done up and down the tree to find the closest language 
that is supported by the AI tool of interest. Closeness is defined 
by a simple count of the number of nodes.  This is a naive solution. 
If someone can suggest a statistic that defines the similarity of a 
language to its children in the hierarchy, a better algorithm can be 
added.

To install this package as a go module:
> go get github.com/garygriswold/lang_tree.git

To use lang_tree in a go program:
```
var tree = NewLanguageTree(ctx context.Context)
err := tree.Load()
languages, distance, err := tree.Search(iso639_3 string, aiTool string)
```

The AI tools supported are as follows:
* espeak
* mms_asr
* mms_lid
* mms_tts
* whisper

The Search func returns a slice of type db.Language, which is as follows:
```
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
```