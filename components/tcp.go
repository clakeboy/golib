package components

import (
	"fmt"
	"log"
	"net"
	"time"
)

//TCP状态常量
const (
	TCPStatusNone = iota
	TCPStatusConnected
	TCPStatusDisconnected
)

//TCP事件常量
const (
	TCPEventConnected = iota
	TCPEventDisconnected
	TCPEventRecv
	TCPEventWritten
)

//TCP连接
type TCPConnect struct {
	conn        net.Conn      //TCP连接
	status      int           //TCP状态
	event       IEventTCP     //TCP事件接口
	writeChan   chan []byte   //写入队列
	closeChan   chan bool     //关闭队列
	debug       bool          //是否DEUBG模式
	readTimeout time.Duration //读取超时时间,单位秒
}

//新建TCP服务连接
func NewTCPConnect(c net.Conn, evt IEventTCP) *TCPConnect {
	return &TCPConnect{
		conn:        c,
		status:      TCPStatusNone,
		event:       evt,
		writeChan:   make(chan []byte, 10),
		closeChan:   make(chan bool),
		readTimeout: time.Second * 30,
	}
}

//设置读取超时
func (tp *TCPConnect) SetReadTimeout(ss int64) {
	tp.readTimeout = time.Second * time.Duration(ss)
}

//运行
func (tp *TCPConnect) Run() {
	go tp.mainRun()
}

//主TCP线程
func (tp *TCPConnect) mainRun() {
	tp.status = TCPStatusConnected
	if tp.debug {
		fmt.Println(time.Now().String(), "TCP connected, with remoto ip:", tp.conn.RemoteAddr().String())
	}
	tp.emit(TCPEventConnected, nil)
	go tp.write()
	tp.read()
}

//读取线程
func (tp *TCPConnect) read() {
	defer func() {
		if tp.debug {
			fmt.Println(time.Now().String(), "TCP Read disconnected, with remoto ip:", tp.conn.RemoteAddr().String())
		}
		tp.Close()
	}()
	for {
		if tp.readTimeout != 0 {
			tp.conn.SetReadDeadline(time.Now().Add(tp.readTimeout))
		}
		buf := make([]byte, 256)
		msg_len, err := tp.conn.Read(buf)
		if err != nil {
			log.Println(time.Now().String(), "Conn read error:", err)
			break
		}
		if tp.debug {
			fmt.Printf("[DEBUG] Read: %x,%d\n", buf[:msg_len], msg_len)
		}
		tp.emit(TCPEventRecv, buf[:msg_len])
	}
}

//写入线程
func (tp *TCPConnect) write() {
	defer func() {
		if tp.debug {
			fmt.Println(time.Now().String(), "TCP Write disconnected, with remoto ip:", tp.conn.RemoteAddr().String())
		}
		tp.Close()
	}()

	for {
		select {
		case data := <-tp.writeChan:
			tp.conn.SetWriteDeadline(time.Now().Add(time.Second * 30))
			msg_len, err := tp.conn.Write(data)
			if err != nil {
				log.Println(time.Now().String(), "Conn write error:", err)
				return
			}
			if tp.debug {
				fmt.Printf("[DEBUG] Write: %x %d\n", data, msg_len)
			}
			tp.emit(TCPEventWritten, data)
		case <-tp.closeChan:
			return
		default:
			time.Sleep(time.Second)
		}
	}
}

//写入数据到TCP
func (tp *TCPConnect) WriteData(data []byte) {
	tp.writeChan <- data
}

//关闭连接
func (tp *TCPConnect) Close() {
	if tp.status == TCPStatusDisconnected {
		return
	}
	tp.emit(TCPEventDisconnected, nil)
	tp.status = TCPStatusDisconnected
	tp.conn.Close()
	tp.closeChan <- true
}

//触发一个事件
func (tp *TCPConnect) emit(event_type int, data interface{}) {
	evt := &TCPConnEvent{
		EventType: event_type,
		Conn:      tp,
		Data:      data,
	}
	switch event_type {
	case TCPEventConnected:
		go tp.event.OnConnected(evt)
	case TCPEventDisconnected:
		go tp.event.OnDisconnected(evt)
	case TCPEventRecv:
		go tp.event.OnRecv(evt)
	case TCPEventWritten:
		go tp.event.OnWritten(evt)
	}
}

//得到当前TCP状态
func (tp *TCPConnect) Status() int {
	return tp.status
}

//设置DEBUG状态
func (tp *TCPConnect) SetDebug(yes bool) {
	tp.debug = yes
}

//事件类型结构数据
type TCPConnEvent struct {
	EventType int
	Conn      *TCPConnect
	Data      interface{}
}

type IEventTCP interface {
	OnConnected(evt *TCPConnEvent)
	OnDisconnected(evt *TCPConnEvent)
	OnRecv(evt *TCPConnEvent)
	OnWritten(evt *TCPConnEvent)
}
