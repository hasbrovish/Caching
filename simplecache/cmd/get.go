package cmd

import (
	"fmt"

	"github.com/hasbrovish/Caching/simplecache/cache"
	"github.com/spf13/cobra"
)

var getKey string

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVarP(&getKey, "key", "k", "", "Key to get")
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a key from the cache",
	Run: func(cmd *cobra.Command, args []string) {
		c := cache.NewCache()
		cacheFile := "cache.db"

		err := c.LoadCache(cacheFile)
		if err != nil {
			fmt.Printf("Error loading cache: %v\n", err)
			return
		}

		value, found := c.GetCache(getKey)
		if found {
			fmt.Printf("Found key %s with value %v\n", getKey, value)
		} else {
			fmt.Printf("Cannot find key %s in cache\n", getKey)
		}
	},
}
