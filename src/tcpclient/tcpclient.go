package tcpclient

import (
	"context"
	"fmt"
	"net"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type TCPClient struct {
	ctx   *context.Context
	conn  net.Conn
	ip    string
	port  int
	ready bool
}

func NewTCPClient(ctx *context.Context) *TCPClient {
	return &TCPClient{
		ctx:   ctx,
		conn:  nil,
		ready: false,
	}
}

// 클라이언트 연결 시작
func (c *TCPClient) StartClient(ip string, port int) bool {
	c.ip = ip
	c.port = port

	err := c.connectToServer()

	if err != nil {
		runtime.MessageDialog(*c.ctx, runtime.MessageDialogOptions{
			Type:          runtime.ErrorDialog,
			Title:         "Error",
			Message:       "Could not connect to the server",
			Buttons:       nil,
			DefaultButton: "",
			CancelButton:  "",
		})
		return false
	}

	c.ready = true
	runtime.MessageDialog(*c.ctx, runtime.MessageDialogOptions{
		Type:          runtime.InfoDialog,
		Title:         "Connected",
		Message:       "Connected to the server",
		Buttons:       nil,
		DefaultButton: "",
		CancelButton:  "",
	})
	return true
}

func (c *TCPClient) connectToServer() error {
	serverAddress := fmt.Sprintf("%s:%d", c.ip, c.port)
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return err
	}
	c.conn = conn

	// 클라이언트와 서버 연결 성공
	return nil
}

func (c *TCPClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
