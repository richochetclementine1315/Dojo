import { useState, useEffect } from 'react';
import { Navbar } from '@/components/layout/Navbar';
import { Card, CardContent } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import {
  problemsService,
  type Problem,
  type CFProblem,
  type ProblemFilters,
  type CFProblemsFilters,
} from '@/services/problemsService';
import {
  Loader2,
  ExternalLink,
  ChevronLeft,
  ChevronRight,
  Code2,
  Search,
  RefreshCw,
  CheckCircle2,
  Circle,
  Zap,
} from 'lucide-react';
import Antigravity from '@/components/effects/Antigravity';

// ─── Rating → colour helpers ────────────────────────────────────────────────

function getCFRatingColor(rating: number): string {
  if (rating === 0) return 'text-gray-400';
  if (rating < 1200) return 'text-gray-300';        // newbie
  if (rating < 1400) return 'text-green-400';       // pupil
  if (rating < 1600) return 'text-cyan-400';        // specialist
  if (rating < 1900) return 'text-blue-400';        // expert
  if (rating < 2100) return 'text-violet-400';      // candidate master
  if (rating < 2300) return 'text-orange-400';      // master
  if (rating < 2400) return 'text-orange-500';      // international master
  if (rating < 3000) return 'text-red-500';         // grandmaster
  return 'text-red-600';                            // legendary grandmaster
}

function getDifficultyColor(difficulty: string): string {
  switch (difficulty.toLowerCase()) {
    case 'easy':   return 'text-green-400 bg-green-400/10 border-green-400/20';
    case 'medium': return 'text-yellow-400 bg-yellow-400/10 border-yellow-400/20';
    case 'hard':   return 'text-red-400 bg-red-400/10 border-red-400/20';
    default:       return 'text-gray-400 bg-gray-400/10 border-gray-400/20';
  }
}

function getPlatformColor(platform: string): string {
  switch (platform.toLowerCase()) {
    case 'leetcode':   return 'bg-orange-500/10 text-orange-400 border-orange-500/20';
    case 'codeforces': return 'bg-blue-500/10 text-blue-400 border-blue-500/20';
    case 'codechef':   return 'bg-purple-500/10 text-purple-400 border-purple-500/20';
    case 'gfg':        return 'bg-green-500/10 text-green-400 border-green-500/20';
    default:           return 'bg-gray-500/10 text-gray-400 border-gray-500/20';
  }
}

// ─── CF tag list (top tags from problemset.problems) ────────────────────────
const CF_POPULAR_TAGS = [
  'implementation', 'math', 'greedy', 'dp', 'data structures',
  'brute force', 'constructive algorithms', 'graphs', 'sortings',
  'binary search', 'dfs and similar', 'trees', 'strings',
  'number theory', 'geometry', 'bitmasks', 'two pointers',
  'combinatorics', 'divide and conquer',
];

// ─── Component ───────────────────────────────────────────────────────────────

type TabType = 'all' | 'codeforces';

