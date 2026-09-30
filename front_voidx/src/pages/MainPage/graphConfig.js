import { seedNodes, seedPositions } from '../../data/graph/nodes'
import { seedSlots } from '../../data/graph/contacts'

const POSITIONS_KEY = 'voidx_node_positions'
const SLOTS_KEY = 'voidx_node_slots'
const EMPTY_SLOTS = [null, null, null, null, null, null]

function readJSON(key, fallback) {
  const raw = localStorage.getItem(key)
  return raw ? JSON.parse(raw) : fallback
}

export function getNodes() {
  const positions = readJSON(POSITIONS_KEY, {})
  return seedNodes.map((node) => {
    const pos = positions[node.id] ?? seedPositions[node.id] ?? { x: 0, y: 0 }
    return { ...node, x: pos.x, y: pos.y }
  })
}

export function saveNodePosition(id, x, y) {
  const positions = readJSON(POSITIONS_KEY, {})
  positions[id] = { x: Math.round(x), y: Math.round(y) }
  localStorage.setItem(POSITIONS_KEY, JSON.stringify(positions))
}

export function updateNodePosition(nodes, id, x, y) {
  return nodes.map((n) => (n.id === id ? { ...n, x, y } : n))
}

export function getNodeDepth(nodes, id) {
  let depth = 0
  let current = nodes.find((n) => n.id === id)
  while (current && current.parentId) {
    depth += 1
    current = nodes.find((n) => n.id === current.parentId)
  }
  return depth
}

export function getAncestorChain(nodes, id) {
  const chain = []
  let current = nodes.find((n) => n.id === id)
  while (current) {
    chain.unshift(current.id)
    current = current.parentId ? nodes.find((n) => n.id === current.parentId) : null
  }
  return chain
}

export function getEdges(nodes) {
  return nodes
    .filter((n) => n.parentId)
    .map((n) => {
      const parent = nodes.find((p) => p.id === n.parentId)
      return { id: `${n.parentId}-${n.id}`, parentId: n.parentId, childId: n.id, x1: parent.x, y1: parent.y, x2: n.x, y2: n.y }
    })
}

export function getNodeSize(depth) {
  if (depth === 0) return 96
  if (depth === 1) return 80
  return 68
}

export function getSlots(nodeId) {
  const allSlots = readJSON(SLOTS_KEY, {})
  return allSlots[nodeId] ?? seedSlots[nodeId] ?? [...EMPTY_SLOTS]
}

export function setSlot(nodeId, slotIndex, chatId) {
  const allSlots = readJSON(SLOTS_KEY, {})
  const current = allSlots[nodeId] ?? seedSlots[nodeId] ?? [...EMPTY_SLOTS]
  const updated = [...current]
  updated[slotIndex] = chatId
  allSlots[nodeId] = updated
  localStorage.setItem(SLOTS_KEY, JSON.stringify(allSlots))
  return updated
}
