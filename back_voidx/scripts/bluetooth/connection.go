package bluetooth

import (
	"fmt"

	"tinygo.org/x/bluetooth"
)

/*
Поток подключения (роль central, типичный клиент к периферии voidx):

  1. Scan — собираешь bluetooth.ScanResult (MAC/имя/RSSI из рекламы).
  2. StopScan — перед Connect обязательно останови сканирование: иначе BlueZ
     часто нестабилен (конкурирующие HCI-команды).
  3. Enable — один раз на процесс (повторный вызов обычно безвреден).
  4. Connect — ACL-линк к выбранному пиру (tinygo: adapter.Connect).
  5. DiscoverServices — дождаться кэша GATT у BlueZ (на Linux это обёртка
     вокруг ServicesResolved).
  6. DiscoverCharacteristics — найти нужные UUID под твой протокол.
  7. Read / WriteWithoutResponse / EnableNotifications — обмен данными.
  8. Disconnect — явно рвёшь сессию или выходишь из процесса.

Пока у voidx нет своего GATT UUID в коде — шаги 5–7 вызываешь сам, когда
задашь сервис и характеристики (как в examples/nusclient библиотеки).
*/

// AddressFromMAC строит bluetooth.Address из строки вида "AA:BB:CC:DD:EE:FF"
// (как в твоём JSON устройств). IsRandom=false — публичный адрес.
func AddressFromMAC(mac string) (bluetooth.Address, error) {
	m, err := bluetooth.ParseMAC(mac)
	if err != nil {
		return bluetooth.Address{}, fmt.Errorf("parse MAC: %w", err)
	}
	return bluetooth.Address{MACAddress: bluetooth.MACAddress{MAC: m}}, nil
}

// StopScanBestEffort останавливает активное сканирование; если скана нет —
// ошибка игнорируется (удобно вызывать перед Connect из любого состояния).
func StopScanBestEffort() {
	_ = adapter.StopScan()
}

// ConnectPeer подключается к периферии по MAC-строке.
// Перед вызовом рекомендуется StopScanBestEffort().
func ConnectPeer(mac string) (bluetooth.Device, error) {
	if err := adapter.Enable(); err != nil {
		return bluetooth.Device{}, fmt.Errorf("enable adapter: %w", err)
	}
	addr, err := AddressFromMAC(mac)
	if err != nil {
		return bluetooth.Device{}, err
	}
	dev, err := adapter.Connect(addr, bluetooth.ConnectionParams{})
	if err != nil {
		return bluetooth.Device{}, fmt.Errorf("connect: %w", err)
	}
	return dev, nil
}

// ConnectPeerFromScan подключается, используя адрес из рекламного пакета.
// Предпочитай этот вариант сразу после Scan: адрес совпадает с тем, что видел стек.
func ConnectPeerFromScan(result bluetooth.ScanResult) (bluetooth.Device, error) {
	if err := adapter.Enable(); err != nil {
		return bluetooth.Device{}, fmt.Errorf("enable adapter: %w", err)
	}
	dev, err := adapter.Connect(result.Address, bluetooth.ConnectionParams{})
	if err != nil {
		return bluetooth.Device{}, fmt.Errorf("connect: %w", err)
	}
	return dev, nil
}

// DiscoverAllServices возвращает все сервисы периферии (uuids == nil в API).
func DiscoverAllServices(dev bluetooth.Device) ([]bluetooth.DeviceService, error) {
	services, err := dev.DiscoverServices(nil)
	if err != nil {
		return nil, fmt.Errorf("discover services: %w", err)
	}
	return services, nil
}
