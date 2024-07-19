package cmd

import (
	"fmt"
	"log"

	"github.com/hasbrovish/Caching/simplecache/cache"
	"github.com/spf13/cobra"
)

var deleteKey string

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVarP(&deleteKey, "key", "k", "", "Key to delete")
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a key from the cache",
	Run: func(cmd *cobra.Command, args []string) {
		c := cache.NewCache()
		cacheFile := "cache.db"

		err := c.LoadCache(cacheFile)
		if err != nil {
			fmt.Printf("Error loading cache: %v\n", err)
			return
		}

		c.DeleteCache(deleteKey)
		_, found := c.GetCache(deleteKey)
		if found {
			log.Printf("Failed to delete key %s\n", deleteKey)
		} else {
			log.Printf("Key %s is successfully deleted\n", deleteKey)
		}

		err = c.SaveToFile(cacheFile)
		if err != nil {
			fmt.Printf("Error saving cache: %v\n", err)
		}
	},
}
