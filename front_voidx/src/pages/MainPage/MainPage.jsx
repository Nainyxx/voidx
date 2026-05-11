import { useState, useEffect } from 'react'
import Graph from '../../components/Graph/Graph'
import { getGraphCoordinates } from './graphConfig'

export default function MainPage() {
  const [graphs, setGraphs] = useState([])
  const [dragging, setDragging] = useState(null)
  const [dragOffset, setDragOffset] = useState({ dx: 0, dy: 0 })

  // Initialize graphs from config
  useEffect(() => {
    const initialGraphs = getGraphCoordinates()
    setGraphs(initialGraphs)
  }, [])

  // Save graphs to localStorage (simulating config.json update)
  const saveGraphs = (updatedGraphs) => {
    setGraphs(updatedGraphs)
    localStorage.setItem('graphs', JSON.stringify(updatedGraphs))
    console.log('Graphs saved:', updatedGraphs)
  }

  // Handle graph mouse down
  const handleGraphMouseDown = (id, e) => {
    e.stopPropagation()
    const graph = graphs.find(g => g.id === id)
    if (graph) {
      setDragging(id)
      setDragOffset({
        dx: e.clientX - graph.x,
        dy: e.clientY - graph.y,
      })
    }
  }

  // Handle canvas mouse move
  const handleCanvasMouseMove = (e) => {
    if (dragging) {
      const updatedGraphs = graphs.map(g =>
        g.id === dragging
          ? {
              ...g,
              x: e.clientX - dragOffset.dx,
              y: e.clientY - dragOffset.dy,
            }
          : g
      )
      setGraphs(updatedGraphs)
    }
  }

  // Handle canvas mouse up
  const handleCanvasMouseUp = () => {
    if (dragging) {
      setDragging(null)
      const finalGraphs = graphs.map(g =>
        g.id === dragging
          ? { ...g, x: Math.round(g.x), y: Math.round(g.y) }
          : g
      )
      saveGraphs(finalGraphs)
    }
  }

  // Handle canvas click to create new graph
  const handleCanvasClick = (e) => {
    if (e.target === e.currentTarget && !dragging) {
      const rect = e.currentTarget.getBoundingClientRect()
      const x = e.clientX - rect.left - 60
      const y = e.clientY - rect.top - 60

      const newId = Math.max(...graphs.map(g => parseInt(g.id) || 0), 0) + 1
      const newGraph = {
        id: String(newId),
        x: Math.round(x),
        y: Math.round(y),
      }

      saveGraphs([...graphs, newGraph])
    }
  }

  return (
    <div
      style={{
        position: 'relative',
        width: '100%',
        height: '100vh',
        background: '#ffffff',
        overflow: 'hidden',
        cursor: dragging ? 'grabbing' : 'crosshair',
      }}
      onClick={handleCanvasClick}
      onMouseMove={handleCanvasMouseMove}
      onMouseUp={handleCanvasMouseUp}
      onMouseLeave={handleCanvasMouseUp}
    >
      {graphs.map((graph) => (
        <Graph
          key={graph.id}
          id={graph.id}
          x={graph.x}
          y={graph.y}
          onMouseDown={(e) => handleGraphMouseDown(graph.id, e)}
        />
      ))}
    </div>
  )
}
