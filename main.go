package main

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/joshdcuneo/gitstats/gitstats"
)

// TODO implement some useful aggregations over commits
// TODO expose these things via a cobra CLI
// TODO expose these things via a gin api

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

func main() {
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: "https://github.com/go-git/go-git",
	})

	CheckIfError(err)

	ref, err := r.Head()
	CheckIfError(err)

	commitIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	CheckIfError(err)

	authorIter := gitstats.AuthorIterFromCommitIter(commitIter)
	authorIter.SortByCommitCount()

	authorIter.ForEach(func(a gitstats.Author) error {
		fmt.Printf("author: %s, commits: %d\n", a.Email, a.CommitCount)
		return nil
	})

	fmt.Printf("total commits: %d\n", authorIter.CommitCount())
}
