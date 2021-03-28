package main

import (
	"fmt"
	"testing"
	"time"
)

func TestQualication(t *testing.T) {
	user := &TableUser{
		UUID:       "12ac02901cdda9ef",
		TicketName: "2392C9223426F1232D57974ED46D4EC3FF0BBC4B900150535E754B1D6963D287A477E447EBB54D905BFFF7F2E594CDD0E41F8E20870080B77D7A7AFC2241A9C89268D121263B1CEB67E54686D54B488C924B80F7C79A3D26790CDBF7F7C51696E26E3097119757C25C536B8384FEDB48887B141D6543B0129BD9AD32309742F6",
	}

	qinfo, err := qualification(user)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v \n", qinfo)

}

func TestTomorrow(t *testing.T) {
	y, m, d := time.Now().Add(24 * time.Hour).Date()
	shipment := fmt.Sprintf("%d-%d-%d", y, m, d)
	fmt.Println("shipment:", shipment)
}
