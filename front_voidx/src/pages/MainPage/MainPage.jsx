import { getGraphCount } from './graphConfig'
import { renderGraphs } from './graphRenderer'

export default function MainPage() {
  const graphCount = getGraphCount()

  return (
    <main style={{ padding: '24px' }}>
      <h1>Main Page</h1>
      <p>Всего графов: {graphCount}</p>
      <div
        style={{
          position: 'relative',
          width: '100%',
          minHeight: '600px',
          border: '1px solid #ddd',
          borderRadius: '12px',
          background: '#fafafa',
          overflow: 'hidden',
        }}
      >
        {renderGraphs()}
      </div>
    </main>
  )
}
