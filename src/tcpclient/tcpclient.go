package tcpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"go_file_sync/src/logs"
	"net"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// TCPClient는 서버에 대한 클라이언트 연결을 관리합니다.
type TCPClient struct {
	ctx          *context.Context
	conn         net.Conn
	ip           string
	port         int
	connectState bool
}

type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

// NewTCPClient는 새 TCPClient 인스턴스를 생성합니다.
func NewTCPClient(ctx *context.Context) *TCPClient {
	return &TCPClient{
		ctx: ctx,
	}
}

func (c *TCPClient) connectToServer(ip string, port int) (net.Conn, error) {
	serverAddress := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// 서버에 연결을 시도하고 클라이언트를 초기화
func (c *TCPClient) StartClient(ip string, port int) bool {
	c.ip = ip
	c.port = port

	conn, err := c.connectToServer(ip, port)
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

	if c.conn != nil {
		return true
	}

	c.conn = conn

	// 클라이언트와 서버 연결 성공
	logs.PrintMsgLog("서버에 연결 성공")
	c.connectState = true

	go c.ReceiveMessages() // 클라이언트가 메시지를 받을 수 있도록 고루틴 시작
	return true
}

func (c *TCPClient) handleMessage(buffer []byte, n int) {
	var message Message

	err := json.Unmarshal(buffer[:n], &message)
	if err != nil {
		runtime.MessageDialog(*c.ctx, runtime.MessageDialogOptions{
			Type:          runtime.ErrorDialog,
			Title:         "Error",
			Message:       "데이터 수신에 실패하였습니다.",
			Buttons:       nil,
			DefaultButton: "",
			CancelButton:  "",
		})
		logs.PrintMsgLog(fmt.Sprintf("데이터 수신에 실패하였습니다.: %s\n", err.Error()))
		return
	}

	logs.PrintMsgLog(fmt.Sprintf("서버로부터 받은 헤더: %s\n", message.Type))
	switch message.Type {
	case "close server":
		logs.PrintMsgLog("서버 닫힘, 연결 끊기")
		runtime.EventsEmit(*c.ctx, "client_server_disconnect", true)
		c.Close()
	}
}

func (c *TCPClient) ReceiveMessages() {
	for {
		buffer := make([]byte, 1024)

		n, err := c.conn.Read(buffer)
		if err != nil {
			logs.PrintMsgLog(fmt.Sprintf("메시지 받기 실패 에러: %s\n", err.Error()))
			c.Close()
			return
		}

		c.handleMessage(buffer, n)
	}
}

func (c *TCPClient) Close() {
	if c.conn != nil {
		c.connectState = false
		c.conn.Close()
	}
}

// 클라이언트가 서버에 연결한 이후 서버에 현재 PC에서 실행중인 포트를 보냄
func (c *TCPClient) SendAutoConnectServer(port int) {
	message := Message{
		Type:    "auto connect",
		Content: port,
	}

	// JSON 직렬화
	writeData, err := json.Marshal(message)
	if err != nil {
		runtime.MessageDialog(*c.ctx, runtime.MessageDialogOptions{
			Type:          runtime.ErrorDialog,
			Title:         "Error",
			Message:       "데이터 송신에 실패하였습니다.",
			Buttons:       nil,
			DefaultButton: "",
			CancelButton:  "",
		})
		logs.PrintMsgLog(fmt.Sprintf("데이터 송신에 실패하였습니다.: %s\n", err.Error()))
	}

	_, err = c.conn.Write(writeData)
	if err != nil {
		logs.PrintMsgLog(fmt.Sprintf("Error sending close signal: %s\n", err.Error()))
	}
}
