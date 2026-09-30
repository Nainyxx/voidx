import { useRef, useState } from 'react'
import Graph from '../../components/Graph/Graph'
import FolderWindows from '../../components/FolderWindows/FolderWindows'
import SlotPicker from '../../components/SlotPicker/SlotPicker'
import SettingsButton from '../../components/SettingsButton/SettingsButton'
import { seedChats } from '../../data/graph/contacts'
import {
  getNodes,
  saveNodePosition,
  updateNodePosition,
  getAncestorChain,
  getEdges,
  getNodeDepth,
  getNodeSize,
  getSlots,
  setSlot,
} from './graphConfig'
import './MainPage.css'

export default function MainPage() {
  const [nodes, setNodes] = useState(() => getNodes())
  const [draggingId, setDraggingId] = useState(null)
  const [selectedId, setSelectedId] = useState(null)
  const [slotsByNode, setSlotsByNode] = useState(() => {
    const initial = {}
    getNodes().forEach((n) => { initial[n.id] = getSlots(n.id) })
    return initial
  })
  const [chats, setChats] = useState(() => ({ ...seedChats }))
  const [picker, setPicker] = useState(null)

  const dragOffset = useRef({ dx: 0, dy: 0 })
  const dragStart = useRef({ x: 0, y: 0 })
  const hasDragged = useRef(false)

  const handleNodeMouseDown = (id, e) => {
    e.stopPropagation()
    const node = nodes.find((n) => n.id === id)
    if (!node) return
    setDraggingId(id)
    hasDragged.current = false
    dragStart.current = { x: e.clientX, y: e.clientY }
    dragOffset.current = { dx: e.clientX - node.x, dy: e.clientY - node.y }
  }

  const handleCanvasMouseMove = (e) => {
    if (!draggingId) return
    const movedX = e.clientX - dragStart.current.x
    const movedY = e.clientY - dragStart.current.y
    if (Math.abs(movedX) > 4 || Math.abs(movedY) > 4) hasDragged.current = true
    setNodes((prev) =>
      updateNodePosition(prev, draggingId, e.clientX - dragOffset.current.dx, e.clientY - dragOffset.current.dy)
    )
  }

  const handleCanvasMouseUp = () => {
    if (!draggingId) return
    if (hasDragged.current) {
      const node = nodes.find((n) => n.id === draggingId)
      const x = Math.round(node.x)
      const y = Math.round(node.y)
      setNodes((prev) => updateNodePosition(prev, draggingId, x, y))
      saveNodePosition(draggingId, x, y)
    } else {
      setSelectedId((prev) => (prev === draggingId ? null : draggingId))
      setPicker(null)
    }
    setDraggingId(null)
  }

  const handleCanvasClick = (e) => {
    if (e.target === e.currentTarget) {
      setSelectedId(null)
      setPicker(null)
    }
  }

  const handleOpenPicker = (slotIndex, pos) => {
    const x = Math.min(pos.x, window.innerWidth - 300)
    const y = Math.min(Math.max(pos.y, 20), window.innerHeight - 420)
    setPicker({ nodeId: selectedId, slotIndex, x, y })
  }

  const handleAssignChat = (chatId) => {
    if (!picker) return
    const updated = setSlot(picker.nodeId, picker.slotIndex, chatId)
    setSlotsByNode((prev) => ({ ...prev, [picker.nodeId]: updated }))
    setPicker(null)
  }

  const handleRemoveChat = (nodeId, slotIndex) => {
    const updated = setSlot(nodeId, slotIndex, null)
    setSlotsByNode((prev) => ({ ...prev, [nodeId]: updated }))
  }

  const handleSendMessage = (chatId, text) => {
    setChats((prev) => ({
      ...prev,
      [chatId]: { ...prev[chatId], messages: [...prev[chatId].messages, { from: 'me', text }] },
    }))
  }

  const edges = getEdges(nodes)
  const ancestorChain = selectedId ? getAncestorChain(nodes, selectedId) : []
  const highlightedEdges = new Set(
    ancestorChain.slice(0, -1).map((id, i) => `${id}-${ancestorChain[i + 1]}`)
  )

  const selectedNode = selectedId && selectedId !== draggingId ? nodes.find((n) => n.id === selectedId) : null
  const selectedSlots = selectedNode ? (slotsByNode[selectedNode.id] ?? [null, null, null, null, null, null]) : null

  const pickerContacts = picker
    ? Object.values(chats).filter((c) => !slotsByNode[picker.nodeId]?.includes(c.id))
    : []

  return (
    <div
      className={`canvas${draggingId ? ' canvas--dragging' : ''}`}
      onClick={handleCanvasClick}
      onMouseMove={handleCanvasMouseMove}
      onMouseUp={handleCanvasMouseUp}
      onMouseLeave={handleCanvasMouseUp}
    >
      <svg className="canvas-edges">
        {edges.map((edge) => (
          <line
            key={edge.id}
            x1={edge.x1}
            y1={edge.y1}
            x2={edge.x2}
            y2={edge.y2}
            stroke={highlightedEdges.has(edge.id) ? 'var(--accent)' : 'var(--border-default)'}
            strokeWidth={highlightedEdges.has(edge.id) ? 2.5 : 2}
          />
        ))}
      </svg>

      {nodes.map((node) => {
        const depth = getNodeDepth(nodes, node.id)
        return (
          <Graph
            key={node.id}
            id={node.id}
            label={node.label}
            x={node.x}
            y={node.y}
            size={getNodeSize(depth)}
            isRoot={depth === 0}
            isSelected={node.id === selectedId}
            dragging={node.id === draggingId}
            onMouseDown={handleNodeMouseDown}
          />
        )
      })}

      {selectedNode && (
        <FolderWindows
          node={selectedNode}
          nodeSize={getNodeSize(getNodeDepth(nodes, selectedNode.id))}
          slots={selectedSlots}
          chats={chats}
          onOpenPicker={handleOpenPicker}
          onSend={handleSendMessage}
          onRemove={(slotIndex) => handleRemoveChat(selectedNode.id, slotIndex)}
        />
      )}

      <SettingsButton />

      <div className="canvas-hint">
        {selectedId ? 'Перетащите узел, чтобы переместить папку' : 'Нажмите на папку — рядом появятся чаты'}
      </div>

      {picker && (
        <SlotPicker
          x={picker.x}
          y={picker.y}
          contacts={pickerContacts}
          onSelect={handleAssignChat}
          onClose={() => setPicker(null)}
        />
      )}
    </div>
  )
}
