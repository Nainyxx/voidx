package bluetooth

import (
	"fmt"
	"time"
	messagesscripts "voidx/scripts/messagesScripts"
	"voidx/scripts/usersScripts"
	"voidx/structs/message"
)

func MainBluetoothComponent() {
	// Пример использования: сканирование, подключение и отправка JSON

	// Сначала верификация пользователя
	usersScripts.UserVerify()

	// Сканируем устройства
	fmt.Println("Сканирование устройств...")
	err := ScanDevices(10 * time.Second)
	if err != nil {
		fmt.Printf("Ошибка сканирования: %v\n", err)
		return
	}

	// Предположим, мы подключаемся к первому найденному устройству
	// В реальности нужно выбрать по имени или MAC
	devices, err := ReadDevicesFromJSON()
	if err != nil {
		fmt.Printf("Ошибка загрузки устройств: %v\n", err)
		return
	}
	if len(devices) == 0 {
		fmt.Println("Нет найденных устройств")
		return
	}

	targetMAC := devices[0].MAC
	fmt.Printf("Подключение к устройству: %s\n", targetMAC)

	StopScanBestEffort()
	dev, err := ConnectPeer(targetMAC)
	if err != nil {
		fmt.Printf("Ошибка подключения: %v\n", err)
		return
	}
	defer dev.Disconnect()

	// Генерируем тестовое сообщение
	msg, err := messagesscripts.GenerateMessage("Тестовое сообщение", targetMAC)
	if err != nil {
		fmt.Printf("Ошибка генерации сообщения: %v\n", err)
		return
	}

	// Отправляем сообщение
	err = SendMessage(dev, msg)
	if err != nil {
		fmt.Printf("Ошибка отправки сообщения: %v\n", err)
		return
	}
	fmt.Println("Сообщение отправлено")

	// Генерируем feedback
	fb := message.Feedback{
		ID:                msg.ID + "-f",
		Status:            200,
		FeedbackSenderMAC: "00:11:22:33:44:55", // Пример MAC отправителя feedback
	}

	// Отправляем feedback
	err = SendFeedback(dev, fb)
	if err != nil {
		fmt.Printf("Ошибка отправки feedback: %v\n", err)
		return
	}
	fmt.Println("Feedback отправлен")

	// Для демонстрации получения, включаем уведомления
	msgChan, err := ReceiveMessage(dev)
	if err != nil {
		fmt.Printf("Ошибка включения уведомлений для message: %v\n", err)
	} else {
		go func() {
			for msg := range msgChan {
				fmt.Printf("Получено сообщение: %+v\n", msg)
			}
		}()
	}

	fbChan, err := ReceiveFeedback(dev)
	if err != nil {
		fmt.Printf("Ошибка включения уведомлений для feedback: %v\n", err)
	} else {
		go func() {
			for fb := range fbChan {
				fmt.Printf("Получен feedback: %+v\n", fb)
			}
		}()
	}

	// Ждем немного для демонстрации
	time.Sleep(5 * time.Second)
	fmt.Println("Завершение демонстрации")
}
