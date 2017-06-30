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

type TCPConnect struct {
	conn net.Conn
	status int
	event IEventTCP
}

func NewTCPConnect(c net.Conn,evt IEventTCP) *TCPConnect {
	return &TCPConnect{conn:c,status:TCPStatus_None,event:evt}
}

func (tp *TCPConnect) Run() {
	go tp.mainRun()
}

func (tp *TCPConnect) mainRun() {
	tp.status = TCPStatus_Connected
	fmt.Println(time.Now().String(),"TCP connected, with remoto ip:",tp.conn.RemoteAddr().String())
	tp.read()
}

func (tp *TCPConnect) read() {
	defer func() {
		fmt.Println(time.Now().String(),"TCP disconnected, with remoto ip:",tp.conn.RemoteAddr().String())
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