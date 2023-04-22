package bingchat

import "time"

type HandShake struct {
	Protocol string `json:"protocol"`
	Version  int    `json:"version"`
}

type ConversationSession struct {
	ConversationID        string `json:"conversationId"`
	ClientID              string `json:"clientId"`
	ConversationSignature string `json:"conversationSignature"`
	Result                struct {
		Value   string      `json:"value"`
		Message interface{} `json:"message"`
	} `json:"result"`
}

type SendMessage struct {
	Arguments []struct {
		Source              string   `json:"source"`
		OptionsSets         []string `json:"optionsSets"`
		AllowedMessageTypes []string `json:"allowedMessageTypes"`
		SliceIds            []string `json:"sliceIds"`
		Verbosity           string   `json:"verbosity"`
		TraceId             string   `json:"traceId"`
		IsStartOfSession    bool     `json:"isStartOfSession"`
		Message             struct {
			Locale        string `json:"locale"`
			Market        string `json:"market"`
			Region        string `json:"region"`
			Location      string `json:"location"`
			LocationHints []struct {
				Country           string `json:"country"`
				Timezoneoffset    int    `json:"timezoneoffset"`
				CountryConfidence int    `json:"countryConfidence"`
				Center            struct {
					Latitude  float64 `json:"Latitude"`
					Longitude float64 `json:"Longitude"`
				} `json:"Center"`
				RegionType int `json:"RegionType"`
				SourceType int `json:"SourceType"`
			} `json:"locationHints"`
			Timestamp   time.Time `json:"timestamp"`
			Author      string    `json:"author"`
			InputMethod string    `json:"inputMethod"`
			Text        string    `json:"text"`
			MessageType string    `json:"messageType"`
		} `json:"message"`
		ConversationSignature string `json:"conversationSignature"`
		Participant           struct {
			Id string `json:"id"`
		} `json:"participant"`
		ConversationId string `json:"conversationId"`
		// PreviousMessages []struct {
		// 	Text               string        `json:"text"`
		// 	Author             string        `json:"author"`
		// 	AdaptiveCards      []interface{} `json:"adaptiveCards"`
		// 	ContentOrigin      string        `json:"contentOrigin"`
		// 	SuggestedResponses []struct {
		// 		Text          string `json:"text"`
		// 		ContentOrigin string `json:"contentOrigin"`
		// 		MessageType   string `json:"messageType"`
		// 		MessageId     string `json:"messageId"`
		// 		Offense       string `json:"offense"`
		// 	} `json:"suggestedResponses"`
		// 	MessageId   string `json:"messageId"`
		// 	MessageType string `json:"messageType"`
		// } `json:"previousMessages"`
	} `json:"arguments"`
	InvocationId string `json:"invocationId"`
	Target       string `json:"target"`
	Type         int    `json:"type"`
}

type MessageResp struct {
	Type   int    `json:"type"`
	Target string `json:"target"`
	Item   struct {
		Messages []struct {
			MessageType string `json:"messageType,omitempty"`
		} `json:"messages"`
	} `json:"item"`
	Arguments []struct {
		Cursor struct {
			J string `json:"j"`
			P int    `json:"p"`
		} `json:"cursor"`
		Messages []struct {
			Text          string    `json:"text"`
			Author        string    `json:"author"`
			CreatedAt     time.Time `json:"createdAt"`
			Timestamp     time.Time `json:"timestamp"`
			MessageId     string    `json:"messageId"`
			Offense       string    `json:"offense"`
			MessageType   string    `json:"messageType"`
			AdaptiveCards []struct {
				Type    string `json:"type"`
				Version string `json:"version"`
				Body    []struct {
					Type string `json:"type"`
					Text string `json:"text"`
					Wrap bool   `json:"wrap"`
					Size string `json:"size"`
				} `json:"body"`
			} `json:"adaptiveCards"`
			SourceAttributions []interface{} `json:"sourceAttributions"`
			Feedback           struct {
				Tag       interface{} `json:"tag"`
				UpdatedOn interface{} `json:"updatedOn"`
				Type      string      `json:"type"`
			} `json:"feedback"`
			ContentOrigin      string      `json:"contentOrigin"`
			Privacy            interface{} `json:"privacy"`
			SuggestedResponses []struct {
				Text        string    `json:"text"`
				Author      string    `json:"author"`
				CreatedAt   time.Time `json:"createdAt"`
				Timestamp   time.Time `json:"timestamp"`
				MessageId   string    `json:"messageId"`
				MessageType string    `json:"messageType"`
				Offense     string    `json:"offense"`
				Feedback    struct {
					Tag       interface{} `json:"tag"`
					UpdatedOn interface{} `json:"updatedOn"`
					Type      string      `json:"type"`
				} `json:"feedback"`
				ContentOrigin string      `json:"contentOrigin"`
				Privacy       interface{} `json:"privacy"`
			} `json:"suggestedResponses"`
		} `json:"messages"`
		RequestId string `json:"requestId"`
	} `json:"arguments"`
}
