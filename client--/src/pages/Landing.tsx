import { Link } from 'react-router-dom';
import { Navbar } from '../components/layout/Navbar';
import { Button } from '../components/ui/Button';
import { useAuthStore } from '../store/authStore';
import { Sword, Code2, Users, Trophy, ArrowRight, ChevronLeft, ChevronRight, Calendar, Clock, ExternalLink } from 'lucide-react';
import { useEffect, useState } from 'react';
import { contestService } from '../services/contestService';
import type { Contest } from '../types';
import Antigravity from '../components/effects/Antigravity';
import Balatro from '../components/effects/Balatro';

export default function Landing() {
  const { isAuthenticated } = useAuthStore();
  const [contests, setContests] = useState<Contest[]>([]);
  const [currentSlide, setCurrentSlide] = useState(0);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    fetchContests();
  }, []);

  const fetchContests = async () => {
    try {
      setIsLoading(true);
      const data = await contestService.getContests({ upcoming: true, limit: 12 });
      setContests(Array.isArray(data) ? data : []);
    } catch (error) {
      console.error('Failed to fetch contests:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const contestsPerSlide = 3;
  const totalSlides = Math.ceil(contests.length / contestsPerSlide);

  const nextSlide = () => {
    setCurrentSlide((prev) => (prev + 1) % totalSlides);
  };

  const prevSlide = () => {
    setCurrentSlide((prev) => (prev - 1 + totalSlides) % totalSlides);
  };

  const getCurrentContests = () => {
    const start = currentSlide * contestsPerSlide;
    return contests.slice(start, start + contestsPerSlide);
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
  };

  const formatTime = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' });
  };

  const formatDuration = (seconds: number) => {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    return `${hours}h ${minutes}m`;
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
      
      {/* Hero */}
      <section className="relative min-h-[90vh] flex items-center justify-center overflow-hidden">
        {/* Balatro Shader Background */}
        <div className="absolute inset-0 z-0">
          <Balatro
            isRotate={false}
            mouseInteraction
            pixelFilter={745}
            color1="#DE443B"
            color2="#185d8c"
            color3="#271716"
          />
        </div>
        <div className="container mx-auto px-4 relative z-10">
          <div className="flex flex-col items-center text-center space-y-8 py-20">
            <img 
              src="/dojo-logo.png" 
              alt="Dojo" 
              className="h-64 w-64 object-contain" 
            />

            <h2 className="text-5xl md:text-6xl font-bold text-white max-w-4xl">
              Master Your Coding Journey
            </h2>

            <p className="text-xl md:text-2xl text-gray-300 max-w-2xl">
              Collaborative platform for competitive programming and problem solving
            </p>

            <div className="flex flex-col sm:flex-row gap-4 pt-8">
              {isAuthenticated ? (
                <Link to="/dashboard">
                  <Button size="lg" className="dojo-border-glow">
                    Go to Dashboard
                    <ArrowRight className="ml-2 h-5 w-5" />
                  </Button>
                </Link>
              ) : (
                <>
                  <Link to="/register">
                    <Button size="lg" className="dojo-border-glow">
                      Start Your Journey
                      <ArrowRight className="ml-2 h-5 w-5" />
                    </Button>
                  </Link>
                  <Link to="/login">
                    <Button variant="outline" size="lg">
                      Sign In
                    </Button>
                  </Link>
                </>
              )}
            </div>
          </div>
        </div>
      </section>

      {/* Upcoming Contests */}
      <section className="py-20 bg-gradient-to-b from-dojo-black-800 to-dojo-black-900">
        <div className="container mx-auto px-4">
          <h2 className="text-4xl font-bold text-center mb-12 dojo-text-gradient">
            Upcoming Contests
          </h2>

          {isLoading ? (
            <div className="flex justify-center items-center py-20">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-dojo-red-500"></div>
            </div>
          ) : contests.length === 0 ? (
            <div className="text-center py-20">
              <Trophy className="h-16 w-16 text-gray-600 mx-auto mb-4" />
              <p className="text-gray-400 text-lg">No upcoming contests at the moment</p>
            </div>
          ) : (
            <div className="relative">
              {/* Carousel */}
              <div className="grid md:grid-cols-3 gap-6 mb-8">
                {getCurrentContests().map((contest) => (
                  <div key={contest.id} className="dojo-card p-6 rounded-lg hover:scale-105 transition-transform">
                    <div className="flex items-start justify-between mb-4">
                      <span className={`px-3 py-1 rounded-full text-xs font-semibold ${
                        contest.platform === 'leetcode' ? 'bg-yellow-500/20 text-yellow-400' :
                        contest.platform === 'codeforces' ? 'bg-blue-500/20 text-blue-400' :
                        'bg-purple-500/20 text-purple-400'
                      }`}>
                        {contest.platform.toUpperCase()}
                      </span>
                    </div>

                    <h3 className="text-lg font-semibold text-white mb-3 line-clamp-2">
                      {contest.name}
                    </h3>

                    <div className="space-y-2 text-sm text-gray-400 mb-4">
                      <div className="flex items-center">
                        <Calendar className="h-4 w-4 mr-2 text-dojo-red-500" />
                        {formatDate(contest.start_time)}
                      </div>
                      <div className="flex items-center">
                        <Clock className="h-4 w-4 mr-2 text-dojo-red-500" />
                        {formatTime(contest.start_time)} • {formatDuration(contest.duration_seconds)}
                      </div>
                    </div>

                    {contest.contest_url && (
                      <a
                        href={contest.contest_url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="inline-flex items-center text-dojo-red-500 hover:text-dojo-red-400 text-sm font-medium"
                      >
                        View Contest
                        <ExternalLink className="ml-1 h-4 w-4" />
                      </a>
                    )}
                  </div>
                ))}
              </div>

              {/* Navigation */}
              {totalSlides > 1 && (
                <div className="flex items-center justify-center gap-4">
                  <button
                    onClick={prevSlide}
                    className="p-2 rounded-full bg-dojo-black-700 hover:bg-dojo-black-600 text-white transition-colors"
                    aria-label="Previous slide"
                  >
                    <ChevronLeft className="h-6 w-6" />
                  </button>

                  <div className="flex gap-2">
                    {Array.from({ length: totalSlides }).map((_, idx) => (
                      <button
                        key={idx}
                        onClick={() => setCurrentSlide(idx)}
                        className={`h-2 rounded-full transition-all ${
                          idx === currentSlide ? 'w-8 bg-dojo-red-500' : 'w-2 bg-gray-600'
                        }`}
                        aria-label={`Go to slide ${idx + 1}`}
                      />
                    ))}
                  </div>

                  <button
                    onClick={nextSlide}
                    className="p-2 rounded-full bg-dojo-black-700 hover:bg-dojo-black-600 text-white transition-colors"
                    aria-label="Next slide"
                  >
                    <ChevronRight className="h-6 w-6" />
                  </button>
                </div>
              )}

              <div className="text-center mt-8">
                <Link to="/contests">
                  <Button variant="outline">
                    View All Contests
                    <ArrowRight className="ml-2 h-4 w-4" />
                  </Button>
                </Link>
              </div>
            </div>
          )}
        </div>
      </section>

      {/* Features */}
      <section className="py-20 bg-dojo-black-900">
        <div className="container mx-auto px-4">
          <h2 className="text-4xl font-bold text-center mb-16 dojo-text-gradient">
            Train Like a Master
          </h2>

          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
            <div className="dojo-card p-6 rounded-lg hover:scale-105 transition-transform">
              <Code2 className="h-12 w-12 text-dojo-red-500 mb-4" />
              <h3 className="text-xl font-semibold text-white mb-2">Problem Solving</h3>
              <p className="text-gray-400">Track problems from LeetCode, Codeforces, and more</p>
            </div>

            <div className="dojo-card p-6 rounded-lg hover:scale-105 transition-transform">
              <Trophy className="h-12 w-12 text-dojo-red-500 mb-4" />
              <h3 className="text-xl font-semibold text-white mb-2">Contest Tracking</h3>
              <p className="text-gray-400">Never miss a contest with automated sync</p>
            </div>

            <div className="dojo-card p-6 rounded-lg hover:scale-105 transition-transform">
              <Users className="h-12 w-12 text-dojo-red-500 mb-4" />
              <h3 className="text-xl font-semibold text-white mb-2">Collaborative Rooms</h3>
              <p className="text-gray-400">Real-time code collaboration with video and chat</p>
            </div>

            <div className="dojo-card p-6 rounded-lg hover:scale-105 transition-transform">
              <Sword className="h-12 w-12 text-dojo-red-500 mb-4" />
              <h3 className="text-xl font-semibold text-white mb-2">Social Learning</h3>
              <p className="text-gray-400">Follow friends and climb the leaderboard</p>
            </div>
          </div>
        </div>
      </section>
      </div>
    </div>
  );
}