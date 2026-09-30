export default function EmptySlot({ x, y, onClick }) {
  return (
    <button
      type="button"
      className="mini-chat mini-chat--empty"
      style={{ left: `${x}px`, top: `${y}px` }}
      onMouseDown={(e) => e.stopPropagation()}
      onClick={onClick}
    >
      <div className="mini-chat-empty-icon">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="var(--accent)" strokeWidth="2" strokeLinecap="round">
          <line x1="12" y1="5" x2="12" y2="19" />
          <line x1="5" y1="12" x2="19" y2="12" />
        </svg>
      </div>
      <span className="mini-chat-empty-label">Добавить чат</span>
    </button>
  )
}
