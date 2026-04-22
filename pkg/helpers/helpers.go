// Package helpers 存放辅助方法
package helpers

import (
	"crypto/rand"
	"fmt"
	"io"
	"math"
	mathrand "math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"
)

// Empty 类似于 PHP 的 empty() 函数
func Empty(val interface{}) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
}

// MicrosecondsStr 将 time.Duration 类型（nano seconds 为单位）
// 输出为小数点后 3 位的 ms （microsecond 毫秒，千分之一秒）
func MicrosecondsStr(elapsed time.Duration) string {
	return fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)
}

// RandomNumber 生成长度为 length 随机数字字符串
func RandomNumber(length int) string {
	table := [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

// FirstElement 安全地获取 args[0]，避免 panic: runtime error: index out of range
func FirstElement(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

// RandomString 生成长度为 length 的随机字符串
func RandomString(length int) string {
	mathrand.Seed(time.Now().UnixNano())
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[mathrand.Intn(len(letters))]
	}
	return string(b)
}

// 判断文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	// 如果是其他错误（例如权限问题），默认认为文件存在
	return true
}

func StringContains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val || strings.HasPrefix(val, item) {
			return true
		}
	}
	return false
}

func IsPathAllowed(requestPath string, permissionPath string) bool {
	// 通配符 * 表示匹配任意结尾
	if strings.HasSuffix(permissionPath, "*") {
		prefix := strings.TrimSuffix(permissionPath, "*")
		return strings.HasPrefix(requestPath, prefix)
	}
	// 严格匹配
	return requestPath == permissionPath
}

func GetFileExt(fileName string) string {
	fileOriExt := filepath.Ext(fileName) // 获取文件扩展名 这里包含了 .
	fileExt := fileOriExt[1:]
	return fileExt
}

func GetRandomNumber(digits int) int {
	if digits <= 0 {
		return 0
	}

	// 计算最小值和最大值
	min := int(math.Pow10(digits - 1))
	max := int(math.Pow10(digits)) - 1

	// 初始化随机数生成器的种子
	mathrand.Seed(time.Now().UnixNano())

	// 生成指定位数的随机整数
	return mathrand.Intn(max-min+1) + min
}

// 查找字符串元素在切片中的索引，如果找不到则返回 -1
func FindElement(slice []string, element string) int {
	for i, v := range slice {
		if v == element {
			return i
		}
	}
	return -1
}

type Router struct {
	mu          sync.RWMutex
	ignorePaths []string
}

// 初始值直接写死
var globalRouter = &Router{
	ignorePaths: []string{"/admin/auth/login", "/admin/auth/captcha", "/admin/upload", "/admin/version", "/admin/test"}, // 忽略的路径无需验证,
}

// 只暴露方法，不暴露结构体
func GetIgnorePaths() []string {
	globalRouter.mu.RLock()
	defer globalRouter.mu.RUnlock()

	res := make([]string, len(globalRouter.ignorePaths))
	copy(res, globalRouter.ignorePaths)
	return res
}

func AppendIgnorePaths(v []string) {
	globalRouter.mu.Lock()
	defer globalRouter.mu.Unlock()

	globalRouter.ignorePaths = append(globalRouter.ignorePaths, v...)
}
