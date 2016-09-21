package check

import (
	"errors"
	"github.com/golang/mock/gomock"
	"testing"
)

func Test_checkBrokerExists_WhenBrokerExistsInZookeeper(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	check, zookeeper := newZkTestCheck(ctrl)
	zookeeper.EXPECT().Connect([]string{"localhost:2181"}, gomock.Any()).Return(nil, nil).MaxTimes(0)
	zookeeper.EXPECT().Exists("/brokers/ids/1").Return(true, nil, nil)

	res := check.checkBrokerZookeeperStatus()

	if res != "green" {
		t.Error("ZooKeeper should have broker available")
	}
}

func Test_checkBrokerExists_WhenBrokerDoesNotExistInZookeeper(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	check, zookeeper := newZkTestCheck(ctrl)
	zookeeper.EXPECT().Connect([]string{"localhost:2181"}, gomock.Any()).Return(nil, nil).MaxTimes(0)
	zookeeper.EXPECT().Exists("/brokers/ids/1").Return(false, nil, nil)

	res := check.checkBrokerZookeeperStatus()

	if res != "red" {
		t.Error("ZooKeeper should not have broker")
	}
}

func Test_checkBrokerExists_WhenZookeeperReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	check, zookeeper := newZkTestCheck(ctrl)
	zookeeper.EXPECT().Connect([]string{"localhost:2181"}, gomock.Any()).Return(nil, nil).MaxTimes(0)
	zookeeper.EXPECT().Exists("/brokers/ids/1").Return(true, nil, errors.New("zk: connection closed"))

	res := check.checkBrokerZookeeperStatus()

	if res != "red" {
		t.Error("ZooKeeper should not have broker")
	}
}
