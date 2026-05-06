# VoidX Быстрый старт

## Что было добавлено?

### 1. **Криптография сообщений** (`scripts/cryptography/encryption.go`)
- AES-256-GCM шифрование для всех сообщений
- Функции `EncryptMessage()` и `DecryptMessage()`
- Автоматическое кодирование в Base64

### 2. **Обмен сообщениями через Bluetooth** (`scripts/bluetooth/messageExchange.go`)
- `SendMessageToPeer()` - отправка зашифрованного сообщения
- `ReceiveAndProcessMessage()` - приём и дешифровка сообщений
- Структура `MessageTransfer` для передачи данных
- Верификация отправителя по публичному ключу

### 3. **Хранилище чатов** (`scripts/messagesScripts/chatStorage.go`)
- `SaveMessage()` - сохранение сообщений в JSON
- `GetMessages()` - получение чата с конкретным пользователем
- `GetAllMessages()` - получение всех чатов
- `DeleteMessage()` - удаление сообщений

### 4. **REST API** (`api.go`)
- `/api/user` - информация о текущем пользователе
- `/api/users` - список всех пользователей
- `/api/messages?user_uid=XXX` - сообщения из чата
- `/api/message/send` - отправка сообщения (локально)

### 5. **Основной цикл** (`main.go`)
- Инициализация пользователей и чатов
- Запуск API сервера на порту 8080
- Готовность к подключению Bluetooth компонента

---

## Запуск

```bash
cd back_voidx
go mod tidy
go run main.go
```

**Результат:**
```
[MAIN] Запуск voidx децентрализованного мессенджера...
[INFO] Найдено пользователей: 1
  1. watson (UID: [178 171 59 178 209 9 202 187 234 123 161 23])
[API] Сервер запущен на :8080
```

---

## Быстрые тесты API

### 1. Получить текущего пользователя
```bash
curl http://localhost:8080/api/user
```

### 2. Получить всех пользователей
```bash
curl http://localhost:8080/api/users
```

### 3. Получить сообщения
```bash
curl "http://localhost:8080/api/messages?user_uid=watson"
```

### 4. Отправить сообщение
```bash
curl -X POST http://localhost:8080/api/message/send \
  -H "Content-Type: application/json" \
  -d '{"receiver_uid":"sanya","content":"Привет!"}'
```

---

## Архитектура потока сообщений

```
┌─────────────┐
│   Sender    │
└──────┬──────┘
       │
       ▼ [1. GenerateMessage]
   ┌─────────────┐
   │ Message ID  │  (sender_key-receiver_uid-random_id)
   └──────┬──────┘
          │
          ▼ [2. EncryptMessage]
   ┌──────────────────┐
   │  AES-256-GCM     │  (base64 encoded)
   └──────┬───────────┘
          │
          ▼ [3. CreateTransfer]
   ┌──────────────────┐
   │  MessageTransfer │  (JSON)
   └──────┬───────────┘
          │
    ┌─────▼──────┐
    │  Bluetooth  │  (Peer to Peer)
    └─────┬──────┘
          │
          ▼ [4. ReceiveData]
   ┌──────────────────┐
   │ MessageTransfer  │
   └──────┬───────────┘
          │
          ▼ [5. DecryptMessage]
   ┌──────────────────┐
   │ Plaintext Content│
   └──────┬───────────┘
          │
          ▼ [6. SaveMessage]
   ┌───────────────────┐
   │   Chat JSON File  │
   └────────┬──────────┘
            │
            ▼
        ┌─────────┐
        │ Receiver│
        └─────────┘
```

---

## Структура данных

### В памяти (Go):
```go
// Сообщение
type Message struct {
    ID        string    // Уникальный ID
    Content   string    // Зашифрованное (base64)
    Timestamp time.Time // Когда было создано
}

// Отправка по Bluetooth
type MessageTransfer struct {
    MessageID         string // ID сообщения
    SenderUID         string // UID отправителя
    SenderKey         string // Публичный ключ
    ReceiverUID       string // UID получателя
    EncryptedContent  string // AES-256-GCM (base64)
    Timestamp         string // ISO format
}
```

### На диске (JSON):
```
data/
├── users.json
├── actual_devices_around.json
└── chats/
    ├── watson/
    │   └── chat.json
    └── sanya/
        └── chat.json
```

---

## Следующие шаги

### Backend:
- [ ] Реализовать GATT сервисы для полного Bluetooth обмена
- [ ] Добавить webocket для real-time уведомлений
- [ ] Реализовать mesh-сеть для ретрансляции сообщений
- [ ] Добавить систему верификации устройств

### Frontend (React):
- [ ] Компонент списка чатов
- [ ] Компонент окна чата с сообщениями
- [ ] Форма отправки сообщений
- [ ] Интеграция с API (fetch/axios)
- [ ] Панель пользователя
- [ ] Обнаружение Bluetooth устройств

---

## Файлы с документацией

- [API_DOCUMENTATION.md](API_DOCUMENTATION.md) - Полное описание REST API
- [scripts/cryptography/](scripts/cryptography/) - Шифрование сообщений
- [scripts/bluetooth/messageExchange.go](scripts/bluetooth/messageExchange.go) - P2P обмен
- [scripts/messagesScripts/chatStorage.go](scripts/messagesScripts/chatStorage.go) - Хранилище чатов

---

## Безопасность

⚠️ **ВАЖНО**:
1. Все сообщения шифруются AES-256-GCM
2. Приватные ключи хранятся локально (защитите `data/users.json`)
3. Публичные ключи используются для верификации отправителя
4. Bluetooth требует ручной верификации устройств перед доверием
