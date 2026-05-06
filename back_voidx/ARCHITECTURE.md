# VoidX Architecture

## Высокоуровневая архитектура

```
┌─────────────────────────────────────────────────────────────────┐
│                     Frontend (React + Vite)                     │
│  ┌──────────────┐  ┌──────────────┐  ┌────────────────────┐   │
│  │  Chat List   │  │  Chat Window │  │  User Profile      │   │
│  └──────┬───────┘  └──────┬───────┘  └──────┬─────────────┘   │
│         │                 │                  │                  │
│         └─────────────────┼──────────────────┘                  │
│                           │                                      │
│                   ┌───────▼────────┐                            │
│                   │   HTTP API     │                            │
│                   │  (fetch/axios) │                            │
│                   └───────┬────────┘                            │
└───────────────────────────┼────────────────────────────────────┘
                            │ HTTP (8080)
                            │
┌───────────────────────────▼────────────────────────────────────┐
│                   Backend (Go)                                  │
│                                                                 │
│  ┌────────────────────────────────────────────────────────┐   │
│  │                 REST API Server (api.go)              │   │
│  │  GET  /api/user                                       │   │
│  │  GET  /api/users                                      │   │
│  │  GET  /api/messages?user_uid=XXX                     │   │
│  │  POST /api/message/send                              │   │
│  └────────────────────────────────────────────────────────┘   │
│                           ▲                                     │
│                           │                                     │
│  ┌────────────────────────┼────────────────────────────────┐  │
│  │  Message Processing    │                               │  │
│  │  ┌──────────────────────▼────────────────────────┐    │  │
│  │  │  Message Generation & Encryption             │    │  │
│  │  │  ┌─────────────────────────────────────────┐ │    │  │
│  │  │  │ GenerateID() (randomGen)                │ │    │  │
│  │  │  │ EncryptMessage() (cryptography)         │ │    │  │
│  │  │  │ Create MessageTransfer JSON             │ │    │  │
│  │  │  └─────────────────────────────────────────┘ │    │  │
│  │  └──────────────────┬─────────────────────────────┘    │  │
│  │                     │                                   │  │
│  │  ┌──────────────────▼─────────────────────────────┐    │  │
│  │  │  Message Routing                              │    │  │
│  │  │  ┌─────────────────────────────────────────┐  │    │  │
│  │  │  │ Bluetooth (P2P to peer devices)         │  │    │  │
│  │  │  │ API (local storage + frontend)          │  │    │  │
│  │  │  └─────────────────────────────────────────┘  │    │  │
│  │  └──────────────────┬─────────────────────────────┘    │  │
│  │                     │                                   │  │
│  │  ┌──────────────────▼─────────────────────────────┐    │  │
│  │  │  Storage & Persistence                        │    │  │
│  │  │  ┌─────────────────────────────────────────┐  │    │  │
│  │  │  │ SaveMessage() (chatStorage)             │  │    │  │
│  │  │  │ JSON Files (data/chats/)                │  │    │  │
│  │  │  └─────────────────────────────────────────┘  │    │  │
│  │  └──────────────────────────────────────────────┘    │  │
│  └────────────────────────────────────────────────────────┘  │
│                           │                                    │
│  ┌────────────────────────▼─────────────────────────────┐    │
│  │  Bluetooth Component (scripts/bluetooth/)            │    │
│  │  ┌───────────────────────────────────────────────┐  │    │
│  │  │ blueScan.go: ScanDevices()                   │  │    │
│  │  │ connection.go: ConnectPeer(), DiscoverServices│  │    │
│  │  │ messageExchange.go: SendMessage/Receive      │  │    │
│  │  │ mainBluetoothComponent.go: Orchestration     │  │    │
│  │  └───────────────────────────────────────────────┘  │    │
│  └────────────────────────────────────────────────────────┘  │
│                           │                                    │
└───────────────────────────┼────────────────────────────────────┘
                            │ Bluetooth (Low Energy)
                            │
        ┌───────────────────┼───────────────────┐
        │                   │                   │
        ▼                   ▼                   ▼
    ┌────────┐          ┌────────┐          ┌────────┐
    │ Device │          │ Device │          │ Device │
    │   1    │          │   2    │          │   3    │
    └────────┘          └────────┘          └────────┘
   (Local App)        (Local App)         (Local App)
```

