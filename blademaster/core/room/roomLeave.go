package room

import (
	"net"

	//. "github.com/KouKouChan/CSO2-Server/blademaster/core/message"
	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnLeaveRoom(client net.Conn, end bool) {
	//找到玩家
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		//DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "try to leave room but not in server !")
		return
	}
	//找到玩家的房间
	rm := GetRoomFromID(uPtr.GetUserChannelServerID(),
		uPtr.GetUserChannelID(),
		uPtr.GetUserRoomID())
	if rm == nil ||
		rm.Id <= 0 {
		//DebugInfo(2, "Error : User", string(uPtr.Username), "try to leave a null room !")
		return
	}
	//检查玩家游戏状态，准备情况下并且开始倒计时了，那么就不允许离开房间
	if uPtr.IsUserReady() &&
		rm.IsGlobalCountdownInProgress() {
		DebugInfo(2, "Error : User", string(uPtr.UserName), "try to leave room but is started !")
		return
	}
	//房间移除玩家
	rm.RoomRemoveUser(uPtr.Userid)
	//检查房间是否为空
	if rm.NumPlayers <= 0 {
		DelChannelRoom(rm.Id,
			uPtr.GetUserChannelID(),
			uPtr.GetUserChannelServerID())

	} else {
		//向其他玩家发送离开信息
		SentUserLeaveMes(uPtr, rm)
	}
	// //扣除1000points
	// if uPtr.CurrentIsIngame {
	// 	uPtr.PunishPoints()
	// 	OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox, GAME_ROOM_LEAVE_EARLY)
	// 	//UserInfo部分
	// 	rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeUserInfo), BuildUserInfo(0XFFFFFFFF, NewUserInfo(uPtr), uPtr.Userid, true))
	// 	SendPacket(rst, uPtr.CurrentConnection)
	// }
	//设置玩家状态
	uPtr.QuitRoom()
	//发送房间列表给玩家
	if !end {
		OnBroadcastRoomList(uPtr.GetUserChannelServerID(), uPtr.GetUserChannelID(), uPtr)
	}
	//房间状态
	rm.CheckIngameStatus()
	DebugInfo(2, "User", string(uPtr.UserName), "id", uPtr.Userid, "left room", string(rm.Setting.RoomName), "id", rm.Id)
}
func SentUserLeaveMes(uPtr *User, rm *Room) {
	//发送离开消息
	if rm.HostUserID == uPtr.Userid {
		//选出新房主
		for _, v := range rm.Users {
			rm.SetRoomHost(v)
			DebugInfo(2, "Set User", string(v.UserName), "id", v.Userid, "to host in room", string(rm.Setting.RoomName), "id", rm.Id)
			if !v.CurrentIsIngame {
				v.SetUserStatus(UserNotReady)
				temp := BuildUserReadyStatus(v)
				for _, k := range rm.Users {
					rst := BytesCombine(BuildHeader(k.CurrentSequence, PacketTypeRoom), temp)
					SendPacket(rst, k.CurrentConnection)
				}
			}
			break
		}
		numInGame := 0
		leave := BuildUserLeave(uPtr.Userid)
		sethost := BuildSetHost(rm.HostUserID)

		for _, v := range rm.Users {
			if v.CurrentIsIngame {
				numInGame++
			}
			rst1 := append(BuildHeader(v.CurrentSequence, PacketTypeRoom), OUTPlayerLeave)
			rst1 = BytesCombine(rst1, leave)
			rst2 := append(BuildHeader(v.CurrentSequence, PacketTypeRoom), OUTSetHost)
			rst2 = BytesCombine(rst2, sethost)

			SendPacket(rst1, v.CurrentConnection)
			SendPacket(rst2, v.CurrentConnection)

			// hostu := GetUserFromID(rm.HostUserID)
			// rst := BytesCombine(BuildHeader(hostu.CurrentSequence, PacketTypeHost), BuildGameContinue(rm.HostUserID))
			// SendPacket(rst, hostu.CurrentConnection)

			// rst = UDPBuild(hostu.CurrentSequence, 0, v.Userid, v.NetInfo.ExternalIpAddress, v.NetInfo.ExternalClientPort)
			// SendPacket(rst, hostu.CurrentConnection)
		}

		if numInGame == 0 {
			rm.SetStatus(StatusWaiting)
			// setting := BuildRoomSetting(rm, 0x404000)
			// for _, v := range rm.Users {
			// 	rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeRoom), setting)
			// 	SendPacket(rst, v.CurrentConnection)
			// }
		}
		DebugInfo(2, "Sent a set roomHost packet to other users")
	} else {
		leave := BuildUserLeave(uPtr.Userid)
		for _, v := range rm.Users {
			rst1 := append(BuildHeader(v.CurrentSequence, PacketTypeRoom), OUTPlayerLeave)
			rst1 = BytesCombine(rst1, leave)
			SendPacket(rst1, v.CurrentConnection)
		}
		DebugInfo(2, "Sent a leave room packet to other users")
	}
}
func BuildUserLeave(id uint32) []byte {
	buf := make([]byte, 4)
	offset := 0
	WriteUint32(&buf, id, &offset)
	return buf
}
func BuildSetHost(id uint32) []byte {
	buf := make([]byte, 5)
	offset := 0
	WriteUint32(&buf, id, &offset)
	buf[4] = 1
	return buf
}

func BuildGameContinue(id uint32) []byte {
	buf := make([]byte, 5)
	offset := 0
	WriteUint8(&buf, GameStart, &offset)
	WriteUint32(&buf, id, &offset)
	return buf[:offset]
}
