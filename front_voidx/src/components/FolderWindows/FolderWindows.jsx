import MiniChatWindow from '../MiniChatWindow/MiniChatWindow'
import EmptySlot from '../MiniChatWindow/EmptySlot'

const SLOT_W = 220
const SLOT_H = 200
const GAP_X = 40
const GAP_Y = 40

export function getSlotPositions(node, nodeSize) {
  const radius = nodeSize / 2
  const topY = node.y - radius - GAP_Y - SLOT_H
  const bottomY = node.y + radius + GAP_Y
  const xs = [
    node.x - SLOT_W / 2 - (SLOT_W + GAP_X),
    node.x - SLOT_W / 2,
    node.x - SLOT_W / 2 + (SLOT_W + GAP_X),
  ]
  return [
    { x: xs[0], y: topY },
    { x: xs[1], y: topY },
    { x: xs[2], y: topY },
    { x: xs[0], y: bottomY },
    { x: xs[1], y: bottomY },
    { x: xs[2], y: bottomY },
  ]
}

export default function FolderWindows({ node, nodeSize, slots, chats, onOpenPicker, onSend, onRemove }) {
  const positions = getSlotPositions(node, nodeSize)

  return (
    <>
      {slots.map((chatId, i) => {
        const pos = positions[i]
        if (!chatId || !chats[chatId]) {
          return (
            <EmptySlot
              key={i}
              x={pos.x}
              y={pos.y}
              onClick={() => onOpenPicker(i, pos)}
            />
          )
        }
        return (
          <MiniChatWindow
            key={i}
            chat={chats[chatId]}
            x={pos.x}
            y={pos.y}
            onSend={(text) => onSend(chatId, text)}
            onRemove={() => onRemove(i)}
          />
        )
      })}
    </>
  )
}
