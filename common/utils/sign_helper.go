package utils

import (
	"TFLanHttpDesktop/common/define"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

type Signature struct {
	Salt      string // 随机盐值
	Timestamp int64  // 生成时间戳（Unix秒）
	Sign      string // 最终签名
}

// GenerateSignature 生成字符串的随机签名
func GenerateSignature(data string) (string, error) {
	// 1. 生成随机盐值
	salt, err := generateRandomSalt(define.SignSaltLength)
	if err != nil {
		return "", fmt.Errorf("生成盐值失败: %v", err)
	}

	// 2. 获取当前时间戳（用于过期校验）
	timestamp := time.Now().Unix()

	// 3. 组合待签名的数据（盐值 + 时间戳 + 原始数据 + 密钥）
	signStr := fmt.Sprintf("%s:%d:%s:%s", salt, timestamp, data, define.SignSecretKey)

	// 4. 使用 HMAC-SHA256 计算哈希
	mac := hmac.New(sha256.New, []byte(define.SignSecretKey))
	_, err = mac.Write([]byte(signStr))
	if err != nil {
		return "", fmt.Errorf("计算哈希失败: %v", err)
	}
	sign := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	// 5. 组合盐值、时间戳和签名为一个字符串（便于传输）
	signature := Signature{
		Salt:      salt,
		Timestamp: timestamp,
		Sign:      sign,
	}

	// 6. 序列化为JSON并Base64编码（避免特殊字符问题）
	signBytes, err := json.Marshal(signature)
	if err != nil {
		return "", fmt.Errorf("序列化签名失败: %v", err)
	}

	return base64.URLEncoding.EncodeToString(signBytes), nil
}

// VerifySignature 校验签名是否有效
func VerifySignature(data, signatureStr string) (bool, error) {
	// 1. 解码签名字符串
	signBytes, err := base64.URLEncoding.DecodeString(signatureStr)
	if err != nil {
		return false, fmt.Errorf("签名解码失败: %v", err)
	}

	// 2. 反序列化为Signature结构
	var signature Signature
	if err := json.Unmarshal(signBytes, &signature); err != nil {
		return false, fmt.Errorf("签名反序列化失败: %v", err)
	}

	// 3. 校验签名是否过期（如果设置了有效期）

	now := time.Now().Unix()
	if now-signature.Timestamp > define.SignExpiresIn {
		return false, fmt.Errorf("签名已过期（有效期 %d 秒）", define.SignExpiresIn)
	}

	// 4. 重新计算签名并比对
	signStr := fmt.Sprintf("%s:%d:%s:%s", signature.Salt, signature.Timestamp, data, define.SignSecretKey)
	mac := hmac.New(sha256.New, []byte(define.SignSecretKey))
	_, err = mac.Write([]byte(signStr))
	if err != nil {
		return false, fmt.Errorf("校验时计算哈希失败: %v", err)
	}
	computedSign := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	// 5. 比较两个签名是否一致（使用常数时间比较，防止时序攻击）
	if !hmac.Equal([]byte(computedSign), []byte(signature.Sign)) {
		return false, fmt.Errorf("签名不匹配")
	}

	return true, nil
}

// 生成指定长度的随机盐值
func generateRandomSalt(length int) (string, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt) // 使用加密安全的随机数生成器
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(salt), nil
}
