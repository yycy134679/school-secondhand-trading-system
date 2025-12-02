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
	// 允许的图片类型
	AllowedImageTypes = "jpg,jpeg,png"
)

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

// SaveImage 保存图片到存储目录并返回可访问的URL
func SaveImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	// 校验文件大小
	if header.Size > MaxFileSize {
		return "", errors.New("文件大小超过限制，最大允许2MB")
	}

	// 校验文件类型
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == ".jpeg" {
		ext = ".jpg"
	}
	if ext != ".jpg" && ext != ".png" {
		return "", errors.New("不支持的图片格式，仅支持JPG和PNG")
	}

	// 生成唯一的文件名
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	filename := fmt.Sprintf("%d%s", timestamp, ext)

	// 创建年月子目录
	now := time.Now()
	subDir := fmt.Sprintf("%d/%02d", now.Year(), now.Month())
	fullDir := filepath.Join(FileStorageDir, subDir)

	// 确保子目录存在
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	// 创建目标文件
	dst := filepath.Join(fullDir, filename)
	dstFile, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer dstFile.Close()

	// 将上传的文件内容复制到目标文件
	if _, err = io.Copy(dstFile, file); err != nil {
		return "", fmt.Errorf("保存文件失败: %w", err)
	}

	// 生成可访问的URL
	// 在Windows环境中，需要将路径分隔符从\转换为/
	relativePath := filepath.Join(subDir, filename)
	relativePath = strings.ReplaceAll(relativePath, "\\", "/")
	url := fmt.Sprintf("%s/uploads/%s", BaseURL, relativePath)

	return url, nil
}

// DeleteFile 删除文件
func DeleteFile(filePath string) error {
	// 从URL中提取文件路径
	// 假设URL格式为: http://localhost:8080/uploads/2023/12/123456.jpg
	parts := strings.Split(filePath, "/uploads/")
	if len(parts) < 2 {
		return errors.New("无效的文件路径")
	}

	// 构建完整的本地文件路径
	// 确保正确处理Windows路径分隔符
	relativePath := parts[1]
	relativePath = strings.ReplaceAll(relativePath, "/", string(filepath.Separator))
	fullPath := filepath.Join(FileStorageDir, relativePath)

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return errors.New("文件不存在")
	}

	// 删除文件
	return os.Remove(fullPath)
}
