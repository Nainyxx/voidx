package messagesscripts

import (
	"fmt"
	"os"
	"path/filepath"
	"voidx/constants"
)

func CheckChat(name string) (bool, error) {
	path := filepath.Join(constants.CHATS_DATA_PATH, name)
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("ошибка при проверке папки: %w", err)
	}
	return info.IsDir(), nil
}

func InitChat(name string) error {
	chatPath := filepath.Join(constants.CHATS_DATA_PATH, name)

	if err := os.MkdirAll(chatPath, 0755); err != nil {
		return fmt.Errorf("ошибка создания папки чата: %w", err)
	}

	jsonPath := filepath.Join(chatPath, "chat.json")
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		if err := os.WriteFile(jsonPath, []byte("[]"), 0644); err != nil {
			return fmt.Errorf("ошибка создания JSON: %w", err)
		}
	}

	return nil
}
