package models

import (
	"QuickPass/pkg/constant"
	"QuickPass/pkg/errors"
	"QuickPass/pkg/util"
	"fmt"
	"github.com/jinzhu/gorm"
)

type Order struct {
	//
	Id int64 `json:"id" gorm:"primary_key" gorm:"column:id"`
	// 代理
	Agency string `json:"agency" gorm:"column:agency"`
	// 订单号
	OrderNo string `json:"order_no" gorm:"column:order_no"`
	// 订单类型（1：转账，2：买入，3：卖出）
	OrderType int64 `json:"order_type" gorm:"column:order_type"`
	// 订单状态（1：创建，2：待支付，3：待放行，4：已取消  5：已失败  6：已完成）
	OrderStatus int64 `json:"order_status" gorm:"column:order_status"`
	// 异常单号
	AbnormalOrderNo string `json:"abnormal_order_no" gorm:"column:abnormal_order_no"`
	// 发起者账号
	FromUserName string `json:"from_user_name" gorm:"column:from_user_name"`
	// 商家账号
	MerchantUserName string `json:"merchant_user_name" gorm:"column:merchant_user_name"`
	// 接受者账号
	ToUserName string `json:"to_user_name" gorm:"column:to_user_name"`
	// 订单所匹配的承兑人昵称
	AcceptorNickname string `json:"acceptor_nickname" gorm:"column:acceptor_nickname"`
	// 卡号
	CardNo string `json:"card_no" gorm:"column:card_no"`
	// 卡账户名
	CardAccount string `json:"card_account" gorm:"column:card_account"`
	// 银行名称
	CardBank string `json:"card_bank" gorm:"column:card_bank"`
	// 银行支行
	CardSubBank string `json:"card_sub_bank" gorm:"column:card_sub_bank"`
	// 收款码
	CardImg string `json:"card_img" gorm:"column:card_img"`
	// 金额
	Amount int64 `json:"amount" gorm:"column:amount"`
	// 提交类型=（1：直接提交，2：接口提交）
	SubmitType int `json:"submit_type" gorm:"column:submit_type"`
	// 通道类型= (BANK_CARD,ALIPAY,WECHAT)
	ChannelType string `json:"channel_type" gorm:"column:channel_type"`
	// 摘要
	AppendInfo string `json:"append_info" gorm:"column:append_info"`
	// ip地址
	ClientIp string `json:"client_ip" gorm:"column:client_ip"`
	// 订单完成时间
	FinishTime *util.JSONTime `json:"finish_time" gorm:"column:finish_time"`
	// 创建时间
	CreateTime util.JSONTime `json:"create_time" gorm:"column:create_time"`
	// 更新时间
	UpdateTime util.JSONTime `json:"update_time" gorm:"column:update_time"`

	Session *Session `json:"-" gorm:"-"`

	// 商家平台传过来的订单号
	MerchantOrderNo string `json:"merchant_order_no" gorm:"-"`
	// 回调结果（1、未处理，2、回调成功，3、回调失败）
	CallbackStatus int `json:"callback_status" gorm:"-"`
	// 回调失败的信息
	CallbackInfo string `json:"callback_info" gorm:"-"`
	// 回调地址
	CallbackUrl string `json:"callback_url" gorm:"-"`
}

// 设置Order的表名为`order`
func (Order) TableName() string {
	return "order"
}

func NewOrderModel(session *Session) *Order {
	return &Order{Session: session}
}

//订单待支付
//id 订单ID
//toUsername 接受者账号
//cardId 卡ID
func (a *Order) UpdateOrderWaitPay(agency, orderNo, toUsername, acceptorNickname, cardNo, cardAccount, cardBank, cardSubBank, cardImg string) error {
	tx := GetSessionTx(a.Session)

	order, err := a.GetOrderByOrderNo(agency, orderNo)
	if err != nil {
		return err
	}

	if order.OrderStatus != constant.ORDER_STATUS_CREATE {
		return errors.New("订单状态异常")
	}

	data := map[string]interface{}{
		"order_status":      constant.ORDER_STATUS_WAIT_PAY,
		"to_user_name":      toUsername,
		"acceptor_nickname": acceptorNickname,
		"card_no":           cardNo,
		"card_account":      cardAccount,
		"card_bank":         cardBank,
		"card_sub_bank":     cardSubBank,
		"card_img":          cardImg,
	}
	tx = tx.Model(&Order{}).Where("agency = ? and order_no = ? and order_status = ?", agency, orderNo, order.OrderStatus).Updates(data)
	if tx.Error != nil {
		return tx.Error
	}
	if !(tx.RowsAffected > 0) {
		return errors.New("更新订单状态失败")
	}
	return nil
}

