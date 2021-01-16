package main

import (
	"fmt"
	"math/rand"
)

func GenUniqueID() string {
	b := make([]byte, 32)
	rand.Read(b)

	return fmt.Sprintf("%x", b)
}

// GenUUID 等同于设备号 16位
func GenUUID() string {
	// f549f16478802fb3
	// 904e17cdf1dae0b6
	return GenUniqueID()[:16]
}

// GenCID 19位
func GenCID() string {
	// 170976fa8a2449b3b88
	// 1507bfd3f765c393abe
	return GenUniqueID()[:19]
}

// GenSessionID 32位
func GenSessionID() string {
	// edb4e7eb75e740bbb7da4f62e24df4f4
	// 5bcff61bff624b66bbc8a4016c32385e
	return GenUniqueID()[:32]
}

func RandDeviceName() string {
	index := rand.Int31n(int32(len(PhoneModelDataBase)))
	return PhoneModelDataBase[index]
}
