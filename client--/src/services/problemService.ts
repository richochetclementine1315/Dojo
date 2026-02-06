import api from '@/lib/api';
import type { Problem, ApiResponse } from '@/types';

export const problemService = {
  async getProblems(params?: {
    difficulty?: string;
    platform?: string;
    tags?: string[];
    search?: string;
    page?: number;
    limit?: number;
  }) {
    const queryParams = new URLSearchParams();
    
    if (params?.difficulty) queryParams.append('difficulty', params.difficulty);
    if (params?.platform) queryParams.append('platform', params.platform);
    if (params?.tags?.length) queryParams.append('tags', params.tags.join(','));
    if (params?.search) queryParams.append('search', params.search);
    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.limit) queryParams.append('limit', params.limit.toString());

    const response = await api.get<ApiResponse<Problem[]>>(
      `/problems?${queryParams.toString()}`
    );
    return response.data.data;
  },

  async getProblem(id: string) {
    const response = await api.get<ApiResponse<Problem>>(`/problems/${id}`);
    return response.data.data;
  },

  async createProblem(data: Omit<Problem, 'id' | 'created_at' | 'updated_at'>) {
    const response = await api.post<ApiResponse<Problem>>('/problems', data);
    return response.data.data;
  },

  async updateProblem(id: string, data: Partial<Problem>) {
    const response = await api.put<ApiResponse<Problem>>(`/problems/${id}`, data);
    return response.data.data;
  },

  async deleteProblem(id: string) {
    const response = await api.delete<ApiResponse<void>>(`/problems/${id}`);
    return response.data;
  },
};
