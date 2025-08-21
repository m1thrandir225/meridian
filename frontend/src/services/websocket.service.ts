import config from '@/lib/config'
import { useAuthStore } from '@/stores/auth'
import type { SendMessagePayload, WebSocketMessage } from '@/types/websocket'

class WebsocketService {
  private ws: WebSocket | null = null
  private reconnectAttempts = 0
  private reconnectDelay = 1000
  private maxReconnectAttempts = 5
  private pingInterval: ReturnType<typeof setInterval> | null = null
  private isConnecting = false
  private messageHandlers: Map<string, ((payload: unknown) => void)[]> = new Map()

  connect() {
    if (this.ws?.readyState === WebSocket.OPEN || this.isConnecting) {
      return
    }

    this.isConnecting = true
    const authStore = useAuthStore()

    if (!authStore.accessToken || !authStore.user) {
      console.error('Cannot connect: no access token')
      this.isConnecting = false
      return
    }

    const wsUrl = `${config.baseUrl.replace('http', 'ws')}/api/v1/messages/ws?token=${encodeURIComponent(authStore.accessToken)}`
    this.ws = new WebSocket(wsUrl)

    this.ws.onopen = () => {
      console.log('WebSocket connected')
      this.isConnecting = false
      this.reconnectAttempts = 0
      this.startPingInterval()
      this.emit('connected', {})
    }

    this.ws.onmessage = (event) => {
      try {
        const message: WebSocketMessage = JSON.parse(event.data)
        this.handleMessage(message)
      } catch (error) {
        console.error('Error parsing WebSocket message:', error)
      }
    }
    this.ws.onclose = (event) => {
      console.log('WebSocket closed:', event)
      this.isConnecting = false
      this.stopPingInterval()
      this.emit('disconnected', { code: event.code, reason: event.reason })
      if (event.code !== 1000 && this.reconnectAttempts < this.maxReconnectAttempts) {
        this.scheduleReconnect()
      }
    }

    this.ws.onerror = (event) => {
      console.error('WebSocket error:', event)
      this.isConnecting = false
    }
  }
  private scheduleReconnect() {
    this.reconnectAttempts++
    const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts)
    setTimeout(() => {
      console.log(`Reconnecting in ${delay}ms...`)
      this.connect()
    }, delay)
  }

  private startPingInterval() {
    this.pingInterval = setInterval(() => {
      this.send({ type: 'ping', payload: {} })
    }, 30000)
  }

  private stopPingInterval() {
    if (this.pingInterval) {
      clearInterval(this.pingInterval)
      this.pingInterval = null
    }
  }
  private handleMessage(message: WebSocketMessage) {
    const handlers = this.messageHandlers.get(message.type)
    if (handlers) {
      handlers.forEach((handler) => handler(message.payload))
    }
  }
  send(message: WebSocketMessage) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message))
    } else {
      console.error('WebSocket is not connected')
    }
  }
  sendMessage(payload: SendMessagePayload) {
    this.send({
      type: 'message',
      payload,
    })
  }

  sendTypingStart(channelId: string) {
    this.send({
      type: 'typing_start',
      payload: { channel_id: channelId },
    })
  }

  sendTypingStop(channelId: string) {
    this.send({
      type: 'typing_stop',
      payload: { channel_id: channelId },
    })
  }

  sendAddReaction(channelId: string, messageId: string, reactionType: string) {
    this.send({
      type: 'add_reaction',
      payload: {
        channel_id: channelId,
        message_id: messageId,
        reaction_type: reactionType,
      },
    })
  }

  sendRemoveReaction(channelId: string, messageId: string, reactionType: string) {
    this.send({
      type: 'remove_reaction',
      payload: {
        channel_id: channelId,
        message_id: messageId,
        reaction_type: reactionType,
      },
    })
  }

  on(event: string, handler: (payload: unknown) => void) {
    if (!this.messageHandlers.has(event)) {
      this.messageHandlers.set(event, [])
    }
    this.messageHandlers.get(event)!.push(handler)
  }

  off(event: string, handler: (payload: unknown) => void) {
    const handlers = this.messageHandlers.get(event)
    if (handlers) {
      const index = handlers.indexOf(handler)
      if (index > -1) {
        handlers.splice(index, 1)
      }
    }
  }

  private emit(event: string, payload: unknown) {
    const handlers = this.messageHandlers.get(event)
    if (handlers) {
      handlers.forEach((handler) => handler(payload))
    }
  }

  disconnect() {
    this.stopPingInterval()
    if (this.ws) {
      this.ws.close(1000, 'User initiated disconnect')
      this.ws = null
    }
  }

  isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN
  }
}

const websocketService = new WebsocketService()
export default websocketService
