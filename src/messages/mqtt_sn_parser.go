package message

import (
	"encoding/binary"
)

func UnMarshall(packet []byte) (msg MqttSnMessage) {
	h := MqttSnHeader{}
	h.UnMarshall(packet)
	switch h.MsgType {
	case MQTTSNT_ADVERTISE:
		msg = NewAdvertise()
		msg.UnMarshall(packet)
	case MQTTSNT_SEARCHGW:
		msg = NewSearchGw()
		msg.UnMarshall(packet)
	case MQTTSNT_GWINFO:
		msg = NewGwInfo()
		msg.UnMarshall(packet)
	case MQTTSNT_CONNECT:
		msg = NewConnect()
		msg.UnMarshall(packet)
	case MQTTSNT_CONNACK:
		msg = NewConnAck()
		msg.UnMarshall(packet)
	case MQTTSNT_WILLTOPICREQ:
		msg = NewWillTopicReq()
		msg.UnMarshall(packet)
	case MQTTSNT_WILLTOPIC:
		msg = NewWillTopic()
		msg.UnMarshall(packet)
	case MQTTSNT_WILLMSGREQ:
		msg = NewWillMsgReq()
		msg.UnMarshall(packet)
	case MQTTSNT_WILLMSG:
		msg = NewWillMsg()
		msg.UnMarshall(packet)
	case MQTTSNT_REGISTER:
		msg = NewRegister(0, 0, "")
		msg.UnMarshall(packet)
	case MQTTSNT_REGACK:
		msg = NewRegAck(0, 0, MQTTSN_RC_ACCEPTED)
		msg.UnMarshall(packet)
	case MQTTSNT_PUBLISH:
		msg = NewPublish()
		msg.UnMarshall(packet)
	case MQTTSNT_PUBACK:
		msg = NewPubAck(0, 0, 0)
		msg.UnMarshall(packet)
	case MQTTSNT_PUBCOMP:
		msg = NewPubComp()
		msg.UnMarshall(packet)
	case MQTTSNT_PUBREC:
		msg = NewPubRec()
		msg.UnMarshall(packet)
	case MQTTSNT_PUBREL:
		msg = NewPubRel()
		msg.UnMarshall(packet)
	case MQTTSNT_SUBSCRIBE:
		msg = NewSubscribe()
		msg.UnMarshall(packet)
	case MQTTSNT_SUBACK:
		msg = NewSubAck(0, 0, MQTTSN_RC_ACCEPTED)
		msg.UnMarshall(packet)
	case MQTTSNT_UNSUBSCRIBE:
		msg = NewUnSubscribe()
		msg.UnMarshall(packet)
	case MQTTSNT_UNSUBACK:
		msg = NewUnSubAck()
		msg.UnMarshall(packet)
	case MQTTSNT_PINGREQ:
		msg = NewPingReq()
		msg.UnMarshall(packet)
	case MQTTSNT_PINGRESP:
		msg = NewPingResp()
		msg.UnMarshall(packet)
	case MQTTSNT_DISCONNECT:
		msg = NewDisConnect(0)
		msg.UnMarshall(packet)
	case MQTTSNT_WILLTOPICUPD:
		msg = NewWillTopicUpd()
		msg.UnMarshall(packet)
	case MQTTSNT_WILLTOPICRESP:
		msg = NewWillTopicResp()
		msg.UnMarshall(packet)
	case MQTTSNT_WILLMSGUPD:
		msg = NewWillMsgUpd()
		msg.UnMarshall(packet)
	case MQTTSNT_WILLMSGRESP:
		msg = NewWillMsgResp()
		msg.UnMarshall(packet)
	}

	return msg
}

func MqttSnFlags(dup uint8, qos uint8, retain bool, will uint8, csession uint8, tidType uint8) uint8 {
	if retain == false {
		return dup << 7 & qos << 5 & 1 << 4 & will << 3 & csession << 2 & tidType
	} else {
		return dup << 7 & qos << 5 & will << 3 & csession << 2 & tidType
	}
}

func DupFlag(flags uint8) bool {
	dup := false
	if (flags >> 7) > 0 {
		dup = true
	}
	return dup
}

func QosFlag(flags uint8) uint8 {
	return flags >> 5 & 0x03
}

