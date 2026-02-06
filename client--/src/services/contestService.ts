import api from '@/lib/api';
import type { Contest, ApiResponse } from '@/types';

export const contestService = {
  async getContests(params?: {
    platform?: string;
    upcoming?: boolean;
    limit?: number;
  }) {
    const queryParams = new URLSearchParams();
    
    if (params?.platform) queryParams.append('platform', params.platform);
    if (params?.upcoming !== undefined) queryParams.append('upcoming', params.upcoming.toString());
    if (params?.limit) queryParams.append('limit', params.limit.toString());

    const response = await api.get<ApiResponse<{ contests: Contest[]; total: number; page: number; limit: number }>>(
      `/contests?${queryParams.toString()}`
    );
    return response.data.data.contests;
  },

  async getContest(id: string) {
    const response = await api.get<ApiResponse<Contest>>(`/contests/${id}`);
    return response.data.data;
  },

  async syncContests() {
    const response = await api.post<ApiResponse<any>>('/contests/sync');
    return response.data.data;
  },
};
