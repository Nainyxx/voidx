import { useState } from 'react'
import "./SlotPicker.css"

export default function SlotPicker({ x, y, contacts, onSelect, onClose }) {
  const [query, setQuery] = useState('')

  const filtered = contacts.filter((c) =>
    c.name.toLowerCase().includes(query.trim().toLowerCase())
  )

  return (
    <div className="slot-picker-backdrop" onMouseDown={onClose}>
      <div
        className="slot-picker"
        style={{ left: `${x}px`, top: `${y}px` }}
        onMouseDown={(e) => e.stopPropagation()}
      >
        <div className="slot-picker-header">
          <div className="slot-picker-title">Выбрать чат</div>
          <button type="button" className="slot-picker-close" onClick={onClose} aria-label="Закрыть">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
              <line x1="18" y1="6" x2="6" y2="18" />
              <line x1="6" y1="6" x2="18" y2="18" />
            </svg>
          </button>
        </div>

        <div className="slot-picker-search">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="var(--text-tertiary)" strokeWidth="2" strokeLinecap="round">
            <circle cx="11" cy="11" r="7" />
            <line x1="21" y1="21" x2="16.65" y2="16.65" />
          </svg>
          <input
            autoFocus
            placeholder="Поиск чата или контакта"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />
        </div>

        <div className="slot-picker-list">
          {filtered.map((c) => (
            <button
              key={c.id}
              type="button"
              className="slot-picker-row"
              onClick={() => onSelect(c.id)}
            >
              <div className="slot-picker-avatar">{c.initial}</div>
              <div className="slot-picker-row-text">
                <div className="slot-picker-row-name">{c.name}</div>
                <div className="slot-picker-row-preview">{c.messages[c.messages.length - 1]?.text}</div>
              </div>
            </button>
          ))}
          {filtered.length === 0 && (
            <div className="slot-picker-empty">Ничего не найдено</div>
          )}
        </div>
      </div>
    </div>
  )
}
