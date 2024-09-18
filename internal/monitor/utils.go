package monitor

import "os"

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		// 如果返回错误，检查是否是“文件或目录不存在”的错误
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
