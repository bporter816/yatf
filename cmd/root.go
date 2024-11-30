package cmd

import (
	"fmt"
	//"io/fs"
	"os"
	//"path/filepath"
	"strings"

	//"github.com/bporter816/yatf/internal"
	"github.com/bporter816/yatf/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "yatf",
	Short: "yatf is an opinionated formatter for Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		//internal.LintTrailingCommas(nil)
		dir, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		var files []string
		for _, e := range entries {
			fmt.Fprintln(os.Stdout, e.Name(), e.IsDir())
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".tf") {
				files = append(files, e.Name())
			}
		}

		l := internal.NewLinter(files)
		l.Lint()

		/*
			filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				if strings.HasSuffix(path, ".tf") {
					fmt.Fprintln(os.Stdout, path)
				}
				return nil
			})
		*/
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
