/*struct of user packet {
	basePacket 			4 bytes
	type 				1 byte
	lenOfNexonUsername 	1 byte
	nexonUsername 		len bytes
	lenOfGameUsername 	1 byte
	gameUsername 		len bytes
	unknown01			1 byte
	lenOfPassWd			2 bytes
	PassWd				len bytes
	HddHwid				16 bytes
	netCafeId 			4 bytes
	unknown02			4 bytes
	userSn 				8 bytes
	lenOfUnknownString	2 bytes
	UnknownString03		len bytes
	unknown04 			1 byte
	isLeague 			1 byte
	{ always null ignore it /  包里好像不存在，忽视
		unk04 				1 byte
		lenOfUnknown05 		1 byte
		UnknownString05		len bytes
		lenOfUnknown06 		1 byte
		UnknownString06		len bytes
	}
	lenOfString			1 bytes
	String				len bytes
}*/

package main

import (
	"log"
	"net"
)

const (
	//MAXUSERNUM 最大用户数
	MAXUSERNUM = 1024 //房间状态

	UserNotReady = 0
	UserIngame   = 1
	UserReady    = 2
)

//全局用户管理
type userManager struct {
	userNum int
	users   []user
}

type user struct {
	//个人信息
	userid               uint32
	loginName            []byte
	username             []byte
	password             []byte
	level                uint16
	rank                 uint8
	rankFrame            uint8
	points               uint64
	currentExp           uint64
	maxExp               uint64
	playedMatches        uint32
	wins                 uint32
	kills                uint32
	headshots            uint32
	deaths               uint32
	assists              uint32
	accuracy             uint16
	secondsPlayed        uint32
	netCafeName          []byte
	cash                 uint32
	clanName             []byte
	clanMark             uint32
	worldRank            uint32
	mpoints              uint32
	titleId              uint16
	unlockedTitles       []byte
	signature            []byte
	bestGamemode         uint32
	bestMap              uint32
	unlockedAchievements []byte
	avatar               uint16
	unlockedAvatars      []byte
	vipLevel             uint8
	vipXp                uint32
	skillHumanCurXp      uint64
	skillHumanMaxXp      uint64
	skillHumanPoints     uint8
	skillZombieCurXp     uint64
	skillZombieMaxXp     uint64
	skillZombiePoints    uint8
	//连接
	currentConnection net.Conn
	//频道房间信息
	currentChannelServerIndex uint8
	currentChannelIndex       uint8
	currentRoomId             uint16
	currentTeam               uint8
	currentstatus             uint8
	currentIsIngame           bool

	//仓库信息
	//inventory userInventory
}

func addUser(src *user) bool {
	if (*src).userid == 0 {
		log.Fatalln("ID of User", (*src).username, "is illegal !")
		return false
	}
	if UserManager.userNum > MAXUSERNUM {
		log.Fatalln("Online users is too more to login !")
		return false
	}
	for _, v := range UserManager.users {
		if v.userid == (*src).userid {
			log.Fatalln("User is already logged in !")
			return false
		}
	}
	UserManager.userNum++
	UserManager.users = append(UserManager.users, *src)
	return true
}

func delUser(src *user) bool {
	if (*src).userid == 0 {
		log.Fatalln("ID of User", (*src).username, "is illegal !")
		return false
	}
	if UserManager.userNum == 0 {
		log.Fatalln("There is no online user !")
		return false
	}
	for i, v := range UserManager.users {
		if v.userid == (*src).userid {
			UserManager.users = append(UserManager.users[:i], UserManager.users[i+1:]...)
			UserManager.userNum--
			return true
		}
	}
	return false
}

func delUserWithConn(con net.Conn) bool {
	if UserManager.userNum == 0 {
		log.Fatalln("There is no online user !")
		return false
	}
	for i, v := range UserManager.users {
		if v.currentConnection == con {
			UserManager.users = append(UserManager.users[:i], UserManager.users[i+1:]...)
			UserManager.userNum--
			return true
		}
	}
	return false
}
func getNewUserID() uint32 {
	if UserManager.userNum > MAXUSERNUM {
		log.Fatalln("Online users is too much , unable to get a new id !")
		//ID=0 是非法的
		return 0
	}
	var intbuf [MAXUSERNUM + 2]uint32
	//哈希思想
	for i := 0; i < int(UserManager.userNum); i++ {
		intbuf[UserManager.users[i].userid] = 1
	}
	//找到空闲的ID
	for i := 1; i < int(MAXUSERNUM+2); i++ {
		if intbuf[i] == 0 {
			//找到了空闲ID
			return uint32(i)
		}
	}
	return 0
}

//假定nexonUsername是唯一
func getUserByLogin(p loginPacket) user {
	u := findOnlineUserByName(p.gameUsername)
	if u.userid <= 0 {
		return getUserFromDatabase(p)
	}
	return *u
}

