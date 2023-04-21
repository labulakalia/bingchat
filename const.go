package bingchat

import "encoding/json"

const balanceStyle = `{"arguments":[{"source":"cib","optionsSets":["nlu_direct_response_filter","deepleo","disable_emoji_spoken_text","responsible_ai_policy_235","enablemm","galileo","rcallowlist","responseos","jb090","jbfv202","dv3sugg"],"allowedMessageTypes":["Chat","InternalSearchQuery","InternalSearchResult","Disengaged","InternalLoaderMessage","RenderCardRequest","AdsQuery","SemanticSerp","GenerateContentQuery","SearchQuery"],"sliceIds":["contctxp2tf","delayglobjscf","0417bicunivs0","ssoverlap50","sspltop5","sswebtop1","audseq","sbsvgopt","nopreloadsstf","winlongmsg2tf","perfimpcomb","sugdivdis","sydnoinputt","wpcssopt","414suggs0","scctl","418glpv6ps0","417rcallow","321slocs0","407pgparsers0","0329resp","asfixescf","udscahrfoncf","414jbfv202"],"verbosity":"verbose","traceId":"6441f4712452428ab53b745af65c5089","isStartOfSession":true,"message":{"locale":"zh-CN","market":"zh-CN","region":"WW","location":"lat:47.639557;long:-122.128159;re=1000m;","locationHints":[{"country":"Singapore","timezoneoffset":8,"countryConfidence":8,"Center":{"Latitude":1.2929,"Longitude":103.8547},"RegionType":2,"SourceType":1}],"timestamp":"2023-04-21T10:27:06+08:00","author":"user","inputMethod":"Keyboard","text":"我需要帮助制定计划","messageType":"Chat"},"conversationSignature":"AdctQf6lU2LbVhUuyTMDCcchrXUbEaX2jyBbeQ2iXuY=","participant":{"id":"914798353003051"},"conversationId":"51D|BingProd|92EBE954F84335CBF36EE4E86BF8A028765E8922A38BC8A1E248457AF7342CA3"}],"invocationId":"","target":"chat","type":4}`
const createStyle = `{"arguments":[{"source":"cib","optionsSets":["nlu_direct_response_filter","deepleo","disable_emoji_spoken_text","responsible_ai_policy_235","enablemm","h3imaginative","rcallowlist","responseos","jb090","jbfv202","dv3sugg","clgalileo","gencontentv3"],"allowedMessageTypes":["Chat","InternalSearchQuery","InternalSearchResult","Disengaged","InternalLoaderMessage","RenderCardRequest","AdsQuery","SemanticSerp","GenerateContentQuery","SearchQuery"],"sliceIds":["contctxp2tf","delayglobjscf","0417bicunivs0","ssoverlap50","sspltop5","sswebtop1","audseq","sbsvgopt","nopreloadsstf","winlongmsg2tf","perfimpcomb","sugdivdis","sydnoinputt","wpcssopt","414suggs0","scctl","418glpv6ps0","417rcallow","321slocs0","407pgparsers0","0329resp","asfixescf","udscahrfoncf","414jbfv202"],"verbosity":"verbose","traceId":"6441f4712452428ab53b745af65c5089","isStartOfSession":true,"message":{"locale":"zh-CN","market":"zh-CN","region":"WW","location":"lat:47.639557;long:-122.128159;re=1000m;","locationHints":[{"country":"Singapore","timezoneoffset":8,"countryConfidence":8,"Center":{"Latitude":1.2929,"Longitude":103.8547},"RegionType":2,"SourceType":1}],"timestamp":"2023-04-21T10:27:06+08:00","author":"user","inputMethod":"Keyboard","text":"告诉我的星座","messageType":"Chat"},"conversationSignature":"uk3kLopdE2Zb8nTXHFx/smV2IWyec3G11B0y8ehSC4k=","participant":{"id":"914798353003051"},"conversationId":"51D|BingProd|29F20DF6A2946BAD80F2B98E87138C4E14A6D8C5D29E95056F41D5BE3539D4B4"}],"invocationId":"","target":"chat","type":4}`
const preciseStyle = `{"arguments":[{"source":"cib","optionsSets":["nlu_direct_response_filter","deepleo","disable_emoji_spoken_text","responsible_ai_policy_235","enablemm","h3precise","rcallowlist","responseos","jb090","jbfv202","dv3sugg","clgalileo"],"allowedMessageTypes":["Chat","InternalSearchQuery","InternalSearchResult","Disengaged","InternalLoaderMessage","RenderCardRequest","AdsQuery","SemanticSerp","GenerateContentQuery","SearchQuery"],"sliceIds":["contctxp2tf","delayglobjscf","0417bicunivs0","ssoverlap50","sspltop5","sswebtop1","audseq","sbsvgopt","nopreloadsstf","winlongmsg2tf","perfimpcomb","sugdivdis","sydnoinputt","wpcssopt","414suggs0","scctl","418glpv6ps0","417rcallow","321slocs0","407pgparsers0","0329resp","asfixescf","udscahrfoncf","414jbfv202"],"verbosity":"verbose","traceId":"6441f4712452428ab53b745af65c5089","isStartOfSession":true,"message":{"locale":"zh-CN","market":"zh-CN","region":"WW","location":"lat:47.639557;long:-122.128159;re=1000m;","locationHints":[{"country":"Singapore","timezoneoffset":8,"countryConfidence":8,"Center":{"Latitude":1.2929,"Longitude":103.8547},"RegionType":2,"SourceType":1}],"timestamp":"2023-04-21T10:27:06+08:00","author":"user","inputMethod":"Keyboard","text":"我需要帮助做研究","messageType":"Chat"},"conversationSignature":"1F8e/oVRPtqkMq+/hrKWphxvXbc5DTQTsItUsoaxedE=","participant":{"id":"914798353003051"},"conversationId":"51D|BingProd|23B9F05272D0D7471D94F332A995F6996B9E192846CB8B7007092B4B6DE6FDEC"}],"invocationId":"14","target":"chat","type":4}`
const DELIMITER = "\x1e"

type ConversationStyle uint8

const (
	ConversationCreateStyle ConversationStyle = iota + 1
	ConversationBalanceStyle
	ConversationPreciseStyle
)

func (c ConversationStyle) String() string {
	switch c {
	case ConversationBalanceStyle:
		return "Balance"
	case ConversationCreateStyle:
		return "Create"
	case ConversationPreciseStyle:
		return "Precise"
	}
	return ""
}

func (c ConversationStyle) TmpMessage() *SendMessage {
	var data string
	switch c {
	case ConversationBalanceStyle:
		data = balanceStyle
	case ConversationCreateStyle:
		data = createStyle
	case ConversationPreciseStyle:
		data = preciseStyle
	}
	msg := SendMessage{}
	json.Unmarshal([]byte(data), &msg)
	return &msg
}
