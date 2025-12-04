// Package util 提供通用工具函数
package util

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 文件存储配置常量
const (
	// 最大文件大小：2MB
	MaxFileSize = 2 * 1024 * 1024
)

// 支持的图片格式
var allowedImageFormats = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

// 全局变量存储目录，可从环境变量或配置中获取
var (
	// 文件存储根目录
	FileStorageDir string
	// 可访问的基础URL
	BaseURL string
)

// 初始化文件存储目录
func init() {
	// 默认存储目录
	FileStorageDir = getEnvOrDefault("FILE_STORAGE_DIR", "./uploads")
	// 默认基础URL
	BaseURL = getEnvOrDefault("BASE_URL", "http://localhost:8080")

	// 确保存储目录存在
	if err := os.MkdirAll(FileStorageDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create storage directory: %v", err))
	}
}

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// SaveImage 将图片保存到FILE_STORAGE_DIR下并返回可访问URL
// 参数：
// - file: 上传的图片文件
// - header: 文件头信息
// 返回值：
// - url: 可访问的图片URL
// - err: 错误信息
func SaveImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	// 1. 检查文件大小
	if header.Size > MaxFileSize {
		return "", errors.New("文件大小超过限制，最大支持2MB")
	}

	// 2. 检查文件格式
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == ".jpeg" {
		ext = ".jpg"
	}
	if !allowedImageFormats[ext] {
		return "", errors.New("仅支持JPG和PNG格式的图片")
	}

	// 3. 生成唯一文件名
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	fileName := fmt.Sprintf("%d%s", timestamp, ext)

	// 4. 创建年月子目录
	now := time.Now()
	subDir := fmt.Sprintf("%d/%02d", now.Year(), now.Month())
	fullDir := filepath.Join(FileStorageDir, subDir)

	// 5. 确保存储目录存在
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return "", fmt.Errorf("创建存储目录失败: %w", err)
	}

	// 6. 创建目标文件
	filePath := filepath.Join(fullDir, fileName)
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	// 7. 复制文件内容
	if _, err = io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("保存文件失败: %w", err)
	}

	// 8. 构建可访问的URL路径
	// 在Windows环境中，需要将路径分隔符从\转换为/
	relativePath := filepath.Join(subDir, fileName)
	relativePath = strings.ReplaceAll(relativePath, "\\", "/")
	url := fmt.Sprintf("%s/uploads/%s", BaseURL, relativePath)

	return url, nil
}

// DeleteFile 删除指定路径的文件
// 参数：
// - filePath: 文件路径
// 返回值：
// - err: 错误信息
func DeleteFile(filePath string) error {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // 文件不存在视为删除成功
	}

	// 删除文件
	return os.Remove(filePath)
}

// GetFileSize 获取文件大小
// 参数：
// - file: 上传的文件
// 返回值：
// - size: 文件大小（字节）
// - err: 错误信息
func GetFileSize(file multipart.File) (int64, error) {
	// 获取文件当前位置
	currentPos, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	// 移动到文件末尾获取大小
	size, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	// 恢复文件位置
	_, err = file.Seek(currentPos, io.SeekStart)
	if err != nil {
		return 0, err
	}

	return size, nil
}