// 理论上来讲是要从数据库里面读取用户
// 但是目前暂时先不考虑数据库，以后加进，
// 目前先发送新用户数据getNewUser()
func findOnlineUserByName(name []byte) *user {
	l := len(name)
	if l <= 0 {
		log.Fatalln("User name is illegal !")
		u := getNewUser()
		return &u
	}
	for _, v := range UserManager.users {
		if string(v.username) == string(name) {
			return &v
		}
	}
	u := getNewUser()
	return &u
}

func (u user) isVIP() bool {
	if u.vipLevel <= 0 {
		return false
	}
	return true
}

func (u *user) setID(id uint32) {
	(*u).userid = id
}

func (u *user) setUserName(name []byte) {
	(*u).loginName = name
	(*u).username = name
}

func (u *user) setUserChannelServer(id uint8) {
	(*u).currentChannelServerIndex = id
}

func (u *user) setUserChannel(id uint8) {
	(*u).currentChannelIndex = id
}

func (u *user) setUserRoom(id uint16) {
	(*u).currentRoomId = id
}

func (u *user) quitChannel() {
	(*u).currentChannelIndex = 0
}

func (u *user) quitRoom() {
	(*u).currentRoomId = 0
}

func getNewUser() user {
	return user{
		0,
		[]byte{},        //loginname
		[]byte{},        //username,looks can change it to another name
		[]byte{},        //passwd
		1,               //level
		0,               //rank
		0,               //rankframe
		0x7AF3,          //points
		0,               //curEXP
		1000,            //maxEXP
		0x0A,            //playermatchs
		0,               //wins
		0,               //kills
		0,               //headshots
		0,               //deaths
		0,               //assists
		0x0A,            // accuracy
		0x290C,          // secondsPlayed
		newNullString(), // netCafeName
		0,               // cash
		newNullString(), // clanName
		0,               // clanMark
		0,               // worldRank
		0,               // mpoints
		0,               // titleId
		[]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // unlockedTitles
		newNullString(), // signature
		5,               // bestGamemode
		9,               // bestMap
		[]uint8{0x00, 0x00, 0x18, 0x08, 0x00, 0x00, 0x00, 0x00, 0x42, 0x02,
			0x18, 0xC0, 0x09, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0xC0, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0xC8, 0xB7, 0x08, 0x00, 0x00, 0x04, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // unlockedAchievements
		1006, // avatar
		[]uint8{0x00, 0x00, 0x18, 0x08, 0x00, 0x00, 0x00, 0x00, 0x42, 0x02,
			0x18, 0xC0, 0x09, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0xC0, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0xC8, 0xB7, 0x08, 0x00, 0x00, 0x04, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3E, 0x00, 0x00}, // unlockedAvatars
		0,      //viplevel
		0,      //vipXp
		0x02FB, //skillHumanCurXp
		0x19AC, //skillHumanMaxXp
		0,      //skillHumanPoints
		0,      //skillZombieCurXp
		0x16F6, //skillZombieMaxXp
		0,      //skillZombiePoints
		nil,    //connection
		1,      //serverid
		0,      //channelid
		0,      //roomid
		0,      //currentTeam
		0,      //currentstatus
		false,  //currentIsIngame
	}
}

//暂定功能
//从数据库中读取用户数据
//如果是新用户则保存到数据库中
func getUserFromDatabase(p loginPacket) user {
	u := getNewUser()
	u.setID(getNewUserID())
	u.setUserName(p.gameUsername)
	u.password = p.PassWd
	return u
}

//通过连接获取用户
func getUserFromConnection(client net.Conn) *user {
	if UserManager.userNum <= 0 {
		return nil
	}
	for k, v := range UserManager.users {
		if v.currentConnection == client {
			return &UserManager.users[k]
		}
	}
	return nil
}

//通过ID获取用户
func getUserFromID(id uint32) *user {
	if UserManager.userNum <= 0 {
		return nil
	}
	for k, v := range UserManager.users {
		if v.userid == id {
			return &UserManager.users[k]
		}
	}
	return nil
}

//获取用户所在分区服务器ID
func (u user) getUserChannelServerID() uint8 {
	if u.userid <= 0 {
		return 0
	}
	return u.currentChannelServerIndex
}

//获取用户所在频道ID
func (u user) getUserChannelID() uint8 {
	if u.userid <= 0 {
		return 0
	}
	return u.currentChannelIndex
}

//获取用户所在房间ID
func (u user) getUserRoomID() uint16 {
	if u.userid <= 0 {
		return 0
	}
	return u.currentRoomId
}

func (u user) getUserTeam() uint8 {
	return u.currentTeam
}