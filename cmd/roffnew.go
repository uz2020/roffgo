/*
Copyright © 2022 Zhj Rong <rongzhj2020@163.com>

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
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// roffnewCmd represents the roffnew command
var roffnewCmd = &cobra.Command{
	Use:   "new",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: roffnewRun,
}

func roffnewRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("post name pls")
		return
	}

	baseDir := "content"
	for _, name := range args {
		dirName := baseDir + "/" + filepath.Dir(name)
		pathName := baseDir + "/" + name + ".ms"
		err := os.MkdirAll(dirName, 0755)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = os.Create(pathName)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	fmt.Println("roffnew called")
}

func init() {
	rootCmd.AddCommand(roffnewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// roffnewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// roffnewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
