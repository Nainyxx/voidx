package message

import "time"

/*
ID - уникальный идентифактор любого сообщения
ID состоит из 3 частей:
1. PublicKey отправителя
2. PublicKey получателя
3. Контекстный ID (который генерируется случайным образом)
Представляется в виде строки: "senderPublicKey-receiverPublicKey-randomcontextID"

Content - содержимое сообщения

Timestamp - время отправки сообщения
Представляется в виде строки: "YYYY-MM-DD HH:MM:SS"

Все это вместе шифруется приватным ключом отправителя
*/

type Message struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	SenderMAC string    `json:"sender_mac"`
}

// ID фидбека формируется из ID сообщения + флага "-f", что означает -feedback
// При получении любого сообщения/фидбека человек сохраняет MAC адрес отправителя

type Feedback struct {
	ID                string `json:"id"`
	Status            int    `json:"status"` // status 200 и тд
	FeedbackSenderMAC string `json:"feedback_sender_mac"`
}