export default function Problems() {
  // ── Shared state ──────────────────────────────────────────────────────────
  const [activeTab, setActiveTab] = useState<TabType>('all');

  // ── "All problems" tab state ──────────────────────────────────────────────
  const [problems, setProblems]             = useState<Problem[]>([]);
  const [isLoading, setIsLoading]           = useState(true);
  const [currentPage, setCurrentPage]       = useState(1);
  const [totalProblems, setTotalProblems]   = useState(0);
  const [hasMore, setHasMore]               = useState(false);
  const [selectedPlatform, setSelectedPlatform] =
    useState<'All' | 'leetcode' | 'codeforces' | 'codechef' | 'gfg'>('All');
  const [selectedDifficulty, setSelectedDifficulty] =
    useState<'All' | 'easy' | 'medium' | 'hard'>('All');
  const [searchQuery, setSearchQuery]       = useState('');
  const [authError, setAuthError]           = useState(false);
  const [isSyncing, setIsSyncing]           = useState(false);
  const [syncSuccess, setSyncSuccess]       = useState('');

  // ── "Browse Codeforces" tab state ──────────────────────────────────────────
  const [cfProblems, setCFProblems]         = useState<CFProblem[]>([]);
  const [cfLoading, setCFLoading]           = useState(false);
  const [cfPage, setCFPage]                 = useState(1);
  const [cfTotal, setCFTotal]               = useState(0);
  const [cfHasMore, setCFHasMore]           = useState(false);
  const [cfSearch, setCFSearch]             = useState('');
  const [cfMinRating, setCFMinRating]       = useState<string>('');
  const [cfMaxRating, setCFMaxRating]       = useState<string>('');
  const [cfSelectedTags, setCFSelectedTags] = useState<string[]>([]);
  const [cfError, setCFError]               = useState('');

  const pageSize = 20;

  // ── Effects ───────────────────────────────────────────────────────────────

  useEffect(() => {
    if (activeTab === 'all') fetchProblems();
  }, [currentPage, selectedPlatform, selectedDifficulty, activeTab]);

  useEffect(() => {
    if (activeTab === 'codeforces') fetchCFProblems();
  }, [cfPage, activeTab]);

  // ── All-problems fetch ────────────────────────────────────────────────────

  const fetchProblems = async () => {
    try {
      setIsLoading(true);
      setAuthError(false);
      const filters: ProblemFilters = { page: currentPage, limit: pageSize };
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

  const handleSyncProblems = async (platform: 'leetcode' | 'codeforces') => {
    try {
      setIsSyncing(true);
      const result = await problemsService.syncProblems(platform, 100);
      setSyncSuccess(`Successfully synced ${result.imported} ${platform} problems!`);
      setTimeout(() => setSyncSuccess(''), 3000);
      fetchProblems();
    } catch (err: any) {
      console.error('Failed to sync problems:', err);
    } finally {
      setIsSyncing(false);
    }
  };

  const handleMarkSolved = async (problemId: string, currentlymark: boolean) => {
    try {
      await problemsService.markProblemSolved(problemId, !currentlymark);
      setProblems(problems.map(p =>
        p.id === problemId ? { ...p, is_solved: !currentlymark } : p
      ));
    } catch (err: any) {
      console.error('Failed to mark problem:', err);
    }
  };

  // ── CF-browse fetch ───────────────────────────────────────────────────────

  const fetchCFProblems = async () => {
    try {
      setCFLoading(true);
      setCFError('');
      const filters: CFProblemsFilters = { page: cfPage, limit: pageSize };
      if (cfSearch) filters.search = cfSearch;
      if (cfMinRating) filters.min_rating = parseInt(cfMinRating);
      if (cfMaxRating) filters.max_rating = parseInt(cfMaxRating);
      if (cfSelectedTags.length > 0) filters.tags = cfSelectedTags;

      const data = await problemsService.getCodeforcesProblems(filters);
      setCFProblems(data.problems);
      setCFTotal(data.total);
      setCFHasMore(data.has_more);
    } catch (err: any) {
      setCFError('Failed to load Codeforces problems. Please try again.');
      setCFProblems([]);
    } finally {
      setCFLoading(false);
    }
  };

  const handleCFSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setCFPage(1);
    fetchCFProblems();
  };

  const toggleCFTag = (tag: string) => {
    setCFSelectedTags(prev =>
      prev.includes(tag) ? prev.filter(t => t !== tag) : [...prev, tag]
    );
  };

  // ── Render ────────────────────────────────────────────────────────────────

  return (
    <div className="min-h-screen bg-dojo-black-900 relative">
      {/* Background */}
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

          {/* Tab switcher */}
          <div className="flex gap-2 mb-6">
            <Button
              variant={activeTab === 'all' ? 'default' : 'outline'}
              size="sm"
              className="rounded-full gap-2"
              onClick={() => setActiveTab('all')}
            >
              <Code2 className="h-4 w-4" />
              All Problems
            </Button>
            <Button
              variant={activeTab === 'codeforces' ? 'default' : 'outline'}
              size="sm"
              className="rounded-full gap-2"
              onClick={() => setActiveTab('codeforces')}
            >
              <Zap className="h-4 w-4 text-blue-400" />
              Browse Codeforces
            </Button>
          </div>

          {/* ─── ALL PROBLEMS TAB ─────────────────────────────────────────── */}
          {activeTab === 'all' && (
            <>
              {/* Search */}
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
                    <Button type="submit" size="sm">Search</Button>
                  </form>
                </CardContent>
              </Card>

              {/* Filters */}
              <div className="mb-6 space-y-3">
                <div>
                  <label className="text-xs md:text-sm text-gray-400 mb-2 block">Platform</label>
                  <div className="flex flex-wrap gap-2">
                    {['All', 'leetcode', 'codeforces', 'codechef', 'gfg'].map((platform) => (
                      <Button
                        key={platform}
                        variant={selectedPlatform === platform ? 'default' : 'outline'}
                        onClick={() => { setSelectedPlatform(platform as any); setCurrentPage(1); }}
                        size="sm"
                        className="rounded-full capitalize"
                      >
                        {platform === 'gfg' ? 'GFG' : platform}
                      </Button>
                    ))}
                  </div>
                </div>
                <div>
                  <label className="text-xs md:text-sm text-gray-400 mb-2 block">Difficulty</label>
                  <div className="flex flex-wrap gap-2">
                    {['All', 'easy', 'medium', 'hard'].map((difficulty) => (
                      <Button
                        key={difficulty}
                        variant={selectedDifficulty === difficulty ? 'default' : 'outline'}
                        onClick={() => { setSelectedDifficulty(difficulty as any); setCurrentPage(1); }}
                        size="sm"
                        className="rounded-full capitalize"
                      >
                        {difficulty}
                      </Button>
                    ))}
                  </div>
                </div>
                <div>
                  <label className="text-xs md:text-sm text-gray-400 mb-2 block">Sync Problems</label>
                  <div className="flex flex-wrap gap-2">
                    <Button
                      variant="outline" size="sm"
                      onClick={() => handleSyncProblems('leetcode')}
                      disabled={isSyncing}
                      className="rounded-full gap-2"
                    >
                      <RefreshCw className={`h-4 w-4 ${isSyncing ? 'animate-spin' : ''}`} />
                      LeetCode
                    </Button>
                    <Button
                      variant="outline" size="sm"
                      onClick={() => handleSyncProblems('codeforces')}
                      disabled={isSyncing}
                      className="rounded-full gap-2"
                    >
                      <RefreshCw className={`h-4 w-4 ${isSyncing ? 'animate-spin' : ''}`} />
                      Codeforces
                    </Button>
                  </div>
                </div>
              </div>

              {syncSuccess && (
                <div className="bg-green-500/10 border border-green-500/50 text-green-500 px-4 py-2 rounded-lg mb-4">
                  {syncSuccess}
                </div>
              )}

              {/* Problem list */}
              {authError ? (
                <Card>
                  <CardContent className="py-12 text-center">
                    <Code2 className="h-12 w-12 text-gray-600 mx-auto mb-4" />
                    <p className="text-gray-400 text-lg mb-2">Authentication Required</p>
                    <p className="text-gray-500 text-sm mb-4">Please log in to view problems</p>
                    <Button onClick={() => window.location.href = '/login'}>Go to Login</Button>
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
                                <button
                                  onClick={() => handleMarkSolved(problem.id, problem.is_solved || false)}
                                  className="flex items-center gap-1.5 hover:opacity-80 transition-opacity"
                                  title={problem.is_solved ? 'Mark as unsolved' : 'Mark as solved'}
                                >
                                  {problem.is_solved
                                    ? <CheckCircle2 className="h-5 w-5 text-green-500" />
                                    : <Circle className="h-5 w-5 text-gray-500" />}
                                </button>
                                <span className={`px-2 md:px-3 py-1 rounded-full text-xs font-medium border ${getPlatformColor(problem.platform)}`}>
                                  {problem.platform === 'gfg' ? 'GFG' : problem.platform.charAt(0).toUpperCase() + problem.platform.slice(1)}
                                </span>
                                <span className={`px-2 md:px-3 py-1 rounded-full text-xs font-medium border ${getDifficultyColor(problem.difficulty)}`}>
                                  {problem.difficulty.charAt(0).toUpperCase() + problem.difficulty.slice(1)}
                                </span>
                                {/* Show CF rating badge for codeforces problems */}
                                {problem.platform === 'codeforces' && (problem.cf_rating ?? 0) > 0 && (
                                  <span className={`px-2 py-1 rounded-full text-xs font-semibold bg-blue-900/30 border border-blue-700/30 ${getCFRatingColor(problem.cf_rating!)}`}>
                                    ★ {problem.cf_rating}
                                  </span>
                                )}
                              </div>
                              <h3 className="text-base md:text-lg font-semibold text-white mb-2">
                                {problem.title}
                              </h3>
                              {problem.tags && problem.tags.length > 0 && (
                                <div className="flex flex-wrap gap-1.5 md:gap-2">
                                  {problem.tags.slice(0, 5).map((tag, idx) => (
                                    <span key={idx} className="px-2 py-0.5 md:py-1 bg-dojo-black-800 text-gray-400 text-xs rounded-md">
                                      {tag}
                                    </span>
                                  ))}
                                </div>
                              )}
                            </div>
                            <Button
                              variant="outline" size="sm"
                              className="rounded-full gap-2 w-full md:w-auto"
                              onClick={() => window.open(problem.problem_url, '_blank')}
                            >
                              Solve <ExternalLink className="h-4 w-4" />
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

                  {problems.length > 0 && (
                    <div className="mt-6 md:mt-8 flex flex-col sm:flex-row items-center justify-center gap-3 md:gap-4">
                      <Button
                        variant="outline" size="sm"
                        onClick={() => setCurrentPage(p => Math.max(1, p - 1))}
                        disabled={currentPage === 1}
                        className="rounded-full gap-2 w-full sm:w-auto"
                      >
                        <ChevronLeft className="h-4 w-4" /> Previous
                      </Button>
                      <span className="text-sm md:text-base text-gray-400">
                        Page {currentPage} {totalProblems > 0 && `of ${Math.ceil(totalProblems / pageSize)}`}
                      </span>
                      <Button
                        variant="outline" size="sm"
                        onClick={() => setCurrentPage(p => p + 1)}
                        disabled={!hasMore}
                        className="rounded-full gap-2 w-full sm:w-auto"
                      >
                        Next <ChevronRight className="h-4 w-4" />
                      </Button>
                    </div>
                  )}
                </>
              )}
            </>
          )}

          {/* ─── BROWSE CODEFORCES TAB ────────────────────────────────────── */}
          {activeTab === 'codeforces' && (
            <>
              <Card className="mb-4 md:mb-6">
                <CardContent className="p-4 space-y-4">
                  {/* Search row */}
                  <form onSubmit={handleCFSearch} className="flex gap-2">
                    <div className="relative flex-1">
                      <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
                      <Input
                        type="text"
                        placeholder="Search by problem name..."
                        value={cfSearch}
                        onChange={(e) => setCFSearch(e.target.value)}
                        className="pl-9"
                      />
                    </div>
                    <Button type="submit" size="sm">Search</Button>
                  </form>

                  {/* Rating range */}
                  <div className="flex gap-3 items-center">
                    <label className="text-sm text-gray-400 whitespace-nowrap">Rating range</label>
                    <Input
                      type="number"
                      placeholder="Min (e.g. 800)"
                      value={cfMinRating}
                      onChange={(e) => setCFMinRating(e.target.value)}
                      className="w-36"
                    />
                    <span className="text-gray-500">–</span>
                    <Input
                      type="number"
                      placeholder="Max (e.g. 2000)"
                      value={cfMaxRating}
                      onChange={(e) => setCFMaxRating(e.target.value)}
                      className="w-36"
                    />
                    <Button
                      size="sm" variant="outline"
                      onClick={() => { setCFPage(1); fetchCFProblems(); }}
                      className="rounded-full gap-1"
                    >
                      <RefreshCw className="h-3.5 w-3.5" /> Apply
                    </Button>
                  </div>

                  {/* Tag filter */}
                  <div>
                    <label className="text-xs text-gray-400 mb-2 block">Filter by tag</label>
                    <div className="flex flex-wrap gap-2">
                      {CF_POPULAR_TAGS.map((tag) => (
                        <button
                          key={tag}
                          onClick={() => { toggleCFTag(tag); }}
                          className={`px-2 py-1 rounded-md text-xs border transition-colors ${
                            cfSelectedTags.includes(tag)
                              ? 'bg-blue-600/30 border-blue-500/60 text-blue-300'
                              : 'bg-dojo-black-800 border-gray-700 text-gray-400 hover:border-gray-500'
                          }`}
                        >
                          {tag}
                        </button>
                      ))}
                    </div>
                    {cfSelectedTags.length > 0 && (
                      <button
                        onClick={() => setCFSelectedTags([])}
                        className="mt-2 text-xs text-gray-500 hover:text-gray-300 underline"
                      >
                        Clear tags
                      </button>
                    )}
                  </div>

                  <Button
                    className="w-full sm:w-auto gap-2"
                    onClick={() => { setCFPage(1); fetchCFProblems(); }}
                  >
                    <Zap className="h-4 w-4" /> Fetch Problems
                  </Button>
                </CardContent>
              </Card>

              {/* Info banner */}
              <div className="bg-blue-900/20 border border-blue-700/30 rounded-lg px-4 py-2 mb-4 text-xs text-blue-300">
                Problems are fetched live from the Codeforces API and cached for 6 hours. Click "Solve" to open on Codeforces.
              </div>

              {cfError && (
                <div className="bg-red-500/10 border border-red-500/30 text-red-400 px-4 py-2 rounded-lg mb-4 text-sm">
                  {cfError}
                </div>
              )}

              {cfLoading ? (
                <div className="flex items-center justify-center py-12">
                  <Loader2 className="h-8 w-8 animate-spin text-blue-500" />
                </div>
              ) : (
                <>
                  {cfProblems.length > 0 && (
                    <p className="text-sm text-gray-500 mb-3">
                      Showing {cfProblems.length} of {cfTotal.toLocaleString()} problems
                    </p>
                  )}

                  <div className="grid gap-3 md:gap-4">
                    {cfProblems.map((problem) => (
                      <Card
                        key={`${problem.contest_id}-${problem.index}`}
                        className="hover:border-blue-500/50 transition-colors"
                      >
                        <CardContent className="p-4 md:p-6">
                          <div className="flex flex-col md:flex-row md:items-center justify-between gap-3 md:gap-4">
                            <div className="flex-1">
                              <div className="flex flex-wrap items-center gap-2 mb-2">
                                {/* CF rating badge */}
                                {problem.rating > 0 ? (
                                  <span className={`px-2 py-1 rounded-full text-xs font-semibold bg-blue-900/30 border border-blue-700/30 ${getCFRatingColor(problem.rating)}`}>
                                    ★ {problem.rating}
                                  </span>
                                ) : (
                                  <span className="px-2 py-1 rounded-full text-xs font-medium bg-gray-800 border border-gray-700 text-gray-500">
                                    Unrated
                                  </span>
                                )}

                                {/* Contest + Index */}
                                <span className="px-2 py-1 rounded-full text-xs bg-blue-500/10 border border-blue-500/20 text-blue-400">
                                  {problem.contest_id}{problem.index}
                                </span>

                                {/* Solved count */}
                                {problem.solved_count > 0 && (
                                  <span className="text-xs text-gray-500">
                                    {problem.solved_count.toLocaleString()} solved
                                  </span>
                                )}
                              </div>

                              <h3 className="text-base md:text-lg font-semibold text-white mb-2">
                                {problem.name}
                              </h3>

                              {problem.tags && problem.tags.length > 0 && (
                                <div className="flex flex-wrap gap-1.5 md:gap-2">
                                  {problem.tags.map((tag, idx) => (
                                    <span
                                      key={idx}
                                      className="px-2 py-0.5 bg-dojo-black-800 text-gray-400 text-xs rounded-md"
                                    >
                                      {tag}
                                    </span>
                                  ))}
                                </div>
                              )}
                            </div>

                            <Button
                              variant="outline" size="sm"
                              className="rounded-full gap-2 w-full md:w-auto border-blue-700/40 hover:border-blue-500"
                              onClick={() => window.open(problem.problem_url, '_blank')}
                            >
                              Solve <ExternalLink className="h-4 w-4" />
                            </Button>
                          </div>
                        </CardContent>
                      </Card>
                    ))}
                  </div>

                  {cfProblems.length === 0 && !cfLoading && (
                    <div className="text-center py-12">
                      <Zap className="h-12 w-12 text-gray-600 mx-auto mb-4" />
                      <p className="text-gray-500">
                        Click "Fetch Problems" to load Codeforces problems, or adjust your filters.
                      </p>
                    </div>
                  )}

                  {cfProblems.length > 0 && (
                    <div className="mt-6 md:mt-8 flex flex-col sm:flex-row items-center justify-center gap-3 md:gap-4">
                      <Button
                        variant="outline" size="sm"
                        onClick={() => setCFPage(p => Math.max(1, p - 1))}
                        disabled={cfPage === 1}
                        className="rounded-full gap-2 w-full sm:w-auto"
                      >
                        <ChevronLeft className="h-4 w-4" /> Previous
                      </Button>
                      <span className="text-sm md:text-base text-gray-400">
                        Page {cfPage} {cfTotal > 0 && `of ${Math.ceil(cfTotal / pageSize)}`}
                      </span>
                      <Button
                        variant="outline" size="sm"
                        onClick={() => setCFPage(p => p + 1)}
                        disabled={!cfHasMore}
                        className="rounded-full gap-2 w-full sm:w-auto"
                      >
                        Next <ChevronRight className="h-4 w-4" />
                      </Button>
                    </div>
                  )}
                </>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  );
}