func HasRetainFlag(flags uint8) bool {
	retain := false
	if (flags >> 4 & 0x01) > 0 {
		retain = true
	}
	return retain
}

func HasWillFlag(flags uint8) bool {
	will := false
	if (flags >> 3 & 0x01) > 0 {
		will = true
	}
	return will
}

func HasCleanSessionFlag(flags uint8) bool {
	csession := false
	if (flags >> 2 & 0x01) > 0 {
		csession = true
	}
	return csession
}

func TopicIdType(flags uint8) uint8 {
	return flags & 0x03
}

/*********************************************/
/* MqttSnHeader                              */
/*********************************************/
// func NewMqttSnHeader(l uint16, t uint8) *MqttSnHeader {
// 	h := &MqttSnHeader{l, t}
// 	return h
// }

func (m *MqttSnHeader) Marshall() []byte {
	packet := make([]byte, m.Size())
	index := 0
	if m.Size() == 2 {
		packet[index] = (uint8)(m.Length)
		index += 1
	} else if m.Size() == 4 {
		packet[index] = 0x01
		index++
		binary.BigEndian.PutUint16(packet[index:], m.Length)
		index += 2
	}
	packet[index] = m.MsgType

	return packet
}

func (m *MqttSnHeader) UnMarshall(packet []byte) {
	index := 0
	length := packet[index]
	if length == 0x01 {
		index++
		m.Length = binary.BigEndian.Uint16(packet[index:])
		index += 2
	} else {
		m.Length = (uint16)(length)
		index++
	}
	m.MsgType = packet[index]
}

func (m *MqttSnHeader) Size() int {
	size := 4
	// if length is smaller than 256
	if m.Length < 256 {
		size = 2
	}

	return (int)(size)
}

/*********************************************/
/* Advertise                                 */
/*********************************************/
func NewAdvertise() *Advertise {
	m := &Advertise{
		Header: MqttSnHeader{5, MQTTSNT_ADVERTISE}}
	return m
}

func (m *Advertise) Marshall() []byte {
	index := 0
	packet := make([]byte, m.Size())

	hPacket := m.Header.Marshall()
	copy(packet[index:], hPacket)
	index += m.Header.Size()

	packet[index] = m.GwId
	binary.BigEndian.PutUint16(packet[index:], m.Duration)

	return packet
}

func (m *Advertise) UnMarshall(packet []byte) {
	index := 0
	m.Header.UnMarshall(packet[index:])
	index += m.Size()

	m.GwId = packet[index]
	index++

	m.Duration = binary.BigEndian.Uint16(packet[index:])
}

func (m *Advertise) Size() int {
	return 5
}

/*********************************************/
/* SearchGw                                  */
/*********************************************/
func NewSearchGw() *SearchGw {
	m := &SearchGw{
		Header: MqttSnHeader{3, MQTTSNT_SEARCHGW}}
	return m
}

func (m *SearchGw) Marshall() []byte {
	index := 0
	packet := make([]byte, m.Size())

	hPacket := m.Header.Marshall()
	copy(packet[index:], hPacket)
	index += m.Header.Size()

	packet[index] = m.Radius

	return packet
}

func (m *SearchGw) UnMarshall(packet []byte) {
	index := 0
	m.Header.UnMarshall(packet[index:])
	index += m.Header.Size()

	m.Radius = packet[index]
}

func (m *SearchGw) Size() int {
	return 3
}

/*********************************************/
/* GwInfo                                    */
/*********************************************/
func NewGwInfo() *GwInfo {
	m := &GwInfo{
		Header: MqttSnHeader{3, MQTTSNT_GWINFO}}
	return m
}

func (m *GwInfo) Marshall() []byte {
	index := 0
	packet := make([]byte, m.Size())

	hPacket := m.Header.Marshall()
	copy(packet[index:], hPacket)
	index += m.Size()

	packet[index] = m.GwId
	index++

	copy(packet[index:], m.GwAdd)

	return packet
}

func (m *GwInfo) UnMarshall(packet []byte) {
	index := 0
	m.Header.UnMarshall(packet[index:])
	index += m.Header.Size()

	m.GwId = packet[index]
	index++

	addrLen := m.Header.Length - 3
	m.GwAdd = make([]uint8, addrLen)
	copy(m.GwAdd, packet[index:])
}

