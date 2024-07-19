package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hasbrovish/Caching/simplecache/cache"
)

func main() {
	cacheFile := "cache.db"
	c := cache.NewCache()

	err := c.LoadCache(cacheFile)
	if err != nil {
		fmt.Printf("Error in loading cache %v \n", err)
	}
	// Here we will use flag for creating set, get, and delete commands
	// When we are making map
	setCmd := flag.NewFlagSet("set", flag.ExitOnError)
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)

	setKey := setCmd.String("key", "", "Key to set")
	setValue := setCmd.String("value", "", "Value to set")
	setDuration := setCmd.Int("duration", 10, "Duration to set in seconds")

	getKey := getCmd.String("key", "", "Key to get")

	deleteKey := deleteCmd.String("key", "", "Key to delete")

	if len(os.Args) < 2 {
		fmt.Println("Expected 'set', 'get' or 'delete' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "set":
		setCmd.Parse(os.Args[2:])
		if *setKey == "" || *setValue == "" {
			fmt.Println("Key and Value should be provided")
			os.Exit(1)
		}
		c.SetCache(*setKey, *setValue, time.Duration(*setDuration)*time.Second)
		fmt.Printf("Set key %s with value %s for %d seconds\n", *setKey, *setValue, *setDuration)
		err := c.SaveToFile(cacheFile)
		if err != nil {
			fmt.Printf("Error saving cache: %v\n", err)
		}

	case "get":
		getCmd.Parse(os.Args[2:])
		if *getKey == "" {
			fmt.Println("Key should be provided")
			os.Exit(1)
		}
		value, found := c.GetCache(*getKey)
		if found {
			fmt.Printf("Found key %s with value %v\n", *getKey, value)
		} else {
			fmt.Printf("Cannot find key %s in cache\n", *getKey)
		}
	case "delete":
		deleteCmd.Parse(os.Args[2:])
		if *deleteKey == "" {
			fmt.Println("Key is required to delete")
			os.Exit(1)
		}
		c.DeleteCache(*deleteKey)
		_, found := c.GetCache(*deleteKey)
		if found {
			panic("Key is not deleted")
		} else {
			log.Printf("Key %s is successfully deleted\n", *deleteKey)
		}
		// Save the cache to the file
		err := c.SaveToFile(cacheFile)
		if err != nil {
			fmt.Printf("Error saving cache: %v\n", err)
		}

	default:
		fmt.Println("Expected 'set', 'get' or 'delete' subcommands")
		os.Exit(1)
	}
}
