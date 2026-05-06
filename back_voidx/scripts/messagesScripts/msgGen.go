package messagesscripts

import (
	"fmt"
	"time"
	"voidx/constants"
	"voidx/scripts/randomGen"
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