func (m *GwInfo) Size() int {
	size := 3
	size += len(m.GwAdd)
	return size
}

/*********************************************/
/* Connect                                   */
/*********************************************/
func NewConnect() *Connect {
	m := &Connect{
		Header: MqttSnHeader{6, MQTTSNT_CONNECT}}
	return m
}

func (m *Connect) Marshall() []byte {
	index := 0
	packet := make([]byte, m.Size())

	hPacket := m.Header.Marshall()
	copy(packet[index:], hPacket)
	index += m.Header.Size()

	packet[index] = m.Flags
	index++

	packet[index] = m.ProtocolId
	index++

	binary.BigEndian.PutUint16(packet[index:], m.Duration)
	index += 2

	copy(packet[index:], []byte(m.ClientId))

	return packet
}

func (m *Connect) UnMarshall(packet []byte) {
	index := 0
	m.Header.UnMarshall(packet[index:])
	index += m.Header.Size()

	m.Flags = packet[index]
	index++

	m.ProtocolId = packet[index]
	index++

	m.Duration = binary.BigEndian.Uint16(packet[index:])
	index += 2

	m.ClientId = string(packet[index:m.Header.Length])
}

func (m *Connect) Size() int {
	size := 6
	size += len(m.ClientId)
	return size
}

/*********************************************/
/* ConnAck                                   */
/*********************************************/
func NewConnAck() *ConnAck {
	m := &ConnAck{
		Header: MqttSnHeader{3, MQTTSNT_CONNACK}}
	return m
}

func (m *ConnAck) Marshall() []byte {
	index := 0
	packet := make([]byte, m.Size())

	hPacket := m.Header.Marshall()
	copy(packet[index:], hPacket)
	index += m.Header.Size()

	packet[index] = m.ReturnCode

	return packet
}

func (m *ConnAck) UnMarshall(packet []byte) {
	index := 0
	m.Header.UnMarshall(packet[index:])
	index += m.Size()

	m.ReturnCode = packet[index]
}

func (m *ConnAck) Size() int {
	return 3
}

/*********************************************/
/* WillTopicReq                              */
/*********************************************/
func NewWillTopicReq() *MqttSnHeader {
	m := &MqttSnHeader{2, MQTTSNT_WILLTOPICREQ}
	return m
}

/*********************************************/
/* WillTopic                                 */
/*********************************************/
func NewWillTopic() *WillTopic {
	return nil
}

func (m *WillTopic) Marshall() []byte {
	return nil
}

func (m *WillTopic) UnMarshall(packet []byte) {
}

func (m *WillTopic) Size() int {
	return 0
}

/*********************************************/
/* WillMsgReq                                */
/*********************************************/
func NewWillMsgReq() *MqttSnHeader {
	return nil
}

/*********************************************/
/* WillMsg                                   */
/*********************************************/
func NewWillMsg() *WillMsg {
	return nil
}

func (m *WillMsg) Marshall() []byte {
	return nil
}

func (m *WillMsg) UnMarshall(packet []byte) {
}

func (m *WillMsg) Size() int {
	return 0
}

/*********************************************/
/* Register                                  */
/*********************************************/
func NewRegister(topicId uint16, msgId uint16, topicName string) *Register {
	m := &Register{
		MqttSnHeader{uint16(6 + len(topicName)), MQTTSNT_REGISTER},
		topicId,
		msgId,
		topicName}
	return m
}

func (m *Register) Marshall() []byte {
	index := 0
	packet := make([]byte, m.Size())

	hPacket := m.Header.Marshall()
	copy(packet[index:], hPacket)
	index += m.Header.Size()

	binary.BigEndian.PutUint16(packet[index:], m.TopicId)
	index += 2

	binary.BigEndian.PutUint16(packet[index:], m.MsgId)
	index += 2

	copy(packet[index:], m.TopicName)
	index += 2

	return packet
}

func (m *Register) UnMarshall(packet []byte) {
	index := 0
	m.Header.UnMarshall(packet[index:])
	index += m.Header.Size()

	m.TopicId = binary.BigEndian.Uint16(packet[index:])
	index += 2

	m.MsgId = binary.BigEndian.Uint16(packet[index:])
	index += 2

	m.TopicName = string(packet[index:m.Header.Length])
}

