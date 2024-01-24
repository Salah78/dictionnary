package dictionary

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

const redisAddr = "localhost:6379" 

type Entry struct {
	Word       string `json:"word"` 
	Definition string `json:"definition"`
}

func (e Entry) String() string {
	return e.Definition
}

type Dictionary struct {
	redisConn redis.Conn
}

func New() (*Dictionary, error) {
	conn, err := redis.Dial("tcp", redisAddr)
	if err != nil {
		return nil, err
	}

	d := &Dictionary{redisConn: conn}
	return d, nil
}

func (d *Dictionary) Add(word string, definition string) {
	_, err := d.redisConn.Do("SET", word, definition)
	if err != nil {
		fmt.Println("Error adding entry to Redis:", err)
	}
}

func (d *Dictionary) Get(word string) (Entry, error) {
	definition, err := redis.String(d.redisConn.Do("GET", word))
	if err != nil {
		return Entry{}, fmt.Errorf("word '%s' not found", word)
	}
	return Entry{Word: word, Definition: definition}, nil
}

func (d *Dictionary) Remove(word string) {
	_, err := d.redisConn.Do("DEL", word)
	if err != nil {
		fmt.Println("Error removing entry from Redis:", err)
	}
}

func (d *Dictionary) List() ([]string, map[string]Entry) {
	keys, err := redis.Strings(d.redisConn.Do("KEYS", "*"))
	if err != nil {
		fmt.Println("Error listing entries from Redis:", err)
		return nil, nil
	}

	entries := make(map[string]Entry)
	for _, key := range keys {
		definition, _ := redis.String(d.redisConn.Do("GET", key))
		entries[key] = Entry{Word: key, Definition: definition}
	}

	return keys, entries
}

func (d *Dictionary) SaveToFile(filePath string) error {
	return nil
}

func (d *Dictionary) LoadFromFile(filePath string) error {
	return nil
}
