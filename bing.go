package bingchat

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const conversationSessionHeader = `    "accept": "application/json",
"accept-language": "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
"sec-ch-ua": "\"Chromium\";v=\"112\", \"Microsoft Edge\";v=\"112\", \"Not:A-Brand\";v=\"99\"",
"sec-ch-ua-arch": "\"x86\"",
"sec-ch-ua-bitness": "\"64\"",
"sec-ch-ua-full-version": "\"112.0.1722.48\"",
"sec-ch-ua-full-version-list": "\"Chromium\";v=\"112.0.5615.121\", \"Microsoft Edge\";v=\"112.0.1722.48\", \"Not:A-Brand\";v=\"99.0.0.0\"",
"sec-ch-ua-mobile": "?0",
"sec-ch-ua-model": "\"\"",
"sec-ch-ua-platform": "\"Windows\"",
"sec-ch-ua-platform-version": "\"15.0.0\"",
"sec-fetch-dest": "empty",
"sec-fetch-mode": "cors",
"sec-fetch-site": "same-origin",
"sec-ms-gec": "673E82A42CAB0AF8C4F97398D164CA4F1F69BEC0D5E41226FD5375F12B17F341",
"sec-ms-gec-version": "1-112.0.1722.48",
"x-ms-client-request-id": "8f0c8a85-28bb-49eb-9c5d-c35dfc5112dd",
"x-ms-useragent": "azsdk-js-api-client-factory/1.0.0-beta.1 core-rest-pipeline/1.10.0 OS/Win32"
"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.48"}`

const wsHeader = `{
    "authority": "edgeservices.bing.com",
    "accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
    "accept-language": "en-US,en;q=0.9",
    "cache-control": "max-age=0",
    "sec-ch-ua": '"Chromium";v="110", "Not A(Brand";v="24", "Microsoft Edge";v="110"',
    "sec-ch-ua-arch": '"x86"',
    "sec-ch-ua-bitness": '"64"',
    "sec-ch-ua-full-version": '"110.0.1587.69"',
    "sec-ch-ua-full-version-list": '"Chromium";v="110.0.5481.192", "Not A(Brand";v="24.0.0.0", "Microsoft Edge";v="110.0.1587.69"',
    "sec-ch-ua-mobile": "?0",
    "sec-ch-ua-model": '""',
    "sec-ch-ua-platform": '"Windows"',
    "sec-ch-ua-platform-version": '"15.0.0"',
    "sec-fetch-dest": "document",
    "sec-fetch-mode": "navigate",
    "sec-fetch-site": "none",
    "sec-fetch-user": "?1",
    "upgrade-insecure-requests": "1",
    "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36 Edg/110.0.1587.69",
    "x-edge-shopping-flag": "1"
  }`

var Timeout = time.Second * 30

type BingChatHub struct {
	sync.Mutex
	wsConn            *websocket.Conn
	cookies           []*http.Cookie
	client            *http.Client
	chatSession       *ConversationSession
	invocationId      int
	sendMessage       *SendMessage
	conversationStyle ConversationStyle
}

func (b *BingChatHub) parseHeader(headerData string) http.Header {
	data := map[string]string{}
	json.Unmarshal([]byte(headerData), &data)
	header := http.Header{}
	for key, value := range data {
		header.Add(key, value)
	}
	return header
}

type IBingChat interface {
	Reset(style ...ConversationStyle)
	SendMessage(msg string) (*MsgResp, error)
	Style() ConversationStyle
}

func NewBingChat(cookiesJson []byte, style ConversationStyle) (IBingChat, error) {
	var cookies []*http.Cookie
	json.Unmarshal(cookiesJson, &cookies)
	return &BingChatHub{
		cookies: cookies,
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
			Timeout: Timeout,
		},
		conversationStyle: style,
	}, nil
}

func (b *BingChatHub) Reset(style ...ConversationStyle) {
	if len(style) > 0 {
		fmt.Println("Switch Style", style[0])
		b.conversationStyle = style[0]
	}
	b.wsConn.Close()
	b.chatSession = nil
	b.invocationId = 0
	b.sendMessage = nil
}

func (b *BingChatHub) Style() ConversationStyle {
	return b.conversationStyle
}

