package v2

type GetPayInformationBackData struct {
	TxHash      string `json:"txHash" binding:"required"`
	BlockNumber int    `json:"blockNumber" binding:"required"`
	Timestamp   int64  `json:"timestamp" binding:"required"`
	From        string `json:"from" binding:"required"`
	To          string `json:"to" binding:"required"`
	Amount      int    `json:"amount" binding:"required"`
	Token       string `json:"token" binding:"required"`
	UserID      string `json:"userId" binding:"required"`
}

type CreatePrepaidPhoneOrdersData struct {
	ThreeOrder        string `json:"ThreeOrder" binding:"required"`
	CollectionAddress string `json:"CollectionAddress" binding:"required"`
	RechargeAddress   string `json:"RechargeAddress" binding:"required"`
	Username          string `json:"Username" binding:"required"`
	AccountOrders     int    `json:"AccountOrders"  binding:"required"`
	RechargeType      string `json:"RechargeType"  binding:"required"`
}



