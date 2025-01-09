package implhelper

import (
	"sync"
)

type FileRepoEntry struct {
	SourceFile string
	Duplicates []string
}

func NewFileRepoEntry(f string) FileRepoEntry {
	return FileRepoEntry{
		SourceFile: f,
		Duplicates: make([]string, 0),
	}
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

func (f *FileRepo) InitFromSource(sourceDir string, blackList []string, whiteList []string) error {
	resp := make(chan TraversResponse)
	var err error
	f.mutex.Lock()
	defer f.mutex.Unlock()
	go TraversDir(sourceDir, blackList, whiteList, resp)
	for r := range resp {
		if r.err != nil {
			err = r.err
		} else {
			if v, exist := f.repo[r.md5]; exist {
				// Entry with that md5 is already there
				v.Duplicates = append(v.Duplicates, r.file)
			} else {
				// save new entry
				f.repo[r.md5] = NewFileRepoEntry(r.file)
			}
		}
	}
	return err
}

func (f *FileRepo) HasEntry(md5Str string) bool {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	_, e := f.repo[md5Str]
	return e
}

func (f *FileRepo) GetEntry(md5Str string) (FileRepoEntry, bool) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	v, e := f.repo[md5Str]
	return v, e
}

func (f *FileRepo) Size() int {
	return len(f.repo)
}

func (f *FileRepo) Repo() map[string]FileRepoEntry {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	ret := make(map[string]FileRepoEntry)
	for k, v := range f.repo {
		ret[k] = v
	}
	return ret
}
