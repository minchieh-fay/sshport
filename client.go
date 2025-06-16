package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
)

// quic客户端

type QuicClient struct {
	config   *ConfigInfo
	quicconn quic.Connection
	mtx      sync.Mutex
}

func NewQuicClient(config *ConfigInfo) *QuicClient {
	return &QuicClient{
		config: config,
	}
}

func (c *QuicClient) Start() {
	// 连上quicserver
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // 确保释放资源
	quicconn, err := quic.DialAddr(ctx, c.config.ServerAddress, GetSimpleTLSConfig(), GetQuicConfig())
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
		os.Exit(1)
	}
	c.quicconn = quicconn

	// 坚听一个tcp端口
	port := c.config.Port
	if debug == true {
		port = port + 1
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Failed to accept: %v", err)
		}
		go c.handleConnection(conn)
	}

}

func (c *QuicClient) getQuicStream() (quic.Stream, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if c.quicconn == nil {
		return nil, errors.New("quicconn is nil")
	}
	stream, err := c.quicconn.OpenStreamSync(context.Background())
	if err != nil {
		c.quicconn.CloseWithError(quic.ApplicationErrorCode(quic.NoError), "quicconn is nil")
		c.quicconn = nil
		// 重新dial
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		quicconn, err := quic.DialAddr(ctx, c.config.ServerAddress, GetSimpleTLSConfig(), GetQuicConfig())
		if err != nil {
			return nil, err
		}
		c.quicconn = quicconn
		stream, err = c.quicconn.OpenStreamSync(context.Background())
		if err != nil {
			return nil, err
		}
	}
	return stream, nil
}

func (c *QuicClient) handleConnection(conn net.Conn) {
	stream, err := c.getQuicStream()
	if err != nil {
		log.Printf("Failed to get quic stream: %v", err)
		conn.Close()
		return
	}

	// 检查是否配置了加密密钥
	if c.config.Key == "" {
		log.Println("警告: 未配置加密密钥，数据将不加密传输")
		// 不加密的数据传输
		done := make(chan struct{})
		go func() {
			io.Copy(conn, stream)
			if debug == true {
				log.Println("11111")
			}
			done <- struct{}{}
		}()
		go func() {
			io.Copy(stream, conn)
			if debug == true {
				log.Println("22222")
			}
			done <- struct{}{}
		}()
		<-done
	} else {
		// 加密的数据传输
		c.handleEncryptedConnection(conn, stream)
	}

	stream.Close()
	conn.Close()
}

func (c *QuicClient) handleEncryptedConnection(conn net.Conn, stream quic.Stream) {
	// 获取加密密钥
	key, _ := GetKeyAndIV(c.config.Key)

	// 创建 AES-GCM 加密器
	encryptor, err := NewAESGCMEncryptor(key)
	if err != nil {
		log.Printf("创建加密器失败: %v", err)
		return
	}

	done := make(chan struct{})

	// TCP -> QUIC Stream (加密)
	go func() {
		defer func() { done <- struct{}{} }()

		buffer := make([]byte, 4096)
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				if err != io.EOF {
					log.Printf("TCP读取失败: %v", err)
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

	// QUIC Stream -> TCP (解密)
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

			// 写入解密数据到TCP连接
			if _, err := conn.Write(decrypted); err != nil {
				log.Printf("写入解密数据失败: %v", err)
				break
			}
		}
	}()

	<-done
	conn.Close()
	stream.Close()
}
