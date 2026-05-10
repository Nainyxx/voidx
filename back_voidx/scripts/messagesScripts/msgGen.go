package messagesscripts

import (
	"fmt"
	"time"
	"voidx/constants"
	"voidx/scripts/randomGen"
	"voidx/structs/message"
)

func GenerateMsgID(senderPubKey string, readerPubKey string) (string, error) {

	contextMsgID, err := randomGen.GenerateString(40, constants.ALPHABET)
	if err != nil {
		return "", fmt.Errorf("не удалось сгенерировать контекстный ID сообщения: %w", err)
	}
	msgID := senderPubKey + "-" + readerPubKey + "-" + contextMsgID
	return msgID, nil
}

func GenerateMsg(sender string, reader string, content string, timestamp time.Time) (message.Message, error) {
	
	Resulted_Message := {};

	contextMsgID, err := GenerateMsgID(sender, reader);
	if (err != nil) {
		return "", fmt.Errorf("Ошибка компановки сообщентя: %w", err)
	}
	Resulted_Message.ID = contextMsgID;
	Resulted_Message.content = content;
	Resulted_Message.SenderMAC = "HUYHUYHUY";

	return Resulted_Message, nil;
}

func GenerateMessage(content string, receiverMAC string) (message.Message, error) {
	// Для простоты используем MAC как ключи
	senderPubKey := "senderKey" // В реальности брать из пользователя
	readerPubKey := "readerKey" // В реальности брать из получателя

	id, err := GenerateMsgID(senderPubKey, readerPubKey)
	if err != nil {
		return message.Message{}, err
	}

	msg := message.Message{
		ID:        id,
		Content:   content,
		Timestamp: time.Now(),
		SenderMAC: receiverMAC, // В примере используем receiverMAC как senderMAC для теста
	}
	return msg, nil
}
