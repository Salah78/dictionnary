package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
   "strings"

   "estiam/dictionary"
)

func main() {
   d := dictionary.New()
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
         actionAdd(d, reader)
      case 2:
         actionDefine(d, reader)
      case 3:
         actionRemove(d, reader)
      case 4:
         actionList(d)
      case 5:
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

func actionAdd(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Enter the word: ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	fmt.Print("Enter the definition: ")
	definition, _ := reader.ReadString('\n')
	definition = strings.TrimSpace(definition)

	d.Add(word, definition)
	fmt.Printf("Word '%s' added with definition '%s'.\n", word, definition)
}

func actionDefine(d *dictionary.Dictionary, reader *bufio.Reader) {
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

func actionRemove(d *dictionary.Dictionary, reader *bufio.Reader) {
	fmt.Print("Enter the word to remove: ")
	word, _ := reader.ReadString('\n')
	word = strings.TrimSpace(word)

	d.Remove(word)
	fmt.Printf("Word '%s' removed.\n", word)
}

func actionList(d *dictionary.Dictionary) {
	words, entries := d.List()
	fmt.Println("Words in the dictionary:")
	for _, word := range words {
		fmt.Printf("%s: %s\n", word, entries[word])
	}
}
