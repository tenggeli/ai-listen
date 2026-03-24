package providercenter

type WorkStatusRequest struct {
	WorkStatus int `json:"workStatus"`
}

type AccountRequest struct {
	AccountType int    `json:"accountType"`
	AccountName string `json:"accountName"`
	AccountNo   string `json:"accountNo"`
	BankName    string `json:"bankName"`
}

type WithdrawRequest struct {
	AccountID uint64 `json:"accountId"`
	Amount    int64  `json:"amount"`
}
