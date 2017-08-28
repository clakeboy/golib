package components

import (
	"net"
	"log"
	"fmt"
	"time"
)

const (
	TCPStatus_None = iota
	TCPStatus_Connected
	TCPStatus_Disconnected
)

const (
	TCPEvent_Connected = iota
	TCPEvent_Disconnected
	TCPEvent_Recv
)
//TCP连接
type TCPConnect struct {
	conn net.Conn
	status int
	event IEventTCP
}

//新建TCP服务连接
func NewTCPConnect(c net.Conn,evt IEventTCP) *TCPConnect {
	return &TCPConnect{conn:c,status:TCPStatus_None,event:evt}
}
//运行
func (tp *TCPConnect) Run() {
	go tp.mainRun()
}
//主TCP线程
func (tp *TCPConnect) mainRun() {
	tp.status = TCPStatus_Connected
	fmt.Println(time.Now().String(),"TCP connected, with remoto ip:",tp.conn.RemoteAddr().String())
	tp.read()
}
//读取线程
func (tp *TCPConnect) read() {
	defer func() {
		fmt.Println(time.Now().String(),"TCP Read disconnected, with remoto ip:",tp.conn.RemoteAddr().String())
		tp.conn.Close()
		tp.status = TCPStatus_Disconnected
	}()
	for {
		buf := make([]byte,256)
		msg_len, err := tp.conn.Read(buf)
		if err != nil {
			log.Println(time.Now().String(),"Conn read error:", err)
			break
		}

		tp.emit(TCPEvent_Recv,buf[:msg_len])
	}
}
//写入线程
func (tp *TCPConnect) write() {
	defer func() {
		fmt.Println(time.Now().String(),"TCP Write disconnected, with remoto ip:",tp.conn.RemoteAddr().String())
		tp.conn.Close()
		tp.status = TCPStatus_Disconnected
	}()


}

//触发一个事件
func (tp *TCPConnect) emit(event_type int,data []byte) {
	evt := &TCPConnEvent{
		EventType:TCPEvent_Connected,
		Conn:tp,
		Data:data,
	}
	switch event_type {
	case TCPEvent_Connected:
		go tp.event.OnConnected(evt)
	case TCPEvent_Disconnected:
		go tp.event.OnDisconnected(evt)
	case TCPEvent_Recv:
		go tp.event.OnRecv(evt)
	}
}

//得到当前TCP状态
func (tp *TCPConnect) Status() int {
	return tp.status
}

type TCPConnEvent struct {
	EventType int
	Conn      *TCPConnect
	Data      []byte
}

type IEventTCP interface {
	OnConnected(evt *TCPConnEvent)
	OnDisconnected(evt *TCPConnEvent)
	OnRecv(evt *TCPConnEvent)
}