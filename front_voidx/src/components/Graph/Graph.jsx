import "./Graph.css";

export default function Graph({ id, x = 0, y = 0, onMouseDown }) {
  const style = {
    position: 'absolute',
    left: `${x}px`,
    top: `${y}px`,
  };

  return (
    <div
      className="graph"
      style={style}
      onMouseDown={onMouseDown}
    >
      <span className="graph-label">Graph {id}</span>
    </div>
  );
}