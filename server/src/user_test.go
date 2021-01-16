package main

import (
	"fmt"
	"testing"
)

func TestNewUser(t *testing.T) {
	p := "15197910055"
	u := NewUser(p)
	if u.Phone != p {
		t.Error("新建用户错误")
	}

	fmt.Printf("%+v \n", u)
}
