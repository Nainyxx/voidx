import "./Graph.css"

function FolderIcon({ color }) {
  return (
    <svg width="45%" height="45%" viewBox="0 0 24 24" fill="none">
      <path
        d="M4 7C4 6.44772 4.44772 6 5 6H9.17157C9.70201 6 10.2107 6.21071 10.5858 6.58579L11.4142 7.41421C11.7893 7.78929 12.298 8 12.8284 8H19C19.5523 8 20 8.44772 20 9V17C20 17.5523 19.5523 18 19 18H5C4.44772 18 4 17.5523 4 17V7Z"
        stroke={color}
        strokeWidth="1.75"
        strokeLinejoin="round"
      />
    </svg>
  )
}

function ChatIcon() {
  return (
    <svg width="32%" height="32%" viewBox="0 0 24 24" fill="none">
      <path
        d="M4 5.5C4 4.67157 4.67157 4 5.5 4H18.5C19.3284 4 20 4.67157 20 5.5V15.5C20 16.3284 19.3284 17 18.5 17H9L5 20.5V17H5.5C4.67157 17 4 16.3284 4 15.5V5.5Z"
        stroke="white"
        strokeWidth="1.75"
        strokeLinejoin="round"
      />
    </svg>
  )
}

export default function Graph({ id, label, x, y, size, isRoot, isSelected, dragging, onMouseDown }) {
  const circleClass = [
    'graph-circle',
    isRoot && 'graph-circle--root',
    isSelected && 'graph-circle--selected',
  ].filter(Boolean).join(' ')

  const labelClass = [
    'graph-label',
    isRoot && 'graph-label--root',
    isSelected && 'graph-label--selected',
  ].filter(Boolean).join(' ')

  return (
    <div
      className={`graph${dragging ? ' graph--dragging' : ''}`}
      style={{ left: `${x}px`, top: `${y}px` }}
      onMouseDown={(e) => onMouseDown(id, e)}
    >
      <div className={circleClass} style={{ width: size, height: size }}>
        {isRoot ? <ChatIcon /> : <FolderIcon color="var(--accent)" />}
      </div>
      <span className={labelClass}>{label}</span>
    </div>
  )
}
