const (
	HEART_BEAT = 5005
	RECV_TXT_MSG = 1
	RECV_PIC_MSG = 3
	USER_LIST = 5000
	GET_USER_LIST_SUCCSESS = 5001
	GET_USER_LIST_FAIL = 5002
	TXT_MSG = 555
	PIC_MSG = 500
	AT_MSG = 550
	CHATROOM_MEMBER = 5010
	CHATROOM_MEMBER_NICK = 5020
	PERSONAL_INFO = 6500
	DEBUG_SWITCH = 6000
	PERSONAL_DETAIL = 6550
	DESTROY_ALL = 9999
)

type TextMSG struct {
    Id string  `json:"id"`
	Type int   `json:"type"`
	Content string `json:"content"`
	Wxid    string 'json:"wxid"'
}

type TextMSG struct {
    Id string  `json:"id"`
	Type int   `json:"type"`
	Content string `json:"content"`
	Wxid    string 'json:"wxid"'
}