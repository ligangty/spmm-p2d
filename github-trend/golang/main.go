package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// run first one
	run(3)
	TIME_FORMAT := "15:04:05"
	INTERVAL := 30 * time.Second
	ticker := time.NewTicker(INTERVAL)
	next := time.Now().Add(INTERVAL)
	fmt.Printf("Next invoke will happen at %v\n\n", next.Format(TIME_FORMAT))
	for t := range ticker.C {
		next := time.Now().Add(INTERVAL)
		fmt.Println("Invoked at ", t.Format(TIME_FORMAT))
		run(3)
		fmt.Printf("Next invoke will happen at %v\n\n", next.Format(TIME_FORMAT))
	}
}

func run(top int) {
	htmlContent := getGithubTrending()
	reposInfo := parseContent(htmlContent, top)
	for _, repoInfo := range reposInfo {
		fmt.Printf("Repo:        %s\n", repoInfo["url"])
		fmt.Printf("Description: %s\n", repoInfo["description"])
		fmt.Printf("Language:    %s\n", repoInfo["language"])
		fmt.Printf("Star:        %s\n", repoInfo["star"])
		fmt.Printf("Fork:        %s\n", repoInfo["fork"])
		fmt.Println()
	}

	gitURL := fmt.Sprintf("https://github.com/%s", reposInfo[0]["url"])
	repoLocalPath := DownloadRepo(gitURL)
	defer os.RemoveAll(repoLocalPath)

	readme := ""
	filepath.Walk(repoLocalPath, func(path string, info os.FileInfo, err error) error {
		pathElems := strings.Split(path, "/")
		fileName := pathElems[len(pathElems)-1]
		if strings.HasPrefix(strings.ToLower(fileName), "readme") {
			readme = path
		}
		return nil
	})

	f, err := os.Open(readme)
	if err != nil {
		fmt.Printf("Can not open file: %s", err)
		os.Exit(0)
	}
	defer f.Close()

	content, _ := ioutil.ReadAll(f)
	fmt.Println(string(content))
	fmt.Print("\n\n")
}
