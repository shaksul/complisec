import { api } from './client'

export interface ChatMessage {
  role: 'user' | 'assistant' | 'system'
  content: string
}

export interface SendChatMessageRequest {
  provider_id: string
  messages: ChatMessage[]
  model?: string
  stream?: boolean
}

export interface SendChatMessageResponse {
  message: ChatMessage
  model: string
}

export const aiChatApi = {
  /**
   * Отправить сообщение в AI чат
   */
  async sendMessage(
    providerId: string,
    messages: ChatMessage[],
    model?: string
  ): Promise<SendChatMessageResponse> {
    const response = await api.post('/ai/chat', {
      provider_id: providerId,
      messages,
      model: model || 'llama3.2',
      stream: false,
    })
    return response.data.data
  },
}

