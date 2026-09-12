import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'

const API_URL = import.meta.env.VITE_PYTHON_API || 'http://localhost:5080'

const BATCH_SIZE = 5
const BATCH_INTERVAL = 2 * 60 * 1000  // 2 минуты

const queue = ref([])
let flushTimer = null
let initialized = false

function scheduleFlush() {
  if (flushTimer) return
  flushTimer = setTimeout(() => {
    flush()
    flushTimer = null
  }, BATCH_INTERVAL)
}

export function useTracking() {
  const auth = useAuthStore()

  function track(type, target = {}) {
    if (!auth.profile?.id) return

    const event = {
      user_id: auth.profile.id,
      type,
      target_id: target.id || null,
      tags: target.tags || [],
      author_id: target.author_id || null,
      timestamp: Math.floor(Date.now() / 1000)
    }

    queue.value.push(event)

    if (queue.value.length >= BATCH_SIZE) {
      flush()
    } else {
      scheduleFlush()
    }
  }

  async function flush() {
    if (queue.value.length === 0) return

    const events = [...queue.value]
    queue.value = []

    try {
      await fetch(`${API_URL}/events`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ events }),
        keepalive: true
      })
    } catch (e) {
      console.warn('Failed to send events:', e)
      // не возвращаем в очередь, чтобы не копить мусор
    }
  }

  function flushSync() {
    if (queue.value.length === 0) return

    const events = [...queue.value]
    queue.value = []

    try {
      navigator.sendBeacon(
        `${API_URL}/events`,
        new Blob([JSON.stringify({ events })], { type: 'application/json' })
      )
    } catch (e) {
      console.warn('Beacon failed:', e)
    }
  }

  function init() {
    if (initialized) return
    initialized = true

    window.addEventListener('beforeunload', flushSync)

    document.addEventListener('visibilitychange', () => {
      if (document.visibilityState === 'hidden') {
        flushSync()
      }
    })
  }

  return { track, flush, flushSync, init }
}