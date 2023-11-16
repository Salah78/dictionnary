package dictionary

import "fmt"

type Entry struct {
	Word       string
	Definition string
}

func (e Entry) String() string {
	return e.Definition
}

type Dictionary struct {
   entries map[string]Entry
}

func New() *Dictionary {
   return &Dictionary{
      entries: make(map[string]Entry),
   }
}

func (d *Dictionary) Add(word string, definition string) {
   entry := Entry{Word: word, Definition: definition}
   d.entries[word] = entry
}

func (d *Dictionary) Get(word string) (Entry, error) {
   entry, found := d.entries[word]
   if !found {
      return Entry{}, fmt.Errorf("word '%s' not found", word)
   }
   return entry, nil
}

func (d *Dictionary) Remove(word string) {
   delete(d.entries, word)
}

func (d *Dictionary) List() ([]string, map[string]Entry) {
   wordList := make([]string, 0, len(d.entries))
   for word := range d.entries {
      wordList = append(wordList, word)
   }
   return wordList, d.entries
}
