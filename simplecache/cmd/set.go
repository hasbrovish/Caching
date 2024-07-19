package cmd

import (
	"fmt"
	"time"

	"github.com/hasbrovish/Caching/simplecache/cache"
	"github.com/spf13/cobra"
)

var setKey string
var setValue string
var setDuration string

func init() {
	rootCmd.AddCommand(setCmd)
	setCmd.Flags().StringVarP(&setKey, "key", "k", "", "Key to set")
	setCmd.Flags().StringVarP(&setValue, "value", "v", "", "Value to set")
	setCmd.Flags().StringVarP(&setDuration, "duration", "d", "10s", "Duration to set")
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a key in the cache",
	Run: func(cmd *cobra.Command, args []string) {
		c := cache.NewCache()
		cacheFile := "cache.db"

		err := c.LoadCache(cacheFile)
		if err != nil {
			fmt.Printf("Error loading cache: %v\n", err)
			return
		}

		duration, err := time.ParseDuration(setDuration)
		if err != nil {
			fmt.Println("Invalid duration format. Use Go duration format, e.g., '10s', '1m'")
			return
		}

		c.SetCache(setKey, setValue, duration)
		fmt.Printf("Set key %s with value %s for %s\n", setKey, setValue, duration)

		err = c.SaveToFile(cacheFile)
		if err != nil {
			fmt.Printf("Error saving cache: %v\n", err)
		}
	},
}
