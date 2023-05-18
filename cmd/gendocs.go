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
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// gendocsCmd represents the gendocs command
var gendocsCmd = &cobra.Command{
	Use:   "gen",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: gendocsRun,
}

func readFirstLine(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if text != "" {
			return text
		}
	}
	return ""
}

func orgToTxt(path string) (out []byte, err error) {
	extCmd := `BEGIN{FS=OFS=" "} {gsub(/\*/, "\t", $1)} 1`
	out, err = exec.Command("awk", extCmd, path).Output()
	return
}

func txtToTxt(path string) (out []byte, err error) {
	out, err = ioutil.ReadFile(path)
	return
}

func genIndex(path string) []byte {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err, path)
		return nil
	}

	content := "<div style='margin: 10% 20%'><ul style='list-style-type: none; font-size: 1.5em;'>"
	for _, file := range files {
		name := file.Name()
		link := name
		ext := filepath.Ext(name)
		if name == "index.html" {
			continue
		}
		if file.IsDir() {
			name += "/"
		} else if ext == ".txt" {
			fp := path + "/" + name
			name = readFirstLine(fp)
		} else {
			continue
		}
		content += fmt.Sprintf("<li style='margin: 10px 0;'><a href='./%s'>%s</a></li>\n", link, name)
	}
	content += "</ul></div>"
	return []byte(content)
}

func gendocs(path string) ([]byte, error) {
	var err error
	out := []byte{}

	ext := filepath.Ext(path)
	switch ext {
	case ".ms", ".mm":
		extCmd := fmt.Sprintf("-%s", ext[1:])
		out, err = exec.Command("groff", extCmd, "-Tutf8", "-k", path).Output()
	case ".org":
		out, err = orgToTxt(path)
	case ".txt":
		out, err = txtToTxt(path)
	}
	if err != nil {
		fmt.Println("------------- err", err, path)
	}
	return out, err
}

func gendocsRun(cmd *cobra.Command, args []string) {
	contentDir := "content"

	err := os.Mkdir(contentDir, 0755)
	if err != nil {
		panic(err)
	}
	err = os.Mkdir("docs", 0755)
	if err != nil {
		panic(err)
	}

	dirLen := len(contentDir)
	err = filepath.Walk(contentDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fileName := "docs" + path[dirLen:]
			ext := filepath.Ext(fileName)
			dirName := filepath.Dir(fileName)
			var out []byte
			if !info.IsDir() {
				fileName = fileName[0:len(fileName)-len(ext)] + ".txt"
				out, err = gendocs(path)
				if err != nil {
					fmt.Println("------------", err)
					return err
				}
			}

			if len(out) == 0 {
				return nil
			}

			os.MkdirAll(dirName, 0755)

			err = os.WriteFile(fileName, out, 0644)
			if err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		fmt.Println(err)
	}

	dir := "docs"
	err = filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				return nil
			}

			out := genIndex(path)
			path += "/index.html"
			return os.WriteFile(path, out, 0644)
		})

	if err != nil {
		fmt.Println(err)
	}
}

func init() {
	rootCmd.AddCommand(gendocsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gendocsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gendocsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
