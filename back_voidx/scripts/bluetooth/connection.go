package bluetooth

import (
	"encoding/json"
	"fmt"

	"voidx/constants"
	"voidx/structs/message"

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

// FindServiceByUUID ищет сервис по UUID среди обнаруженных сервисов.
func FindServiceByUUID(services []bluetooth.DeviceService, uuid string) (bluetooth.DeviceService, error) {
	for _, svc := range services {
		if svc.UUID().String() == uuid {
			return svc, nil
		}
	}
	return bluetooth.DeviceService{}, fmt.Errorf("service with UUID %s not found", uuid)
}

// FindCharacteristicByUUID ищет характеристику по UUID в сервисе.
func FindCharacteristicByUUID(svc bluetooth.DeviceService, uuid string) (bluetooth.DeviceCharacteristic, error) {
	chars, err := svc.DiscoverCharacteristics(nil)
	if err != nil {
		return bluetooth.DeviceCharacteristic{}, fmt.Errorf("discover characteristics: %w", err)
	}
	for _, char := range chars {
		if char.UUID().String() == uuid {
			return char, nil
		}
	}
	return bluetooth.DeviceCharacteristic{}, fmt.Errorf("characteristic with UUID %s not found", uuid)
}

// SendMessage отправляет JSON message через BLE характеристику.
func SendMessage(dev bluetooth.Device, msg message.Message) error {
	services, err := DiscoverAllServices(dev)
	if err != nil {
		return err
	}
	svc, err := FindServiceByUUID(services, constants.BLE_SERVICE_UUID)
	if err != nil {
		return err
	}
	char, err := FindCharacteristicByUUID(svc, constants.BLE_MESSAGE_CHAR_UUID)
	if err != nil {
		return err
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}
	_, err = char.WriteWithoutResponse(data)
	return err
}

// SendFeedback отправляет JSON feedback через BLE характеристику.
func SendFeedback(dev bluetooth.Device, fb message.Feedback) error {
	services, err := DiscoverAllServices(dev)
	if err != nil {
		return err
	}
	svc, err := FindServiceByUUID(services, constants.BLE_SERVICE_UUID)
	if err != nil {
		return err
	}
	char, err := FindCharacteristicByUUID(svc, constants.BLE_FEEDBACK_CHAR_UUID)
	if err != nil {
		return err
	}
	data, err := json.Marshal(fb)
	if err != nil {
		return fmt.Errorf("marshal feedback: %w", err)
	}
	_, err = char.WriteWithoutResponse(data)
	return err
}

// ReceiveMessage включает уведомления для характеристики message и возвращает канал для получения JSON.
func ReceiveMessage(dev bluetooth.Device) (<-chan message.Message, error) {
	services, err := DiscoverAllServices(dev)
	if err != nil {
		return nil, err
	}
	svc, err := FindServiceByUUID(services, constants.BLE_SERVICE_UUID)
	if err != nil {
		return nil, err
	}
	char, err := FindCharacteristicByUUID(svc, constants.BLE_MESSAGE_CHAR_UUID)
	if err != nil {
		return nil, err
	}
	ch := make(chan message.Message, 10)
	err = char.EnableNotifications(func(buf []byte) {
		var msg message.Message
		if err := json.Unmarshal(buf, &msg); err != nil {
			fmt.Printf("Error unmarshaling message: %v\n", err)
			return
		}
		ch <- msg
	})
	if err != nil {
		close(ch)
		return nil, fmt.Errorf("enable notifications: %w", err)
	}
	return ch, nil
}

// ReceiveFeedback включает уведомления для характеристики feedback и возвращает канал для получения JSON.
func ReceiveFeedback(dev bluetooth.Device) (<-chan message.Feedback, error) {
	services, err := DiscoverAllServices(dev)
	if err != nil {
		return nil, err
	}
	svc, err := FindServiceByUUID(services, constants.BLE_SERVICE_UUID)
	if err != nil {
		return nil, err
	}
	char, err := FindCharacteristicByUUID(svc, constants.BLE_FEEDBACK_CHAR_UUID)
	if err != nil {
		return nil, err
	}
	ch := make(chan message.Feedback, 10)
	err = char.EnableNotifications(func(buf []byte) {
		var fb message.Feedback
		if err := json.Unmarshal(buf, &fb); err != nil {
			fmt.Printf("Error unmarshaling feedback: %v\n", err)
			return
		}
		ch <- fb
	})
	if err != nil {
		close(ch)
		return nil, fmt.Errorf("enable notifications: %w", err)
	}
	return ch, nil
}