---

## Message Flow Diagram

### 1. Sending Flow

```
┌───────────────────────────────────────────────────────────────┐
│                    SENDING MESSAGE FLOW                        │
└───────┬─────────────────────────────────────────────────────┬─┘
        │                                                       │
        ▼                                                       ▼
    ┌─────────────────┐                          ┌──────────────────┐
    │  Frontend (React) │                          │  Backend (api.go)│
    │                  │                          │                  │
    │ 1. User clicks   │ POST /api/message/send   │ 1. Parse request │
    │    "Send Button" │─────────────────────────>│ 2. GenerateID()  │
    │                  │                          │ 3. SaveMessage() │
    │                  │                          │ 4. Return response│
    │ 2. Display msg   │ {success: true, ...}     │                  │
    │    in chat UI    │<─────────────────────────│                  │
    └──────────────────┘                          └──────────────────┘
                                                           │
                                                  (Bluetooth Ready)
                                                           │
                                                           ▼
                                  ┌────────────────────────────────┐
                                  │  Bluetooth GATT Write           │
                                  │ messageExchange.SendMessageToPeer│
                                  │                                 │
                                  │ 1. GenerateID()               │
                                  │ 2. EncryptMessage()            │
                                  │ 3. Create MessageTransfer      │
                                  │ 4. Serialize to JSON           │
                                  │ 5. Write via GATT              │
                                  └────────────────────────────────┘
```

### 2. Receiving Flow

```
┌───────────────────────────────────────────────────────────────┐
│                   RECEIVING MESSAGE FLOW                       │
└───────┬─────────────────────────────────────────────────────┬─┘
        │                                                       │
        ▼                                                       ▼
    ┌──────────────────┐                          ┌──────────────────┐
    │ Bluetooth Device │                          │  Our Device      │
    │                  │                          │                  │
    │ Sends JSON via   │  GATT Notification      │ 1. Receive JSON  │
    │ GATT Write       │─────────────────────────>│ 2. Parse Transfer│
    │ (MessageTransfer)│                          │ 3. Verify Sender │
    │                  │                          │ 4. DecryptMessage│
    │                  │                          │ 5. SaveMessage() │
    │                  │                          │ 6. Notify API    │
    └──────────────────┘                          └──────────────────┘
                                                           │
                                                           ▼
                                  ┌────────────────────────────────┐
                                  │  Update Frontend (WebSocket)   │
                                  │  OR Poll via API               │
                                  └────────────────────────────────┘
                                                           │
                                                           ▼
                                  ┌────────────────────────────────┐
                                  │  Frontend Shows New Message    │
                                  │  ✓ Message appears in chat     │
                                  └────────────────────────────────┘
```

### 3. Encryption Flow

```
Plain Text:
"Привет, как дела?"

        │
        ▼ EncryptMessage(content, privateKey)

┌─────────────────────────────────────────┐
│  1. Create AES-256 cipher               │
│  2. Generate random nonce (12 bytes)    │
│  3. Use GCM mode for authentication     │
│  4. Encrypt: cipher.Seal(nonce + msg)   │
│  5. Encode to Base64                    │
└─────────────────────────────────────────┘

        │
        ▼ Base64 Encoded Ciphertext:

"xK9L2nM+...pqR3sT4uV5w..." (looks like gibberish)

        │
        ├─> Sent via Bluetooth
        │
        ▼ DecryptMessage(encryptedBase64, publicKey)

┌─────────────────────────────────────────┐
│  1. Decode Base64                       │
│  2. Extract nonce (first 12 bytes)      │
│  3. Create same AES-256 cipher          │
│  4. Decrypt: cipher.Open(nonce, data)   │
│  5. Return plaintext                    │
└─────────────────────────────────────────┘

        │
        ▼ Decrypted:

"Привет, как дела?"
```

---

## File Structure & Data Flow

```
data/
├── users.json
│   └─ Contains: [User{UID, Name, PublicKey, PrivateKey}]
│
├── actual_devices_around.json
│   └─ Contains: [Device{Name, MAC, RSSI, LastSeen}]
│
├── chats/
│   ├── watson/
│   │   └── chat.json
│   │       └─ Contains: [Message{ID, Content, Timestamp}]
│   │
│   └── sanya/
│       └── chat.json
│           └─ Contains: [Message{ID, Content, Timestamp}]
│
└── (future: encrypted keyring)
```

