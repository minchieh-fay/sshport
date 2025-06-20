package help

type ConfigInfo struct {
	Port          int // 如果是服务端  那么，他就是tcp的端口， 如果是客户端那么就是tcp端口
	ServerAddress string
	Key           string
	Help          bool
	SshAddress    string
	Debug         bool
}

func (c *ConfigInfo) GetType() string {
	if c.ServerAddress != "" {
		return "client"
	}
	return "server"
}
