package usersScripts

import (
	"fmt"

	"voidx/constants"
	"voidx/scripts/randomGen"
	"voidx/structs/user"
)

func UserVerify() {
	EnsureJSON(constants.USER_DATA_PATH)
	users, err := ReadUsers()
	if err != nil {
		fmt.Printf("[ERROR] Не удалось считать пользователей.")
	}
	if len(users) == 0 {
		var user_name, public_key, private_key string
		var uid [12]byte

		fmt.Printf("\nОтсустствуют записи о пользователях. Давайте создадим нового. Введите имя:")

		fmt.Scanln(&user_name)

		uid, err = randomGen.GenerateUID()
		if err != nil {
			fmt.Printf("[ERROR] Не удалось сгенерировать UID")

		}
		public_key, err = randomGen.GenerateString(constants.KEYS_LEN, constants.ALPHABET)
		if err != nil {
			fmt.Printf("[ERROR] Не удалось сгенерировать публичный ключ")
		}
		private_key, err = randomGen.GenerateString(constants.KEYS_LEN, constants.ALPHABET)
		if err != nil {
			fmt.Printf("[ERROR] Не удалось сгенерировать приватный ключ")
		}

		new_user := user.User{
			Name:       user_name,
			PrivateKey: private_key,
			PublicKey:  public_key,
			UID:        uid,
		}
		fmt.Println(new_user)

		AddUser(new_user)

	} else {
		fmt.Println("[OK] User found")
	}
}
