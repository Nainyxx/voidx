import Graph from '../../components/Graph/Graph'
import { getGraphCoordinates } from './graphConfig'

export function renderGraphs() {
  return getGraphCoordinates().map(({ id, x, y }) => (
    <Graph key={id} id={id} x={x} y={y} />
  ))
}
