package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int
	//online users
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	//channel to send message
	Message chan string
}

func (this *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

func (this *Server) ListenMessager() {
	for {
		message := <-this.Message
		this.mapLock.Lock()
		for _, user := range this.OnlineMap {
			user.C <- message //note: don't send a channel to a channel, a channel can only be extracted to a value
		}
		this.mapLock.Unlock()
	}
}

func (this *Server) Handler(conn net.Conn) {
	fmt.Println("Server Connects Successful!")
	user := NewUser(conn, this)
	//handle user online flow
	user.Online()
	isLive := make(chan bool)
	//user flow
	go func() {
		buf := make([]byte, 4096)
		for {
			n, error := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if error != nil && error != io.EOF {
				fmt.Println("Content read error: ", error)
			}
			msg := string(buf[:n-1]) //ignore the last /n
			user.DoMessage(msg)
			isLive <- true
		}
	}()
	//block the handler
	for {
		select {
		case <-isLive: //case <-isLive必须在上面，如果一直没有isLive, 10秒就会记着，如果有了，select就会run一次，for loop重新进，所以10秒也会重新记
		case <-time.After(100 * time.Second):
			user.SendMsgToUser("You are kicked out.")
			close(user.C)
			conn.Close()
			return //runtime.GoExit()
		}
	}
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}
func (this *Server) Start() {
	listener, error := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if error != nil {
		fmt.Println("network error: ", error)
		return
	}
	defer listener.Close()
	//start the go routine to listen message
	go this.ListenMessager()

	for {
		conn, error := listener.Accept()
		if error != nil {
			fmt.Println("listener accept error: ", error)
			continue
		}
		//do handler
		go this.Handler(conn)
	}
}
