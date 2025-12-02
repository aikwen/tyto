package hash

import (
	"github.com/cespare/xxhash/v2"
	"os"
	"io"
)

// Hash64 计算给定数据的 xxHash64 哈希值.
func Hash64(data []byte) uint64 {
	return xxhash.Sum64(data)
}

// HashFile 计算指定文件的 xxHash64 哈希值.
func HashFile(filePath string) (uint64, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return 0, err
    }
    defer file.Close()

    hasher := xxhash.New()

    // io.Copy 会自动处理缓冲和循环，将文件内容全部写入 hasher
    // 它会返回写入的字节数和遇到的错误
    if _, err := io.Copy(hasher, file); err != nil {
        return 0, err
    }

    return hasher.Sum64(), nil
}
