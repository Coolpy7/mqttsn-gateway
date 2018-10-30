package mqttsn

import (
	"env"
	"errors"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"messages"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

/**
 *
 */
const (
	MQTTSN_SESSION_ACTIVE = iota
	MQTTSN_SESSION_ASLEEP
	MQTTSN_SESSION_DISCONNECTED
)

/**
 *
 */
type MqttSnSession struct {
	mutex         sync.RWMutex
	ClientId      string
	Conn          *net.UDPConn
	Remote        *net.UDPAddr
	Topics        *TopicMap
	msgId         *ManagedId
	duration      uint16
	state         byte
	tatCalculator *TatCalculator
}

/**
 *
 */
type TransparentSnSession struct {
	mutex sync.RWMutex
	*MqttSnSession
	mqttClient         MQTT.Client
	brokerHost         string
	brokerPort         int
	brokerUser         string
	brokerPassword     string
	sendBuffer         chan message.MqttSnMessage
	shutDown           chan bool
	statisticsReporter *StatisticsReporter
}

/**
 *
 */
func NewMqttSnSession(id string, conn *net.UDPConn, remote *net.UDPAddr) *MqttSnSession {
	s := &MqttSnSession{
		sync.RWMutex{},
		id,
		conn,
		remote,
		NewTopicMap(),
		&ManagedId{},
		0,
		MQTTSN_SESSION_ACTIVE,
		NewTatCalculator(),
	}
	return s
}

/**
 *
 */
func (s *MqttSnSession) StoreTopic(topicName string) uint16 {
	return s.Topics.StoreTopic(topicName)
}

/**
 *
 */
func (s *MqttSnSession) StoreTopicWithId(topicName string, id uint16) bool {
	return s.Topics.StoreTopicWithId(topicName, id)
}

/**
 *
 */
func (s *MqttSnSession) LoadTopic(topicId uint16) (string, bool) {
	return s.Topics.LoadTopic(topicId)
}

/**
 *
 */
func (s *MqttSnSession) LoadTopicId(topicName string) (uint16, bool) {
	return s.Topics.LoadTopicId(topicName)
}

/**
 *
 */
func (s *MqttSnSession) FreeMsgId(i uint16) {
	s.msgId.FreeId(i)
}

/**
 *
 */
func (s *MqttSnSession) NextMsgId() uint16 {
	return s.msgId.NextId()
}

/**
 *
 */
func NewTransparentSnSession(
	id string,
	conn *net.UDPConn,
	remote *net.UDPAddr,
	host string,
	port int,
	user string,
	password string,
	queueSize int,
	statisticsReporter *StatisticsReporter) *TransparentSnSession {
	s := &TransparentSnSession{
		sync.RWMutex{},
		&MqttSnSession{
			sync.RWMutex{},
			id,
			conn,
			remote,
			NewTopicMap(),
			&ManagedId{},
			0,
			MQTTSN_SESSION_ACTIVE,
			NewTatCalculator(),
		},
		nil, host, port, user, password,
		make(chan message.MqttSnMessage, queueSize),
		make(chan bool, 1),
		statisticsReporter}
	return s
}

/**
 *
 */
func (s *TransparentSnSession) connLostHandler(
	c MQTT.Client, err error) {
	log.Println("ERROR : MQTT connection is lost with ", err)
}

/**
 *
 */
func (s *TransparentSnSession) ConnectToBroker(cleanSession bool) error {
	log.Println("Connect to broker")

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// create opts
	addr := "tcp://" + s.brokerHost + ":" + strconv.Itoa(s.brokerPort)
	opts := MQTT.NewClientOptions().AddBroker(addr)
	opts.SetClientID(s.MqttSnSession.ClientId)
	opts.SetUsername(s.brokerUser)
	opts.SetPassword(s.brokerPassword)
	opts.SetAutoReconnect(false)
	opts.SetConnectionLostHandler(
		s.connLostHandler)
	opts.SetCleanSession(cleanSession)

	// create client instance
	s.mqttClient = (MQTT.NewClient(opts))

	// connect, timeout time is 15 sec
	if token := s.mqttClient.Connect(); !token.WaitTimeout(15 * time.Second) {
		return errors.New("Connect to broker is Timeout.")
	} else if token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (s *TransparentSnSession) Disconnect() {
	log.Println("Disconnect from broker")

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.mqttClient.Disconnect(250)
}

/**
 *
 */
func (s *TransparentSnSession) OnPublish(client MQTT.Client, msg MQTT.Message) {
	if env.DEBUG {
		log.Println("on publish. Receive message from broker.")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// get topic
	topic := msg.Topic()

	// TODO: check TransparentSnSession's state.
	// if session is sleep, gateway must buffer the message.

	var m *message.Publish
	msgId := s.NextMsgId()

	// get predef topicid
	topicId, ok := LoadPredefTopicId(topic)
	if ok {
		// process as predef topicid
		// qos, retain, topicId, msgId, data
		// Now, qos is hard coded as 1
		m = message.NewPublishPredefined(
			1, false, topicId, msgId, msg.Payload())
	} else if topicId, ok = s.LoadTopicId(topic); ok {
		// process as fixed topicid
		// qos, retain, topicId, msgId, data
		// Now, qos is hard coded as 1
		m = message.NewPublishNormal(
			1, false, topicId, msgId, msg.Payload())
	} else {
		// TODO: implement wildcarded or short topic route

		log.Println("WARN : topic id was not found for ", topic, ".")
		return
	}

	// send message
	s.Conn.WriteToUDP(m.Marshall(), s.Remote)
}

/**
 *
 */
func (s *TransparentSnSession) sendMqttMessageLoop() {
	for {
		select {
		case msg := <-s.sendBuffer:
			switch msgi := msg.(type) {
			case *message.Publish:
				s.doPublish(msgi)
			case *message.Subscribe:
				s.doSubscribe(msgi)
			}
		case <-s.shutDown:
			// wait 100 ms and disconnect
			if s.mqttClient.IsConnected() {
				s.mqttClient.Disconnect(100)
			}
			return
		}
	}
}

/**
 *
 */
func (s *TransparentSnSession) doPublish(m *message.Publish) {
	if env.DEBUG {
		log.Println("Publish to broker")
	}

	var topicName string
	var ok bool
	if message.TopicIdType(m.Flags) == message.MQTTSN_TIDT_NORMAL {
		// search topic name from topic id
		topicName, ok = s.LoadTopic(m.TopicId)
		if ok == false {
			// error handling
			log.Println("ERROR : topic was not found.")
			puback := message.NewPubAck(
				m.TopicId,
				m.MsgId,
				message.MQTTSN_RC_REJECTED_INVALID_TOPIC_ID)
			// send
			s.Conn.WriteToUDP(puback.Marshall(), s.Remote)

			return
		}
	} else if message.TopicIdType(m.Flags) == message.MQTTSN_TIDT_PREDEFINED {
		// search topic name from topic id
		topicName, ok = LoadPredefTopic(m.TopicId)
		if ok == false {
			// error handling
			log.Println("ERROR : topic was not found.")
			puback := message.NewPubAck(
				m.TopicId,
				m.MsgId,
				message.MQTTSN_RC_REJECTED_INVALID_TOPIC_ID)
			// send
			s.Conn.WriteToUDP(puback.Marshall(), s.Remote)

			return
		}
	} else if message.TopicIdType(m.Flags) == message.MQTTSN_TIDT_SHORT_NAME {
		topicName = m.TopicName
	} else {
		log.Println("ERROR : invalid TopicIdType ", message.TopicIdType(m.Flags))
		return
	}

	// get qos
	qos := message.QosFlag(m.Flags)

	// send Publish to Mqtt Broker.
	// Topic, Qos, Retain, Payload
	var token MQTT.Token
	if s.mqttClient.IsConnected() {
		token = s.mqttClient.Publish(topicName, qos, false, m.Data)
	} else if err := s.ConnectToBroker(false); err != nil {
		log.Println("ERROR : Broker connection is not available")
		return
	} else {
		token = s.mqttClient.Publish(topicName, qos, false, m.Data)
	}

	if qos == 1 {
		// if qos 1
		go waitPubAck(token, s.MqttSnSession, m.TopicId, m.MsgId, s.statisticsReporter)
		s.statisticsReporter.countUpSendPublish()
	} else if qos == 2 {
		// elif qos 2
		// TODO: send PubRec
		// TODO: wait PubComp
	}
}

/**
 *
 */
func (s *TransparentSnSession) doSubscribe(m *message.Subscribe) {
	// if topic include wildcard, set topicId as 0x0000
	// else regist topic to client-session instance and assign topiId
	var topicId uint16
	topicId = uint16(0x0000)

	// return code
	var rc byte = message.MQTTSN_RC_ACCEPTED

	var topicName string
	switch message.TopicIdType(m.Flags) {
	// if TopicIdType is NORMAL, regist it
	case message.MQTTSN_TIDT_NORMAL:
		log.Println("Subscribe from : ", s.Remote.String(), ", TopicID Type : NORMAL, TopicName : ", m.TopicName)

		// check topic is wildcarded or not
		topics := strings.Split(m.TopicName, "/")
		if IsWildCarded(topics) != true {
			topicId = s.StoreTopic(m.TopicName)
		}
		topicName = m.TopicName

	// else if PREDEFINED, get TopicName and Subscribe to Broker
	case message.MQTTSN_TIDT_PREDEFINED:
		log.Println("Subscribe from : ", s.Remote.String(), ", TopicID Type : PREDEFINED, TopicID : ", m.TopicId)

		// Get topicId
		topicId = m.TopicId

		// get topic name and subscribe to broker
		var ok bool
		topicName, ok = LoadPredefTopic(topicId)
		if ok != true {
			// TODO: error handling
			rc = message.MQTTSN_RC_REJECTED_INVALID_TOPIC_ID
		}

	// else if SHORT_NAME, subscribe to broker
	case message.MQTTSN_TIDT_SHORT_NAME:
		log.Println("WARN : Subscribe SHORT Topic is not implemented")
		// TODO: implement
	}

	// send subscribe
	qos := message.QosFlag(m.Flags)
	var token MQTT.Token
	if s.mqttClient.IsConnected() {
		token = s.mqttClient.Subscribe(topicName, qos, s.OnPublish)
	} else if s.ConnectToBroker(false) != nil {
		log.Println("ERROR : Broker connection is not available")
		token = nil
		rc = message.MQTTSN_RC_REJECTED_CONGESTION
	} else {
		token = s.mqttClient.Subscribe(topicName, qos, s.OnPublish)
	}

	if token != nil && (!token.WaitTimeout(15*time.Second) || token.Error() != nil) {
		log.Println("ERROR : Broker connection is not available")
	}

	// send subscribe to broker
	suback := message.NewSubAck(
		topicId,
		m.MsgId,
		rc)

	// send suback
	s.Conn.WriteToUDP(suback.Marshall(), s.Remote)
}
