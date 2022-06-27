package V2

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
	PlatformOrder   string  `json:"PlatformOrder" binding:"required"`   //平台订单号
	//RechargeAddress string  `json:"RechargeAddress" binding:"required"` //充值地址
	Username        string  `json:"Username" binding:"required"`        //充值用户名
	AccountOrders   float64 `json:"AccountOrders"  binding:"required"`  //充值金额
	RechargeType    string  `json:"RechargeType"  binding:"required"`   //充值类型
	BackUrl         string  `json:"BackUrl"  binding:"required" `       // 回调地址
	//ThreeOrder      string  `json:"three_order" binding:"required"`   三方订单号 不需要 这个是我自己生成的!
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
	Data struct {
		Addr    string `json:"addr"`
		Balance string `json:"balance"`
		Seq     int64  `json:"seq"`
		User    string `json:"user"`
	} `json:"data"`
	Type string `json:"type"`
}

type UsernameAddress struct {
	Username string `json:"username"`
}


