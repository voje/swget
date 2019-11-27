package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type FileInfo struct {
	Name string
	Date string
	// Folder size will be "-"
	Size    string
	Version string
}

// ListFiles takes an url and reads the index.html file provided by the url.
// It returns a []FileInfo array.
func ListFiles(url string) ([]FileInfo, error) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed getting file: %v\n", err)
		return nil, err
	}
	defer res.Body.Close()

	reader := bufio.NewReader(res.Body)
	if err != nil {
		fmt.Printf("Failed reading file: %v\n", err)
	}

	reHyperLink := regexp.MustCompile("^<a")
	reFileName := regexp.MustCompile(">(.*)</a>")
	reVersion := regexp.MustCompile("[0-9]+.[0-9]+.[0-9]+")
	var files []FileInfo
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		// Skip lines that don't represent file links.
		if reHyperLink.Find(line) == nil {
			continue
		}

		// Create and fill new FileInfo.
		var fi FileInfo

		fld := strings.Fields(string(line))
		matches := reFileName.FindStringSubmatch(fld[1])
		if matches == nil {
			continue
		}
		fi.Name = matches[1]
		fi.Date = fmt.Sprintf("%s %s", fld[2], fld[3])
		fi.Size = fld[4]

		// Parse file version.
		matches = reVersion.FindStringSubmatch(fi.Name)
		if matches != nil {
			fi.Version = matches[0]
		}

		files = append(files, fi)
	}
	return files, nil
}

func InteractiveSearch(files []FileInfo, matchName string) {
	var filtered []FileInfo
	for _, f := range files {
		if strings.Contains(f.Name, matchName) {
			filtered = append(filtered, f)
		}
	}
	if len(filtered) == 0 {
		fmt.Println("No file matches.")
		return
	}
	fmt.Println("Pick a file:")
	for i, f := range filtered {
		fmt.Printf("%d) %s\n", i, f.Name)
	}
	// TODO read user input... also figure out what to do with copying structs around.
}

func main() {
	var url string
	var file string
	var exact, highestVersion bool
	flag.StringVar(&url, "url", "", "Uniform resource locator.")
	flag.StringVar(&file, "file", "", "File name.")
	flag.BoolVar(&exact, "exact", false, "Find exact match for filename, else return a non-zero exit code.")
	flag.BoolVar(&highestVersion, "highest-version", false, "Find file that matches the filename and has the highest three-number-version.")
	flag.Parse()

	// url = "http://files.k-vm-repo-server.docker.iskratel.mak"

	// Test print files.
	files, _ := ListFiles(url)
	for _, f := range files {
		fmt.Printf("%+v\n", f)
	}

	InteractiveSearch(files, "itkf")

}