//买入订单支付待放行
//orderNo 订单号
func (a *Order) UpdateOrderWaitDischarge(agency, orderNo string) error {
	tx := GetSessionTx(a.Session)

	order, err := a.GetOrderByOrderNo(agency, orderNo)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("order not found")
	}

	if order.OrderStatus != constant.ORDER_STATUS_WAIT_PAY {
		return errors.New("订单状态异常")
	}

	data := map[string]interface{}{
		"order_status": constant.ORDER_STATUS_WAIT_DISCHARGE,
	}
	tx = tx.Model(&Order{}).Where("agency = ? and order_no = ? and order_status = ?", agency, orderNo, order.OrderStatus).Updates(data)
	if tx.Error != nil {
		return tx.Error
	}
	if !(tx.RowsAffected > 0) {
		return errors.New("更新订单状态失败")
	}
	return nil
}

// 唤醒已取消的订单
func (a *Order) WakeUpCancelOrder(agency, orderNo string) error {
	tx := GetSessionTx(a.Session)

	order, err := a.GetOrderByOrderNo(agency, orderNo)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("order not found")
	}

	if order.OrderStatus != constant.ORDER_STATUS_CANCEL {
		return errors.New("订单状态异常")
	}

	data := map[string]interface{}{
		"order_status": constant.ORDER_STATUS_WAIT_DISCHARGE,
	}
	tx = tx.Model(&Order{}).Where("agency = ? and order_no = ? and order_status = ?", agency, orderNo, order.OrderStatus).Updates(data)
	if tx.Error != nil {
		return tx.Error
	}
	if !(tx.RowsAffected > 0) {
		return errors.New("更新订单状态失败")
	}
	return nil
}

//卖出订单支付待放行
//id 订单ID
func (a *Order) UpdateSellOrderWaitDischarge(agency, orderNo string) error {
	tx := GetSessionTx(a.Session)

	order, err := a.GetOrderByOrderNo(agency, orderNo)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("order not found")
	}

	if order.OrderStatus != constant.ORDER_STATUS_WAIT_PAY {
		return errors.New("订单状态异常")
	}

	data := map[string]interface{}{
		"order_status": constant.ORDER_STATUS_WAIT_DISCHARGE,
	}

	tx = tx.Model(&Order{}).Where("agency = ? and order_no = ? and order_status = ?", agency, orderNo, order.OrderStatus).Updates(data)
	if tx.Error != nil {
		return tx.Error
	}

	if !(tx.RowsAffected > 0) {
		return errors.New("更新订单状态失败")
	}
	return nil
}

//更新订单被取消
//orderNo 订单号
func (a *Order) UpdateOrderCancel(agency, orderNo string) error {
	tx := GetSessionTx(a.Session)

	order, err := a.GetOrderByOrderNo(agency, orderNo)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New(fmt.Sprintf("orderNo[%s] info not found", orderNo))
	}

	// 可取消订单状态标记
	enableCancelFlag := true
	switch order.OrderType {
	case constant.ORDER_TYPE_BUY:
		// 买单只有待支付和待放行可以取消
		enableCancelFlag = (order.OrderStatus == constant.ORDER_STATUS_WAIT_PAY) ||
			(order.OrderStatus == constant.ORDER_STATUS_WAIT_DISCHARGE)
	case constant.ORDER_TYPE_SELL:
		// 卖单只有待支付可以取消
		enableCancelFlag = order.OrderStatus == constant.ORDER_STATUS_WAIT_PAY
	}

	if !enableCancelFlag {
		return errors.New("UpdateOrderCancel [订单状态异常]")
	}

	data := map[string]interface{}{
		"order_status": constant.ORDER_STATUS_CANCEL,
	}
	tx = tx.Model(&Order{}).Where("agency = ? and order_no = ? and order_status = ?", agency, orderNo, order.OrderStatus).Updates(data)
	if tx.Error != nil {
		return tx.Error
	}
	if !(tx.RowsAffected > 0) {
		return errors.New("更新订单状态失败")
	}
	return nil
}

//订单失败
//id 订单ID
func (a *Order) UpdateOrderFailed(agency, orderNo string) error {
	tx := GetSessionTx(a.Session)

	order, err := a.GetOrderByOrderNo(agency, orderNo)
	if err != nil {
		return err
	}

	if order == nil {
		return errors.New("not found order")
	}
	if order.OrderStatus != constant.ORDER_STATUS_CREATE {
		return errors.New("订单状态异常")
	}

	data := map[string]interface{}{
		"order_status": constant.ORDER_STATUS_FAILED,
	}
	tx = tx.Model(&Order{}).Where("agency = ?  and order_no = ? and order_status = ?", agency, orderNo, order.OrderStatus).Updates(data)
	if tx.Error != nil {
		return tx.Error
	}
	if !(tx.RowsAffected > 0) {
		return errors.New("更新订单状态失败")
	}

	return nil
}

