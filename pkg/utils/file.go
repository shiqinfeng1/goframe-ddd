package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// nextFilename 在文件重名时，通过在原有文件名后加数字,生成下一个文件名，
func NextFileName(path, newDir string) string {
	newPath := filepath.Join(newDir, filepath.Base(path))
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		if _, err := os.Stat(newPath + ".downloading"); os.IsNotExist(err) {
			return newPath
		}
	}

	ext := filepath.Ext(newPath)
	base := strings.TrimSuffix(filepath.Base(newPath), ext)

	for i := 1; ; i++ {
		newPath := filepath.Join(newDir, fmt.Sprintf("%s(%d)%s", base, i, ext))
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			if _, err := os.Stat(newPath + ".downloading"); os.IsNotExist(err) {
				return newPath
			}
		}
	}
}
