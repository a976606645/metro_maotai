package main

import (
	"fmt"
	"testing"
)

func TestGenUUID(t *testing.T) {
	fmt.Println(len("f549f16478802fb3"))
}

func TestGenCID(t *testing.T) {
	fmt.Println(len("170976fa8a2449b3b88"))
}

func TestGenSessionID(t *testing.T) {
	result := GenSessionID()
	if len(result) != len("5bcff61bff624b66bbc8a4016c32385e") {
		t.Error("sessionid 长度不一致:", result)
	}

	fmt.Println("session:", result)
}

func TestRandDeviceName(t *testing.T) {
	result := RandDeviceName()
	fmt.Println("device:", result)
}
