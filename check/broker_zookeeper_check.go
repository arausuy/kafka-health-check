package check

import (
	"log"
	"strconv"
)

const (
	available   = "green"
	unavailable = "red"
)

func (check *HealthCheck) checkBrokerZookeeperStatus() string {
	brokerId := int(check.config.brokerID)
	brokerPath := "/brokers/ids/" + strconv.Itoa(brokerId)

	return check.getStatus(brokerPath)
}

func (check *HealthCheck) getStatus(brokerPath string) string {
	nodeExists, _, err := check.zookeeper.Exists(brokerPath)

	if err != nil {
		log.Printf("Broker path %s cannot be found in zookeeper"+
			" so failing health check", brokerPath)

		return unavailable
	}

	if nodeExists == false {
		log.Printf("Failing health check for %s - "+
			"Node not found in zookeeper, but no error thrown", brokerPath)

		return unavailable
	}

	log.Printf("Broker path %s can be found in zookeeper", brokerPath)

	return available

}
