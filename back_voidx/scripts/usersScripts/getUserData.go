package usersScripts

import ()

func GetUserData() (user.User, error) {
	user := os.ReadFile(constants.USER_DATA_PATH)
	fmt.Printf(user)
} 