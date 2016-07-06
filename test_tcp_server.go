package main

////////////////////////////////////////////////////
// Time        : 2016/6/17 11:42
// Author      : wilsonloo21@163.com
// File        : test_tcp_server.go
// Software    : PyCharm

// Description : 测试tcp 服务器
////////////////////////////////////////////////////

import (
	"fmt"
	"asio"
	"time"
	"asio_based_framework"
)

////////////////////////////////////////////////////////////////////////////////////////////////
// 消息分派
type TCPPacketHandleType func(*asio.TCPSession, asio.Message) error

type TCPEventDispatcher struct {
	// 消息处理器
	packet_handlers map[uint16] TCPPacketHandleType
}

var (
	g_tcp_event_dispatcher TCPEventDispatcher
	g_tcp_server *asio.TCPServer
)

////////////////////////////////////////////////////////////////////////////////////////////////
// 事件处理器
type MyServerEventController struct {

}

// 有新连接到本服务器
func (this *MyServerEventController)OnNewConnectionToThisServer (client_session *asio.TCPSession) {
	// 独立协程处理心跳
	go server_send_heartbeat_packet(client_session)
}

// 连接断开回调
func (this *MyServerEventController)OnDisconnected (*asio.TCPSession) {

}

// PacketHandler 消息处理器
func (this *MyServerEventController)OnPacketValid (session *asio.TCPSession, msg asio.Message) error {

	var err error

	packet := msg.(*wgnet.LenLeadingMessage)

	// 获取消息处理器
	packet_handler, ok := g_tcp_event_dispatcher.packet_handlers[packet.Cmd()]
	if !ok {
		// 消息没有被注册
		err = this.OnRawMessage(session, packet)
	} else {
		//消息已经被注册
		err = packet_handler(session, msg)
	}

	// 先回收消息
	wgnet.FreeLenLeadingMessage(packet)

	if err != nil {
		// 最后一次推送给用户
		this.OnPacketHandleFailed(session, err)

		//回调返回值不为空，则关闭连接
		fmt.Println("err001:", err.Error())
		session.Close()
	}

	return err
}

// packet 处理失败回调
func (this *MyServerEventController)OnPacketHandleFailed (*asio.TCPSession, error) {

}

//RawMessageCallback 没有注册的消息的回调
func (this *MyServerEventController)OnRawMessage (conn *asio.TCPSession, msg asio.Message) error {
	return nil
}

// 向客户端发送心跳包
func server_send_heartbeat_packet(client_session *asio.TCPSession) {

	// 每个一段时间发送心跳
	timer := time.NewTimer(30 * time.Second)
	for {
		select {
		case <- timer.C:
			// todo 发送心跳包
		}
	}
}

// 处理客户端的心跳反馈
func on_server_receive_hearbeat_feedback_packet(client_session *asio.TCPSession) {
	client_session.OnHeartBeat()
}

func TestServer() {
	var err error

	io_service := asio.NewIoService()
	g_tcp_server = io_service.NewTCPServer()

	// 绑定事件处理器
	var event_control MyEventController
	g_tcp_server.BindEventController(&event_control)

	// 绑定协议适配器
	var protocol_apapter wgnet.LenLeadingProtocolProcessor
	g_tcp_server.BindProtocolAdapter(&protocol_apapter)

	err = g_tcp_server.Start("127.0.0.1:22111")
	if err != nil {
		fmt.Errorf("failed to start server:%s\n", err.Error());
		return
	}

	io_service.Run()
}
