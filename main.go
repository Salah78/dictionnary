package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"os"
	"github.com/gorilla/mux"
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

	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/add", addHandler).Methods("POST")
	router.HandleFunc("/define/{word}", defineHandler).Methods("PUT")
	router.HandleFunc("/remove/{word}", removeHandler).Methods("DELETE")
	router.HandleFunc("/list", listHandler).Methods("GET")
	router.HandleFunc("/exit", exitHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	word := r.FormValue("word")
	definition := r.FormValue("definition")

	d.Add(word, definition)

	fmt.Fprintf(w, "Word '%s' added with definition '%s'.\n", word, definition)
}

func defineHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    word := vars["word"]

    // Récupérer la nouvelle définition du corps JSON
    var input map[string]string
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Erreur de décodage JSON", http.StatusBadRequest)
        return
    }

    newDefinition, ok := input["definition"]
    if !ok {
        http.Error(w, "Champ 'definition' manquant dans le corps JSON", http.StatusBadRequest)
        return
    }

    // Mettre à jour la définition existante
    d.Add(word, newDefinition)

    fmt.Fprintf(w, "Définition de '%s' mise à jour avec '%s'.\n", word, newDefinition)
}

func removeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	word := vars["word"]

	d.Remove(word)

	fmt.Fprintf(w, "Word '%s' removed.\n", word)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	words, entries := d.List()
	fmt.Fprintln(w, "Words in the dictionary:")
	for _, word := range words {
		fmt.Fprintf(w, "%s: %s\n", word, entries[word])
	}
}

func exitHandler(w http.ResponseWriter, r *http.Request) {
	saveDictionary()
	fmt.Fprintln(w, "Exiting the program.")
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

	fmt.Println("Loading dictionary from file:", dictionaryFilePath)

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
