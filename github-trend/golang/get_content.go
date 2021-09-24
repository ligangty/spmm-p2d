package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const GITHUB_TREANDING = "https://github.com/trending"

func getGithubTrending() string {
	fmt.Printf("Start requesting %s", GITHUB_TREANDING)
	response, err := http.Get(GITHUB_TREANDING)
	if err != nil {
		fmt.Printf("Error to access github trending page due to %s", err)
		os.Exit(1)
	}
	defer response.Body.Close()
	htmlContentByte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error to get trending page content due to %s", err)
		os.Exit(1)
	}
	htmlContent := string(htmlContentByte)
	return htmlContent
}
