package gitstats

import (
	"sort"

	"github.com/go-git/go-git/v5/plumbing/object"
)

// TODO implement our own commit and commit iters

type Author struct {
	Name  string
	Email string
	// TODO store commits
	CommitCount int
}

type AuthorIter struct {
	authors []Author
}

func NewAuthorIter(authors []Author) AuthorIter {
	return AuthorIter{authors: authors}
}

func AuthorIterFromCommitIter(commitIter object.CommitIter) AuthorIter {
	authorsByEmail := make(map[string]Author)
	commitIter.ForEach(func(c *object.Commit) error {
		authorsByEmail[c.Author.Email] = Author{
			Name:        c.Author.Name,
			Email:       c.Author.Email,
			CommitCount: authorsByEmail[c.Author.Email].CommitCount + 1,
		}
		return nil
	})

	authors := make([]Author, 0, len(authorsByEmail))
	for _, a := range authorsByEmail {
		authors = append(authors, a)
	}

	return NewAuthorIter(authors)
}

func (ai *AuthorIter) SortByCommitCount() {
	sort.SliceStable(ai.authors, func(i, j int) bool {
		return ai.authors[i].CommitCount > ai.authors[j].CommitCount
	})
}

// TODO this should be on the commit iter once we implement it
func (ai *AuthorIter) CommitCount() int {
	count := 0
	for _, a := range ai.authors {
		count += a.CommitCount
	}
	return count
}

func (ai *AuthorIter) ForEach(fn func(a Author) error) error {
	for _, a := range ai.authors {
		if err := fn(a); err != nil {
			return err
		}
	}
	return nil
}
