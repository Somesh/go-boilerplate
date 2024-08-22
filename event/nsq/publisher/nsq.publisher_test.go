package nsqPublisher

import (
	"fmt"
	"testing"
	"time"

	"github.com/Somesh/go-boilerplate/common/config"
	"github.com/nsqio/go-nsq"
)

func mockPublish(producer *nsq.Producer, topic string, data []byte) error {
	return nil
}
func mockPublishErr(producer *nsq.Producer, topic string, data []byte) error {
	return fmt.Errorf("Error in publishing")
}

func mockDeferredPublish(producer *nsq.Producer, topic string, data []byte, delay time.Duration) error {
	return nil
}
func mockDeferredPublishErr(producer *nsq.Producer, topic string, data []byte, delay time.Duration) error {
	return fmt.Errorf("Error in publishing")
}

func TestPublish(t *testing.T) {

	n, err := NewNSQPublisher(config.NSQCfg{
		LookUpAddress: []string{"127.0.0.1:4150"},
		Prefix:        "omscart_",
	}, nil)
	if err != nil {
		t.Error(err)
	}

	// {"source":"webhook","message":{"guest_id":"24befcc7-0b11-488b-928e-fbc7f192c474","request_id":"","user_id":,"payment_status":"","ride_status":"ready","total_fare":"$9.00","currency_code":"USD"}}
	publishData := struct {
		GuestId       string `json:"guest_id"`
		RequestId     string `json:"request_id"`
		UserId        int64  `json:"user_id"`
		PaymentStatus string `json:"payment_status"`
		RideStatus    string `json:"ride_status"`
		TotalFare     string `json:"total_fare"`
		CurrencyCode  string `json:"currency_code"`
	}{
		GuestId:       "24befcc7-0b11-488b-928e-fbc7f192c474",
		RequestId:     "ab18dbad-b44d-4cb4-af05-0d0e991fc0c4",
		UserId:        3045010,
		PaymentStatus: "",
		RideStatus:    "ready",
		TotalFare:     "$9.00",
		CurrencyCode:  "USD",
	}

	p := n.producer
	n.producer = nil
	err = n.Publish("payment", QMessage{Source: "webhook", Message: publishData})
	if err == nil {
		t.Errorf("Expected Error. Got Nil")
	}

	err = n.PublishWithoutPrefix("payment", QMessage{Source: "webhook", Message: publishData})
	if err == nil {
		t.Errorf("Expected Error. Got Nil")
	}

	n.producer = p
	n.publish = mockPublish
	err = n.Publish("payment", QMessage{Source: "webhook", Message: publishData})
	if err != nil {
		t.Error(err)
	}

	err = n.Publish("payment", nil)
	if err == nil {
		t.Errorf("Expected Error. Got Nil")
	}
	n.publish = mockPublishErr
	err = n.Publish("payment", QMessage{Source: "webhook", Message: publishData})
	if err == nil {
		t.Errorf("Expected Error. Got Nil")
	}

}

func TestDeferredPublish(t *testing.T) {

	n, err := NewNSQPublisher(config.NSQCfg{
		LookUpAddress: []string{"127.0.0.1:4150"},
		Prefix:        "omscart_",
	}, nil)
	if err != nil {
		t.Log(err)
	}

	publishData := struct {
		GuestId       string `json:"guest_id"`
		RequestId     string `json:"request_id"`
		UserId        int64  `json:"user_id"`
		PaymentStatus string `json:"payment_status"`
		RideStatus    string `json:"ride_status"`
		TotalFare     string `json:"total_fare"`
		CurrencyCode  string `json:"currency_code"`
	}{
		GuestId:       "oiyuioiy-qcasfa24234-2342",
		RequestId:     "1234-5678-91ab",
		UserId:        9065826,
		PaymentStatus: "SUCCESS",
		RideStatus:    "COMPLETED",
		TotalFare:     "INR 10",
		CurrencyCode:  "INR",
	}

	p := n.producer
	n.producer = nil
	err = n.DeferredPublish("payment", QMessage{Source: "webhook", Message: publishData}, 1*time.Second)
	if err == nil {
		t.Errorf("Expected Error. Got Nil")
	}

	err = n.DeferredPublishWithoutPrefix("payment", QMessage{Source: "webhook", Message: publishData}, 1*time.Second)
	if err == nil {
		t.Errorf("Expected Error. Got Nil")
	}

	n.producer = p
	n.deferredPublish = mockDeferredPublish
	err = n.DeferredPublish("payment", QMessage{Source: "webhook", Message: publishData}, 1*time.Second)
	if err != nil {
		t.Error(err)
	}

	err = n.Publish("payment", nil)
	if err == nil {
		t.Errorf("Expected Error. Got Nil")
	}

	n.deferredPublish = mockDeferredPublishErr
	err = n.DeferredPublish("payment", QMessage{Source: "webhook", Message: publishData}, 1*time.Second)
	if err == nil {
		t.Errorf("Expected Error. Got Nil")
	}
}
