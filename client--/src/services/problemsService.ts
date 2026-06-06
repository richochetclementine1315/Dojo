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
  cf_rating?: number;
  solved_count?: number;
  problem_url: string;
  description?: string;
  constraints?: string;
  examples?: any;
  hints?: any;
  is_solved?: boolean;
  created_at?: string;
}

// A Codeforces problem returned by the live browse endpoint (not stored in DB)
interface CFProblem {
  contest_id: number;
  index: string;
  name: string;
  type: string;
  rating: number;
  tags: string[];
  solved_count: number;
  problem_url: string;
}

interface ProblemFilters {
  platform?: 'leetcode' | 'codeforces' | 'codechef' | 'gfg';
  difficulty?: 'easy' | 'medium' | 'hard';
  tags?: string[];
  search?: string;
  page?: number;
  limit?: number;
}

interface CFProblemsFilters {
  tags?: string[];
  min_rating?: number;
  max_rating?: number;
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

interface CFProblemsResponse {
  problems: CFProblem[];
  total: number;
  page: number;
  limit: number;
  has_more: boolean;
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
      if (error.response?.status === 401) {
        throw new Error('AUTHENTICATION_REQUIRED');
      }
      return { problems: [], total: 0, page: 1, limit: 20, hasMore: false };
    }
  }

  /**
   * Fetch Codeforces problems live from the CF API (cached 6h on backend).
   * Supports tag, rating range, and text search filters.
   */
  async getCodeforcesProblems(filters: CFProblemsFilters = {}): Promise<CFProblemsResponse> {
    const params = new URLSearchParams();

    if (filters.search) params.append('search', filters.search);
    if (filters.min_rating) params.append('min_rating', filters.min_rating.toString());
    if (filters.max_rating) params.append('max_rating', filters.max_rating.toString());
    if (filters.page) params.append('page', filters.page.toString());
    if (filters.limit) params.append('limit', filters.limit.toString());
    if (filters.tags && filters.tags.length > 0) {
      filters.tags.forEach(tag => params.append('tags', tag));
    }

    const response = await axios.get(
      `${API_URL}/problems/codeforces/browse?${params.toString()}`,
      { headers: this.getHeaders() }
    );

    return response.data.data as CFProblemsResponse;
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

  /**
   * Sync problems from external platforms (imports into DB)
   */
  async syncProblems(platform: 'leetcode' | 'codeforces', limit: number = 100): Promise<{imported: number}> {
    const response = await axios.post(
      `${API_URL}/problems/sync`,
      { platform, limit },
      { headers: this.getHeaders() }
    );
    return response.data.data;
  }

  /**
   * Mark a problem as solved or unsolved
   */
  async markProblemSolved(problemId: string, isSolved: boolean): Promise<void> {
    await axios.post(
      `${API_URL}/problems/${problemId}/solve`,
      { is_solved: isSolved },
      { headers: this.getHeaders() }
    );
  }

  /**
   * Get the count of problems solved by the current user
   */
  async getSolvedCount(): Promise<number> {
    try {
      const response = await axios.get(`${API_URL}/problems/solved/count`, {
        headers: this.getHeaders(),
      });
      return response.data.data.count || 0;
    } catch (error) {
      console.error('Failed to fetch solved count:', error);
      return 0;
    }
  }
}

export const problemsService = new ProblemsService();
export type { Problem, CFProblem, ProblemFilters, CFProblemsFilters, ProblemsResponse, CFProblemsResponse };