func (m *Register) Size() int {
	return 6 + len(m.TopicName)
}

/*********************************************/
/* RegAck                                    */
/*********************************************/
func NewRegAck(topicId uint16, msgId uint16, rc uint8) *RegAck {
	m := &RegAck{
		MqttSnHeader{7, MQTTSNT_REGACK},
		topicId,
		msgId,
		rc}
	return m
}

func (m *RegAck) Marshall() []byte {
	index := 0
	packet := make([]byte, m.Size())

	hPacket := m.Header.Marshall()
	copy(packet[index:], hPacket)
	index += m.Header.Size()

	binary.BigEndian.PutUint16(packet[index:], m.TopicId)
	index += 2

	binary.BigEndian.PutUint16(packet[index:], m.MsgId)
	index += 2

	packet[index] = m.ReturnCode

	return packet
}

func (m *RegAck) UnMarshall(packet []byte) {
	index := 0
	m.Header.UnMarshall(packet[index:])
	index += m.Header.Size()

	m.TopicId = binary.BigEndian.Uint16(packet[index:])
	index += 2

	m.MsgId = binary.BigEndian.Uint16(packet[index:])
	index += 2

	m.ReturnCode = packet[index]
}

func (m *RegAck) Size() int {
	return 7
}

/*********************************************/
/* Publish                                   */
/*********************************************/
func NewPublish() *Publish {
	return NewPublishNormal(0, false, 0, 0, nil)
}

func NewPublishNormal(qos uint8, retain bool, topicId uint16, msgId uint16, data []uint8) *Publish {
	// make flags
	// dup, qos, retain, will, cleansession, topicIdType
	flags := MqttSnFlags(0, qos, retain, 0, 0, MQTTSN_TIDT_NORMAL)

	// make Publish
	// Header, Flags, TopicId, TopicName, MsgId, Data
	m := &Publish{
		Header:  MqttSnHeader{7 + uint16(len(data)), MQTTSNT_PUBLISH},
		Flags:   flags,
		TopicId: topicId,
		MsgId:   msgId,
		Data:    data}
	return m
}

func NewPublishShortName(qos uint8, retain bool, topic string, msgId uint16, data []uint8) *Publish {
	// make flags
	// dup, qos, retain, will, cleansession, topicIdType
	flags := MqttSnFlags(0, qos, retain, 0, 0, MQTTSN_TIDT_NORMAL)

	// make Publish
	// Header, Flags, TopicId, TopicName, MsgId, Data
	m := &Publish{
		Header:    MqttSnHeader{7 + uint16(len(data)), MQTTSNT_PUBLISH},
		Flags:     flags,
		TopicName: topic,
		MsgId:     msgId,
		Data:      data}
	return m
}

func NewPublishPredefined(qos uint8, retain bool, topicId uint16, msgId uint16, data []uint8) *Publish {
	// make flags
	// dup, qos, retain, will, cleansession, topicIdType
	flags := MqttSnFlags(0, qos, retain, 0, 0, MQTTSN_TIDT_PREDEFINED)

	// make Publish
	// Header, Flags, TopicId, TopicName, MsgId, Data
	m := &Publish{
		Header:  MqttSnHeader{7 + uint16(len(data)), MQTTSNT_PUBLISH},
		Flags:   flags,
		TopicId: topicId,
		MsgId:   msgId,
		Data:    data}
	return m
}

func (m *Publish) Marshall() []byte {
	index := 0
	packet := make([]byte, m.Size())

	hPacket := m.Header.Marshall()
	copy(packet[index:], hPacket)
	index += m.Header.Size()

	packet[index] = m.Flags
	index++

	if (TopicIdType(m.Flags) == MQTTSN_TIDT_PREDEFINED) ||
		(TopicIdType(m.Flags) == MQTTSN_TIDT_NORMAL) {
		// if topic type is TopicId
		binary.BigEndian.PutUint16(packet[index:], m.TopicId)
	} else if TopicIdType(m.Flags) == MQTTSN_TIDT_SHORT_NAME {
		// else if topic type is Short TopicName
		copy(packet[index:], m.TopicName)
	}

	index += 2

	binary.BigEndian.PutUint16(packet[index:], m.MsgId)
	index += 2

	copy(packet[index:], m.Data[0:])

	return packet
}

