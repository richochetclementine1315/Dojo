import api from '@/lib/api';
import type { UserProfileStats, ApiResponse } from '@/types';

export const profileService = {
  async getProfileStats() {
    // Profile stats are included in the profile endpoint
    const response = await api.get<ApiResponse<any>>('/users/profile');
    const platformStats = response.data.data?.profile?.platform_stats || [];
    
    // Transform array of platform stats to object keyed by platform
    const stats: UserProfileStats = {
      leetcode: undefined,
      codeforces: undefined,
      codechef: undefined,
      geeksforgeeks: undefined,
    };
    
    platformStats.forEach((stat: any) => {
      const platformKey = stat.platform === 'gfg' ? 'geeksforgeeks' : stat.platform;
      if (platformKey === 'leetcode') {
        stats.leetcode = {
          username: response.data.data?.profile?.leetcode_username || '',
          ranking: stat.global_rank || 0,
          total_solved: stat.solved_count || 0,
          easy_solved: stat.easy_solved || 0,
          medium_solved: stat.medium_solved || 0,
          hard_solved: stat.hard_solved || 0,
          acceptance_rate: 0,
          contribution_points: 0,
          reputation: 0,
        };
      } else if (platformKey === 'codeforces') {
        stats.codeforces = {
          username: response.data.data?.profile?.codeforces_username || '',
          rating: stat.rating || 0,
          max_rating: stat.max_rating || 0,
          rank: '',
          max_rank: '',
          contribution: 0,
          friend_of_count: 0,
          contests_participated: stat.contests_attended || 0,
        };
      } else if (platformKey === 'codechef') {
        stats.codechef = {
          username: response.data.data?.profile?.codechef_username || '',
          rating: stat.rating || 0,
          global_rank: stat.global_rank || 0,
          country_rank: 0,
          stars: '',
          problems_solved: stat.solved_count || 0,
          contests_participated: stat.contests_attended || 0,
        };
      } else if (platformKey === 'geeksforgeeks') {
        stats.geeksforgeeks = {
          username: response.data.data?.profile?.gfg_username || '',
          institute_rank: 0,
          current_streak: 0,
          max_streak: 0,
          total_problems_solved: stat.solved_count || 0,
          monthly_score: 0,
          overall_score: stat.global_rank || 0,
        };
      }
    });
    
    return stats;
  },

  async updateProfile(data: {
    bio?: string;
    location?: string;
    website?: string;
    leetcode_username?: string;
    codeforces_username?: string;
    codechef_username?: string;
    gfg_username?: string;
  }) {
    const response = await api.put<ApiResponse<any>>('/users/profile', data);
    return response.data.data;
  },

  async syncPlatformStats(platforms: string[]) {
    const response = await api.post<ApiResponse<any>>('/users/sync-stats', {
      platforms,
    });
    
    return response.data.data;
  },

  async syncLeetCode(username: string) {
    const response = await api.post<ApiResponse<any>>('/users/sync/leetcode', { username });
    return response.data.data;
  },

  async syncCodeforces(username: string) {
    const response = await api.post<ApiResponse<any>>('/users/sync/codeforces', { username });
    return response.data.data;
  },

  async syncCodeChef(username: string) {
    const response = await api.post<ApiResponse<any>>('/users/sync/codechef', { username });
    return response.data.data;
  },

  async syncGeeksForGeeks(username: string) {
    const response = await api.post<ApiResponse<any>>('/users/sync/geeksforgeeks', { username });
    return response.data.data;
  },
};
