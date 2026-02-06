export interface User {
  id: string;
  email: string;
  username: string;
  avatar_url?: string;
  created_at: string;
  leetcode_username?: string;
  codeforces_username?: string;
  codechef_username?: string;
  gfg_username?: string;
  profile?: {
    bio?: string;
    location?: string;
    website?: string;
    leetcode_username?: string;
    codeforces_username?: string;
    codechef_username?: string;
    gfg_username?: string;
    total_solved?: number;
    easy_solved?: number;
    medium_solved?: number;
    hard_solved?: number;
    platform_stats?: Array<{
      platform: string;
      rating: number;
      max_rating: number;
      solved_count: number;
      global_rank: number;
      last_synced_at?: string;
    }>;
  };
}

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface RegisterData {
  email: string;
  username: string;
  password: string;
  leetcode_username?: string;
  codeforces_username?: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  token_type: string;
  expires_in: number;
}

export interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
}

export interface Problem {
  id: string;
  title: string;
  description: string;
  difficulty: 'easy' | 'medium' | 'hard';
  platform: 'leetcode' | 'codeforces' | 'custom';
  problem_url?: string;
  tags: string[];
  acceptance_rate?: number;
  total_submissions?: number;
  created_at: string;
  updated_at: string;
}

export interface ProblemFilters {
  difficulty?: 'easy' | 'medium' | 'hard';
  platform?: 'leetcode' | 'codeforces' | 'custom';
  tags?: string[];
  search?: string;
}

export interface LeetCodeStats {
  username: string;
  ranking: number;
  total_solved: number;
  easy_solved: number;
  medium_solved: number;
  hard_solved: number;
  acceptance_rate: number;
  contribution_points: number;
  reputation: number;
}

export interface CodeforcesStats {
  username: string;
  rating: number;
  max_rating: number;
  rank: string;
  max_rank: string;
  contribution: number;
  friend_of_count: number;
  contests_participated: number;
}

export interface CodeChefStats {
  username: string;
  rating: number;
  global_rank: number;
  country_rank: number;
  stars: string;
  problems_solved: number;
  contests_participated: number;
}

export interface GeeksForGeeksStats {
  username: string;
  institute_rank: number;
  current_streak: number;
  max_streak: number;
  total_problems_solved: number;
  monthly_score: number;
  overall_score: number;
}

export interface UserProfileStats {
  leetcode?: LeetCodeStats;
  codeforces?: CodeforcesStats;
  codechef?: CodeChefStats;
  geeksforgeeks?: GeeksForGeeksStats;
}

export interface Contest {
  id: string;
  name: string;
  platform: 'leetcode' | 'codeforces';
  start_time: string;
  duration: number;
  duration_seconds: number;
  url: string;
  contest_url: string;
  phase?: string;
  is_registered?: boolean;
  created_at: string;
}

export interface Sheet {
  id: string;
  name: string;
  description: string;
  user_id: string;
  is_public: boolean;
  problems: Problem[];
  created_at: string;
  updated_at: string;
}

export interface Room {
  id: string;
  name: string;
  room_code: string;
  description?: string;
  max_participants: number;
  current_participants: number;
  created_by: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}
