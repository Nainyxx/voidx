import graphConfig from '../../data/config/config.json'

export function getGraphCount() {
  return graphConfig.graph_counts ?? Object.keys(graphConfig.graph_cords ?? {}).length
}

export function getGraphCoordinates() {
  const cords = graphConfig.graph_cords ?? {}
  return Object.entries(cords).map(([id, { x, y }]) => ({ id, x, y }))
}