//订单待支付
//id 订单ID
func (a *Order) UpdateOrderFinished(agency, orderNo string) error {
	tx := GetSessionTx(a.Session)

	order, err := a.GetOrderByOrderNo(agency, orderNo)
	if err != nil {
		return err
	}

	// 可放行订单状态标记
	enableDischargeFlag := true
	switch order.OrderType {
	case constant.ORDER_TYPE_BUY:
		// 买单只有待支付和待放行可以放行完成
		enableDischargeFlag = (order.OrderStatus == constant.ORDER_STATUS_WAIT_PAY) ||
			(order.OrderStatus == constant.ORDER_STATUS_WAIT_DISCHARGE)
	case constant.ORDER_TYPE_SELL:
		// 卖单只有待放行可以放行完成
		enableDischargeFlag = order.OrderStatus == constant.ORDER_STATUS_WAIT_DISCHARGE
	}
	if !enableDischargeFlag {
		return errors.New("订单状态异常")
	}

	data := map[string]interface{}{
		"order_status": constant.ORDER_STATUS_FINISHED,
		"finish_time":  util.JSONTimeNow(),
	}
	tx = tx.Model(&Order{}).Where("agency = ? and order_no = ? and order_status = ?", agency, orderNo, order.OrderStatus).Updates(data)
	if tx.Error != nil {
		return tx.Error
	}
	if !(tx.RowsAffected > 0) {
		return errors.New("更新订单状态失败")
	}
	return nil
}

// ExistOrderByID checks if an order exists based on ID
func (a *Order) ExistOrderByID(id int64) (bool, error) {
	order := new(Order)
	err := a.Session.db.Select("id").Where("id = ? ", id).First(order).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if order.Id > 0 {
		return true, nil
	}

	return false, nil
}

// GetOrderTotal gets the total number of orders based on the constraints
func (a *Order) GetOrderTotal(params *OrderListQueryParams, maps interface{}) (int64, error) {
	var count int64
	query := getOrderListParams(a, params)
	err := query.Where(maps).Count(&count).Error
	return count, err
}

// GetOrders gets a list of orders based on paging constraints
func (a *Order) GetOrders(params *OrderListQueryParams, pageNum int, pageSize int, maps interface{}) ([]*Order, error) {
	var orders []*Order
	query := getOrderListParams(a, params)
	err := query.Where(maps).Offset(pageNum).Limit(pageSize).Find(&orders).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return orders, nil
}

//动态组合查询条件
func getOrderListParams(a *Order, params *OrderListQueryParams) *gorm.DB {
	query := a.Session.db.Model(&Order{})
	if params.StartTime != "" {
		query = query.Where("create_time >= ?", params.StartTime)
	}
	if params.EndTime != "" {
		query = query.Where("create_time <= ?", params.EndTime)
	}
	if params.FinishStartTime != "" {
		query = query.Where("finish_time >= ?", params.FinishStartTime)
	}
	if params.FinishEndTime != "" {
		query = query.Where("finish_time <= ?", params.FinishEndTime)
	}
	if params.MinAmount != 0 {
		query = query.Where("amount >= ?", params.MinAmount)
	}
	if params.MaxAmount != 0 {
		query = query.Where("amount <= ?", params.MaxAmount)
	}
	if params.FinishStartTime != "" || params.FinishEndTime != "" {
		query = query.Where("order_status = ?", constant.ORDER_STATUS_FINISHED)
	}
	orderBy := fmt.Sprintf("FIELD(order_status, %d, %d) desc, create_time desc ", constant.ORDER_STATUS_WAIT_DISCHARGE, constant.ORDER_STATUS_WAIT_PAY)
	query = query.Order(orderBy)
	return query
}

// GetOrder Get a single order based on ID
func (a *Order) GetOrder(id int64) (*Order, error) {
	order := new(Order)
	err := a.Session.db.Where("id = ?", id).First(order).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return order, nil
}

// EditOrder modify a single order
func (a *Order) EditOrder(id int64, data interface{}) error {
	tx := GetSessionTx(a.Session)
	return tx.Model(&Order{}).Where("id = ?", id).Updates(data).Error
}

// AddOrder add a single order
func (a *Order) AddOrder(order *Order) error {
	tx := GetSessionTx(a.Session)
	return tx.Create(order).Error
}

// DeleteOrder delete a single order
func (a *Order) DeleteOrder(id int64) error {
	tx := GetSessionTx(a.Session)
	return tx.Where("id = ?", id).Delete(Order{}).Error
}

