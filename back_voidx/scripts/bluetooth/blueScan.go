package bluetooth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"voidx/constants"
	"voidx/structs/device"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

// saveDevicesToJSON сохраняет устройства в JSON-файл, используя структуру из пакета device
func saveDevicesToJSON(currentDevices map[string]bluetooth.ScanResult) error {
	filePath := filepath.Join(constants.DATA_PATH, "actual_devices_around.json")

	// Читаем существующие устройства из JSON
	existingDevices := make(map[string]device.Device)
	if data, err := os.ReadFile(filePath); err == nil {
		var devices []device.Device
		if err := json.Unmarshal(data, &devices); err == nil {
			for _, d := range devices {
				existingDevices[d.MAC] = d
			}
		}
	}

	// Создаём новый список с актуальными устройствами
	var updatedDevices []device.Device

	// Добавляем текущие устройства
	for mac, d := range currentDevices {
		updatedDevices = append(updatedDevices, device.Device{
			Name: d.LocalName(),
			MAC:  mac,
			RSSI: int(d.RSSI),
			Time: time.Now().Format("2006-01-02 15:04:05"),
		})
		// Удаляем из existingDevices те, которые ещё есть
		delete(existingDevices, mac)
	}

	// existingDevices теперь содержит только пропавшие устройства
	if len(existingDevices) > 0 {
		fmt.Printf("[BLE] Удалено пропавших устройств: %d\n", len(existingDevices))
		for mac := range existingDevices {
			fmt.Printf("[BLE] - Удалено: %s\n", mac)
		}
	}

	// Сохраняем обновлённый список в JSON
	data, err := json.MarshalIndent(updatedDevices, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации JSON: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("ошибка записи файла: %w", err)
	}

	fmt.Printf("[BLE] Сохранено устройств: %d\n", len(updatedDevices))
	return nil
}

func ScanDevices(duration time.Duration) error {
	if err := adapter.Enable(); err != nil {
		return fmt.Errorf("не удалось включить Bluetooth: %w", err)
	}

	if err := os.MkdirAll(constants.DATA_PATH, 0755); err != nil {
		return fmt.Errorf("не удалось создать папку данных: %w", err)
	}

	deviceMap := make(map[string]bluetooth.ScanResult)
	stopScan := make(chan struct{})

	// Таймер для сохранения (каждые 5 секунд)
	saveTicker := time.NewTicker(5 * time.Second)
	defer saveTicker.Stop()

	// Горутина для периодического сохранения
	go func() {
		for range saveTicker.C {
			if err := saveDevicesToJSON(deviceMap); err != nil {
				fmt.Printf("[BLE] Ошибка сохранения: %v\n", err)
			} else {
				fmt.Printf("[BLE] Сохранено %d устройств\n", len(deviceMap))
			}
		}
	}()

	// Запускаем сканирование
	err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		name := device.LocalName()
		if name == "" {
			return
		}
		mac := device.Address.String()

		if _, exists := deviceMap[mac]; !exists {
			deviceMap[mac] = device
			fmt.Printf("[BLE] Найдено: %s [%s] RSSI: %d\n", name, mac, device.RSSI)
		}
	})
	if err != nil {
		return fmt.Errorf("ошибка сканирования: %w", err)
	}

	// Останавливаем сканирование через duration
	time.AfterFunc(duration, func() {
		adapter.StopScan()
		close(stopScan)
	})

	<-stopScan

	// Финальное сохранение
	saveTicker.Stop()
	if err := saveDevicesToJSON(deviceMap); err != nil {
		return fmt.Errorf("ошибка финального сохранения: %w", err)
	}

	fmt.Printf("[BLE] Сканирование завершено. Сохранено %d устройств\n", len(deviceMap))
	return nil
}

// ReadDevicesFromJSON читает сохранённые устройства из JSON
func ReadDevicesFromJSON() ([]device.Device, error) {
	filePath := filepath.Join(constants.DATA_PATH, "actual_devices_around.json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %w", err)
	}
	var devices []device.Device
	if err := json.Unmarshal(data, &devices); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}
	return devices, nil
}

/*
func Connect([]devices) {
	return 1
}

func SendMessage(msg message.Message) (result string, error string) {
	if err := adapter.Enable(); err != nil {
		return fmt.Errorf("не удалось включить Bluetooth: %w", err)
	}

	data, err = json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("ошибка сериализации сообщения из стринг в json. Проверьте формат стринга. %w", err)
	}

}
*/
