package engine

import (
	"encoding/xml"
	"os"
	"sync"

	"github.com/huichen/sego"
)

type Documents struct {
	Title    string `xml:"title"`
	URL      string `xml:"url"`
	Abstract string `xml:"abstract"`
	ID       uint64
}

type Index map[string][]uint64

func (idx Index) Add(docs []Documents) {
	for _, doc := range docs {
		for _, token := range analyze(doc.Abstract) {
			ids := idx[token]
			if ids != nil && ids[len(ids)-1] == doc.ID {
				continue
			}
			idx[token] = append(ids, doc.ID)
		}
	}
}

func LoadDocument(path string) []Documents {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	type S struct {
		Documents []Documents `xml:"doc"`
	}

	dp := sync.Pool{
		New: func() interface{} {
			return new(S)
		},
	}
	dump := dp.Get().(*S)
	dec := xml.NewDecoder(f)
	if err = dec.Decode(&dump); err != nil {
		return nil
	}
	for i := 1; i < len(dump.Documents); i++ {
		dump.Documents[i].ID = uint64(i)
	}

	return dump.Documents
}

func (idx Index) Search(text string) []uint64 {
	var r []uint64
	for _, token := range analyze(text) {
		if ids, ok := idx[token]; ok {
			if r == nil {
				r = ids
			} else {
				r = intersection(r, ids)
			}
		} else {
			// Token doesn't exist.
			return nil
		}
	}
	return r
}

var sm sego.Segmenter

func LoadDict() {
	sm.LoadDictionary("/Users/songjian.li/Desktop/dictionary.txt")
}

func analyze(text string) []string {
	segments := sm.Segment([]byte(text))
	return sego.SegmentsToSlice(segments, false)
}

func intersection(a []uint64, b []uint64) []uint64 {
	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}
	r := make([]uint64, 0, maxLen)
	var i, j int
	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			i++
		} else if a[i] > b[j] {
			j++
		} else {
			r = append(r, a[i])
			i++
			j++
		}
	}
	return r
}
