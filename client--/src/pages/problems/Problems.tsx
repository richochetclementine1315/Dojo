import { useState, useEffect } from 'react';
import { Navbar } from '@/components/layout/Navbar';
import { Card, CardContent } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { problemsService, type Problem, type ProblemFilters } from '@/services/problemsService';
import { Loader2, ExternalLink, ChevronLeft, ChevronRight, Code2, Search } from 'lucide-react';
import Antigravity from '@/components/effects/Antigravity';

export default function Problems() {
  const [problems, setProblems] = useState<Problem[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalProblems, setTotalProblems] = useState(0);
  const [hasMore, setHasMore] = useState(false);
  const [selectedPlatform, setSelectedPlatform] = useState<'All' | 'leetcode' | 'codeforces' | 'codechef' | 'gfg'>('All');
  const [selectedDifficulty, setSelectedDifficulty] = useState<'All' | 'easy' | 'medium' | 'hard'>('All');
  const [searchQuery, setSearchQuery] = useState('');
  const [authError, setAuthError] = useState(false);
  
  const pageSize = 20;

  useEffect(() => {
    fetchProblems();
  }, [currentPage, selectedPlatform, selectedDifficulty]);

  const fetchProblems = async () => {
    try {
      setIsLoading(true);
      setAuthError(false);
      
      const filters: ProblemFilters = {
        page: currentPage,
        limit: pageSize,
      };
      
      if (selectedPlatform !== 'All') filters.platform = selectedPlatform;
      if (selectedDifficulty !== 'All') filters.difficulty = selectedDifficulty;
      if (searchQuery) filters.search = searchQuery;
      
      const data = await problemsService.getProblems(filters);
      setProblems(data.problems);
      setTotalProblems(data.total);
      setHasMore(data.hasMore);
    } catch (err: any) {
      if (err.response?.status === 401 || err.message === 'AUTHENTICATION_REQUIRED') {
        setAuthError(true);
      }
      setProblems([]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setCurrentPage(1);
    fetchProblems();
  };

  const getDifficultyColor = (difficulty: string) => {
    switch (difficulty.toLowerCase()) {
      case 'easy': return 'text-green-400 bg-green-400/10 border-green-400/20';
      case 'medium': return 'text-yellow-400 bg-yellow-400/10 border-yellow-400/20';
      case 'hard': return 'text-red-400 bg-red-400/10 border-red-400/20';
      default: return 'text-gray-400 bg-gray-400/10 border-gray-400/20';
    }
  };

  const getPlatformColor = (platform: string) => {
    switch (platform.toLowerCase()) {
      case 'leetcode': return 'bg-orange-500/10 text-orange-400 border-orange-500/20';
      case 'codeforces': return 'bg-blue-500/10 text-blue-400 border-blue-500/20';
      case 'codechef': return 'bg-purple-500/10 text-purple-400 border-purple-500/20';
      case 'gfg': return 'bg-green-500/10 text-green-400 border-green-500/20';
      default: return 'bg-gray-500/10 text-gray-400 border-gray-500/20';
    }
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
      
      <div className="container mx-auto px-4 py-6 md:py-8">
        {/* Header */}
        <div className="mb-6 md:mb-8">
          <div className="flex items-center gap-3 mb-3 md:mb-4">
            <Code2 className="h-7 w-7 md:h-8 md:w-8 text-dojo-red-500" />
            <h1 className="text-2xl md:text-4xl font-bold text-white">
              Practice Problems
            </h1>
          </div>
          <p className="text-sm md:text-base text-gray-400">
            Solve problems from LeetCode, Codeforces, CodeChef, and GeeksforGeeks
          </p>
        </div>

        {/* Search Bar */}
        <Card className="mb-4 md:mb-6">
          <CardContent className="p-4">
            <form onSubmit={handleSearch} className="flex gap-2">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 md:h-5 md:w-5 text-gray-400" />
                <Input
                  type="text"
                  placeholder="Search problems..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-9 md:pl-10"
                />
              </div>
              <Button type="submit" size="sm">
                Search
              </Button>
            </form>
          </CardContent>
        </Card>

        {/* Filters */}
        <div className="mb-6 space-y-3">
          {/* Platform Filter */}
          <div>
            <label className="text-xs md:text-sm text-gray-400 mb-2 block">Platform</label>
            <div className="flex flex-wrap gap-2">
              {['All', 'leetcode', 'codeforces', 'codechef', 'gfg'].map((platform) => (
                <Button
                  key={platform}
                  variant={selectedPlatform === platform ? 'default' : 'outline'}
                  onClick={() => {
                    setSelectedPlatform(platform as any);
                    setCurrentPage(1);
                  }}
                  size="sm"
                  className="rounded-full capitalize"
                >
                  {platform === 'gfg' ? 'GFG' : platform}
                </Button>
              ))}
            </div>
          </div>

          {/* Difficulty Filter */}
          <div>
            <label className="text-xs md:text-sm text-gray-400 mb-2 block">Difficulty</label>
            <div className="flex flex-wrap gap-2">
              {['All', 'easy', 'medium', 'hard'].map((difficulty) => (
                <Button
                  key={difficulty}
                  variant={selectedDifficulty === difficulty ? 'default' : 'outline'}
                  onClick={() => {
                    setSelectedDifficulty(difficulty as any);
                    setCurrentPage(1);
                  }}
                  size="sm"
                  className="rounded-full capitalize"
                >
                  {difficulty}
                </Button>
              ))}
            </div>
          </div>
        </div>

        {/* Problems List */}
        {authError ? (
          <Card>
            <CardContent className="py-12 text-center">
              <Code2 className="h-12 w-12 text-gray-600 mx-auto mb-4" />
              <p className="text-gray-400 text-lg mb-2">Authentication Required</p>
              <p className="text-gray-500 text-sm mb-4">Please log in to view problems</p>
              <Button onClick={() => window.location.href = '/login'}>
                Go to Login
              </Button>
            </CardContent>
          </Card>
        ) : isLoading ? (
          <div className="flex items-center justify-center py-12">
            <Loader2 className="h-8 w-8 animate-spin text-dojo-red-500" />
          </div>
        ) : (
          <>
            <div className="grid gap-3 md:gap-4">
              {problems.map((problem) => (
                <Card key={problem.id} className="hover:border-dojo-red-500/50 transition-colors">
                  <CardContent className="p-4 md:p-6">
                    <div className="flex flex-col md:flex-row md:items-center justify-between gap-3 md:gap-4">
                      <div className="flex-1">
                        <div className="flex flex-wrap items-center gap-2 mb-2">
                          <span className={`px-2 md:px-3 py-1 rounded-full text-xs font-medium border ${getPlatformColor(problem.platform)}`}>
                            {problem.platform === 'gfg' ? 'GFG' : problem.platform.charAt(0).toUpperCase() + problem.platform.slice(1)}
                          </span>
                          <span className={`px-2 md:px-3 py-1 rounded-full text-xs font-medium border ${getDifficultyColor(problem.difficulty)}`}>
                            {problem.difficulty.charAt(0).toUpperCase() + problem.difficulty.slice(1)}
                          </span>
                        </div>
                        <h3 className="text-base md:text-lg font-semibold text-white mb-2">
                          {problem.title}
                        </h3>
                        {problem.tags && problem.tags.length > 0 && (
                          <div className="flex flex-wrap gap-1.5 md:gap-2">
                            {problem.tags.slice(0, 5).map((tag, idx) => (
                              <span
                                key={idx}
                                className="px-2 py-0.5 md:py-1 bg-dojo-black-800 text-gray-400 text-xs rounded-md"
                              >
                                {tag}
                              </span>
                            ))}
                          </div>
                        )}
                      </div>
                      <Button
                        variant="outline"
                        size="sm"
                        className="rounded-full gap-2 w-full md:w-auto"
                        onClick={() => window.open(problem.problem_url, '_blank')}
                      >
                        Solve
                        <ExternalLink className="h-4 w-4" />
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>

            {problems.length === 0 && (
              <div className="text-center py-12">
                <Code2 className="h-12 w-12 text-gray-600 mx-auto mb-4" />
                <p className="text-gray-500">No problems found. Try adjusting your filters.</p>
              </div>
            )}

            {/* Pagination */}
            {problems.length > 0 && (
              <div className="mt-6 md:mt-8 flex flex-col sm:flex-row items-center justify-center gap-3 md:gap-4">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setCurrentPage(p => Math.max(1, p - 1))}
                  disabled={currentPage === 1}
                  className="rounded-full gap-2 w-full sm:w-auto"
                >
                  <ChevronLeft className="h-4 w-4" />
                  Previous
                </Button>
                <span className="text-sm md:text-base text-gray-400">
                  Page {currentPage} {totalProblems > 0 && `of ${Math.ceil(totalProblems / pageSize)}`}
                </span>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setCurrentPage(p => p + 1)}
                  disabled={!hasMore}
                  className="rounded-full gap-2 w-full sm:w-auto"
                >
                  Next
                  <ChevronRight className="h-4 w-4" />
                </Button>
              </div>
            )}
          </>
        )}
      </div>
      </div>
    </div>
  );
}