func (a *Order) GetOrderByOrderNo(agency, orderNo string) (*Order, error) {
	order := new(Order)
	err := a.Session.db.Where("agency = ? and order_no = ?", agency, orderNo).First(order).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return order, nil
}

func (a *Order) GetOrderTimeoutList(minute int64) ([]*Order, error) {
	var orders []*Order
	err := a.Session.db.
		// 仅限于买单
		Where("order_type = ?", constant.ORDER_TYPE_BUY).
		// 仅限于(2：待支付，3：待放行)
		Where("order_status in (?, ?)", constant.ORDER_STATUS_WAIT_PAY, constant.ORDER_STATUS_WAIT_DISCHARGE).
		// 仅限于正常下单(排除异常单)
		Where("abnormal_order_no = ''").
		//
		Where("update_time <= DATE_SUB(NOW(), INTERVAL ? MINUTE)", minute).
		Find(&orders).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return orders, nil
}

// 需要回调的订单列表
func (a *Order) GetOrderCallbackList() ([]*Order, error) {
	var orders []*Order

	err := a.Session.db.Table(a.TableName()).
		Joins("JOIN merchant_order ON merchant_order.order_no = `order`.order_no").
		// 回调次数少于3次
		Where("merchant_order.callback_url != ''").
		// 回调未请求的
		Where("merchant_order.callback_lock = ?", constant.CallbackUnLock).
		// 回调次数少于3次
		Where("merchant_order.callback_times < ?", 3).
		// 限于回调失败和未回调的
		Where("merchant_order.callback_status in (?, ?)", constant.CallbackUnprocessed, constant.CallbackFailed).
		// 限于订单完结的
		Where("`order`.order_status in (?, ?, ?)", constant.ORDER_STATUS_CANCEL, constant.ORDER_STATUS_FAILED, constant.ORDER_STATUS_FINISHED).
		Find(&orders).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return orders, nil
}

func (a *Order) GetSellOrderAutoDischargeList(minute int64) ([]*Order, error) {
	var orders []*Order
	err := a.Session.db.
		// 仅限于买单
		Where("order_type = ?", constant.ORDER_TYPE_SELL).
		// 仅限于(3：待放行)
		Where("order_status = ?", constant.ORDER_STATUS_WAIT_DISCHARGE).
		// 仅限于正常下单(排除异常单)
		Where("abnormal_order_no = ''").
		//
		Where("update_time <= DATE_SUB(NOW(), INTERVAL ? MINUTE)", minute).
		Find(&orders).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return orders, nil
}

func (a Order) GetAcceptorGroup(agency string) ([]string, error) {
	var orders []*Order
	err := a.Session.db.Table(a.TableName()).
		Select("acceptor_nickname").
		Where("acceptor_nickname != '' and agency = ?", agency).
		Group("acceptor_nickname").Find(&orders).Error

	acceptorGroup := make([]string, 0, len(orders))
	for _, item := range orders {
		acceptorGroup = append(acceptorGroup, item.AcceptorNickname)
	}
	return acceptorGroup, err
}

func (a Order) GetCurrentOrder(agency, merchant, member string) (*Order, error) {
	order := new(Order)
	err := a.Session.db.Table(a.TableName()).
		Where("agency = ? and merchant_user_name = ? and from_user_name = ?", agency, merchant, member).
		// 仅限于(2：待支付，3：待放行)
		Where("order_status in (?, ?)", constant.ORDER_STATUS_WAIT_PAY, constant.ORDER_STATUS_WAIT_DISCHARGE).
		Find(order).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return order, nil
}

func (a Order) GetCurrentIPOrder(agency, ip string) ([]Order, error) {
	orders := make([]Order, 0)
	err := a.Session.db.Table(a.TableName()).
		Where("agency = ? and client_ip = ?", agency, ip).
		// 仅限于(2：待支付，3：待放行)
		Where("order_type = ? and order_status in (?, ?)", constant.ORDER_TYPE_BUY, constant.ORDER_STATUS_WAIT_PAY, constant.ORDER_STATUS_WAIT_DISCHARGE).
		Find(&orders).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return orders, nil
}

type OrderListQueryParams struct {
	// 最大金额
	MaxAmount int64 `form:"max_amount" json:"max_amount"`
	// 最小金额
	MinAmount int64 `form:"min_amount" json:"min_amount"`
	// 开始时间
	StartTime string `form:"start_time" json:"start_time"`
	// 结束时间
	EndTime string `form:"end_time" json:"end_time"`

	// 完成开始时间
	FinishStartTime string `form:"finish_start_time" json:"finish_start_time"`
	// 完成结束时间
	FinishEndTime string `form:"finish_end_time" json:"finish_end_time"`
}
