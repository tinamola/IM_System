package main

import (
	"net"
	"strings"
)

type User struct {
	Name       string
	Addr       string
	C          chan string
	conn       net.Conn
	UserServer *Server
}

// create a User api
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	newUser := &User{
		Name:       userAddr,
		Addr:       userAddr,
		C:          make(chan string),
		conn:       conn,
		UserServer: server,
	}
	//listen for current channel's message with go routine
	go newUser.ListenMessage()
	return newUser
}

// listen for current chan, obtain message and print out
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}

// user log in  task
func (this *User) Online() {
	//add user to online user map
	this.UserServer.mapLock.Lock()
	this.UserServer.OnlineMap[this.Name] = this
	this.UserServer.mapLock.Unlock()
	//  the message that user is online, include the user themselves.
	this.UserServer.Broadcast(this, "is online!")
}

// user log off  task
func (this *User) Offline() {
	//add user to online user map
	this.UserServer.mapLock.Lock()
	delete(this.UserServer.OnlineMap, this.Name)
	this.UserServer.mapLock.Unlock()
	//broadcast the message that user is online, include the user themselves.
	this.UserServer.Broadcast(this, "is offline!")
}

// user send message task
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		//check all online users details
		this.UserServer.mapLock.Lock()
		for _, user := range this.UserServer.OnlineMap {
			onlinemsg := "[" + user.Addr + "]" + user.Name + " is online.\n"
			this.SendMsgToUser(onlinemsg)
		}
		this.UserServer.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		//change user's own name with format rename|xxx
		newName := strings.Split(msg, "|")[1]
		_, ok := this.UserServer.OnlineMap[newName]
		if ok {
			this.SendMsgToUser("The user name is taken!\n")
		} else {
			this.UserServer.mapLock.Lock()
			delete(this.UserServer.OnlineMap, this.Name)
			this.UserServer.OnlineMap[newName] = this
			this.UserServer.mapLock.Unlock()

			this.Name = newName
			this.SendMsgToUser("You have successfully changed your user name to: " + newName + "\n")
		}
	} else if len(msg) > 3 && msg[:3] == "to|" {
		//Send private msg to somebody
		//get their user name
		toName := strings.Split(msg, "|")[1]
		if toName == "" {
			this.SendMsgToUser("No target name. Should be like to|tina|hello!\n")
			return
		}
		//check if they are online
		remoteUser, ok := this.UserServer.OnlineMap[toName]

		if !ok {
			this.SendMsgToUser("Target user does not exist sorryy\n")
			return
		}
		//send message to that person
		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.SendMsgToUser("No content. Should be like to|tina|hello!\n")
			return
		}
		remoteUser.SendMsgToUser(this.Name + " sent to you: " + content + "\n")
	} else {
		this.UserServer.Broadcast(this, "Sent message: "+msg)
	}
}

func (this *User) SendMsgToUser(msg string) {
	this.conn.Write([]byte(msg))
}
