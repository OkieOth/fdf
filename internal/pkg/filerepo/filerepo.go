package filerepo

import (
	"sync"
)

type FileRepoEntry struct {
	sourceFile string
	duplicates []string
}

type FileRepo struct {
	mutex sync.RWMutex
	repo  map[string]FileRepoEntry
}

func NewFileRepo() *FileRepo {
	return &FileRepo{
		repo: make(map[string]FileRepoEntry),
	}
}
