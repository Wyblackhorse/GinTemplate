package v2

//type GetPayInformationBackData struct {
//
//	TxHash      string `json:"txHash" binding:"required"`
//	BlockNumber int    `json:"blockNumber" binding:"required"`
//	Timestamp   int64  `json:"timestamp" binding:"required"`
//	From        string `json:"from" binding:"required"`
//	To          string `json:"to" binding:"required"`
//	Amount      int    `json:"amount" binding:"required"`
//	Token       string `json:"token" binding:"required"`
//	UserID      string `json:"userId" binding:"required"`
//}

type CreatePrepaidPhoneOrdersData struct {
	PlatformOrder string `json:"PlatformOrder" binding:"required"` //平台订单号
	//CollectionAddress string `json:"CollectionAddress" binding:"required"`
	RechargeAddress string  `json:"RechargeAddress" binding:"required"`
	Username        string  `json:"Username" binding:"required"`
	AccountOrders   float64 `json:"AccountOrders"  binding:"required"`
	RechargeType    string  `json:"RechargeType"  binding:"required"`
}

type GetPayInformationBackData struct {
	Type string `json:"type" binding:"required"`
	Data Data   `json:"data" binding:"required"`
	Sign string `json:"sign" binding:"required"`
}
type Data struct {
	TxHash      string `json:"txHash" binding:"required"`
	BlockNumber int    `json:"blockNumber" binding:"required"`
	Timestamp   int64  `json:"timestamp" binding:"required"`
	From        string `json:"from" binding:"required"`
	To          string `json:"to" binding:"required"`
	Amount      int    `json:"amount" binding:"required"`
	Token       string `json:"token" binding:"required"`
	UserID      string `json:"userId" binding:"required"`
	Balance     string `json:"balance" binding:"required"`
}

type ReturnBase64 struct {
	Data string `json:"data"`
	Sign string `json:"sign"`
}

type BalanceType struct {
	Type string `json:"type"`
	Data Data   `json:"data"`
}
type DataTwo struct {
	User    string `json:"user"`
	Addr    string `json:"addr"`
	Balance string `json:"balance"`
	Seq     int    `json:"seq"`
}
