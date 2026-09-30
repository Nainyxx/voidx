export const seedChats = {
  mark: {
    id: 'mark',
    name: 'Марк',
    initial: 'М',
    online: true,
    messages: [
      { from: 'them', text: 'Го в кино?' },
      { from: 'me', text: 'давай в 19' },
      { from: 'them', text: 'договорились' },
    ],
  },
  anya: {
    id: 'anya',
    name: 'Аня',
    initial: 'А',
    online: true,
    messages: [
      { from: 'them', text: 'видел макет?' },
      { from: 'me', text: 'да, красиво' },
      { from: 'them', text: 'рада что зашло' },
    ],
  },
  familyChat: {
    id: 'familyChat',
    name: 'Семейный чат',
    initial: 'С',
    online: false,
    messages: [
      { from: 'them', text: 'ужин в 8' },
      { from: 'me', text: 'буду' },
      { from: 'them', text: 'не опаздывай' },
    ],
  },
  mom: {
    id: 'mom',
    name: 'Мама',
    initial: 'М',
    online: true,
    messages: [
      { from: 'them', text: 'как дела?' },
      { from: 'me', text: 'всё хорошо' },
      { from: 'them', text: 'скоро позвоню' },
    ],
  },
  teamlead: {
    id: 'teamlead',
    name: 'Тимлид',
    initial: 'Т',
    online: false,
    messages: [
      { from: 'them', text: 'созвон в 15:00' },
    ],
  },
  designChat: {
    id: 'designChat',
    name: 'Дизайн-чат',
    initial: 'Д',
    online: true,
    messages: [
      { from: 'them', text: '3 новых сообщения' },
    ],
  },
  igor: {
    id: 'igor',
    name: 'Игорь',
    initial: 'И',
    online: false,
    messages: [
      { from: 'them', text: 'был в сети недавно' },
    ],
  },
}

// 6 slots per folder: [topLeft, topCenter, topRight, bottomLeft, bottomCenter, bottomRight]
export const seedSlots = {
  family: ['mark', 'anya', null, null, 'familyChat', 'mom'],
}
