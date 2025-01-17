package implhelper

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

const MEGA_BYTE = 1024 * 1024

func GetFileSize(filePath string) (int64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()
	fileStat, err := file.Stat()
	if err != nil {
		return 0, fmt.Errorf("error while query file stats: %v", err)
	}
	return fileStat.Size(), nil
}

func GetMd5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}
	// io.Copy(hash, strings.NewReader(filePath))
	md5Sum := hash.Sum(nil)

	return hex.EncodeToString(md5Sum), nil
}

func GetCommonPrefix(s1 string, s2 string) string {
	var ret string
	runeSlice1 := []rune(s1)
	runeSlice2 := []rune(s2)
	l := len(runeSlice1)
	if l2 := len(runeSlice2); l2 < l {
		l = l2
	}
	for i := 0; i < l; i++ {
		r1 := runeSlice1[i]
		r2 := runeSlice2[i]
		if r1 == r2 {
			ret += string(r1)
		} else {
			break
		}
	}
	return ret
}

func AdjustCommonPrefix(prefix string, s string) string {
	var ret string
	runeSlice1 := []rune(prefix)
	runeSlice2 := []rune(s)
	l := len(runeSlice1)
	if l2 := len(runeSlice2); l2 < l {
		l = l2
	}
	for i := 0; i < l; i++ {
		r1 := runeSlice1[i]
		r2 := runeSlice2[i]
		if r1 == r2 {
			ret += string(r1)
		} else {
			break
		}
	}
	return ret
}

func GetMaxPathPrefixLen(repo map[string]FileRepoEntry) int {
	var ret int
	var prefix string
	prefixNotSet := true
	for _, v := range repo {
		if len(v.Duplicates) == 0 {
			continue
		}
		if !prefixNotSet {
			prefix = AdjustCommonPrefix(prefix, v.SourceFile)
		}
		for i := 1; i < len(v.Duplicates); i++ {
			if prefixNotSet {
				prefix = GetCommonPrefix(v.SourceFile, v.Duplicates[i])
				prefixNotSet = false
			} else {
				prefix = AdjustCommonPrefix(prefix, v.Duplicates[i])
			}
		}
	}
	ret = len(prefix)
	return ret
}

func GetMaxTextLen(repo map[string]FileRepoEntry) int {
	var ret int
	for _, v := range repo {
		if len(v.Duplicates) == 0 {
			continue
		}
		l := len([]rune(v.SourceFile))
		if l > ret {
			ret = l
		}
		for i := 1; i < len(v.Duplicates); i++ {
			l := len([]rune(v.Duplicates[i]))
			if l > ret {
				ret = l
			}
		}
	}
	return ret
}
