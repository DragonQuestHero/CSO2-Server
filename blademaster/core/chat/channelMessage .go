package chat

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnChatChannelMessage(p *InChatPacket, client net.Conn) {
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent ChannelMessage but not in server !")
		return
	}
	//找到对应频道
	chlsrv := GetChannelServerWithID(uPtr.GetUserChannelServerID())
	if chlsrv == nil || chlsrv.ServerIndex <= 0 {
		DebugInfo(2, "Error : User", string(uPtr.IngameName), "sent ChannelMessage but not in channelserver !")
		return
	}
	chl := GetChannelWithID(uPtr.GetUserChannelID(), chlsrv)
	if chl == nil || chl.ChannelID <= 0 {
		DebugInfo(2, "Error : User", string(uPtr.IngameName), "sent ChannelMessage but not in channel !")
		return
	}
	//发送数据
	msg := BuildChatMessage(uPtr, p, ChatChannel)
	for _, v := range UsersManager.Users {
		if !v.CurrentIsIngame && v.GetUserChannelServerID() == chlsrv.ServerIndex && v.GetUserChannelID() == chl.ChannelID {
			//DebugInfo(2, v.Userid)
			SendPacket(BytesCombine(BuildHeader(v.CurrentSequence, PacketTypeChat), msg), v.CurrentConnection)
		}
	}
	DebugInfo(1, "User", string(uPtr.IngameName), "say <", string(p.Message), "> in channel", chl.ChannelID, "channelserver", chlsrv.ServerIndex)
}
