package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"reflect"
	"strconv"
	"time"
	"unsafe"

	"github.com/dgrijalva/jwt-go"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

// redis Key
var (
	SEND_CODE_KEY = "send_code:"
	CHAT_KEY      = "chat:"
)

func Int64ToUint64(i int64) uint64 {
	return uint64(i)
}

func StringToInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func StringToUint64(s string) uint64 {
	i, _ := strconv.ParseUint(s, 10, 64)
	return i
}

func Uint64ToString(i uint64) string {
	return fmt.Sprintf("%d", i)
}

// md5
func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// 6位随机数
func GetRandomNum() string {
	return fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
}

func GetRandomString(length int) string {
	// 字符集
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const charsetLength = len(charset)

	// 初始化随机种子
	rand.NewSource(time.Now().UnixNano())

	// 生成随机字符串
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(charsetLength)]
	}
	return string(result)
}

func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func getFastReqClient() *fasthttp.Client {
	reqClient := &fasthttp.Client{
		ReadTimeout:                   time.Second * 20,
		WriteTimeout:                  time.Second * 20,
		MaxIdleConnDuration:           time.Second * 20,
		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}
	return reqClient
}

func FastWithDo(result any, method, uri string, params map[string]string, body any, headers map[string]string) error {
	client := getFastReqClient()
	req, resp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()
	req.Header.SetMethod(method)
	req.SetRequestURI(uri)
	// 设置参数
	var arg fasthttp.Args
	for k, v := range params {
		arg.Add(k, v)
	}
	req.URI().SetQueryString(arg.String())
	// 设置header信息
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	// 设置body
	if body != nil {
		marshal, _ := json.Marshal(body)
		fmt.Println(string(marshal))
		req.SetBody(marshal)
	}
	// 设置Cookie信息
	// req.Header.SetCookie("key", "val")
	// 发起请求
	if err := client.Do(req, resp); err != nil {
		return err
	}
	// 读取结果
	fmt.Println(string(resp.Body()))
	if err := json.Unmarshal(resp.Body(), result); err != nil {
		result = Bytes2String(resp.Body())
	}
	return nil
}

func FastWithUrlencodeDo(result any, method, uri string, params map[string]string, body string, headers map[string]string) error {
	client := getFastReqClient()
	req, resp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()
	req.Header.SetMethod(method)
	req.SetRequestURI(uri)
	// 设置参数
	var arg fasthttp.Args
	for k, v := range params {
		arg.Add(k, v)
	}
	req.URI().SetQueryString(arg.String())
	// 设置header信息
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	// 设置body
	if len(body) > 0 {
		req.SetBodyString(body)
	}
	// 设置Cookie信息
	// req.Header.SetCookie("key", "val")
	// 发起请求
	if err := client.Do(req, resp); err != nil {
		return err
	}
	// 读取结果
	fmt.Println(string(resp.Body()))
	if err := json.Unmarshal(resp.Body(), result); err != nil {
		result = Bytes2String(resp.Body())
	}
	return nil
}

type JwtPayload struct {
	Domain  string
	UserId  uint64
	AdminId uint64
}

// @secretKey: JWT 加解密密钥
// @iat: 时间戳
// @seconds: 过期时间，单位秒
// @payload: 数据载体
func EncodeJwtToken(secretKey string, iat, seconds int64, payload *JwtPayload) (string, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	payloadStr, _ := json.MarshalToString(payload)
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["payload"] = payloadStr
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

func DecodeJwtToken(payload any, data any) error {
	payloadStr := payload.(string)
	return json.Unmarshal([]byte(payloadStr), data)
}

// 星期转中文
func WeekdayInChinese(weekday time.Weekday) string {
	weekdays := []string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}
	return weekdays[weekday]
}

// 文件转byte
func FileToBytesWithBuffer(file multipart.File) ([]byte, error) {
	file.Seek(0, io.SeekStart)

	var buf bytes.Buffer
	_, err := io.Copy(&buf, file)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// 删除数组中的某个元素
func RemoveItem[T comparable](arr []T, item T) []T {
	for i, v := range arr {
		if v == item {
			return append(arr[:i], arr[i+1:]...)
		}
	}
	return arr
}

// 读取文件内容
func ReadFileToString(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}
	return string(data), nil
}
