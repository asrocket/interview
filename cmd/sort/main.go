package main

import (
	"asrocket/interview/sort"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var root = &cobra.Command{
	Use:   "sort <filepath>",
	Short: "Постоковая сортировка файла",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("need file path argument")
		}

		file, err := os.Open(args[0])
		if err != nil {
			return err
		}
		defer file.Close()

		newFile, err := os.Create(args[0] + "_sorted")
		if err != nil {
			return err
		}
		defer newFile.Close()

		sorter := sort.NewLineSorter(file, newFile)
		return sorter.Sort()
	},
}

func main() {
	if err := root.Execute(); err != nil {
		panic(err)
	}
}
