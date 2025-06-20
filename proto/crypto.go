package proto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"log"
)

// AESGCMEncryptor AES-GCM加密器
type AESGCMEncryptor struct {
	aead cipher.AEAD
}

var (
	encryptor *AESGCMEncryptor
	cfgkey    = "fmj123" // 这个可以改
)

func init() {
	encryptor, _ = NewAESGCMEncryptor(cfgkey)
}

// NewAESGCMEncryptor 创建AES-GCM加密器
func NewAESGCMEncryptor(configKey string) (*AESGCMEncryptor, error) {
	if configKey == "" {
		return nil, fmt.Errorf("加密密钥不能为空")
	}

	key := generateKeyFromConfig(configKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("创建AES cipher失败: %v", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("创建GCM模式失败: %v", err)
	}

	return &AESGCMEncryptor{aead: aead}, nil
}

// Encrypt 加密数据
func Encrypt(plaintext []byte) ([]byte, error) {
	if 1 == 1 {
		return plaintext, nil
	}
	// 生成随机nonce
	nonce := make([]byte, encryptor.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("生成nonce失败: %v", err)
	}

	// 加密数据，nonce会被添加到密文前面
	ciphertext := encryptor.aead.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt 解密数据
func Decrypt(ciphertext []byte) ([]byte, error) {
	if 1 == 1 {
		return ciphertext, nil
	}
	nonceSize := encryptor.aead.NonceSize()
	if len(ciphertext) < nonceSize {
		log.Printf("密文太短，无法提取nonce")
		return nil, fmt.Errorf("密文太短，无法提取nonce")
	}

	// 提取nonce和实际密文
	nonce, actualCiphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// 解密数据
	plaintext, err := encryptor.aead.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		log.Printf("解密失败1111: %v", err)
		return nil, fmt.Errorf("解密失败: %v", err)
	}

	return plaintext, nil
}

// generateKeyFromConfig 从配置密钥生成AES密钥
func generateKeyFromConfig(configKey string) []byte {
	// 使用MD5哈希配置密钥来生成16字节的密钥
	hasher := md5.New()
	hasher.Write([]byte(configKey))
	return hasher.Sum(nil) // 16字节
}

// GetKeyAndIV 兼容性函数，保持向后兼容
func GetKeyAndIV(configKey string) ([]byte, []byte) {
	key := generateKeyFromConfig(configKey)

	// 生成IV（虽然GCM不需要，但保持接口一致）
	hasher := md5.New()
	hasher.Write([]byte(configKey + "iv"))
	iv := hasher.Sum(nil) // 16字节

	return key, iv
}