---

## Security Layers

```
┌─────────────────────────────────────────────────────────┐
│        Layer 1: Device Verification (Bluetooth)         │
│   - Whitelist trusted devices by MAC address           │
│   - Verify peer capabilities (GATT UUID)               │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│     Layer 2: Message Authentication (Public Key)        │
│   - Verify sender identity via sender's public key     │
│   - Check message ID format matches expected pattern   │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│   Layer 3: Message Encryption (AES-256-GCM)           │
│   - Each message encrypted with sender's private key   │
│   - GCM provides authentication tag                    │
│   - Random nonce prevents replay attacks              │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│      Layer 4: Transport Security (Bluetooth LE)         │
│   - Bluetooth LE with optional pairing/bonding        │
│   - 40-bit AES encryption at BLE layer                │
└─────────────────────────────────────────────────────────┘
```

---

## Component Dependencies

```
main.go
├── scripts/usersScripts/userInit.go
├── scripts/messagesscripts/chatInit.go
├── api.go
│   ├── scripts/usersScripts/userInit.go
│   ├── scripts/messagesscripts/chatStorage.go
│   └── scripts/messagesscripts/msgGen.go
│
├── scripts/bluetooth/mainBluetoothComponent.go
│   ├── scripts/bluetooth/blueScan.go
│   ├── scripts/bluetooth/connection.go
│   ├── scripts/bluetooth/messageExchange.go
│   │   ├── scripts/cryptography/encryption.go
│   │   ├── scripts/messagesscripts/chatStorage.go
│   │   └── scripts/randomGen/randomGen.go
│   └── scripts/messagesscripts/msgGen.go
│
└── scripts/randomGen/randomGen.go
└── Constants (used everywhere)
```

---

## Data Models

```go
// User - Пользователь системы децентрализованного мессенджера
type User struct {
    UID        [12]byte // 96-bit уникальный ID (без серверов!)
    Name       string   // Человеческое имя
    PublicKey  string   // 40 символов - для идентификации
    PrivateKey string   // 40 символов - для шифрования
}

// Message - Сообщение в чате
type Message struct {
    ID        string    // "sender_pubkey-receiver_uid-random_id"
    Content   string    // Зашифровано: base64(AES256-GCM(...))
    Timestamp time.Time // Когда было отправлено
}

// Device - Bluetooth устройство
type Device struct {
    Name string // "iPhone 12 Pro"
    MAC  string // "AA:BB:CC:DD:EE:FF"
    RSSI int    // -85 (чем выше = ближе)
    Time string // "2026-05-06 14:30:45"
}

// MessageTransfer - Структура для Bluetooth передачи
type MessageTransfer struct {
    MessageID        string // Уникальный ID
    SenderUID        string // "watson_uid_bytes"
    SenderKey        string // Публичный ключ для верификации
    ReceiverUID      string // "sanya_uid_bytes"
    EncryptedContent string // base64(...) 
    Timestamp        string // RFC3339 format
}
```

---

## Future Enhancements

```
┌─────────────────────────────────────────────────────────┐
│         Planned: Mesh Network Relay (v2.0)              │
│                                                         │
│  Device A ──> Device B ──> Device C                    │
│  (sender)    (relay)      (receiver)                    │
│                                                         │
│  - Hop counters to prevent loops                       │
│  - Adaptive routing based on RSSI                      │
│  - TTL (Time To Live) for message expiry               │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│      Planned: Real-time WebSocket (v2.0)               │
│                                                         │
│  Frontend ←──→ Backend ←──→ Bluetooth                  │
│              WebSocket         Devices                  │
│                                                         │
│  - Live notifications for new messages                 │
│  - Device discovery notifications                      │
│  - Typing indicators                                   │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│    Planned: Group Chats (v2.0)                          │
│                                                         │
│  - Multiple recipients per message                     │
│  - Group discovery & management                        │
│  - Admin controls                                      │
└─────────────────────────────────────────────────────────┘
```
