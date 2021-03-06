package typestruct

const (
	MessageCongratulate    = 11
	MessageNoticeImportant = 20
	MessageDialogBox       = 21
	MessageNotice          = 22
	MessageNotice2         = 23
	MessageNotice3         = 24
	MessageDialogBoxExit   = 60
)

var (
	GAME_ROOM_LEAVE_EARLY                = []byte("你因提前退出而被扣除1000Points!")
	GAME_SERVER_ERROR                    = []byte("你与服务端的连接遭遇到不可恢复性错误，请联系服务器管理员并提供报错信息以便管理员查找和修复错误！")
	GAME_LOGIN_ALREADY                   = []byte("您的账号当前已经有人登录！如果你的账号已泄露请联系管理员！")
	GAME_LOGIN_ERROR                     = []byte("登录过程发生错误，请联系服务器管理员并提供报错信息以便管理员查找和修复错误！")
	GAME_ROOM_COUNT_MODE_ERROR           = []byte("无法开始游戏！请检查你的房间设置，比如是否有2人及以上或者开启bot")
	GAME_ROOM_JOIN_ERROR                 = []byte("加入房间发生错误！")
	GAME_ROOM_JOIN_FAILED_CLOSED         = []byte("#CSO2_POPUP_ROOM_JOIN_FAILED_CLOSED")
	GAME_ROOM_JOIN_FAILED_FULL           = []byte("#CSO2_POPUP_ROOM_JOIN_FAILED_FULL")
	GAME_ROOM_JOIN_FAILED_BAD_PASSWORD   = []byte("#CSO2_POPUP_ROOM_JOIN_FAILED_INVALID_PASSWD")
	GAME_ROOM_CHANGETEAM_FAILED          = []byte("#CSO2_POPUP_ROOM_CHANGETEAM_FAILED")
	GAME_ROOM_COUNTDOWN_FAILED_NOENEMIES = []byte("#CSO2_UI_ROOM_COUNTDOWN_FAILED_NOENEMY")
	GAME_LOGIN_BAD_USERINFO              = []byte("#CSO2_LoginAuth_Certify_NoPassport")
	GAME_LOGIN_BAD_PASSWORD              = []byte("#CSO2_LoginAuth_WrongPassword")
	GAME_LOGIN_BAD_USERNAME              = []byte("#CSO2_LoginAuth_UserNotExists")
	GAME_LOGIN_TENTH_FAILED              = []byte("#CSO2_LoginAuth_TempBlockedByLoginFail")
	GAME_LOGIN_INVALID_USERINFO          = []byte("#CSO2_ServerMessage_INVALID_USERINFO")
)
