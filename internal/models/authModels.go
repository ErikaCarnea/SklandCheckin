package models

type LoginResult struct {
	Status  int    `json:"status"`
	Message string `json:"msg"`
	Data    struct {
		Token string `json:"token"`
	} `json:"data"`
}

func (r LoginResult) GetCode() int {
	return r.Status
}
func (r LoginResult) GetMessage() string {
	return r.Message
}

type GrantResult struct {
	Status  int    `json:"status"`
	Message string `json:"msg"`
	Data    struct {
		Code string `json:"code"`
	} `json:"data"`
}

func (g GrantResult) GetCode() int {
	return g.Status
}
func (g GrantResult) GetMessage() string {
	return g.Message
}

type CredResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Cred   string `json:"cred"`
		Token  string `json:"token"`
		UserId string `json:"userId"`
	} `json:"data"`
}

func (r CredResult) GetCode() int {
	return r.Code
}
func (r CredResult) GetMessage() string {
	return r.Message
}

type SendCodeResult struct {
	Status  int    `json:"status"`
	Message string `json:"msg"`
}

func (r SendCodeResult) GetCode() int {
	return r.Status
}

func (r SendCodeResult) GetMessage() string {
	return r.Message
}
