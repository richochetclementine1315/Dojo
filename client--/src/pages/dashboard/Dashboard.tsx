import { useAuthStore } from '../../store/authStore';
import { Navbar } from '../../components/layout/Navbar';
import { Card, CardHeader, CardContent } from '../../components/ui/Card';
import { Button } from '../../components/ui/Button';
import { Link } from 'react-router-dom';
import { useState, useEffect } from 'react';
import { problemService } from '@/services/problemService';
import { contestService } from '@/services/contestService';
import type { Problem, Contest } from '@/types';
import Antigravity from '../../components/effects/Antigravity';
import { 
  Trophy, 
  Target, 
  Flame,
  Code2,
  Users,
  Calendar,
  ArrowRight,
  Clock,
  ExternalLink
} from 'lucide-react';

export default function Dashboard() {
  const { user } = useAuthStore();
  const [recentProblems, setRecentProblems] = useState<Problem[]>([]);
  const [upcomingContests, setUpcomingContests] = useState<Contest[]>([]);
  const [solvedCount, setSolvedCount] = useState<number>(0);

  useEffect(() => {
    fetchDashboardData();
  }, []);

  const fetchDashboardData = async () => {
    try {
      const [problems, contests, solved] = await Promise.all([
        problemService.getProblems({ page: 1, limit: 5 }),
        contestService.getContests({ upcoming: true }),
        problemService.getSolvedCount(),
      ]);
      setRecentProblems(Array.isArray(problems) ? problems : []);
      setUpcomingContests(Array.isArray(contests) ? contests.slice(0, 3) : []);
      setSolvedCount(solved);
    } catch (err) {
      console.error('Failed to fetch dashboard data:', err);
      setRecentProblems([]);
      setUpcomingContests([]);
      setSolvedCount(0);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const getTimeUntil = (startTime: string) => {
    const now = new Date();
    const start = new Date(startTime);
    const diff = start.getTime() - now.getTime();
    
    if (diff < 0) return 'In Progress';
    
    const days = Math.floor(diff / (1000 * 60 * 60 * 24));
    const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
    
    if (days > 0) return `In ${days}d ${hours}h`;
    return `In ${hours}h`;
  };

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
      
      <div className="container mx-auto px-4 py-8 space-y-8">
        {/* Welcome Header */}
        <div>
          <h1 className="text-4xl font-bold text-white">
            Welcome back, <span className="dojo-text-gradient">{user?.username}</span>! 🥋
          </h1>
          <p className="text-gray-400 mt-2">Track your progress and continue your coding journey</p>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <Card className="hover:scale-105 transition-transform">
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <h3 className="text-sm font-medium text-gray-300">Problems Solved</h3>
              <Target className="h-4 w-4 text-dojo-red-500" />
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold text-white">{solvedCount}</div>
              <p className="text-xs text-gray-500 mt-1">
                {solvedCount === 0 ? 'Start solving problems' : 'Keep going!'}
              </p>
            </CardContent>
          </Card>

          <Card className="hover:scale-105 transition-transform">
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <h3 className="text-sm font-medium text-gray-300">Current Streak</h3>
              <Flame className="h-4 w-4 text-orange-500" />
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold text-white">0 days</div>
              <p className="text-xs text-gray-500 mt-1">Solve today to start</p>
            </CardContent>
          </Card>

          <Card className="hover:scale-105 transition-transform">
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <h3 className="text-sm font-medium text-gray-300">Contest Rating</h3>
              <Trophy className="h-4 w-4 text-yellow-500" />
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold text-white">-</div>
              <p className="text-xs text-gray-500 mt-1">Participate in contests</p>
            </CardContent>
          </Card>

          <Card className="hover:scale-105 transition-transform">
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <h3 className="text-sm font-medium text-gray-300">Global Rank</h3>
              <Trophy className="h-4 w-4 text-green-500" />
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold text-white">-</div>
              <p className="text-xs text-gray-500 mt-1">Keep solving to rank</p>
            </CardContent>
          </Card>
        </div>

        {/* Quick Actions */}
        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          <Link to="/problems">
            <Card className="group hover:border-dojo-red-500/50 transition-all cursor-pointer">
              <CardHeader>
                <div className="flex items-center space-x-3">
                  <div className="p-2 bg-dojo-red-500/20 rounded-lg">
                    <Code2 className="h-6 w-6 text-dojo-red-500" />
                  </div>
                  <h3 className="text-xl font-semibold text-white">Browse Problems</h3>
                </div>
              </CardHeader>
              <CardContent>
                <p className="text-gray-400 mb-4">
                  Explore thousands of coding problems from multiple platforms
                </p>
                <Button className="w-full">
                  Start Solving
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Button>
              </CardContent>
            </Card>
          </Link>

          <Link to="/contests">
            <Card className="group hover:border-yellow-500/50 transition-all cursor-pointer">
            <CardHeader>
              <div className="flex items-center space-x-3">
                <div className="p-2 bg-yellow-500/20 rounded-lg">
                  <Calendar className="h-6 w-6 text-yellow-500" />
                </div>
                <h3 className="text-xl font-semibold text-white">Upcoming Contests</h3>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-gray-400 mb-4">
                Never miss a contest with automated tracking
              </p>
              <Button className="w-full">
                View Contests
                <Calendar className="ml-2 h-4 w-4" />
              </Button>
            </CardContent>
          </Card>
          </Link>

          <Link to="/rooms">
            <Card className="group hover:border-blue-500/50 transition-all cursor-pointer">
              <CardHeader>
                <div className="flex items-center space-x-3">
                  <div className="p-2 bg-blue-500/20 rounded-lg">
                    <Users className="h-6 w-6 text-blue-500" />
                  </div>
                  <h3 className="text-xl font-semibold text-white">Collaborative Rooms</h3>
                </div>
              </CardHeader>
              <CardContent>
                <p className="text-gray-400 mb-4">
                  Practice with friends in real-time
                </p>
                <Button className="w-full">
                  Join Room
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Button>
              </CardContent>
            </Card>
          </Link>
        </div>

        {/* Recent Problems & Upcoming Contests */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Recent Problems */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <h3 className="text-xl font-semibold text-white">Recent Problems</h3>
                <Link to="/problems">
                  <Button variant="outline" size="sm">
                    View All
                    <ArrowRight className="ml-2 h-4 w-4" />
                  </Button>
                </Link>
              </div>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                {recentProblems.length === 0 ? (
                  <p className="text-gray-500 text-center py-4">No problems yet</p>
                ) : (
                  recentProblems.map((problem) => (
                    <div
                      key={problem.id}
                      className="p-3 bg-dojo-black-800 rounded-lg hover:bg-dojo-black-700 transition-colors cursor-pointer"
                    >
                      <div className="flex items-start justify-between gap-2">
                        <div className="flex-1">
                          <h4 className="text-white font-medium mb-1">{problem.title}</h4>
                          <div className="flex items-center gap-2 text-xs">
                            <span className={`px-2 py-0.5 rounded ${
                              problem.difficulty === 'easy' 
                                ? 'bg-green-500/20 text-green-400'
                                : problem.difficulty === 'medium'
                                ? 'bg-yellow-500/20 text-yellow-400'
                                : 'bg-red-500/20 text-red-400'
                            }`}>
                              {problem.difficulty}
                            </span>
                            <span className="text-gray-500">{problem.platform}</span>
                          </div>
                        </div>
                        <Button size="sm">Solve</Button>
                      </div>
                    </div>
                  ))
                )}
              </div>
            </CardContent>
          </Card>

          {/* Upcoming Contests */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <h3 className="text-xl font-semibold text-white">Upcoming Contests</h3>
                <Link to="/contests">
                  <Button variant="outline" size="sm">
                    View All
                    <ArrowRight className="ml-2 h-4 w-4" />
                  </Button>
                </Link>
              </div>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                {upcomingContests.length === 0 ? (
                  <p className="text-gray-500 text-center py-4">No upcoming contests</p>
                ) : (
                  upcomingContests.map((contest) => (
                    <div
                      key={contest.id}
                      className="p-3 bg-dojo-black-800 rounded-lg hover:bg-dojo-black-700 transition-colors cursor-pointer"
                      onClick={() => window.open(contest.url, '_blank')}
                    >
                      <div className="flex items-start justify-between gap-2">
                        <div className="flex-1">
                          <div className="flex items-center gap-2 mb-1">
                            <span className={`text-xs px-2 py-0.5 rounded text-white ${
                              contest.platform === 'leetcode' 
                                ? 'bg-orange-500'
                                : 'bg-blue-500'
                            }`}>
                              {contest.platform.toUpperCase()}
                            </span>
                            <span className="text-xs px-2 py-0.5 rounded bg-green-500/20 text-green-400">
                              {getTimeUntil(contest.start_time)}
                            </span>
                          </div>
                          <h4 className="text-white font-medium mb-1">{contest.name}</h4>
                          <div className="flex items-center gap-3 text-xs text-gray-500">
                            <span className="flex items-center gap-1">
                              <Clock className="h-3 w-3" />
                              {formatDate(contest.start_time)}
                            </span>
                          </div>
                        </div>
                        <ExternalLink className="h-4 w-4 text-gray-400 flex-shrink-0 mt-1" />
                      </div>
                    </div>
                  ))
                )}
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
      </div>
    </div>
  );
}