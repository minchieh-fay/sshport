package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/quic-go/quic-go"
)

// quic服务端

type QuicServer struct {
	config *ConfigInfo
}

func NewQuicServer(config *ConfigInfo) *QuicServer {
	return &QuicServer{
		config: config,
	}
}

func (s *QuicServer) Start() {
	// 生成TLS配置
	tlsConfig := GenerateTLSConfig()
	listener, err := quic.ListenAddr(fmt.Sprintf(":%d", s.config.Port), tlsConfig, GetQuicConfig())
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			log.Printf("接受连接失败: %v", err)
			time.Sleep(time.Second * 1)
			continue
		}
		log.Printf("接受了新的Bridge连接: %s", conn.RemoteAddr())

		go s.handleConnection(conn)
	}
}

// 接受stream 并 转发到ssh 并返回
func (s *QuicServer) handleConnection(conn quic.Connection) {
	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			log.Printf("接受stream失败: %v", err)
			return
		}
		log.Printf("接受了新的stream: %s", stream.StreamID())
		go s.handleStream(stream)
	}
}

// 接受stream 并 转发到ssh 并返回
func (s *QuicServer) handleStream(stream quic.Stream) {
	sshconn, err := net.DialTimeout("tcp", s.config.SshAddress, time.Second*2)
	if err != nil {
		log.Printf("连接ssh失败: %v", err)
		stream.Close()
		return
	}
	defer sshconn.Close()

	// 检查是否配置了加密密钥
	if s.config.Key == "" {
		log.Println("警告: 未配置加密密钥，数据将不加密传输")
		// 不加密的数据传输
		done := make(chan struct{})
		go func() {
			io.Copy(sshconn, stream)
			if debug == true {
				log.Println("333333")
			}
			done <- struct{}{}
		}()
		go func() {
			io.Copy(stream, sshconn)
			if debug == true {
				log.Println("44444")
			}
			done <- struct{}{}
		}()
		<-done
	} else {
		// 加密的数据传输
		s.handleEncryptedStream(stream, sshconn)
	}

	stream.Close()
	sshconn.Close()
}

func (s *QuicServer) handleEncryptedStream(stream quic.Stream, sshconn net.Conn) {
	// 获取加密密钥
	key, _ := GetKeyAndIV(s.config.Key)

	// 创建 AES-GCM 加密器
	encryptor, err := NewAESGCMEncryptor(key)
	if err != nil {
		log.Printf("创建加密器失败: %v", err)
		return
	}

	done := make(chan struct{})

	// QUIC Stream -> SSH (解密)
	go func() {
		defer func() { done <- struct{}{} }()

		for {
			// 读取数据长度（4字节）
			lengthBytes := make([]byte, 4)
			if _, err := io.ReadFull(stream, lengthBytes); err != nil {
				if err != io.EOF {
					log.Printf("读取长度失败: %v", err)
				}
				break
			}

			// 解析数据长度
			length := int(lengthBytes[0])<<24 | int(lengthBytes[1])<<16 | int(lengthBytes[2])<<8 | int(lengthBytes[3])

			// 读取加密数据
			encrypted := make([]byte, length)
			if _, err := io.ReadFull(stream, encrypted); err != nil {
				log.Printf("读取加密数据失败: %v", err)
				break
			}

			// 解密数据
			decrypted, err := encryptor.Decrypt(encrypted)
			if err != nil {
				log.Printf("解密失败: %v", err)
				break
			}

			// 写入解密数据到SSH连接
			if _, err := sshconn.Write(decrypted); err != nil {
				log.Printf("写入解密数据失败: %v", err)
				break
			}
		}
	}()

	// SSH -> QUIC Stream (加密)
	go func() {
		defer func() { done <- struct{}{} }()

		buffer := make([]byte, 4096)
		for {
			n, err := sshconn.Read(buffer)
			if err != nil {
				if err != io.EOF {
					log.Printf("SSH读取失败: %v", err)
				}
				break
			}

			// 加密数据
			encrypted, err := encryptor.Encrypt(buffer[:n])
			if err != nil {
				log.Printf("加密失败: %v", err)
				break
			}

			// 写入加密数据长度（4字节）+ 加密数据
			lengthBytes := make([]byte, 4)
			lengthBytes[0] = byte(len(encrypted) >> 24)
			lengthBytes[1] = byte(len(encrypted) >> 16)
			lengthBytes[2] = byte(len(encrypted) >> 8)
			lengthBytes[3] = byte(len(encrypted))

			if _, err := stream.Write(lengthBytes); err != nil {
				log.Printf("写入长度失败: %v", err)
				break
			}

			if _, err := stream.Write(encrypted); err != nil {
				log.Printf("写入加密数据失败: %v", err)
				break
			}
		}
	}()

	<-done
	sshconn.Close()
	stream.Close()
}
