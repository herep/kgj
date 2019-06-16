package types

type ReturnAccount struct {
	Id             int    `json:"id"`
	AccountName    string `json:"account_name"`
	AccountNum     string `json:"account_num"`
	AccountMailbox string `json:"account_mailbox"`
	AccountPhone   string `json:"account_phone"`
	AccountStatus  string `json:"account_status"`
	Role_name      string `json:"role_name"`
	Company_Name   string `json:"company_name"`
	CreateTime     int64  `json:"create_time"`
	UpdateTime     int64  `json:"update_time"`
}
