import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { Navbar } from '@/components/layout/Navbar';
import { Card, CardHeader, CardContent } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { profileService } from '@/services/profileService';
import type { UserProfileStats } from '@/types';
import Antigravity from '@/components/effects/Antigravity';
import { 
  Trophy, 
  Target, 
  TrendingUp, 
  Award,
  Code2,
  Flame,
  Star,
  Users,
  ExternalLink,
  RefreshCw,
  Loader2,
  Settings
} from 'lucide-react';

const platformColors = {
  leetcode: 'from-orange-500 to-yellow-500',
  codeforces: 'from-blue-500 to-cyan-500',
  codechef: 'from-amber-600 to-orange-500',
  geeksforgeeks: 'from-green-500 to-emerald-500',
};

export default function ProfileStats() {
  const [stats, setStats] = useState<UserProfileStats | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSyncing, setIsSyncing] = useState<string | null>(null);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchStats();
  }, []);

  const fetchStats = async () => {
    try {
      setIsLoading(true);
      const data = await profileService.getProfileStats();
      console.log('ProfileStats - Fetched stats:', data);
      setStats(data || {});
      setError('');
    } catch (err: any) {
      console.error('Failed to fetch profile stats:', err);
      setError(err.response?.data?.message || 'Failed to load profile stats');
      // Set empty stats on error
      setStats({});
    } finally {
      setIsLoading(false);
    }
  };

  const handleSync = async (platform: string) => {
    try {
      setIsSyncing(platform);
      // Call sync endpoint based on platform
      await fetchStats();
    } catch (err: any) {
      setError(`Failed to sync ${platform} data`);
    } finally {
      setIsSyncing(null);
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-dojo-black-900">
        <Navbar />
        <div className="flex items-center justify-center h-[calc(100vh-4rem)]">
          <Loader2 className="h-12 w-12 animate-spin text-dojo-red-500" />
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-dojo-black-900 relative">
      {/* Antigravity Background */}
      <div className="fixed inset-0 z-0">
        <Antigravity
          count={300}
          magnetRadius={6}
          ringRadius={7}
          waveSpeed={0.4}
          waveAmplitude={1}
          particleSize={1.5}
          lerpSpeed={0.05}
          color="#d40808"
          autoAnimate
          particleVariance={1}
          rotationSpeed={0}
          depthFactor={1}
          pulseSpeed={3}
          particleShape="capsule"
          fieldStrength={10}
        />
      </div>
      <div className="relative z-10">
        <Navbar />
      
      <div className="container mx-auto px-4 py-8">
        {/* Header */}
        <div className="mb-8 flex items-center justify-between">
          <div>
            <h1 className="text-4xl font-bold text-white mb-2">
              Coding <span className="dojo-text-gradient">Profile Stats</span>
            </h1>
            <p className="text-gray-400">Track your progress across all major competitive programming platforms</p>
          </div>
          <Link to="/settings/platforms">
            <Button variant="outline">
              <Settings className="mr-2 h-4 w-4" />
              Manage Platforms
            </Button>
          </Link>
        </div>

        {error && (
          <div className="mb-6 p-4 bg-red-500/10 border border-red-500/20 rounded-lg text-red-400">
            {error}
          </div>
        )}

        <div className="grid md:grid-cols-2 gap-6">
          {/* LeetCode Card */}
          <Card className="relative overflow-hidden">
            <div className={`absolute top-0 left-0 right-0 h-1 bg-gradient-to-r ${platformColors.leetcode}`} />
            <CardHeader className="flex flex-row items-center justify-between">
              <div className="flex items-center gap-3">
                <div className={`p-2 rounded-lg bg-gradient-to-br ${platformColors.leetcode}`}>
                  <Code2 className="h-6 w-6 text-white" />
                </div>
                <div>
                  <h2 className="text-xl font-bold text-white">LeetCode</h2>
                  {stats?.leetcode && (
                    <p className="text-sm text-gray-400">@{stats.leetcode.username}</p>
                  )}
                </div>
              </div>
              <div className="flex gap-2">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => handleSync('leetcode')}
                  disabled={isSyncing === 'leetcode'}
                >
                  {isSyncing === 'leetcode' ? (
                    <Loader2 className="h-4 w-4 animate-spin" />
                  ) : (
                    <RefreshCw className="h-4 w-4" />
                  )}
                </Button>
                {stats?.leetcode && (
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => window.open(`https://leetcode.com/${stats.leetcode?.username}`, '_blank')}
                  >
                    <ExternalLink className="h-4 w-4" />
                  </Button>
                )}
              </div>
            </CardHeader>
            <CardContent>
              {stats?.leetcode && (stats.leetcode.total_solved > 0 || stats.leetcode.ranking > 0) ? (
                <div className="space-y-4">
                  <div className="grid grid-cols-2 gap-4">
                    <div className="bg-dojo-black-800 p-4 rounded-lg">
                      <div className="flex items-center gap-2 text-gray-400 text-sm mb-1">
                        <Trophy className="h-4 w-4" />
                        Ranking
                      </div>
                      <div className="text-2xl font-bold text-white">
                        #{stats.leetcode.ranking.toLocaleString()}
                      </div>
                    </div>
                    <div className="bg-dojo-black-800 p-4 rounded-lg">
                      <div className="flex items-center gap-2 text-gray-400 text-sm mb-1">
                        <Target className="h-4 w-4" />
                        Total Solved
                      </div>
                      <div className="text-2xl font-bold text-white">
                        {stats.leetcode.total_solved}
                      </div>
                    </div>
                  </div>

                  <div className="space-y-2">
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-400">Easy</span>
                      <span className="text-sm font-medium text-green-400">
                        {stats.leetcode.easy_solved}
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-400">Medium</span>
                      <span className="text-sm font-medium text-yellow-400">
                        {stats.leetcode.medium_solved}
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-400">Hard</span>
                      <span className="text-sm font-medium text-red-400">
                        {stats.leetcode.hard_solved}
                      </span>
                    </div>
                  </div>

                  <div className="pt-3 border-t border-dojo-black-700 grid grid-cols-2 gap-4 text-sm">
                    <div>
                      <span className="text-gray-400">Acceptance Rate</span>
                      <div className="text-white font-medium">{stats.leetcode.acceptance_rate}%</div>
                    </div>
                    <div>
                      <span className="text-gray-400">Reputation</span>
                      <div className="text-white font-medium">{stats.leetcode.reputation}</div>
                    </div>
                  </div>
                </div>
              ) : (
                <div className="text-center py-8">
                  <Code2 className="h-12 w-12 text-gray-600 mx-auto mb-3" />
                  <p className="text-gray-400 mb-3">LeetCode profile not linked</p>
                  <Button size="sm" variant="outline">
                    Link Profile
                  </Button>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Codeforces Card */}
          <Card className="relative overflow-hidden">
            <div className={`absolute top-0 left-0 right-0 h-1 bg-gradient-to-r ${platformColors.codeforces}`} />
            <CardHeader className="flex flex-row items-center justify-between">
              <div className="flex items-center gap-3">
                <div className={`p-2 rounded-lg bg-gradient-to-br ${platformColors.codeforces}`}>
                  <Award className="h-6 w-6 text-white" />
                </div>
                <div>
                  <h2 className="text-xl font-bold text-white">Codeforces</h2>
                  {stats?.codeforces && (
                    <p className="text-sm text-gray-400">@{stats.codeforces.username}</p>
                  )}
                </div>
              </div>
              <div className="flex gap-2">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => handleSync('codeforces')}
                  disabled={isSyncing === 'codeforces'}
                >
                  {isSyncing === 'codeforces' ? (
                    <Loader2 className="h-4 w-4 animate-spin" />
                  ) : (
                    <RefreshCw className="h-4 w-4" />
                  )}
                </Button>
                {stats?.codeforces && (
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => window.open(`https://codeforces.com/profile/${stats.codeforces?.username}`, '_blank')}
                  >
                    <ExternalLink className="h-4 w-4" />
                  </Button>
                )}
              </div>
            </CardHeader>
            <CardContent>
              {stats?.codeforces && (stats.codeforces.rating > 0 || stats.codeforces.max_rating > 0) ? (
                <div className="space-y-4">
                  <div className="grid grid-cols-2 gap-4">
                    <div className="bg-dojo-black-800 p-4 rounded-lg">
                      <div className="flex items-center gap-2 text-gray-400 text-sm mb-1">
                        <TrendingUp className="h-4 w-4" />
                        Current Rating
                      </div>
                      <div className="text-2xl font-bold text-white">
                        {stats.codeforces.rating}
                      </div>
                      <div className="text-xs text-blue-400 mt-1">
                        {stats.codeforces.rank}
                      </div>
                    </div>
                    <div className="bg-dojo-black-800 p-4 rounded-lg">
                      <div className="flex items-center gap-2 text-gray-400 text-sm mb-1">
                        <Trophy className="h-4 w-4" />
                        Max Rating
                      </div>
                      <div className="text-2xl font-bold text-white">
                        {stats.codeforces.max_rating}
                      </div>
                      <div className="text-xs text-purple-400 mt-1">
                        {stats.codeforces.max_rank}
                      </div>
                    </div>
                  </div>

                  <div className="pt-3 border-t border-dojo-black-700 space-y-2">
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-400">Contests</span>
                      <span className="text-sm font-medium text-white">
                        {stats.codeforces.contests_participated}
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-400">Contribution</span>
                      <span className="text-sm font-medium text-white">
                        {stats.codeforces.contribution}
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-400">Friend of</span>
                      <span className="text-sm font-medium text-white">
                        {stats.codeforces.friend_of_count}
                      </span>
                    </div>
                  </div>
                </div>
              ) : (
                <div className="text-center py-8">
                  <Award className="h-12 w-12 text-gray-600 mx-auto mb-3" />
                  <p className="text-gray-400 mb-3">Codeforces profile not linked</p>
                  <Button size="sm" variant="outline">
                    Link Profile
                  </Button>
                </div>
              )}
            </CardContent>
          </Card>

          {/* CodeChef Card */}
          <Card className="relative overflow-hidden">
            <div className={`absolute top-0 left-0 right-0 h-1 bg-gradient-to-r ${platformColors.codechef}`} />
            <CardHeader className="flex flex-row items-center justify-between">
              <div className="flex items-center gap-3">
                <div className={`p-2 rounded-lg bg-gradient-to-br ${platformColors.codechef}`}>
                  <Star className="h-6 w-6 text-white" />
                </div>
                <div>
                  <h2 className="text-xl font-bold text-white">CodeChef</h2>
                  {stats?.codechef && (
                    <p className="text-sm text-gray-400">@{stats.codechef.username}</p>
                  )}
                </div>
              </div>
              <div className="flex gap-2">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => handleSync('codechef')}
                  disabled={isSyncing === 'codechef'}
                >
                  {isSyncing === 'codechef' ? (
                    <Loader2 className="h-4 w-4 animate-spin" />
                  ) : (
                    <RefreshCw className="h-4 w-4" />
                  )}
                </Button>
                {stats?.codechef && (
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => window.open(`https://www.codechef.com/users/${stats.codechef?.username}`, '_blank')}
                  >
                    <ExternalLink className="h-4 w-4" />
                  </Button>
                )}
              </div>
            </CardHeader>
            <CardContent>
              {stats?.codechef && (stats.codechef.rating > 0 || stats.codechef.global_rank > 0 || stats.codechef.problems_solved > 0) ? (
                <div className="space-y-4">
                  <div className="grid grid-cols-2 gap-4">
                    <div className="bg-dojo-black-800 p-4 rounded-lg">
                      <div className="flex items-center gap-2 text-gray-400 text-sm mb-1">
                        <TrendingUp className="h-4 w-4" />
                        Rating
                      </div>
                      <div className="text-2xl font-bold text-white">
                        {stats.codechef.rating}
                      </div>
                      <div className="flex items-center gap-1 mt-1">
                        <Star className="h-3 w-3 text-yellow-400" />
                        <span className="text-xs text-yellow-400">{stats.codechef.stars}</span>
                      </div>
                    </div>
                    <div className="bg-dojo-black-800 p-4 rounded-lg">
                      <div className="flex items-center gap-2 text-gray-400 text-sm mb-1">
                        <Trophy className="h-4 w-4" />
                        Global Rank
                      </div>
                      <div className="text-2xl font-bold text-white">
                        #{stats.codechef.global_rank.toLocaleString()}
                      </div>
                    </div>
                  </div>

                  <div className="pt-3 border-t border-dojo-black-700 space-y-2">
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-400">Country Rank</span>
                      <span className="text-sm font-medium text-white">
                        #{stats.codechef.country_rank.toLocaleString()}
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-400">Problems Solved</span>
                      <span className="text-sm font-medium text-white">
                        {stats.codechef.problems_solved}
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-400">Contests</span>
                      <span className="text-sm font-medium text-white">
                        {stats.codechef.contests_participated}
                      </span>
                    </div>
                  </div>
                </div>
              ) : (
                <div className="text-center py-8">
                  <Star className="h-12 w-12 text-gray-600 mx-auto mb-3" />
                  <p className="text-gray-400 mb-3">CodeChef profile not linked</p>
                  <Button size="sm" variant="outline">
                    Link Profile
                  </Button>
                </div>
              )}
            </CardContent>
          </Card>

          {/* GeeksforGeeks Card */}
          <Card className="relative overflow-hidden">
            <div className={`absolute top-0 left-0 right-0 h-1 bg-gradient-to-r ${platformColors.geeksforgeeks}`} />
            <CardHeader className="flex flex-row items-center justify-between">
              <div className="flex items-center gap-3">
                <div className={`p-2 rounded-lg bg-gradient-to-br ${platformColors.geeksforgeeks}`}>
                  <Flame className="h-6 w-6 text-white" />
                </div>
                <div>
                  <h2 className="text-xl font-bold text-white">GeeksforGeeks</h2>
                  {stats?.geeksforgeeks && (
                    <p className="text-sm text-gray-400">@{stats.geeksforgeeks.username}</p>
                  )}
                </div>
              </div>
              <div className="flex gap-2">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => handleSync('geeksforgeeks')}
                  disabled={isSyncing === 'geeksforgeeks'}
                >
                  {isSyncing === 'geeksforgeeks' ? (
                    <Loader2 className="h-4 w-4 animate-spin" />
                  ) : (
                    <RefreshCw className="h-4 w-4" />
                  )}
                </Button>
                {stats?.geeksforgeeks && (
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => window.open(`https://auth.geeksforgeeks.org/user/${stats.geeksforgeeks?.username}`, '_blank')}
                  >
                    <ExternalLink className="h-4 w-4" />
                  </Button>
                )}
              </div>
            </CardHeader>
            <CardContent>
              {stats?.geeksforgeeks && (stats.geeksforgeeks.total_problems_solved > 0 || stats.geeksforgeeks.overall_score > 0) ? (
                <div className="space-y-4">
                  <div className="grid grid-cols-2 gap-4">
                    <div className="bg-dojo-black-800 p-4 rounded-lg">
                      <div className="flex items-center gap-2 text-gray-400 text-sm mb-1">
                        <Flame className="h-4 w-4" />
                        Current Streak
                      </div>
                      <div className="text-2xl font-bold text-white">
                        {stats.geeksforgeeks.current_streak}
                      </div>
                      <div className="text-xs text-green-400 mt-1">
                        Max: {stats.geeksforgeeks.max_streak} days
                      </div>
                    </div>
                    <div className="bg-dojo-black-800 p-4 rounded-lg">
                      <div className="flex items-center gap-2 text-gray-400 text-sm mb-1">
                        <Target className="h-4 w-4" />
                        Problems Solved
                      </div>
                      <div className="text-2xl font-bold text-white">
                        {stats.geeksforgeeks.total_problems_solved}
                      </div>
                    </div>
                  </div>

                  <div className="pt-3 border-t border-dojo-black-700 space-y-2">
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-400">Institute Rank</span>
                      <span className="text-sm font-medium text-white">
                        #{stats.geeksforgeeks.institute_rank}
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-400">Monthly Score</span>
                      <span className="text-sm font-medium text-white">
                        {stats.geeksforgeeks.monthly_score}
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-sm text-gray-400">Overall Score</span>
                      <span className="text-sm font-medium text-white">
                        {stats.geeksforgeeks.overall_score}
                      </span>
                    </div>
                  </div>
                </div>
              ) : (
                <div className="text-center py-8">
                  <Flame className="h-12 w-12 text-gray-600 mx-auto mb-3" />
                  <p className="text-gray-400 mb-3">GeeksforGeeks profile not linked</p>
                  <Button size="sm" variant="outline">
                    Link Profile
                  </Button>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
      </div>
    </div>
  );
}
