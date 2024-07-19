package main

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/hasbrovish/Caching/simplecache/cache"
)

func main() {
	c := cache.NewCache()
	cacheFile := "cache2.db"

	// Load the cache from the file
	err := c.LoadCache(cacheFile)
	if err != nil {
		fmt.Printf("Error loading cache: %v\n", err)
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "cache-cli> ",
		EOFPrompt:       "exit",
		InterruptPrompt: "^C",
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    nil,
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	fmt.Println("Interactive Cache CLI")
	fmt.Println("Available commands: set, get, delete, exit")

	for {
		line, err := rl.Readline()
		if err != nil { // Handle exit signals
			if err == readline.ErrInterrupt {
				if len(line) == 0 {
					break
				} else {
					continue
				}
			} else if err == io.EOF {
				break
			}
			continue
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		cmd := parts[0]
		args := parts[1:]

		switch cmd {
		case "set":
			if len(args) != 3 {
				fmt.Println("Usage: set <key> <value> <duration>")
				continue
			}
			key := args[0]
			value := args[1]
			duration, err := time.ParseDuration(args[2])
			if err != nil {
				fmt.Println("Invalid duration format. Use Go duration format, e.g., '10s', '1m'")
				continue
			}
			c.SetCache(key, value, duration)
			fmt.Printf("Set key %s with value %s for %s\n", key, value, duration)
			err = c.SaveToFile(cacheFile)
			if err != nil {
				fmt.Printf("Error saving cache: %v\n", err)
			}

		case "get":
			if len(args) != 1 {
				fmt.Println("Usage: get <key>")
				continue
			}
			key := args[0]
			value, found := c.GetCache(key)
			if found {
				fmt.Printf("Found key %s with value %v\n", key, value)
			} else {
				fmt.Printf("Cannot find key %s in cache\n", key)
			}

		case "delete":
			if len(args) != 1 {
				fmt.Println("Usage: delete <key>")
				continue
			}
			key := args[0]
			c.DeleteCache(key)
			_, found := c.GetCache(key)
			if found {
				fmt.Println("Failed to delete key", key)
			} else {
				fmt.Printf("Key %s is successfully deleted\n", key)
			}
			err = c.SaveToFile(cacheFile)
			if err != nil {
				fmt.Printf("Error saving cache: %v\n", err)
			}

		case "exit":
			return

		default:
			fmt.Println("Unknown command:", cmd)
			fmt.Println("Available commands: set, get, delete, exit")
		}
	}
}
