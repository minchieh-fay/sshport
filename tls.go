package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math"
	"math/big"
	"net"
	"time"

	"github.com/quic-go/quic-go"
)

// 从配置密钥生成AES密钥
func GetKeyAndIV(configKey string) ([]byte, []byte) {
	// 使用MD5哈希配置密钥来生成16字节的密钥
	hasher := md5.New()
	hasher.Write([]byte(configKey))
	key := hasher.Sum(nil) // 16字节

	// 使用密钥的另一个哈希作为IV
	hasher.Reset()
	hasher.Write([]byte(configKey + "iv"))
	iv := hasher.Sum(nil) // 16字节

	return key, iv
}

// AES-GCM 加密器
type AESGCMEncryptor struct {
	aead cipher.AEAD
}

// 创建 AES-GCM 加密器
func NewAESGCMEncryptor(key []byte) (*AESGCMEncryptor, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &AESGCMEncryptor{aead: aead}, nil
}

// 加密数据
func (e *AESGCMEncryptor) Encrypt(plaintext []byte) ([]byte, error) {
	// 生成随机 nonce
	nonce := make([]byte, e.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	// 加密数据，nonce 会被添加到密文前面
	ciphertext := e.aead.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// 解密数据
func (e *AESGCMEncryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	nonceSize := e.aead.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// 提取 nonce 和实际密文
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// 解密数据
	plaintext, err := e.aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// 生成TLS配置
func GenerateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		NotBefore:    time.Now().Add(-time.Hour * 24),
		NotAfter:     time.Now().Add(time.Hour * 24 * 1800),
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
		DNSNames:     []string{"localhost"},
	}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"ff-quic-tunnel"},
	}
}

func GetSimpleTLSConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"ff-quic-tunnel"},
	}
}

func GetQuicConfig() *quic.Config {
	return &quic.Config{
		MaxIdleTimeout:  math.MaxInt64,
		KeepAlivePeriod: time.Second * 10,
	}
}
