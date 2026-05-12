import { useState, useEffect } from 'react'
import Graph from '../../components/Graph/Graph'
import {
  getGraphCoordinates,
  saveGraphCoordinates,
  createNewGraph,
  updateGraphPosition,
  roundGraphPositions,
} from './graphConfig'

export default function MainPage() {
  const [graphs, setGraphs] = useState([])
  const [dragging, setDragging] = useState(null)
  const [dragOffset, setDragOffset] = useState({ dx: 0, dy: 0 })

  useEffect(() => {
    setGraphs(getGraphCoordinates())
  }, [])

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

  const handleCanvasMouseMove = (e) => {
    if (dragging) {
      const updatedGraphs = updateGraphPosition(
        graphs,
        dragging,
        e.clientX - dragOffset.dx,
        e.clientY - dragOffset.dy
      )
      setGraphs(updatedGraphs)
    }
  }

  const handleCanvasMouseUp = () => {
    if (dragging) {
      const finalGraphs = roundGraphPositions(graphs, dragging)
      setGraphs(finalGraphs)
      const draggingGraph = finalGraphs.find(g => g.id === dragging)
      if (draggingGraph) {
        saveGraphCoordinates(draggingGraph)
      }
      setDragging(null)
    }
  }

  const handleCanvasClick = (e) => {
    if (e.target === e.currentTarget && !dragging) {
      const rect = e.currentTarget.getBoundingClientRect()
      const x = e.clientX - rect.left - 60
      const y = e.clientY - rect.top - 60
      const newGraph = createNewGraph(graphs, x, y)
      setGraphs([...graphs, newGraph])
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
