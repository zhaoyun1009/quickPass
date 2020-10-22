package mq

import (
	"QuickPass/models"
	"QuickPass/pkg/constant"
	"QuickPass/pkg/logf"
	"encoding/json"
)

// 订单状态推送改用管道通信

var (
	socketMessageChan          chan []byte
	SocketChanReceiverInstance *SocketChanReceiver
)

func SocketChanSetup() {
	socketMessageChan = make(chan []byte, 50)

	SocketChanReceiverInstance = &SocketChanReceiver{}
	SocketChanReceiverInstance.AgencyChanList = make([]*RegisterInfo, 0)
	SocketChanReceiverInstance.AcceptorChanList = make([]*RegisterInfo, 0)
	SocketChanReceiverInstance.MerchantChanList = make([]*RegisterInfo, 0)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				logf.Error(err)
			}
		}()
		for {
			if bytes, ok := <-socketMessageChan; ok {
				err := SocketChanReceiverInstance.consumer(bytes)
				if err != nil {
					logf.Error("order status consumer", err)
				}
			} else {
				break
			}
		}
	}()
}

//订单状态接收监听
// 订单状态是固定的,注册的时候分配到不同的注册列表
type SocketChanReceiver struct {
	// 代理通道注册
	AgencyChanList []*RegisterInfo
	// 承兑人通道注册
	AcceptorChanList []*RegisterInfo
	// 商家通道注册
	MerchantChanList []*RegisterInfo
}

// 通道注册信息
type RegisterInfo struct {
	// 代理
	Agency string
	// 用户名
	Username string
	// 角色
	Role     int8
	DataChan chan []byte
}

// 通道消息
type SocketMessage struct {
	// 消息类型
	Type    int         `json:"type"`
	Message interface{} `json:"message"`
}

func (receiver *SocketChanReceiver) consumer(bytes []byte) error {
	message := new(SocketMessage)
	err := json.Unmarshal(bytes, message)
	if err != nil {
		return err
	}

	messageBytes, err := json.Marshal(message.Message)
	if err != nil {
		return err
	}
	switch message.Type {
	case constant.MessageTypeOrderStatus:
		return orderStatusConsumer(receiver, messageBytes, bytes)
	case constant.MessageTypeAbnormalOrderStatus:
		return abnormalOrderStatusConsumer(receiver, messageBytes, bytes)
	}

	return nil
}

func orderStatusConsumer(receiver *SocketChanReceiver, messageBytes []byte, bytes []byte) error {
	order := new(models.Order)

	err := json.Unmarshal(messageBytes, order)
	if err != nil {
		return err
	}

	// 商家通知
	for _, item := range receiver.MerchantChanList {
		if item.Agency == order.Agency && item.Username == order.MerchantUserName {
			item.DataChan <- bytes
		}
	}
	// 代理通知
	for _, item := range receiver.AgencyChanList {
		if item.Agency == order.Agency {
			item.DataChan <- bytes
		}
	}
	// 承兑人通知
	for _, item := range receiver.AcceptorChanList {
		if item.Agency == order.Agency && item.Username == order.ToUserName {
			item.DataChan <- bytes
		}
	}
	return nil
}

func abnormalOrderStatusConsumer(receiver *SocketChanReceiver, messageBytes []byte, bytes []byte) error {
	abnormalOrder := new(models.AbnormalOrder)

	err := json.Unmarshal(messageBytes, abnormalOrder)
	if err != nil {
		return err
	}
	// 代理通知
	for _, item := range receiver.AgencyChanList {
		if item.Agency == abnormalOrder.Agency {
			item.DataChan <- bytes
		}
	}
	// 承兑人通知
	for _, item := range receiver.AcceptorChanList {
		if item.Agency == abnormalOrder.Agency && item.Username == abnormalOrder.Username {
			item.DataChan <- bytes
		}
	}

	return nil
}

func (receiver *SocketChanReceiver) RegisterChan(item *RegisterInfo) {
	//注册时根据不同角色监听的订单信息不同
	switch item.Role {
	case constant.AGENCY:
		receiver.AgencyChanList = append(receiver.AgencyChanList, item)
	case constant.ACCEPTOR:
		receiver.AcceptorChanList = append(receiver.AcceptorChanList, item)
	case constant.MERCHANT:
		receiver.MerchantChanList = append(receiver.MerchantChanList, item)
	}
}

func (receiver *SocketChanReceiver) RemoveChan(removeItem *RegisterInfo) {
	switch removeItem.Role {
	case constant.AGENCY:
		for index, item := range receiver.AgencyChanList {
			if item == removeItem {
				receiver.AgencyChanList = append(receiver.AgencyChanList[:index], receiver.AgencyChanList[index+1:]...)
			}
		}
	case constant.ACCEPTOR:
		for index, item := range receiver.AcceptorChanList {
			if item == removeItem {
				receiver.AcceptorChanList = append(receiver.AcceptorChanList[:index], receiver.AcceptorChanList[index+1:]...)
			}
		}
	case constant.MERCHANT:
		for index, item := range receiver.MerchantChanList {
			if item == removeItem {
				receiver.MerchantChanList = append(receiver.MerchantChanList[:index], receiver.MerchantChanList[index+1:]...)
			}
		}
	}
}

func SendOrderStatusChange(msg *models.Order) {
	defer func() {
		if err := recover(); err != nil {
			logf.Error(err)
		}
	}()
	message := &SocketMessage{
		Type:    constant.MessageTypeOrderStatus,
		Message: msg,
	}
	marshal, _ := json.Marshal(message)
	socketMessageChan <- marshal
}

func SendAbnormalOrderStatusChange(msg *models.AbnormalOrder) {
	defer func() {
		if err := recover(); err != nil {
			logf.Error(err)
		}
	}()
	message := &SocketMessage{
		Type:    constant.MessageTypeAbnormalOrderStatus,
		Message: msg,
	}
	marshal, _ := json.Marshal(message)
	socketMessageChan <- marshal
}
