package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
)

func DownloadRepo(gitURL string) string {
	dir := fmt.Sprintf("%s/%s", getTempDir(), getGitDir(gitURL))
	return goGit(gitURL, dir)
}

func getGitDir(gitURL string) string {
	segs := strings.Split(gitURL, "/")
	last := segs[len(segs)-1]
	if strings.Contains(last, ".git") {
		splits := strings.Split(last, ".")
		return splits[0]
	}
	return last
}

func getRepo(gitURL, directory string, clone bool) (r *git.Repository, err error) {
	if clone {
		fmt.Printf("Start cloning %s to %s\n", gitURL, directory)
		r, err = git.PlainClone(directory, false, &git.CloneOptions{
			URL:               gitURL,
			Progress:          os.Stdout,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})
	} else {
		fmt.Printf("Open existed repo %s\n", directory)
		r, err = git.PlainOpen(directory)
	}
	return r, err
}

func goGit(gitURL, directory string) string {
	clone := true
	if fileExists(directory) {
		clone = false
	}

	r, err := getRepo(gitURL, directory, clone)
	if err != nil {
		fmt.Printf("Get repo failed due to error: %s", err)
		os.Exit(1)
	}

	// Updating heads
	fetchUpdates(r)

	showHEAD(r)

	return directory
}

func showHEAD(r *git.Repository) {
	// Show HEAD
	ref, err := r.Head()
	checkIfError(err)
	fmt.Printf("Now HEAD is %s\n", ref.Hash())

	commitIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	checkIfError(err)

	commit, err := commitIter.Next()
	checkIfError(err)
	hash := commit.Hash.String()
	line := strings.Split(commit.Message, "\n")
	fmt.Println(hash[:7], line[0])
}

// Fetching updates...
func fetchUpdates(r *git.Repository) {
	fmt.Printf("Fetching Refs....\n")
	err := r.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
		Progress: os.Stdout,
	})
	if err != git.NoErrAlreadyUpToDate {
		checkIfError(err)
	}
}

func getTempDir() string {
	temp := os.Getenv("TMPDIR")
	if strings.TrimSpace(temp) == "" {
		temp = "/tmp"
	}
	return temp
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func checkIfError(err error) {
	if err != nil {
		fmt.Printf("Git operations failed due to error: %s", err)
	}
}