func (b *BingChatHub) createConversation() error {
	req, err := http.NewRequest("GET", "https://www.bing.com/turing/conversation/create", nil)
	if err != nil {
		return err
	}
	req.Header = b.parseHeader(conversationSessionHeader)
	req.Header.Set("x-ms-client-request-id", uuid.New().String())
	for _, cookie := range b.cookies {
		req.AddCookie(cookie)
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	b.chatSession = &ConversationSession{}
	err = json.NewDecoder(resp.Body).Decode(b.chatSession)
	if err != nil {
		return err
	}
	return nil
}

func (b *BingChatHub) initWsConnect() error {
	dial := websocket.DefaultDialer
	dial.Proxy = http.ProxyFromEnvironment
	dial.HandshakeTimeout = Timeout
	dial.EnableCompression = true

	dial.TLSClientConfig = &tls.Config{}
	conn, resp, err := dial.Dial("wss://sydney.bing.com/sydney/ChatHub", b.parseHeader(wsHeader))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusSwitchingProtocols {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	b.wsConn = conn
	err = conn.WriteMessage(websocket.BinaryMessage, []byte(`{"protocol":"json","version":1}`+DELIMITER))
	if err != nil {
		return fmt.Errorf("write json response: %v", err)
	}
	_, _, err = conn.NextReader()
	go func() {
		for {
			b.Lock()
			err := conn.WriteMessage(websocket.BinaryMessage, []byte(`{"type":6}`+DELIMITER))
			b.Unlock()
			if err != nil {
				break
			}
			time.Sleep(time.Second * 5)
		}
	}()
	return err
}

type MsgResp struct {
	Suggest []string
	Notify  chan string
	Msg     string
}

func (b *BingChatHub) SendMessage(msg string) (*MsgResp, error) {
	if b.chatSession == nil {
		err := b.createConversation()
		if err != nil {
			log.Println("create conversation error: ", err)
			return nil, err
		}
	}
	err := b.initWsConnect()
	if err != nil {
		return nil, err
	}
	if b.sendMessage == nil {
		b.sendMessage = b.conversationStyle.TmpMessage()
		b.sendMessage.Arguments[0].ConversationSignature = b.chatSession.ConversationSignature
		b.sendMessage.Arguments[0].Participant.Id = b.chatSession.ClientID
		b.sendMessage.Arguments[0].ConversationId = b.chatSession.ConversationID
	}
	b.sendMessage.Arguments[0].TraceId = b.getTraceId()
	b.sendMessage.Arguments[0].IsStartOfSession = b.invocationId == 0
	b.sendMessage.Arguments[0].Message.Text = msg
	b.sendMessage.Arguments[0].Message.Timestamp = time.Now()
	b.sendMessage.InvocationId = fmt.Sprint(b.invocationId)
	b.invocationId += 1
	msgData, _ := json.Marshal(b.sendMessage)
	b.Lock()
	err = b.wsConn.WriteMessage(websocket.BinaryMessage, append(msgData, []byte(DELIMITER)...))
	b.Unlock()
	if err != nil {
		return nil, err
	}
	msgRespChannel := &MsgResp{
		Notify: make(chan string, 1),
	}
	go func() {
		var startRev bool
		lastMsg := ""
		defer close(msgRespChannel.Notify)
		for {
			_, data, err := b.wsConn.ReadMessage()
			if err != nil {
				log.Println(err)
				b.Reset()
				break
			}
			if len(data) == 0 {
				continue
			}
			spData := bytes.Split(data, []byte(DELIMITER))
			if len(spData) == 0 {
				continue
			}
			data = spData[0]
			resp := MessageResp{}
			_ = json.Unmarshal(data, &resp)

			for _, message := range resp.Item.Messages {
				if message.MessageType == "Disengaged" {
					b.Reset()

					return
				}
			}

			if resp.Type == 1 && len(resp.Arguments) > 0 && resp.Arguments[0].Cursor.J != "" {
				startRev = true
				continue
			}
			if !startRev {
				continue
			}
			if resp.Type == 1 && len(resp.Arguments) > 0 && len(resp.Arguments[0].Messages) > 0 {
				if resp.Arguments[0].Messages[0].SuggestedResponses != nil {
					var suggests []string
					for _, suggest := range resp.Arguments[0].Messages[0].SuggestedResponses {
						suggests = append(suggests, suggest.Text)
					}
					msgRespChannel.Suggest = suggests
				}

				if resp.Arguments[0].Messages[0].MessageType == "Disengaged" {
					b.Reset()

					break
				}
				msg := strings.TrimSpace(resp.Arguments[0].Messages[0].Text)
				msgRespChannel.Msg = msg
				if len(lastMsg) > len(msg) {
					continue
				}
				if msg == "" || msg[len(lastMsg):] == "" {
					continue
				}
				msgRespChannel.Notify <- msg[len(lastMsg):]
				lastMsg = msg
			}
			if resp.Type == 2 {
				b.wsConn.Close()
				break
			}
		}

	}()

	return msgRespChannel, nil
}

func (b *BingChatHub) getTraceId() string {
	rand.Seed(time.Now().UnixNano())
	length := 32
	bytes := make([]byte, length)
	str := "0123456789abcdef"
	for i := 0; i < length; i++ {
		bytes[i] = byte(str[rand.Intn(len(str))])
	}
	return string(bytes)
}
