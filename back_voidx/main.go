package main

import (
	"fmt"
	"time"
	"voidx/scripts/bluetooth"
)

/*
func main() {
	usersScripts.UserVerify()

	devices, err := bluetooth.ScanDevices(3 * time.Second)

	if err != nil {
		fmt.Printf("Ошибка сканирования: %v\n", err)
	} else {
		fmt.Printf("Найдено устройств: %d\n", len(devices))
		for i, d := range devices {
			fmt.Printf("%d. %s [%s] RSSI: %d\n", i+1, d.LocalName(), d.Address.String(), d.RSSI)
		}
	}
}
*/

func main() {

	err := bluetooth.ScanDevices(3 * time.Second)
	if err != nil {
		fmt.Printf("HUYNYA")
	}
}
