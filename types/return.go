package types

type Return struct {
	Status  int
	Message string
	Code    int
	Err     map[string]interface{}
	Token   string
	Tokens  interface{}
}

type Successre struct {
	Status  int
	Message string
	Code    int
}

type SuccessLogin struct {
	Status    int
	Message   string
	Code      int
	Token     string
	AdminName string
}

type Tokeninfo struct {
	Userid int
	Exp    int
	Admin  string
}

type SuccessreInfo struct {
	Status  int
	Message string
	Data    interface{}
	Code    int
}
type Pes struct {
	PsID       int         `json:"ps_id"`
	PsName     string      `json:"ps_name"`
	PsPid      int         `json:"ps_pid"`
	PsC        string      `json:"ps_c"`
	PsA        string      `json:"ps_a"`
	PsLevel    int64       `json:"ps_level"`
	CreateTime int64       `json:"create_time"`
	UpdateTime int64       `json:"update_time"`
	DeleteTime int64       `json:"delete_time"`
	Children   map[int]Per `json:"children"`
}

//
type Per struct {
	PsID       int    `json:"ps_id"`
	PsName     string `json:"ps_name"`
	PsPid      int    `json:"ps_pid"`
	PsC        string `json:"ps_c"`
	PsA        string `json:"ps_a"`
	PsLevel    int64  `json:"ps_level"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
	DeleteTime int64  `json:"delete_time"`
}