func (m *Publish) UnMarshall(packet []byte) {
	index := 0
	m.Header.UnMarshall(packet[index:])
	index += m.Header.Size()

	m.Flags = packet[index]
	index++

	if TopicIdType(m.Flags) == MQTTSN_TIDT_NORMAL ||
		TopicIdType(m.Flags) == MQTTSN_TIDT_PREDEFINED {
		m.TopicId = binary.BigEndian.Uint16(packet[index:])
	} else if TopicIdType(m.Flags) == MQTTSN_TIDT_SHORT_NAME {
		m.TopicName = string(packet[index : index+2])
	}
	index += 2

	m.MsgId = binary.BigEndian.Uint16(packet[index:])
	index += 2

	remain := int(m.Header.Length) - (m.Header.Size() + 5)
	m.Data = make([]uint8, remain)
	copy(m.Data, packet[index:m.Header.Length])
}

func (m *Publish) Size() int {
	return m.Header.Size() + 5 + len(m.Data)
}

/*********************************************/
/* PubAck                                    */
/*********************************************/
func NewPubAck(topicId uint16, msgId uint16, rc uint8) *PubAck {
	m := &PubAck{
		MqttSnHeader{7, MQTTSNT_PUBACK}, topicId, msgId, rc}
	return m
}

func (m *PubAck) Marshall() []byte {
	packet := make([]byte, m.Size())

	index := 0
	h_packet := m.Header.Marshall()
	copy(packet[index:], h_packet)
	index += m.Header.Size()

	binary.BigEndian.PutUint16(packet[index:], m.TopicId)
	index += 2

	binary.BigEndian.PutUint16(packet[index:], m.MsgId)
	index += 2

	packet[index] = m.ReturnCode

	return packet
}

func (m *PubAck) UnMarshall(packet []byte) {
	index := 0
	m.Header.UnMarshall(packet[index:])
	index += m.Header.Size()

	m.TopicId = binary.BigEndian.Uint16(packet[index:])
	index += 2

	m.MsgId = binary.BigEndian.Uint16(packet[index:])
	index += 2

	m.ReturnCode = packet[index]
}

func (m *PubAck) Size() int {
	return 7
}

/*********************************************/
/* PubRec                                    */
/*********************************************/
func NewPubRec() *PubRec {
	return nil
}

func (m *PubRec) Marshall() []byte {
	return nil
}

func (m *PubRec) UnMarshall(packet []byte) {
}

func (m *PubRec) Size() int {
	return 0
}

/*********************************************/
/* PubRel                                    */
/*********************************************/
func NewPubRel() *PubRec {
	return nil
}

/*********************************************/
/* PubComp                                   */
/*********************************************/
func NewPubComp() *PubRec {
	return nil
}

/*********************************************/
/* Subscribe                                 */
/*********************************************/
func NewSubscribe() *Subscribe {
	m := &Subscribe{
		Header: MqttSnHeader{0, MQTTSNT_SUBSCRIBE}}
	return m
}

func (m *Subscribe) Marshall() []byte {
	return nil
}

func (m *Subscribe) UnMarshall(packet []byte) {
	index := 0
	m.Header.UnMarshall(packet)
	index += m.Header.Size()

	m.Flags = packet[index]
	index++

	m.MsgId = binary.BigEndian.Uint16(packet[index:])
	index += 2

	switch TopicIdType(m.Flags) {
	// if Topic is Normal or Short Name
	case MQTTSN_TIDT_NORMAL, MQTTSN_TIDT_SHORT_NAME:
		m.TopicName = string(packet[index:m.Header.Length])
		// else if is Id
	case MQTTSN_TIDT_PREDEFINED:
		m.TopicId = binary.BigEndian.Uint16(packet[index:])
	}
}

func (m *Subscribe) Size() int {
	// TODO: implement
	return 0
}

/*********************************************/
/* SubAck                                 */
/*********************************************/
func NewSubAck(topicId uint16, msgId uint16, rc uint8) *SubAck {
	m := &SubAck{
		Header:     MqttSnHeader{8, MQTTSNT_SUBACK},
		TopicId:    topicId,
		MsgId:      msgId,
		ReturnCode: rc}
	return m
}

