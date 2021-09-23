package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	run(3)
}

func run(top int) {
	htmlContent := getGithubTrending()
	// htmlContent := getGithubTrendingLocal()
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

	os.RemoveAll(repoLocalPath)
}
