package main

import (
	"asrocket/interview/generator"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var root = &cobra.Command{
	Use:   "generate file_name line_count max_line_len",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return fmt.Errorf("need file_name, line_count and max_line_len arguments")
		}

		count, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}

		length, err := strconv.Atoi(args[2])
		if err != nil {
			return err
		}

		file, err := os.Create(args[0])
		if err != nil {
			return err
		}
		defer file.Close()

		gen := generator.NewLineGenerator(length)

		for i := 0; i <= count; i++ {
			if _, err := file.Write(gen.Generate()); err != nil {
				return err
			}
		}
		return nil
	},
}

func main() {
	if err := root.Execute(); err != nil {
		panic(err)
	}
}
