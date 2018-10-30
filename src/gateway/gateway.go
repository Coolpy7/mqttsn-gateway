package mqttsn

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"messages"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

type Gateway interface {
	HandlePacket(*net.UDPConn, *net.UDPAddr, []byte)
	StartUp() error
}

type MessageQueueObject struct {
	remote *net.UDPAddr
	conn   *net.UDPConn
	msg    message.MqttSnMessage
}

type AggregatingGateway struct {
	mutex              sync.RWMutex
	MqttSnSessions     map[string]*MqttSnSession
	Config             *GatewayConfig
	mqttClient         MQTT.Client
	signalChan         chan os.Signal
	recvBuffer         chan *MessageQueueObject
	statisticsReporter *StatisticsReporter
}

type TransparentGateway struct {
	MqttSnSessions     map[string]*TransparentSnSession
	Config             *GatewayConfig
	signalChan         chan os.Signal
	recvBuffer         chan *MessageQueueObject
	statisticsReporter *StatisticsReporter
}

func serverLoop(
	gateway Gateway,
	host string,
	port int,
	readBuff int,
	writeBuff int) error {
	serverStr := host + ":" + strconv.Itoa(port)
	udpAddr, err := net.ResolveUDPAddr("udp", serverStr)
	conn, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		return err
	}

	// set os buffer size
	if readBuff > 0 {
		log.Println("set read buff size : ", readBuff)
		if err := conn.SetReadBuffer(readBuff); err != nil {
			log.Println("failed to set read buff size")
		}
	}
	if writeBuff > 0 {
		log.Println("set write buff size : ", writeBuff)
		if err := conn.SetWriteBuffer(writeBuff); err != nil {
			log.Println("failed to set write buff size")
		}
	}

	buf := make([]byte, 64*1024)
	for {
		_, remote, err := conn.ReadFromUDP(buf)
		if err != nil {
			return err
		}

		// launch gateway
		gateway.HandlePacket(conn, remote, buf)
	}
}

func waitPubAck(token MQTT.Token, s *MqttSnSession, topicId uint16, msgId uint16, reporter *StatisticsReporter) {
	// timeout time is 15 sec
	if !token.WaitTimeout(15 * time.Second) {
		log.Println("ERROR : Wait PubAck is Timeout.")
	} else if token.Error() != nil {
		log.Println("ERROR : ", token.Error())
	} else {
		// send puback
		puback := message.NewPubAck(topicId, msgId, message.MQTTSN_RC_ACCEPTED)
		s.Conn.WriteToUDP(puback.Marshall(), s.Remote)

		// count up statistics information
		tat := s.tatCalculator.CalculateTatAndRemoveMsgId(msgId, time.Now())
		reporter.countUpPubComp()
		reporter.storeTat(tat)
	}
}

func waitPubComp(token MQTT.Token, s *MqttSnSession, topicId uint16, msgId uint16) {
	// TODO: implement
	// token.Wait()
	// pubcomp := message.NewPubComp()
	// s.Conn.WriteToUDP(pubcomp.Marshall(), s.Remote)
}
