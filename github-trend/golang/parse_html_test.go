package main

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseContent(t *testing.T) {
	Convey("TestParseContent", t, func() {
		htmlContent := getGithubTrendingLocal()
		reposInfo := parseContent(htmlContent, 3)
		So(len(reposInfo), ShouldEqual, 3)
		repoInfo := reposInfo[0]
		So(len(repoInfo), ShouldEqual, 5)
		So(repoInfo["url"], ShouldEqual, "public-apis/public-apis")
		So(repoInfo["description"], ShouldEqual, "A collective list of free APIs")
		So(repoInfo["language"], ShouldEqual, "Python")
		So(repoInfo["star"], ShouldEqual, "158,252")
		So(repoInfo["fork"], ShouldEqual, "17,923")
	})
}

func getGithubTrendingLocal() string {
	f, _ := os.Open("./test.html")
	defer f.Close()
	content, _ := ioutil.ReadAll(f)
	return string(content)
}
