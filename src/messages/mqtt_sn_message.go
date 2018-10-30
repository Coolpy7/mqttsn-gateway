package message

type MqttSnMessage interface {
	Marshall() []byte
	UnMarshall([]byte)
	Size() int
}

const (
	MQTTSN_TIDT_NORMAL     = 0x00
	MQTTSN_TIDT_PREDEFINED = 0x01
	MQTTSN_TIDT_SHORT_NAME = 0x02
)

const (
	MQTTSNT_ADVERTISE     = 0x00
	MQTTSNT_SEARCHGW      = 0x01
	MQTTSNT_GWINFO        = 0x02
	MQTTSNT_CONNECT       = 0x04
	MQTTSNT_CONNACK       = 0x05
	MQTTSNT_WILLTOPICREQ  = 0x06
	MQTTSNT_WILLTOPIC     = 0x07
	MQTTSNT_WILLMSGREQ    = 0x08
	MQTTSNT_WILLMSG       = 0x09
	MQTTSNT_REGISTER      = 0x0a
	MQTTSNT_REGACK        = 0x0b
	MQTTSNT_PUBLISH       = 0x0c
	MQTTSNT_PUBACK        = 0x0d
	MQTTSNT_PUBCOMP       = 0x0e
	MQTTSNT_PUBREC        = 0x0f
	MQTTSNT_PUBREL        = 0x10
	MQTTSNT_SUBSCRIBE     = 0x12
	MQTTSNT_SUBACK        = 0x13
	MQTTSNT_UNSUBSCRIBE   = 0x14
	MQTTSNT_UNSUBACK      = 0x15
	MQTTSNT_PINGREQ       = 0x16
	MQTTSNT_PINGRESP      = 0x17
	MQTTSNT_DISCONNECT    = 0x18
	MQTTSNT_WILLTOPICUPD  = 0x1a
	MQTTSNT_WILLTOPICRESP = 0x1b
	MQTTSNT_WILLMSGUPD    = 0x1c
	MQTTSNT_WILLMSGRESP   = 0x1d
)

const (
	MQTTSN_RC_ACCEPTED                  = 0x00
	MQTTSN_RC_REJECTED_CONGESTION       = 0x01
	MQTTSN_RC_REJECTED_INVALID_TOPIC_ID = 0x02
	MQTTSN_RC_REJECTED_NOT_SUPPORTED    = 0x03
)

type MqttSnHeader struct {
	Length  uint16
	MsgType uint8
}

type Advertise struct {
	Header   MqttSnHeader
	GwId     uint8
	Duration uint16
}

type SearchGw struct {
	Header MqttSnHeader
	Radius uint8
}

type GwInfo struct {
	Header MqttSnHeader
	GwId   uint8
	GwAdd  []uint8
}

type Connect struct {
	Header     MqttSnHeader
	Flags      uint8
	ProtocolId uint8
	Duration   uint16
	ClientId   string
}

type ConnAck struct {
	Header     MqttSnHeader
	ReturnCode uint8
}

type WillTopicReq MqttSnHeader

//type WillTopicReq struct {
//	Length  uint8
//	MsgType uint8
//}

type WillTopic struct {
	Header    MqttSnHeader
	Flags     uint8
	WillTopic string
}

type WillMsgReq MqttSnHeader

type WillMsg struct {
	Header  MqttSnHeader
	WillMsg string
}

type Register struct {
	Header    MqttSnHeader
	TopicId   uint16
	MsgId     uint16
	TopicName string
}

type RegAck struct {
	Header     MqttSnHeader
	TopicId    uint16
	MsgId      uint16
	ReturnCode uint8
}

type Publish struct {
	Header    MqttSnHeader
	Flags     uint8
	TopicId   uint16
	TopicName string
	MsgId     uint16
	Data      []uint8
}

type PubAck struct {
	Header     MqttSnHeader
	TopicId    uint16
	MsgId      uint16
	ReturnCode uint8
}

type PubRec struct {
	Header MqttSnHeader
	MsgId  uint16
}

type PubRel PubRec
type PubComp PubRec

type Subscribe struct {
	Header    MqttSnHeader
	Flags     uint8
	MsgId     uint16
	TopicName string
	TopicId   uint16
}

type SubAck struct {
	Header     MqttSnHeader
	Flags      uint8
	TopicId    uint16
	MsgId      uint16
	ReturnCode uint8
}

type UnSubscribe Subscribe

type UnSubAck struct {
	Header MqttSnHeader
	MsgId  uint16
}

type PingReq struct {
	Header   MqttSnHeader
	ClientId []uint8
}

type PingResp MqttSnHeader

// type PingResp struct {
// 	Length  uint8
// 	MsgType uint8
// }

type DisConnect struct {
	Header   MqttSnHeader
	Duration uint16
}

type WillTopicUpd struct {
	Header    MqttSnHeader
	Flags     uint8
	WillTopic string
}

type WillMsgUpd struct {
	Header  MqttSnHeader
	WillMsg string
}

type WillTopicResp struct {
	Header     MqttSnHeader
	ReturnCode uint8
}

type WillMsgResp WillTopicResp
