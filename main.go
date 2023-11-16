package main

import (
	"bufio"
	"encoding/json" 
	"fmt"
	"os"
	"strconv" 
	"strings"
	"sync"

	"estiam/dictionary"
)

const dictionaryFilePath = "dictionary.json"

type SaveData struct {
	Entries map[string]struct {
		Definition string `json:"definition"`
	} `json:"entries"`
}

var d *dictionary.Dictionary
var mu sync.Mutex

func main() {
	d = loadDictionary()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Choose an action:")
		fmt.Println("1. Add")
		fmt.Println("2. Define")
		fmt.Println("3. Remove")
		fmt.Println("4. List")
		fmt.Println("5. Exit")

		choice, err := getUserChoice(reader)
		if err != nil {
			fmt.Println("Error reading choice:", err)
			continue
		}

		switch choice {
		case 1:
			actionAdd(reader)
		case 2:
			actionDefine(reader)
		case 3:
			actionRemove(reader)
		case 4:
			actionList()
		case 5:
			saveDictionary()
			fmt.Println("Exiting the program.")
			return
		default:
			fmt.Println("Invalid choice. Please choose a valid option.")
		}
	}
}

func getUserChoice(reader *bufio.Reader) (int, error) {
	fmt.Print("Enter your choice: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	text = strings.TrimSpace(text)
	choice, err := strconv.Atoi(text)
	if err != nil {
		return 0, err
	}
	return choice, nil
}

func actionAdd(reader *bufio.Reader) {
	fmt.Print("Enter the word: ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	fmt.Print("Enter the definition: ")
	definition, _ := reader.ReadString('\n')
	definition = strings.TrimSpace(definition)

	d.Add(word, definition)

	fmt.Printf("Word '%s' added with definition '%s'.\n", word, definition)
}

func actionDefine(reader *bufio.Reader) {
	fmt.Print("Enter the word: ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	entry, err := d.Get(word)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Definition of '%s': %s\n", entry.Word, entry.Definition)
}

func actionRemove(reader *bufio.Reader) {
	fmt.Print("Enter the word to remove: ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	d.Remove(word)

	fmt.Printf("Word '%s' removed.\n", word)
}

func actionList() {
	words, entries := d.List()
	fmt.Println("Words in the dictionary:")
	for _, word := range words {
		fmt.Printf("%s: %s\n", word, entries[word])
	}
}

func saveDictionary() {
	saveData := SaveData{
		Entries: make(map[string]struct {
			Definition string `json:"definition"`
		}),
	}

	words, entries := d.List()
	for _, word := range words {
		saveData.Entries[word] = struct {
			Definition string `json:"definition"`
		}{Definition: entries[word].Definition}
	}

	data, err := json.MarshalIndent(saveData, "", "  ")
	if err != nil {
		fmt.Println("Error encoding dictionary:", err)
		return
	}

	err = os.WriteFile(dictionaryFilePath, data, 0644)
	if err != nil {
		fmt.Println("Error saving dictionary to file:", err)
	}
}

func loadDictionary() *dictionary.Dictionary {
	d := dictionary.New()

	data, err := os.ReadFile(dictionaryFilePath)
	if err != nil {
		fmt.Println("Error reading dictionary file:", err)
		return d
	}

	if len(data) == 0 {
		return d
	}

	var saveData SaveData
	err = json.Unmarshal(data, &saveData)
	if err != nil {
		fmt.Println("Error decoding dictionary:", err)
		return d
	}

	for word, entry := range saveData.Entries {
		d.Add(word, entry.Definition)
	}

	return d
}
