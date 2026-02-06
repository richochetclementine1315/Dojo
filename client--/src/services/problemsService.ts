import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

interface Problem {
  id: string;
  title: string;
  difficulty: 'easy' | 'medium' | 'hard';
  platform: 'leetcode' | 'codeforces' | 'codechef' | 'gfg';
  platform_problem_id: string;
  slug?: string;
  tags: string[];
  acceptance_rate?: number;
  problem_url: string;
  description?: string;
  constraints?: string;
  examples?: any;
  hints?: any;
  created_at?: string;
}

interface ProblemFilters {
  platform?: 'leetcode' | 'codeforces' | 'codechef' | 'gfg';
  difficulty?: 'easy' | 'medium' | 'hard';
  tags?: string[];
  search?: string;
  page?: number;
  limit?: number;
}

interface ProblemsResponse {
  problems: Problem[];
  total: number;
  page: number;
  limit: number;
  hasMore: boolean;
}

class ProblemsService {
  private getAuthToken(): string | null {
    return localStorage.getItem('access_token');
  }

  private getHeaders() {
    const token = this.getAuthToken();
    return {
      'Content-Type': 'application/json',
      ...(token && { Authorization: `Bearer ${token}` }),
    };
  }

  /**
   * Fetch problems from backend API with filters and pagination
   */
  async getProblems(filters: ProblemFilters = {}): Promise<ProblemsResponse> {
    try {
      const params = new URLSearchParams();
      
      if (filters.platform) params.append('platform', filters.platform);
      if (filters.difficulty) params.append('difficulty', filters.difficulty);
      if (filters.search) params.append('search', filters.search);
      if (filters.page) params.append('page', filters.page.toString());
      if (filters.limit) params.append('limit', filters.limit.toString());
      if (filters.tags && filters.tags.length > 0) {
        filters.tags.forEach(tag => params.append('tags', tag));
      }

      const response = await axios.get(`${API_URL}/problems?${params.toString()}`, {
        headers: this.getHeaders(),
      });

      const { problems = [], total = 0, page = 1, limit = 20 } = response.data.data;

      return {
        problems,
        total,
        page,
        limit,
        hasMore: page * limit < total,
      };
    } catch (error: any) {
      // Check if it's an authentication error
      if (error.response?.status === 401) {
        throw new Error('AUTHENTICATION_REQUIRED');
      }
      
      return {
        problems: [],
        total: 0,
        page: 1,
        limit: 20,
        hasMore: false,
      };
    }
  }

  /**
   * Get a single problem by ID
   */
  async getProblemById(id: string): Promise<Problem | null> {
    try {
      const response = await axios.get(`${API_URL}/problems/${id}`, {
        headers: this.getHeaders(),
      });
      return response.data.data;
    } catch (error) {
      return null;
    }
  }
}

export const problemsService = new ProblemsService();
export type { Problem, ProblemFilters, ProblemsResponse };
