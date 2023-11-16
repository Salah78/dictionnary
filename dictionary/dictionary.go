package dictionary

import (
	"fmt"
	"sync"
)

type Entry struct {
	Word       string
	Definition string
}

func (e Entry) String() string {
	return e.Definition
}

type Dictionary struct {
	entries  map[string]Entry
	updateCh chan dictionaryUpdate
	mutex    sync.Mutex
}

type dictionaryUpdate struct {
	entry Entry
	del   bool
}

func New() *Dictionary {
	d := &Dictionary{
		entries:  make(map[string]Entry),
		updateCh: make(chan dictionaryUpdate),
	}

	go d.startConcurrentOperations()

	return d
}

func (d *Dictionary) startConcurrentOperations() {
	for {
		update := <-d.updateCh

		d.mutex.Lock()

		if update.del {
			delete(d.entries, update.entry.Word)
		} else {
			d.entries[update.entry.Word] = update.entry
		}

		d.mutex.Unlock()
	}
}

func (d *Dictionary) Add(word string, definition string) {
	entry := Entry{Word: word, Definition: definition}
	d.updateCh <- dictionaryUpdate{entry: entry}
}

func (d *Dictionary) Get(word string) (Entry, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	entry, found := d.entries[word]
	if !found {
		return Entry{}, fmt.Errorf("word '%s' not found", word)
	}
	return entry, nil
}

func (d *Dictionary) Remove(word string) {
	d.updateCh <- dictionaryUpdate{entry: Entry{Word: word}, del: true}
}

func (d *Dictionary) List() ([]string, map[string]Entry) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	wordList := make([]string, 0, len(d.entries))
	for word := range d.entries {
		wordList = append(wordList, word)
	}
	return wordList, d.entries
}
