import "./Graph.css";

export default function Graph({ id, x = 0, y = 0 }) {
  const style = {
    position: 'absolute',
    left: `${x}px`,
    top: `${y}px`,
  }

  return (
    <div className="graph" style={style}>
      <span className="graph-label">Graph {id}</span>
    </div>
  )
}