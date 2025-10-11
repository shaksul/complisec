import { useState, useEffect } from 'react'
import {
  Box,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Alert,
  SelectChangeEvent,
  CircularProgress,
} from '@mui/material'
import { Chat } from '@mui/icons-material'
import { PageContainer, PageHeader, SectionCard } from '@/components/common/Page'
import { AIChatBox } from '@/components/ai/AIChatBox'
import { aiChatApi, type ChatMessage } from '@/shared/api/aiChat'
import { getProviders, type AIProvider } from '@/shared/api/ai'

export default function AIChatPage() {
  const [messages, setMessages] = useState<ChatMessage[]>([])
  const [providers, setProviders] = useState<AIProvider[]>([])
  const [selectedProviderId, setSelectedProviderId] = useState('')
  const [selectedModel, setSelectedModel] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [isLoadingProviders, setIsLoadingProviders] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // Загрузка провайдеров при монтировании компонента
  useEffect(() => {
    const loadProviders = async () => {
      try {
        setIsLoadingProviders(true)
        const data = await getProviders()
        const activeProviders = Array.isArray(data) ? data.filter((p: AIProvider) => p.is_active) : []
        setProviders(activeProviders)
        
        // Автоматически выбираем первого провайдера
        if (activeProviders.length > 0) {
          setSelectedProviderId(activeProviders[0].id)
          setSelectedModel(activeProviders[0].default_model || activeProviders[0].models?.[0] || 'llama3.2')
        }
      } catch (err: any) {
        console.error('Ошибка загрузки провайдеров:', err)
        setError('Не удалось загрузить список AI провайдеров')
      } finally {
        setIsLoadingProviders(false)
      }
    }
    
    loadProviders()
  }, [])

  const handleSendMessage = async (content: string) => {
    if (!selectedProviderId) {
      setError('Выберите AI провайдера')
      return
    }

    setError(null)
    
    // Добавляем сообщение пользователя
    const userMessage: ChatMessage = { role: 'user', content }
    const updatedMessages = [...messages, userMessage]
    setMessages(updatedMessages)
    setIsLoading(true)

    try {
      // Отправляем запрос в API с указанным провайдером и моделью
      const response = await aiChatApi.sendMessage(selectedProviderId, updatedMessages, selectedModel)
      
      // Добавляем ответ ассистента
      setMessages([...updatedMessages, response.message])
    } catch (err: any) {
      const errorMessage = err?.response?.data?.error || err?.message || 'Ошибка при отправке сообщения'
      setError(errorMessage)
      console.error('AI Chat error:', err)
    } finally {
      setIsLoading(false)
    }
  }

  const handleProviderChange = (event: SelectChangeEvent) => {
    const providerId = event.target.value
    setSelectedProviderId(providerId)
    
    // Обновляем выбранную модель на дефолтную модель нового провайдера
    const provider = providers.find(p => p.id === providerId)
    if (provider) {
      setSelectedModel(provider.default_model || provider.models?.[0] || 'llama3.2')
    }
  }

  const handleModelChange = (event: SelectChangeEvent) => {
    setSelectedModel(event.target.value)
  }

  const selectedProvider = providers.find(p => p.id === selectedProviderId)

  const handleClearChat = () => {
    setMessages([])
    setError(null)
  }

  if (isLoadingProviders) {
    return (
      <PageContainer>
        <PageHeader
          title="AI Чат"
          subtitle="Общайтесь с AI ассистентом для получения помощи и рекомендаций"
          actions={<Chat color="primary" fontSize="large" />}
        />
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 4 }}>
          <CircularProgress />
        </Box>
      </PageContainer>
    )
  }

  return (
    <PageContainer>
      <PageHeader
        title="AI Чат"
        subtitle="Общайтесь с AI ассистентом для получения помощи и рекомендаций"
        actions={<Chat color="primary" fontSize="large" />}
      />

      {error && (
        <Alert severity="error" onClose={() => setError(null)} sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      {providers.length === 0 && (
        <Alert severity="warning" sx={{ mb: 2 }}>
          Нет доступных AI провайдеров. Добавьте провайдера в разделе "AI-провайдеры".
        </Alert>
      )}

      <SectionCard
        title="Настройки"
        description="Выберите AI провайдера и модель для диалога"
      >
        <Box sx={{ display: 'flex', gap: 2, alignItems: 'center', flexWrap: 'wrap' }}>
          <FormControl sx={{ minWidth: 300 }}>
            <InputLabel>AI Провайдер</InputLabel>
            <Select
              value={selectedProviderId}
              label="AI Провайдер"
              onChange={handleProviderChange}
              disabled={isLoading || providers.length === 0}
            >
              {providers.map((provider) => (
                <MenuItem key={provider.id} value={provider.id}>
                  {provider.name}
                </MenuItem>
              ))}
            </Select>
          </FormControl>

          {selectedProvider && selectedProvider.models && selectedProvider.models.length > 0 && (
            <FormControl sx={{ minWidth: 200 }}>
              <InputLabel>Модель</InputLabel>
              <Select
                value={selectedModel}
                label="Модель"
                onChange={handleModelChange}
                disabled={isLoading}
              >
                {selectedProvider.models.map((model) => (
                  <MenuItem key={model} value={model}>
                    {model}
                    {model === selectedProvider.default_model && ' (по умолчанию)'}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          )}

          {messages.length > 0 && (
            <Box
              component="button"
              onClick={handleClearChat}
              sx={{
                px: 2,
                py: 1,
                borderRadius: 1,
                border: '1px solid',
                borderColor: 'divider',
                bgcolor: 'background.paper',
                cursor: 'pointer',
                '&:hover': {
                  bgcolor: 'action.hover',
                },
              }}
            >
              Очистить чат
            </Box>
          )}
        </Box>
      </SectionCard>

      {selectedProviderId && (
        <SectionCard
          title="Диалог"
          description={`Сообщений: ${messages.length}`}
        >
          <AIChatBox
            messages={messages}
            onSendMessage={handleSendMessage}
            isLoading={isLoading}
          />
        </SectionCard>
      )}
    </PageContainer>
  )
}

