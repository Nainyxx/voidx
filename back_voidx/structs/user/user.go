package user

/*
UID - уникальный набор 12 байт для каждого пользователя (нужно в 
условиях отсутствия сервер чтобы не было конфликтов) 
(число возможных вариантов 256^12 = 2^96 = 7.2 * 10^28)
Name - имя пользователя
PublicKey - публичный ключ пользователя (нужно для передачи сообщения и идентифиации)
PrivateKey - приватный ключ пользователя (нужно для шифрования сообщения)
*/

type User struct {
	UID        [12]byte `json:"uid"`
	Name       string   `json:"name"`
	PublicKey  string   `json:"public_key"`
	PrivateKey string   `json:"private_key"`
}
