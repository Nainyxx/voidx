export const seedNodes = [
  { id: 'root', parentId: null, label: 'Все чаты' },
  { id: 'work', parentId: 'root', label: 'Работа' },
  { id: 'friends', parentId: 'root', label: 'Друзья' },
  { id: 'family', parentId: 'root', label: 'Семья' },
  { id: 'design', parentId: 'work', label: 'Дизайн' },
  { id: 'backend', parentId: 'work', label: 'Backend' },
]

export const seedPositions = {
  root: { x: 720, y: 110 },
  work: { x: 420, y: 360 },
  friends: { x: 720, y: 360 },
  family: { x: 1020, y: 360 },
  design: { x: 300, y: 640 },
  backend: { x: 540, y: 640 },
}
