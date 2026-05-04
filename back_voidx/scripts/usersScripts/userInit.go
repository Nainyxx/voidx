package usersScripts

import (
	"encoding/json"
	"os"
	"path/filepath"

	"voidx/constants"
	"voidx/structs/user"
)

// CREATING USERS JSON FUNCTION (if json not found!)
func EnsureJSON(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	emptyData := []user.User{}
	data, err := json.MarshalIndent(emptyData, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// GETTING USERS FUNCTION
func ReadUsers() ([]user.User, error) {
	path := constants.USER_DATA_PATH
	if err := EnsureJSON(path); err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var users []user.User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// WRITING USERS INTO JSON
func WriteUsers(users []user.User) error {
	path := constants.USER_DATA_PATH
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func AddUser(u user.User) error {
	users, err := ReadUsers()
	if err != nil {
		return err
	}
	users = append(users, u)
	return WriteUsers(users)
}