func (m *SubAck) Marshall() []byte {
	packet := make([]byte, m.Size())
	index := 0

	hPacket := m.Header.Marshall()
	copy(packet[index:], hPacket)
	index += m.Header.Size()

	packet[index] = m.Flags
	index++

	binary.BigEndian.PutUint16(packet[index:], m.TopicId)
	index += 2

	binary.BigEndian.PutUint16(packet[index:], m.MsgId)
	index += 2

	packet[index] = m.ReturnCode

	return packet
}

func (m *SubAck) UnMarshall(packet []byte) {
}

func (m *SubAck) Size() int {
	return 8
}

/*********************************************/
/* UnSubscribe                               */
/*********************************************/
func NewUnSubscribe() *UnSubscribe {
	return nil
}

func (m *UnSubscribe) Marshall() []byte {
	return nil
}

func (m *UnSubscribe) UnMarshall(packet []byte) {
}

func (m *UnSubscribe) Size() int {
	return 0
}

/*********************************************/
/* UnSubAck                                  */
/*********************************************/
func NewUnSubAck() *UnSubAck {
	return nil
}

func (m *UnSubAck) Marshall() []byte {
	return nil
}

func (m *UnSubAck) UnMarshall(packet []byte) {
}

func (m *UnSubAck) Size() int {
	return 0
}

/*********************************************/
/* PingReq                                  */
/*********************************************/
func NewPingReq() *PingReq {
	return nil
}

func (m *PingReq) Marshall() []byte {
	return nil
}

func (m *PingReq) UnMarshall(packet []byte) {
}

func (m *PingReq) Size() int {
	return 0
}

/*********************************************/
/* PingResp                                  */
/*********************************************/
func NewPingResp() *MqttSnHeader {
	return nil
}

/*********************************************/
/* DisConnect                                */
/*********************************************/
func NewDisConnect(duration uint16) *DisConnect {
	length := 2
	if duration > 0 {
		length = 4
	}
	m := &DisConnect{
		Header:   MqttSnHeader{uint16(length), MQTTSNT_DISCONNECT},
		Duration: duration}
	return m
}

func (m *DisConnect) Marshall() []byte {
	index := 0
	packet := make([]byte, m.Size())

	hPacket := m.Header.Marshall()
	copy(packet[index:], hPacket)
	index += m.Header.Size()

	if m.Duration > 0 {
		binary.BigEndian.PutUint16(packet[index:], m.Duration)
	}

	return packet
}

func (m *DisConnect) UnMarshall(packet []byte) {
	index := 0
	m.Header.UnMarshall(packet[index:])
	index += m.Header.Size()

	if m.Header.Length > 2 {
		m.Duration = binary.BigEndian.Uint16(packet[index:])
	}
}

func (m *DisConnect) Size() int {
	size := 2
	if m.Duration > 0 {
		size = 4
	}
	return size
}

/*********************************************/
/* WillTopicUpd                                */
/*********************************************/
func NewWillTopicUpd() *WillTopicUpd {
	return nil
}

func (m *WillTopicUpd) Marshall() []byte {
	return nil
}

func (m *WillTopicUpd) UnMarshall(packet []byte) {
}

func (m *WillTopicUpd) Size() int {
	return 0
}

/*********************************************/
/* WillMsgUpd                                */
/*********************************************/
func NewWillMsgUpd() *WillMsgUpd {
	return nil
}

func (m *WillMsgUpd) Marshall() []byte {
	return nil
}

func (m *WillMsgUpd) UnMarshall(packet []byte) {
}

func (m *WillMsgUpd) Size() int {
	return 0
}

/*********************************************/
/* WillTopicResp                             */
/*********************************************/
func NewWillTopicResp() *WillTopicResp {
	return nil
}

func (m *WillTopicResp) Marshall() []byte {
	return nil
}

func (m *WillTopicResp) UnMarshall(packet []byte) {
}

func (m *WillTopicResp) Size() int {
	return 0
}

/*********************************************/
/* WillMsgResp                             */
/*********************************************/
func NewWillMsgResp() *WillTopicResp {
	return nil
}
