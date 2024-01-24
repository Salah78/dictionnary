package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
	"estiam/dictionary"
	"estiam/middleware"
	"errors"
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

	router.Use(middleware.AuthMiddleware)
	router.Use(middleware.Logger)

	router.HandleFunc("/add", addHandler).Methods("POST")
	router.HandleFunc("/define/{word}", defineHandler).Methods("PUT")
	router.HandleFunc("/remove/{word}", removeHandler).Methods("DELETE")
	router.HandleFunc("/list", listHandler).Methods("GET")
	router.HandleFunc("/exit", exitHandler).Methods("POST")

	router.HandleFunc("/generate-token", generateTokenHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func generateTokenHandler(w http.ResponseWriter, r *http.Request) {
	token := "mon_token_secret"
	fmt.Fprintf(w, "Token généré: %s\n", token)
}

func validateWordAndDefinition(word, definition string) error {
	if len(word) < 3 {
		return errors.New("Le mot doit avoir au moins 3 caractères.")
	}

	if len(definition) < 5 {
		return errors.New("La définition doit avoir au moins 5 caractères.")
	}

	return nil
}

func writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "Erreur: %s\n", message)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	word := r.FormValue("word")
	definition := r.FormValue("definition")

	if err := validateWordAndDefinition(word, definition); err != nil {
		writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	d.Add(word, definition)

	fmt.Fprintf(w, "Word '%s' added with definition '%s'.\n", word, definition)
}

func defineHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	word := vars["word"]

	var input map[string]string
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeErrorResponse(w, "Erreur de décodage JSON", http.StatusBadRequest)
		return
	}

	newDefinition, ok := input["definition"]
	if !ok {
		writeErrorResponse(w, "Champ 'definition' manquant dans le corps JSON", http.StatusBadRequest)
		return
	}

	if err := validateWordAndDefinition(word, newDefinition); err != nil {
		writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

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
	os.Exit(0)
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
	d, err := dictionary.New()
	if err != nil {
		fmt.Println("Error creating dictionary:", err)
		return d
	}

	fmt.Println("Loading dictionary from Redis...")

	// Fetch entries from Redis and populate the dictionary
	wordList, entries := d.List()

	for _, word := range wordList {
		entry, found := entries[word]
		if !found {
			fmt.Println("Error fetching entry from Redis: Entry not found for word", word)
			continue
		}
		d.Add(word, entry.Definition)
	}

	fmt.Println("Dictionary loaded from Redis.")
	return d
}
