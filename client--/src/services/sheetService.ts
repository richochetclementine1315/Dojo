import api from '@/lib/api';
import type { Sheet, Problem, ApiResponse } from '@/types';

export const sheetService = {
  async getSheets() {
    const response = await api.get<ApiResponse<Sheet[]>>('/sheets');
    return response.data.data;
  },

  async getSheet(id: string) {
    const response = await api.get<ApiResponse<Sheet>>(`/sheets/${id}`);
    return response.data.data;
  },

  async createSheet(data: {
    name: string;
    description: string;
    is_public: boolean;
  }) {
    const response = await api.post<ApiResponse<Sheet>>('/sheets', data);
    return response.data.data;
  },

  async updateSheet(id: string, data: Partial<Sheet>) {
    const response = await api.put<ApiResponse<Sheet>>(`/sheets/${id}`, data);
    return response.data.data;
  },

  async deleteSheet(id: string) {
    await api.delete(`/sheets/${id}`);
  },

  async addProblemToSheet(sheetId: string, problemId: string) {
    const response = await api.post<ApiResponse<Sheet>>(
      `/sheets/${sheetId}/problems/${problemId}`
    );
    return response.data.data;
  },

  async removeProblemFromSheet(sheetId: string, problemId: string) {
    await api.delete(`/sheets/${sheetId}/problems/${problemId}`);
  },
};
