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
	"asio_based_framework"
)

////////////////////////////////////////////////////////////////////////////////////////////////
// 消息分派
type PacketHandleType func(*asio.UDPSession, asio.Message) error

type EventDispatcher struct {
	// 消息处理器
	packet_handlers map[uint16] PacketHandleType
}

var (
	g_event_dispatcher EventDispatcher
	g_udp_server *asio.UDPServer
)

////////////////////////////////////////////////////////////////////////////////////////////////
// 事件处理器
type MyUDPServerEventController struct {

}

// 有新连接到本服务器
func (this *MyUDPServerEventController)OnNewConnectionToThisServer (client_session *asio.UDPSession) {

}


// 成功连接到服务器
func (this *MyUDPServerEventController)OnConnectedToServer (*asio.UDPSession) {

}

// 连接断开回调
func (this *MyUDPServerEventController)OnDisconnected (*asio.UDPSession) {

}

// PacketHandler 消息处理器
func (this *MyUDPServerEventController)OnPacketValid (session *asio.UDPSession, msg asio.Message) error {

	var err error

	packet := msg.(*wgnet.LenLeadingMessage)

	// 获取消息处理器
	packet_handler, ok := g_event_dispatcher.packet_handlers[packet.Cmd()]
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
	}

	return err
}

// packet 处理失败回调
func (this *MyUDPServerEventController)OnPacketHandleFailed (*asio.UDPSession, error) {

}

//RawMessageCallback 没有注册的消息的回调
func (this *MyUDPServerEventController)OnRawMessage (conn *asio.UDPSession, msg asio.Message) error {
	return nil
}

func TestUDPServer() {
	var err error

	io_service := asio.NewIoService()

	g_udp_server = io_service.NewUDPServer()
	err = g_udp_server.Start("127.0.0.1:22111")
	if err != nil {
		fmt.Errorf("failed to start server:%s\n", err.Error());
		return
	}

	// 绑定事件处理器
	var event_control MyUDPServerEventController
	g_udp_server.BindEventController(&event_control)

	// 绑定协议适配器
	var protocol_apapter wgnet.LenLeadingProtocolProcessor
	g_udp_server.BindProtocolAdapter(&protocol_apapter)

	io_service.Run()
}
