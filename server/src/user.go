package main

func NewUser(phone string) *TableUser {
	u := &TableUser{
		Phone:     phone,
		UUID:      GenUUID(),
		CID:       GenCID(),
		SessionID: GenSessionID(),
		Device:    RandDeviceName(),
	}

	return u
}
