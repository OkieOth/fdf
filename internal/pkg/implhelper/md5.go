package implhelper

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	//"strings"
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
