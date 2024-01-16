package dictionary

import (
	"testing"
)

func TestAdd(t *testing.T) {
	d := New()
	word := "test"
	definition := "sample definition"

	d.Add(word, definition)

	entry, err := d.Get(word)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if entry.Definition != definition {
		t.Errorf("Expected definition %s, got %s", definition, entry.Definition)
	}
}

func TestGet(t *testing.T) {
	d := New()
	word := "test"
	definition := "sample definition"
	d.Add(word, definition)

	entry, err := d.Get(word)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if entry.Definition != definition {
		t.Errorf("Expected definition %s, got %s", definition, entry.Definition)
	}
}

func TestRemove(t *testing.T) {
	d := New()
	word := "test"
	definition := "sample definition"
	d.Add(word, definition)

	d.Remove(word)

	_, err := d.Get(word)
	if err == nil {
		t.Errorf("Expected an error, got none")
	}
}

func TestList(t *testing.T) {
	d := New()
	word1 := "test1"
	definition1 := "sample definition 1"
	word2 := "test2"
	definition2 := "sample definition 2"
	d.Add(word1, definition1)
	d.Add(word2, definition2)

	words, entries := d.List()

	if len(words) != 2 {
		t.Errorf("Expected 2 words, got %d", len(words))
	}

	if entries[word1].Definition != definition1 {
		t.Errorf("Expected definition %s, got %s", definition1, entries[word1].Definition)
	}

	if entries[word2].Definition != definition2 {
		t.Errorf("Expected definition %s, got %s", definition2, entries[word2].Definition)
	}
}
