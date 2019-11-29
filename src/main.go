package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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
func ListFiles(url *url.URL) ([]FileInfo, error) {
	res, err := http.Get(url.String())
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

// InteractiveSearch searches for files matching a name string.
// It will prompt the user to select one of the matching files.
func InteractiveSearch(files []FileInfo, matchName string) string {
	var filtered []FileInfo
	for _, f := range files {
		if strings.Contains(f.Name, matchName) {
			filtered = append(filtered, f)
		}
	}
	if len(filtered) == 0 {
		fmt.Println("No file matches.")
		return ""
	}
	fmt.Printf("\nPick a file:\n")
	for i, f := range filtered {
		fmt.Printf("%d) %s\n", i, f.Name)
	}
	fmt.Printf("-----\n> ")
	var fileIdx int
	fmt.Scanln(&fileIdx)

	if fileIdx < 0 || fileIdx >= len(filtered) {
		fmt.Printf("Invalid file index: %d\n", fileIdx)
		return ""
	}

	selectedFile := filtered[fileIdx].Name
	fmt.Printf("File selected:\n%d) %s\n", fileIdx, selectedFile)
	return selectedFile
}

func DownloadFile(dir string, url *url.URL) error {
	res, err := http.Get(url.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	path := filepath.Join(filepath.Dir(dir), filepath.Base(url.Path))
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	return err
}

func main() {
	var rurl string
	var file string
	var exact, highestVersion bool
	flag.StringVar(&rurl, "url", "", "Uniform resource locator.")
	flag.StringVar(&file, "file", "", "File name.")
	flag.BoolVar(&exact, "exact", false, "Find exact match for filename, else return a non-zero exit code.")
	flag.BoolVar(&highestVersion, "highest-version", false, "Find file that matches the filename and has the highest three-number-version.")
	flag.Parse()

	if rurl == "" {
		flag.PrintDefaults()
		return
	}

	furl, err := url.Parse(rurl)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Add protocol schemen if omitted in argument.
	if furl.Scheme == "" {
		furl.Scheme = "http"
	}

	// Test print files.
	files, _ := ListFiles(furl)
	for _, f := range files {
		fmt.Printf("%+v\n", f)
	}

	selected := InteractiveSearch(files, "itkf")
	fileURL, err := url.Parse(furl.String() + "/" + selected)
	if err != nil {
		fmt.Println(err)
		return
	}

	DownloadFile("", fileURL)
}
