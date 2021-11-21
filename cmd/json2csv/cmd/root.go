/*
Copyright © 2021 Zaba505

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Zaba505/json2csv"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "json2csv [flags] FILE|-",
	Short: "Format JSON objects as CSV",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		outputName, err := cmd.Flags().GetString("output")
		if err != nil {
			panic(err)
		}

		colFieldMap, err := cmd.Flags().GetStringToString("column")
		if err != nil {
			panic(err)
		}

		fieldColMap := make(map[string]int, len(colFieldMap))
		for colIdxStr, fieldName := range colFieldMap {
			fieldColMap[fieldName], err = strconv.Atoi(colIdxStr)
			if err != nil {
				panic(err)
			}
		}

		skipColumnTitles, err := cmd.Flags().GetBool("no-column-titles")
		if err != nil {
			panic(err)
		}

		src := os.Stdin
		if args[0] != "-" {
			path, err := filepath.Abs(args[0])
			if err != nil {
				panic(err)
			}

			src, err = os.Open(path)
			if err != nil {
				panic(err)
			}
		}

		out := os.Stdout
		if strings.TrimSpace(outputName) != "" {
			out, err = os.Create(outputName)
			if err != nil {
				panic(err)
			}
		}

		opts := make([]json2csv.Option, 0, len(fieldColMap)+1)
		if skipColumnTitles {
			opts = append(opts, json2csv.SkipColumnTitles())
		}
		for field, col := range fieldColMap {
			opts = append(opts, json2csv.MapFieldToColumn(field, col))
		}

		err = json2csv.Convert(out, src, opts...)
		if err != nil {
			panic(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().StringP("output", "o", "", "Filename to write JSON objects as CSV to.")
	rootCmd.Flags().StringToStringP("column", "c", map[string]string{}, "Map column indexes to JSON field names.")
	rootCmd.Flags().Bool("no-column-titles", false, "Don't write JSON field names as first row of values")
}
