import { useState } from "react"
import {
  Container,
  Typography,
  Paper,
  TextField,
  Button,
  Box,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Chip,
} from "@mui/material"
import { Send, Psychology } from "@mui/icons-material"
import { queryAI } from "@/shared/api/ai"

export default function AIQueryPage() {
  const [input, setInput] = useState("")
  const [output, setOutput] = useState("")
  const [role, setRole] = useState("docs")
  const [isLoading, setIsLoading] = useState(false)

  async function handleSend() {
    setIsLoading(true)
    try {
      const res = await queryAI({
        provider_id: "demo",
        role: role,
        input: input,
        context: {}
      })
      setOutput(res.output)
    } catch (error) {
      setOutput("Ошибка: " + (error as any).message)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <Container maxWidth="lg">
      <Box display="flex" alignItems="center" mb={3}>
        <Psychology sx={{ mr: 1 }} />
        <Typography variant="h4">AI Запросы</Typography>
      </Box>

      <Paper sx={{ p: 3 }}>
        <Box mb={3}>
          <FormControl fullWidth sx={{ mb: 2 }}>
            <InputLabel>Роль</InputLabel>
            <Select
              value={role}
              label="Роль"
              onChange={(e) => setRole(e.target.value)}
            >
              <MenuItem value="docs">Анализ документов</MenuItem>
              <MenuItem value="risks">Анализ рисков</MenuItem>
              <MenuItem value="incidents">Анализ инцидентов</MenuItem>
              <MenuItem value="compliance">Соответствие стандартам</MenuItem>
            </Select>
          </FormControl>

          <TextField
            fullWidth
            multiline
            rows={4}
            label="Введите ваш запрос"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="Например: Проанализируй этот документ на соответствие требованиям ИСО 27001..."
            sx={{ mb: 2 }}
          />

          <Button
            variant="contained"
            startIcon={<Send />}
            onClick={handleSend}
            disabled={!input.trim() || isLoading}
            fullWidth
          >
            {isLoading ? "Отправка..." : "Отправить запрос"}
          </Button>
        </Box>

        {output && (
          <Box>
            <Typography variant="h6" gutterBottom>
              Ответ AI:
            </Typography>
            <Paper
              variant="outlined"
              sx={{
                p: 2,
                backgroundColor: "#f5f5f5",
                whiteSpace: "pre-wrap",
                fontFamily: "monospace",
              }}
            >
              {output}
            </Paper>
          </Box>
        )}

        <Box mt={3}>
          <Typography variant="body2" color="text.secondary">
            <strong>Доступные роли:</strong>
          </Typography>
          <Box display="flex" gap={1} mt={1} flexWrap="wrap">
            <Chip label="docs - Анализ документов" size="small" />
            <Chip label="risks - Анализ рисков" size="small" />
            <Chip label="incidents - Анализ инцидентов" size="small" />
            <Chip label="compliance - Соответствие стандартам" size="small" />
          </Box>
        </Box>
      </Paper>
    </Container>
  )
}
