package main

////////////////////////////////////////////////////
// Time        : 2016/6/17 10:11
// Author      : wilsonloo21@163.com
// File        : test_client.go
// Software    : PyCharm
// Description : 测试tcp连接器
////////////////////////////////////////////////////

import (
	"fmt"
	"asio"
	"time"
	"asio_based_framework"
)

type MyEventController struct {

}

// 有新连接到本服务器
func (this *MyEventController)OnNewConnectionToThisServer (*asio.TCPSession) {
	// not in usage
}

// 成功连接到服务器
func (this *MyEventController) OnConnectedToServer(server_session *asio.TCPSession) {
	go client_send_heartbeat_packet(server_session)
}

// 连接断开回调
func (this *MyEventController)OnDisconnected (*asio.TCPSession) {

}

// PacketHandler 消息处理器
func (this *MyEventController)OnPacketValid (*asio.TCPSession, asio.Message) error {
	return nil
}

// packet 处理失败回调
func (this *MyEventController)OnPacketHandleFailed (*asio.TCPSession, error) {

}

//RawMessageCallback 没有注册的消息的回调
func (this *MyEventController)OnRawMessage (conn *asio.TCPSession, msg asio.Message) error {
	return nil
}

// 向服务端发送心跳包
func client_send_heartbeat_packet(server_session *asio.TCPSession) {

	// 每个一段时间发送心跳
	timer := time.NewTimer(30 * time.Second)
	for {
		select {
		case <- timer.C:
			// todo 发送心跳包
		}
	}
}

// 处理服务端的心跳反馈
func on_client_receive_hearbeat_feedback_packet(server_session *asio.TCPSession) {
	server_session.OnHeartBeat()
}

func TestConnector() {
	var err error

	io_service := asio.NewIoService()
	connector_1 := io_service.NewTCPConnector()

	// todo
	connector_1.BindEventController(&MyEventController{})
	connector_1.BindProtocolAdapter(&wgnet.ZlibProtocolProcessor{})

	fmt.Println("starting connecting...")
	err = connector_1.Connect("10.18.12.11:21000")
	if err != nil {
		fmt.Printf("failed to connect server:%s\n", err.Error());
		return
	}
	fmt.Println("starting connecting...OK")

	io_service.Run()
}