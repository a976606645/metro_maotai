package main

const (
	URLValidCode   = "https://appapi.dmall.com/app/passport/validCode"
	URLLogin       = "https://appapi.dmall.com/app/passport/smsLogin"
	URLQualication = "https://presale.dmall.com/maotai/qualification"
	URLJumpQueue   = "https://presale.dmall.com/maotai/jumpQueue"
	URLQueryQueue  = "https://presale.dmall.com/maotai/queryQueue"
	URLQueryStock  = "https://presale.dmall.com/maotai/tradeInfo"
)

const (
	NoErr = iota + 1000
	ErrParam
	ErrReq
	ErrDB
)
