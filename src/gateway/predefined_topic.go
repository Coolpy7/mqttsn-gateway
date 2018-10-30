package mqttsn

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type PredefinedTopicMap map[string]uint16

var PredefinedTopics *TopicMap = NewTopicMap()

/**
 * Initialize Predefined Topic
 */
func InitPredefinedTopic(file string) error {
	topics := PredefinedTopicMap{}

	yamlStr, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println("ERROR : faild to read yaml file.")
		return err
	}

	err = yaml.Unmarshal(yamlStr, &topics)
	if err != nil {
		log.Println("ERROR : faild to unmarshall yaml file.")
		return err
	}

	// read predefined topics
	for key, value := range topics {
		StorePredefTopicWithId(key, value)
	}

	return nil
}

/**
 * Store topic to map
 */
func StorePredefTopic(topicName string) uint16 {
	return PredefinedTopics.StoreTopic(topicName)
}

/**
 * Store topic to map with specified ID
 */
func StorePredefTopicWithId(topicName string, id uint16) bool {
	return PredefinedTopics.StoreTopicWithId(topicName, id)
}

/**
 * Load topic related to topic id
 */
func LoadPredefTopic(topicId uint16) (string, bool) {
	return PredefinedTopics.LoadTopic(topicId)
}

/**
 * Load topic id related to topic name
 */
func LoadPredefTopicId(topicName string) (uint16, bool) {
	return PredefinedTopics.LoadTopicId(topicName)
}
