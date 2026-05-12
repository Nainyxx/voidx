import graphConfig from '../../data/config/config.json'

export function getGraphCount() {
  return graphConfig.graph_counts ?? Object.keys(graphConfig.graph_cords ?? {}).length
}

export function getGraphCoordinates() {
  const savedCords = localStorage.getItem('graph_cords')
  const cords = savedCords ? JSON.parse(savedCords) : (graphConfig.graph_cords ?? {})
  return Object.entries(cords).map(([id, { x, y }]) => ({ id, x, y }))
}

export function saveGraphCoordinates(graph) {
  const savedCords = localStorage.getItem('graph_cords')
  const cords = savedCords ? JSON.parse(savedCords) : (graphConfig.graph_cords ?? {})
  cords[graph.id] = { x: graph.x, y: graph.y }
  localStorage.setItem('graph_cords', JSON.stringify(cords))
  console.log("saved", graph)
}

export function createNewGraph(graphs, x, y) {
  const newId = Math.max(...graphs.map(g => parseInt(g.id) || 0), 0) + 1
  const newGraph = {
    id: String(newId),
    x: Math.round(x),
    y: Math.round(y),
  }
  saveGraphCoordinates(newGraph)
  return newGraph
}

export function updateGraphPosition(graphs, draggingId, x, y) {
  return graphs.map(g =>
    g.id === draggingId
      ? { ...g, x, y }
      : g
  )
}

export function roundGraphPositions(graphs, draggingId) {
  return graphs.map(g =>
    g.id === draggingId
      ? { ...g, x: Math.round(g.x), y: Math.round(g.y) }
      : g
  )
}