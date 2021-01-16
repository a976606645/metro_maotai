package main

type APIValidCode struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

type APILogin struct {
	Code   string `json:"code"`
	Result string `json:"result,omitempty"`
	Data   struct {
		WebUser struct {
			ID         int    `json:"id"`
			LoginID    string `json:"loginId"`
			OtpToken   string `json:"otpToken"`
			Nickname   string `json:"nickname"`
			Realname   string `json:"realName"`
			TicketName string `json:"ticketName"`
			Token      string `json:"token"`
			IDCardNum  string `json:"idCardNum"`
			Birthday   int    `json:"birthday"`
		} `json:"webUser"`
	} `json:"data"`
}

type APIQualification struct {
	Code string `json:"code"`
	Data struct {
		HasQualification bool `json:"hasQualification"`
		RemainCount      int  `json:"remainCount"`
	} `json:"data"`
}

//{"code":"0000","data":true,"msg":"�ɹ�","time":1610449313205}
type APIQueueInfo struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data bool   `json:"data"`
}

type APIStockInfo struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		WareInfo struct {
			Stock int `json:"stock"`
		} `json:"wareInfo"`
	} `json:"data"`
}
