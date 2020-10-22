package order_service

import (
	"QuickPass/models"
	"QuickPass/pkg/util"
)

type MerchantOrder struct {
	//
	Id int64
	// 代理
	Agency string
	// 商家账户名
	Username string
	// 系统订单号
	OrderNo string
	// 商家平台传过来的订单号
	MerchantOrderNo string
	// 下单方式（1、接口，2、非接口）
	SubmitType int
	// 支付成功的返回URL
	ReturnUrl string
	// 回调接口地址
	CallbackUrl string
	// 回调结果（1、未处理，2、回调成功，3、回调失败）
	CallbackStatus int
	// 备注信息（例如：回调失败的信息）
	AppendInfo string
	// 创建时间
	CreateTime util.JSONTime
	// 更新时间
	UpdateTime util.JSONTime

	PageNum  int
	PageSize int
}

func (a *MerchantOrder) GetByOrderNo(agency, orderNo string) (*models.MerchantOrder, error) {
	session := models.NewSession()
	return models.NewMerchantOrderModel(session).GetMerchantOrder(agency, orderNo)
}

func (a *MerchantOrder) GetByMerchantOrderNo(agency, merchantName, merchantOrderNo string) (*models.MerchantOrder, error) {
	session := models.NewSession()
	return models.NewMerchantOrderModel(session).GetOneMerchantOrderByMerchantOrderNo(agency, merchantName, merchantOrderNo)
}

func (a *MerchantOrder) GetAll() ([]*models.MerchantOrder, error) {
	session := models.NewSession()
	return models.NewMerchantOrderModel(session).GetMerchantOrders(a.PageNum, a.PageSize, a.getMaps())
}

func (a *MerchantOrder) Delete() error {
	session := models.NewSession()
	return models.NewMerchantOrderModel(session).DeleteMerchantOrder(a.Id)
}

func (a *MerchantOrder) ExistByID() (bool, error) {
	session := models.NewSession()
	return models.NewMerchantOrderModel(session).ExistMerchantOrderByID(a.Id)
}

func (a *MerchantOrder) Count() (int64, error) {
	session := models.NewSession()
	return models.NewMerchantOrderModel(session).GetMerchantOrderTotal(a.getMaps())
}

func (a *MerchantOrder) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if a.Agency != "" {
		maps["agency"] = a.Agency
	}
	if a.Username != "" {
		maps["username"] = a.Username
	}
	if a.OrderNo != "" {
		maps["order_no"] = a.OrderNo
	}
	if a.MerchantOrderNo != "" {
		maps["merchant_order_no"] = a.MerchantOrderNo
	}
	if a.SubmitType != 0 {
		maps["submit_type"] = a.SubmitType
	}
	if a.CallbackStatus != 0 {
		maps["callback_status"] = a.CallbackStatus
	}

	return maps
}
