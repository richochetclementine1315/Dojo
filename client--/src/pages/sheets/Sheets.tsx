import { useState, useEffect } from 'react';
import { Navbar } from '@/components/layout/Navbar';
import { Card, CardHeader, CardContent } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { sheetService } from '@/services/sheetService';
import type { Sheet } from '@/types';
import Antigravity from '@/components/effects/Antigravity';
import { 
  FileText, 
  Loader2,
  Plus,
  Trash2,
  Lock,
  Globe,
  Edit,
  CheckCircle,
  Circle
} from 'lucide-react';

export default function Sheets() {
  const [sheets, setSheets] = useState<Sheet[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [newSheet, setNewSheet] = useState({
    name: '',
    description: '',
    is_public: false,
  });
  const [isCreating, setIsCreating] = useState(false);

  useEffect(() => {
    fetchSheets();
  }, []);

  const fetchSheets = async () => {
    try {
      setIsLoading(true);
      const data = await sheetService.getSheets();
      setSheets(Array.isArray(data) ? data : []);
      setError('');
    } catch (err: any) {
      console.error('Failed to fetch sheets:', err);
      setError(err.response?.data?.message || 'Failed to load sheets');
      setSheets([]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateSheet = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newSheet.name.trim()) return;

    try {
      setIsCreating(true);
      await sheetService.createSheet(newSheet);
      setNewSheet({ name: '', description: '', is_public: false });
      setShowCreateModal(false);
      await fetchSheets();
    } catch (err: any) {
      setError('Failed to create sheet');
    } finally {
      setIsCreating(false);
    }
  };

  const handleDeleteSheet = async (id: string) => {
    if (!confirm('Are you sure you want to delete this sheet?')) return;

    try {
      await sheetService.deleteSheet(id);
      await fetchSheets();
    } catch (err: any) {
      setError('Failed to delete sheet');
    }
  };

  const getCompletionPercentage = (sheet: Sheet) => {
    if (!sheet.problems || sheet.problems.length === 0) return 0;
    // Note: Would need to track solved problems per user
    // For now, returning 0 as placeholder
    return 0;
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
              Problem <span className="dojo-text-gradient">Sheets</span>
            </h1>
            <p className="text-gray-400">Organize and track your problem-solving progress</p>
          </div>
          <Button onClick={() => setShowCreateModal(true)}>
            <Plus className="mr-2 h-4 w-4" />
            Create Sheet
          </Button>
        </div>

        {error && (
          <div className="mb-6 p-4 bg-red-500/10 border border-red-500/20 rounded-lg text-red-400">
            {error}
          </div>
        )}

        {/* Create Modal */}
        {showCreateModal && (
          <div className="fixed inset-0 bg-black/80 flex items-center justify-center z-50">
            <Card className="w-full max-w-md">
              <CardHeader>
                <h2 className="text-2xl font-bold text-white">Create New Sheet</h2>
              </CardHeader>
              <CardContent>
                <form onSubmit={handleCreateSheet} className="space-y-4">
                  <div>
                    <label className="text-sm text-gray-400 mb-2 block">Sheet Name</label>
                    <Input
                      value={newSheet.name}
                      onChange={(e) => setNewSheet({ ...newSheet, name: e.target.value })}
                      placeholder="e.g., Dynamic Programming Fundamentals"
                      required
                    />
                  </div>

                  <div>
                    <label className="text-sm text-gray-400 mb-2 block">Description</label>
                    <textarea
                      value={newSheet.description}
                      onChange={(e) => setNewSheet({ ...newSheet, description: e.target.value })}
                      placeholder="Brief description of this sheet..."
                      className="w-full px-4 py-2 bg-dojo-black-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:border-dojo-red-500"
                      rows={3}
                    />
                  </div>

                  <div className="flex items-center gap-2">
                    <input
                      type="checkbox"
                      id="is_public"
                      checked={newSheet.is_public}
                      onChange={(e) => setNewSheet({ ...newSheet, is_public: e.target.checked })}
                      className="w-4 h-4 text-dojo-red-500 bg-dojo-black-800 border-gray-700 rounded focus:ring-dojo-red-500"
                    />
                    <label htmlFor="is_public" className="text-sm text-gray-400">
                      Make this sheet public
                    </label>
                  </div>

                  <div className="flex gap-3">
                    <Button
                      type="submit"
                      disabled={isCreating}
                      className="flex-1"
                    >
                      {isCreating ? (
                        <>
                          <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                          Creating...
                        </>
                      ) : (
                        'Create Sheet'
                      )}
                    </Button>
                    <Button
                      type="button"
                      variant="outline"
                      onClick={() => {
                        setShowCreateModal(false);
                        setNewSheet({ name: '', description: '', is_public: false });
                      }}
                    >
                      Cancel
                    </Button>
                  </div>
                </form>
              </CardContent>
            </Card>
          </div>
        )}

        {/* Sheets Grid */}
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <Loader2 className="h-8 w-8 animate-spin text-dojo-red-500" />
          </div>
        ) : sheets.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center">
              <FileText className="h-12 w-12 text-gray-600 mx-auto mb-4" />
              <p className="text-gray-400 text-lg">No sheets yet</p>
              <p className="text-gray-500 text-sm mt-2">Create your first problem sheet to get started</p>
              <Button onClick={() => setShowCreateModal(true)} className="mt-4">
                <Plus className="mr-2 h-4 w-4" />
                Create Sheet
              </Button>
            </CardContent>
          </Card>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {sheets.map((sheet) => {
              const completionPct = getCompletionPercentage(sheet);
              const problemCount = sheet.problems?.length || 0;
              
              return (
                <Card key={sheet.id} className="group hover:border-dojo-red-500/50 transition-all">
                  <CardHeader>
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        <div className="flex items-center gap-2 mb-2">
                          {sheet.is_public ? (
                            <Globe className="h-4 w-4 text-green-400" />
                          ) : (
                            <Lock className="h-4 w-4 text-gray-400" />
                          )}
                          <span className="text-xs text-gray-400">
                            {sheet.is_public ? 'Public' : 'Private'}
                          </span>
                        </div>
                        <h3 className="text-xl font-semibold text-white group-hover:text-dojo-red-400 transition-colors">
                          {sheet.name}
                        </h3>
                      </div>
                      <button
                        onClick={() => handleDeleteSheet(sheet.id)}
                        className="text-gray-400 hover:text-red-400 transition-colors"
                      >
                        <Trash2 className="h-4 w-4" />
                      </button>
                    </div>
                  </CardHeader>
                  
                  <CardContent>
                    <p className="text-gray-400 text-sm mb-4 line-clamp-2">
                      {sheet.description || 'No description provided'}
                    </p>

                    {/* Problem Count */}
                    <div className="flex items-center gap-2 mb-3">
                      <FileText className="h-4 w-4 text-gray-400" />
                      <span className="text-sm text-gray-400">
                        {problemCount} {problemCount === 1 ? 'problem' : 'problems'}
                      </span>
                    </div>

                    {/* Progress Bar */}
                    {problemCount > 0 && (
                      <div className="mb-4">
                        <div className="flex items-center justify-between text-xs text-gray-400 mb-1">
                          <span>Progress</span>
                          <span>{completionPct}%</span>
                        </div>
                        <div className="w-full bg-dojo-black-800 rounded-full h-2">
                          <div
                            className="bg-gradient-to-r from-dojo-red-500 to-orange-500 h-2 rounded-full transition-all"
                            style={{ width: `${completionPct}%` }}
                          />
                        </div>
                      </div>
                    )}

                    <Button
                      variant="outline"
                      className="w-full"
                      onClick={() => window.location.href = `/sheets/${sheet.id}`}
                    >
                      <Edit className="mr-2 h-4 w-4" />
                      View & Edit
                    </Button>
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
