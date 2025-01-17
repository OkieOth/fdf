package implhelper

import (
	"sync"

	"github.com/okieoth/fdf/internal/pkg/progressbar"
)

type FileRepoEntry struct {
	SourceFile string   `json:"sourceFile,omitempty"`
	Duplicates []string `json:"duplicates,omitempty"`
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

func (f *FileRepo) InitFromSource(sourceDir string, blackList []string, whiteList []string, noProgress bool) {
	resp := make(chan TraversResponse, 100000)
	// not needed here, because aligned over the channel
	f.mutex.Lock()
	defer f.mutex.Unlock()
	go TraversDir(sourceDir, blackList, whiteList, resp, false, true)
	for r := range resp {
		if r.err != nil {
			// FIXME
		} else {
			if v, exist := f.repo[r.md5]; exist {
				// Entry with that md5 is already there
				v.Duplicates = append(v.Duplicates, r.file)
				f.repo[r.md5] = v
			} else {
				// save new entry
				f.repo[r.md5] = NewFileRepoEntry(r.file)
			}
		}

		progressbar.ProgressOne()
	}
}

func (f *FileRepo) CheckForDuplicateAndAddInCase(md5Str string, fileName string) {
	if f.HasEntry(md5Str) {
		f.SetEntry(md5Str, fileName)
	}
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

func (f *FileRepo) SetEntry(md5Str string, fileName string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	if v, e := f.repo[md5Str]; e {
		if v.SourceFile != fileName {
			v.Duplicates = append(v.Duplicates, fileName)
			f.repo[md5Str] = v
		}
	}
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

func (f *FileRepo) HasDuplicates() bool {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	for _, v := range f.repo {
		if len(v.Duplicates) > 0 {
			return true
		}
	}
	return false
}
