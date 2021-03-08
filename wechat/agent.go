package main
import (
	"log"
	"time"
	"encoding/json"
	"github.com/gorilla/websocket"
)

type Agent struct {
	Url string 
	Conn *websocket.Conn
}

type TextMSG struct {
    Id string  `json:"id"`
	Type int   `json:"type"`
	Content string `json:"content"`
	Wxid    string `json:"wxid"`
}

func (a *Agent)SendMsg(msg TextMSG)  {
	var msgBytes []byte
	msgBytes, _ = json.Marshal(msg)
	a.Conn.WriteMessage(websocket.TextMessage, msgBytes)
}

func (a *Agent)ReadMsg()  {
	for {
		_, message, err := a.Conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)
	}
}
func (a *Agent)Start()  {
	c, _, err := websocket.DefaultDialer.Dial(a.Url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	a.Conn = c
	go a.ReadMsg()
}

func (a *Agent)Close()  {
	a.Conn.Close()
}

func NewAgent(url string)  *Agent {
	agent := &Agent{
		Url: url,
	}
	return agent
}

func SendLaopo(a *Agent)  {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		msg := TextMSG {
			Id: time.Now().Format("2006-01-02 15:04:05"),
			Type: 555,
			Content: "我爱你",
			Wxid: "wxid_k1jw26novqse21",
		}
		a.SendMsg(msg)
    }
}
func main()  {
	a := NewAgent("ws://192.168.1.5:5555")
	a.Start()

	// msg := TextMSG {
	// 	Id: time.Now().Format("2006-01-02 15:04:05"),
	// 	Type: 5000,
	// 	Content: "user list",
	// 	Wxid: "null",
	// }
	// a.SendMsg(msg)
	// go SendLaopo(a)
	for {}
}