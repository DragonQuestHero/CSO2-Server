package room

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnRoomRequest(p *PacketData, client net.Conn) {
	var pkt InRoomPaket
	if p.PraseRoomPacket(&pkt) {
		switch pkt.InRoomType {
		case NewRoomRequest:
			//log.Println("Recived a new room request from", client.RemoteAddr().String())
			OnNewRoom(p, client)
		case JoinRoomRequest:
			//log.Println("Recived a join room request from", client.RemoteAddr().String())
			OnJoinRoom(p, client)
		case LeaveRoomRequest:
			//log.Println("Recived a leave room request from", client.RemoteAddr().String())
			OnLeaveRoom(client, false)
		case ToggleReadyRequest:
			//log.Println("Recived a ready request from", client.RemoteAddr().String())
			OnToggleReady(p, client)
		case GameStartRequest:
			//log.Println("Recived a start game request from", client.RemoteAddr().String())
			OnGameStart(p, client)
		case UpdateSettings:
			//log.Println("Recived a update room setting request from", client.RemoteAddr().String())
			OnUpdateRoom(p, client)
		case OnCloseResultWindow:
			//log.Println("Recived a close resultWindow request from", client.RemoteAddr().String())
			OnCloseResultRequest(p, client)
		case SetUserTeamRequest:
			//log.Println("Recived a set user team request from", client.RemoteAddr().String())
			OnChangeTeam(p, client)
		case GameStartCountdownRequest:
			//log.Println("Recived a begin start game request from", client.RemoteAddr().String())
			OnGameStartCountdown(p, client)
		case Feedback:

		default:
			DebugInfo(2, "Unknown room packet", pkt.InRoomType, "from", client.RemoteAddr().String(), p.Data)
		}
	} else {
		DebugInfo(2, "Error : Recived a illegal room packet from", client.RemoteAddr().String())
	}
}
