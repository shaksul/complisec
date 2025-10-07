import { apiClient } from './client';

export interface RequestEmailChangeRequest {
  new_email: string;
}

export interface RequestEmailChangeResponse {
  request_id: string;
  message: string;
}

export interface VerifyEmailRequest {
  request_id: string;
  verification_code: string;
}

export interface VerifyEmailResponse {
  message: string;
  status: string;
}

export interface CompleteEmailChangeRequest {
  request_id: string;
}

export interface CompleteEmailChangeResponse {
  message: string;
}

export interface CancelEmailChangeRequest {
  request_id: string;
}

export interface CancelEmailChangeResponse {
  message: string;
}

export interface ResendVerificationCodeRequest {
  request_id: string;
}

export interface ResendVerificationCodeResponse {
  message: string;
}

export interface EmailChangeRequestResponse {
  id: string;
  old_email: string;
  new_email: string;
  status: string;
  expires_at: string;
  created_at: string;
  updated_at: string;
}

export interface EmailChangeStatusResponse {
  has_active_request: boolean;
  request?: EmailChangeRequestResponse;
}

export interface EmailChangeAuditLogResponse {
  id: string;
  old_email: string;
  new_email: string;
  change_type: string;
  ip_address?: string;
  user_agent?: string;
  created_at: string;
}

export interface EmailChangeAuditLogsResponse {
  logs: EmailChangeAuditLogResponse[];
  total: number;
}

export const emailChangeApi = {
  // Создать запрос на смену email
  requestEmailChange: async (data: RequestEmailChangeRequest): Promise<RequestEmailChangeResponse> => {
    const response = await apiClient.post('/email-change/request', data);
    return response.data.data;
  },

  // Подтвердить старый email
  verifyOldEmail: async (data: VerifyEmailRequest): Promise<VerifyEmailResponse> => {
    const response = await apiClient.post('/email-change/verify-old', data);
    return response.data.data;
  },

  // Подтвердить новый email
  verifyNewEmail: async (data: VerifyEmailRequest): Promise<VerifyEmailResponse> => {
    const response = await apiClient.post('/email-change/verify-new', data);
    return response.data.data;
  },

  // Завершить смену email
  completeEmailChange: async (data: CompleteEmailChangeRequest): Promise<CompleteEmailChangeResponse> => {
    const response = await apiClient.post('/email-change/complete', data);
    return response.data.data;
  },

  // Отменить запрос на смену email
  cancelEmailChange: async (data: CancelEmailChangeRequest): Promise<CancelEmailChangeResponse> => {
    const response = await apiClient.post('/email-change/cancel', data);
    return response.data.data;
  },

  // Повторно отправить код подтверждения
  resendVerificationCode: async (data: ResendVerificationCodeRequest): Promise<ResendVerificationCodeResponse> => {
    const response = await apiClient.post('/email-change/resend', data);
    return response.data.data;
  },

  // Получить статус активного запроса
  getEmailChangeStatus: async (): Promise<EmailChangeStatusResponse> => {
    const response = await apiClient.get('/email-change/status');
    return response.data.data;
  },

  // Получить аудит-лог изменений email
  getAuditLogs: async (limit = 20, offset = 0): Promise<EmailChangeAuditLogsResponse> => {
    const response = await apiClient.get(`/email-change/audit-logs?limit=${limit}&offset=${offset}`);
    return response.data.data;
  },
};

