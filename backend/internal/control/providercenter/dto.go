package providercenter

type WorkStatusRequest struct {
	WorkStatus int `json:"workStatus"`
}

type ApplyRequest struct {
	RealName string `json:"realName"`
	IDCardNo string `json:"idCardNo"`
}

type UpdateProfileRequest struct {
	DisplayName string   `json:"displayName"`
	Intro       string   `json:"intro"`
	Tags        []string `json:"tags"`
}

type ProviderServiceItemRequest struct {
	ServiceItemID uint64 `json:"serviceItemId"`
	PriceAmount   int64  `json:"priceAmount"`
	PriceUnit     string `json:"priceUnit"`
}

type UpdateServiceItemsRequest struct {
	Items []ProviderServiceItemRequest `json:"items"`
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
