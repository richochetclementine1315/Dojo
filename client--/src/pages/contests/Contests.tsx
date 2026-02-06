import { useState, useEffect } from 'react';
import { Navbar } from '@/components/layout/Navbar';
import { Card, CardHeader, CardContent } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { contestService } from '@/services/contestService';
import type { Contest } from '@/types';
import Antigravity from '@/components/effects/Antigravity';
import { 
  Calendar, 
  Clock, 
  ExternalLink,
  Loader2,
  Trophy,
  RefreshCw,
  Filter
} from 'lucide-react';

const platformColors = {
  leetcode: 'bg-gradient-to-r from-orange-500 to-yellow-500',
  codeforces: 'bg-gradient-to-r from-blue-500 to-cyan-500',
};

export default function Contests() {
  const [contests, setContests] = useState<Contest[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isSyncing, setIsSyncing] = useState(false);
  const [error, setError] = useState('');
  const [filter, setFilter] = useState<'all' | 'upcoming' | 'past'>('upcoming');
  const [platformFilter, setPlatformFilter] = useState<'all' | 'leetcode' | 'codeforces'>('all');

  useEffect(() => {
    fetchContests();
  }, [filter, platformFilter]);

  const fetchContests = async () => {
    try {
      setIsLoading(true);
      const data = await contestService.getContests({
        platform: platformFilter === 'all' ? undefined : platformFilter,
        upcoming: filter === 'upcoming' ? true : filter === 'past' ? false : undefined,
      });
      setContests(Array.isArray(data) ? data : []);
      setError('');
    } catch (err: any) {
      console.error('Failed to fetch contests:', err);
      setError(err.response?.data?.message || 'Failed to load contests');
      setContests([]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleSync = async () => {
    try {
      setIsSyncing(true);
      await contestService.syncContests();
      await fetchContests();
    } catch (err: any) {
      setError('Failed to sync contests');
    } finally {
      setIsSyncing(false);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const formatDuration = (minutes: number) => {
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    return `${hours}h ${mins}m`;
  };

  const getTimeUntil = (startTime: string) => {
    const now = new Date();
    const start = new Date(startTime);
    const diff = start.getTime() - now.getTime();
    
    if (diff < 0) return 'In Progress / Ended';
    
    const days = Math.floor(diff / (1000 * 60 * 60 * 24));
    const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
    const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
    
    if (days > 0) return `In ${days}d ${hours}h`;
    if (hours > 0) return `In ${hours}h ${minutes}m`;
    return `In ${minutes}m`;
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
      
      <div className="container mx-auto px-4 py-8">
        {/* Header */}
        <div className="flex items-center justify-between mb-8">
          <div>
            <h1 className="text-4xl font-bold text-white mb-2">
              Contest <span className="dojo-text-gradient">Calendar</span>
            </h1>
            <p className="text-gray-400">Track and participate in coding contests</p>
          </div>
          <Button
            onClick={handleSync}
            disabled={isSyncing}
            variant="outline"
          >
            {isSyncing ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Syncing...
              </>
            ) : (
              <>
                <RefreshCw className="mr-2 h-4 w-4" />
                Sync Contests
              </>
            )}
          </Button>
        </div>

        {/* Filters */}
        <Card className="mb-6">
          <CardContent className="p-6">
            <div className="flex flex-wrap gap-6">
              {/* Time Filter */}
              <div>
                <label className="text-sm text-gray-400 mb-2 flex items-center gap-2">
                  <Filter className="h-4 w-4" />
                  Time
                </label>
                <div className="flex gap-2">
                  {(['upcoming', 'all', 'past'] as const).map((f) => (
                    <button
                      key={f}
                      onClick={() => setFilter(f)}
                      className={`px-4 py-1.5 rounded-full text-sm font-medium transition-all ${
                        filter === f
                          ? 'bg-dojo-red-500 text-white'
                          : 'bg-dojo-black-800 text-gray-400 hover:bg-dojo-black-700'
                      }`}
                    >
                      {f.charAt(0).toUpperCase() + f.slice(1)}
                    </button>
                  ))}
                </div>
              </div>

              {/* Platform Filter */}
              <div>
                <label className="text-sm text-gray-400 mb-2 block">Platform</label>
                <div className="flex gap-2">
                  {(['all', 'leetcode', 'codeforces'] as const).map((p) => (
                    <button
                      key={p}
                      onClick={() => setPlatformFilter(p)}
                      className={`px-4 py-1.5 rounded-full text-sm font-medium transition-all ${
                        platformFilter === p
                          ? 'bg-dojo-red-500 text-white'
                          : 'bg-dojo-black-800 text-gray-400 hover:bg-dojo-black-700'
                      }`}
                    >
                      {p.charAt(0).toUpperCase() + p.slice(1)}
                    </button>
                  ))}
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {error && (
          <div className="mb-6 p-4 bg-red-500/10 border border-red-500/20 rounded-lg text-red-400">
            {error}
          </div>
        )}

        {/* Contests List */}
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <Loader2 className="h-8 w-8 animate-spin text-dojo-red-500" />
          </div>
        ) : contests.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center">
              <Calendar className="h-12 w-12 text-gray-600 mx-auto mb-4" />
              <p className="text-gray-400 text-lg">No contests found</p>
              <p className="text-gray-500 text-sm mt-2">Try syncing or adjusting filters</p>
            </CardContent>
          </Card>
        ) : (
          <div className="space-y-4">
            {contests.map((contest) => {
              const isPast = new Date(contest.start_time) < new Date();
              
              return (
                <Card key={contest.id} className="group hover:border-dojo-red-500/50 transition-all">
                  <CardContent className="p-6">
                    <div className="flex items-start justify-between gap-4">
                      <div className="flex-1">
                        <div className="flex items-center gap-3 mb-2">
                          <div className={`px-3 py-1 rounded-full text-xs font-medium text-white ${
                            platformColors[contest.platform]
                          }`}>
                            {contest.platform.toUpperCase()}
                          </div>
                          {!isPast && (
                            <div className="px-3 py-1 rounded-full text-xs font-medium bg-green-500/20 text-green-400">
                              {getTimeUntil(contest.start_time)}
                            </div>
                          )}
                          {isPast && (
                            <div className="px-3 py-1 rounded-full text-xs font-medium bg-gray-500/20 text-gray-400">
                              Ended
                            </div>
                          )}
                        </div>

                        <h3 className="text-xl font-semibold text-white mb-3 group-hover:text-dojo-red-400 transition-colors">
                          {contest.name}
                        </h3>

                        <div className="flex flex-wrap gap-4 text-sm text-gray-400">
                          <div className="flex items-center gap-2">
                            <Calendar className="h-4 w-4" />
                            <span>{formatDate(contest.start_time)}</span>
                          </div>
                          <div className="flex items-center gap-2">
                            <Clock className="h-4 w-4" />
                            <span>{formatDuration(contest.duration)}</span>
                          </div>
                          {contest.phase && (
                            <div className="flex items-center gap-2">
                              <Trophy className="h-4 w-4" />
                              <span>{contest.phase}</span>
                            </div>
                          )}
                        </div>
                      </div>

                      <Button
                        onClick={() => window.open(contest.url, '_blank')}
                        className="flex-shrink-0"
                      >
                        View Contest
                        <ExternalLink className="ml-2 h-4 w-4" />
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              );
            })}
          </div>
        )}
      </div>
      </div>
    </div>
  );
}
