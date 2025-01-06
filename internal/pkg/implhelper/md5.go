package implhelper

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func GetMd5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("Error opening file: %v\n", err)
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("Error reading file: %v\n", err)
	}
	md5Sum := hash.Sum(nil)

	return hex.EncodeToString(md5Sum), nil
}
