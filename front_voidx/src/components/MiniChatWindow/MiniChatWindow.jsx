import { useState } from 'react'
import "./MiniChatWindow.css"

export default function MiniChatWindow({ chat, x, y, onSend, onRemove }) {
  const [draft, setDraft] = useState('')

  const handleSend = () => {
    const text = draft.trim()
    if (!text) return
    onSend(text)
    setDraft('')
  }

  return (
    <div className="mini-chat" style={{ left: `${x}px`, top: `${y}px` }}>
      <div className="mini-chat-header">
        <div className="mini-chat-avatar">{chat.initial}</div>
        <div className="mini-chat-name">{chat.name}</div>
        {chat.online && <div className="mini-chat-status" />}
        <button
          type="button"
          className="mini-chat-close"
          onMouseDown={(e) => e.stopPropagation()}
          onClick={onRemove}
          aria-label="Открепить чат"
        >
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
            <line x1="18" y1="6" x2="6" y2="18" />
            <line x1="6" y1="6" x2="18" y2="18" />
          </svg>
        </button>
      </div>

      <div className="mini-chat-messages">
        {chat.messages.map((m, i) => (
          <div key={i} className={`mini-chat-bubble mini-chat-bubble--${m.from === 'me' ? 'out' : 'in'}`}>
            {m.text}
          </div>
        ))}
      </div>

      <div className="mini-chat-input-row" onMouseDown={(e) => e.stopPropagation()}>
        <input
          className="mini-chat-input"
          placeholder="Сообщение…"
          value={draft}
          onChange={(e) => setDraft(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && handleSend()}
        />
        <button type="button" className="mini-chat-send" onClick={handleSend} aria-label="Отправить">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="white" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <line x1="22" y1="2" x2="11" y2="13" />
            <polygon points="22 2 15 22 11 13 2 9 22 2" />
          </svg>
        </button>
      </div>
    </div>
  )
}
