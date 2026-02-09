import { useState, useEffect } from 'react';
import { Navbar } from '@/components/layout/Navbar';
import { Card, CardHeader, CardContent } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { authService } from '@/services/authService';
import { profileService } from '@/services/profileService';
import { 
  Save,
  Loader2,
  CheckCircle,
  XCircle,
  Link2,
  RefreshCw,
  ExternalLink
} from 'lucide-react';

interface PlatformUsernames {
  leetcode_username: string;
  codeforces_username: string;
  codechef_username: string;
  gfg_username: string;
}

interface SyncResult {
  [key: string]: { status?: string; error?: string };
}

export default function PlatformSettings() {
  const [usernames, setUsernames] = useState<PlatformUsernames>({
    leetcode_username: '',
    codeforces_username: '',
    codechef_username: '',
    gfg_username: '',
  });
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [isSyncing, setIsSyncing] = useState<string | null>(null);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [syncResults, setSyncResults] = useState<SyncResult>({});

  useEffect(() => {
    fetchProfile();
  }, []);

  const fetchProfile = async () => {
    try {
      setIsLoading(true);
      const data = await authService.getMe();
      console.log('Fetched profile data:', data);
      console.log('Profile object:', data.profile);
      const fetchedUsernames = {
        leetcode_username: data.profile?.leetcode_username || data.leetcode_username || '',
        codeforces_username: data.profile?.codeforces_username || data.codeforces_username || '',
        codechef_username: data.profile?.codechef_username || data.codechef_username || '',
        gfg_username: data.profile?.gfg_username || data.gfg_username || '',
      };
      console.log('Setting usernames to:', fetchedUsernames);
      setUsernames(fetchedUsernames);
    } catch (err: any) {
      console.error('Failed to fetch profile:', err);
      setError('Failed to load profile');
    } finally {
      setIsLoading(false);
    }
  };

  const handleSave = async () => {
    try {
      setIsSaving(true);
      setError('');
      setSuccess('');
      
      console.log('Saving platform usernames:', usernames);
      await profileService.updateProfile({
        bio: '',
        location: '',
        website: '',
        ...usernames,
      });
      
      // Refresh profile to get updated data
      await fetchProfile();
      
      setSuccess('Platform usernames saved successfully!');
      setTimeout(() => setSuccess(''), 3000);
    } catch (err: any) {
      console.error('Failed to save usernames:', err);
      setError(err.response?.data?.message || 'Failed to save usernames');
    } finally {
      setIsSaving(false);
    }
  };

  const handleSyncPlatform = async (platform: string) => {
    try {
      setIsSyncing(platform);
      setError('');
      setSyncResults({});
      
      console.log(`Syncing platform: ${platform}`);
      console.log(`Current usernames:`, usernames);
      const result = await profileService.syncPlatformStats([platform]);
      console.log(`Sync result for ${platform}:`, JSON.stringify(result, null, 2));
      console.log(`Platform result:`, result[platform]);
      setSyncResults(result);
      
      // Show success/error for this platform
      if (result[platform]?.status === 'success') {
        setSuccess(`${platform.toUpperCase()} stats synced successfully!`);
        setTimeout(() => setSuccess(''), 3000);
      } else if (result[platform]?.error) {
        console.error(`Sync error for ${platform}:`, result[platform].error);
        setError(`${platform.toUpperCase()}: ${result[platform].error}`);
      }
    } catch (err: any) {
      console.error(`Failed to sync ${platform}:`, err);
      console.error(`Error details:`, err.response?.data);
      setError(err.response?.data?.message || `Failed to sync ${platform} stats`);
    } finally {
      setIsSyncing(null);
    }
  };

  const handleSyncAll = async () => {
    const platforms = [];
    if (usernames.leetcode_username) platforms.push('leetcode');
    if (usernames.codeforces_username) platforms.push('codeforces');
    if (usernames.codechef_username) platforms.push('codechef');
    if (usernames.gfg_username) platforms.push('gfg');

    if (platforms.length === 0) {
      setError('Please add at least one platform username first');
      return;
    }

    try {
      setIsSyncing('all');
      setError('');
      setSyncResults({});
      
      const result = await profileService.syncPlatformStats(platforms);
      setSyncResults(result);
      
      const successCount = Object.values(result).filter((r): r is { status: string; error?: string } => 
        r !== null && typeof r === 'object' && 'status' in r && r.status === 'success'
      ).length;
      setSuccess(`Successfully synced ${successCount}/${platforms.length} platforms!`);
      setTimeout(() => setSuccess(''), 3000);
    } catch (err: any) {
      console.error('Failed to sync all platforms:', err);
      setError(err.response?.data?.message || 'Failed to sync platform stats');
    } finally {
      setIsSyncing(null);
    }
  };

  const platforms = [
    {
      name: 'LeetCode',
      key: 'leetcode',
      usernameKey: 'leetcode_username' as keyof PlatformUsernames,
      color: 'from-orange-500 to-yellow-500',
      placeholder: 'username only (e.g., Mrinmoy_1315)',
      verifyUrl: 'https://leetcode.com/',
    },
    {
      name: 'Codeforces',
      key: 'codeforces',
      usernameKey: 'codeforces_username' as keyof PlatformUsernames,
      color: 'from-blue-500 to-cyan-500',
      placeholder: 'username only (e.g., tourist)',
      verifyUrl: 'https://codeforces.com/profile/',
    },
    {
      name: 'CodeChef',
      key: 'codechef',
      usernameKey: 'codechef_username' as keyof PlatformUsernames,
      color: 'from-amber-500 to-orange-500',
      placeholder: 'Enter your CodeChef username',
      verifyUrl: 'https://www.codechef.com/users/',
    },
    {
      name: 'GeeksforGeeks',
      key: 'gfg',
      usernameKey: 'gfg_username' as keyof PlatformUsernames,
      color: 'from-green-500 to-emerald-500',
      placeholder: 'Enter your GFG username',
      verifyUrl: 'https://auth.geeksforgeeks.org/user/',
    },
  ];

  if (isLoading) {
    return (
      <div className="min-h-screen bg-dojo-black-900 flex items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-dojo-red-500" />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-dojo-black-900">
      <Navbar />
      
      <div className="container mx-auto px-4 py-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-white mb-2">
            Platform <span className="dojo-text-gradient">Connections</span>
          </h1>
          <p className="text-gray-400">Link your coding platform accounts to track your stats</p>
        </div>

        {/* Alerts */}
        {error && (
          <div className="mb-6 p-4 bg-red-500/10 border border-red-500/20 rounded-lg flex items-center gap-3">
            <XCircle className="h-5 w-5 text-red-400 flex-shrink-0" />
            <span className="text-red-400">{error}</span>
          </div>
        )}

        {success && (
          <div className="mb-6 p-4 bg-green-500/10 border border-green-500/20 rounded-lg flex items-center gap-3">
            <CheckCircle className="h-5 w-5 text-green-400 flex-shrink-0" />
            <span className="text-green-400">{success}</span>
          </div>
        )}

        {/* Sync All Button */}
        <div className="mb-6">
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="text-lg font-semibold text-white mb-1">Sync All Platforms</h3>
                  <p className="text-sm text-gray-400">
                    Update stats from all connected platforms at once
                  </p>
                </div>
                <Button
                  onClick={handleSyncAll}
                  disabled={isSyncing !== null}
                  className="flex items-center gap-2"
                >
                  {isSyncing === 'all' ? (
                    <>
                      <Loader2 className="h-4 w-4 animate-spin" />
                      Syncing...
                    </>
                  ) : (
                    <>
                      <RefreshCw className="h-4 w-4" />
                      Sync All
                    </>
                  )}
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Platform Cards */}
        <div className="grid gap-6">
          {platforms.map((platform) => {
            const username = usernames[platform.usernameKey];
            const isConnected = username !== '';
            const syncResult = syncResults[platform.key];

            return (
              <Card key={platform.key}>
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <div className={`h-12 w-12 rounded-lg bg-gradient-to-r ${platform.color} flex items-center justify-center`}>
                        <Link2 className="h-6 w-6 text-white" />
                      </div>
                      <div>
                        <h3 className="text-xl font-semibold text-white">{platform.name}</h3>
                        <div className="flex items-center gap-2 mt-1">
                          {isConnected ? (
                            <span className="text-xs px-2 py-0.5 rounded-full bg-green-500/20 text-green-400 flex items-center gap-1">
                              <CheckCircle className="h-3 w-3" />
                              Connected
                            </span>
                          ) : (
                            <span className="text-xs px-2 py-0.5 rounded-full bg-gray-500/20 text-gray-400">
                              Not Connected
                            </span>
                          )}
                        </div>
                      </div>
                    </div>
                    
                    {isConnected && (
                      <div className="flex gap-2">
                        <a
                          href={`${platform.verifyUrl}${username}`}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="p-2 hover:bg-dojo-black-800 rounded-lg transition-colors"
                        >
                          <ExternalLink className="h-5 w-5 text-gray-400" />
                        </a>
                        <Button
                          onClick={() => handleSyncPlatform(platform.key)}
                          disabled={isSyncing !== null}
                          variant="outline"
                          size="sm"
                        >
                          {isSyncing === platform.key ? (
                            <>
                              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                              Syncing...
                            </>
                          ) : (
                            <>
                              <RefreshCw className="mr-2 h-4 w-4" />
                              Sync
                            </>
                          )}
                        </Button>
                      </div>
                    )}
                  </div>
                </CardHeader>
                
                <CardContent>
                  <div className="space-y-4">
                    <div>
                      <label className="text-sm text-gray-400 mb-2 block">Username</label>
                      <Input
                        value={username}
                        onChange={(e) => setUsernames({
                          ...usernames,
                          [platform.usernameKey]: e.target.value,
                        })}
                        placeholder={platform.placeholder}
                      />
                    </div>

                    {syncResult && (
                      <div className={`p-3 rounded-lg ${
                        syncResult.status === 'success'
                          ? 'bg-green-500/10 border border-green-500/20'
                          : 'bg-red-500/10 border border-red-500/20'
                      }`}>
                        <div className="flex items-center gap-2 text-sm">
                          {syncResult.status === 'success' ? (
                            <>
                              <CheckCircle className="h-4 w-4 text-green-400" />
                              <span className="text-green-400">Stats synced successfully!</span>
                            </>
                          ) : (
                            <>
                              <XCircle className="h-4 w-4 text-red-400" />
                              <span className="text-red-400">{syncResult.error}</span>
                            </>
                          )}
                        </div>
                      </div>
                    )}
                  </div>
                </CardContent>
              </Card>
            );
          })}
        </div>

        {/* Save Button */}
        <div className="mt-8 flex justify-end">
          <Button
            onClick={handleSave}
            disabled={isSaving}
            size="lg"
            className="flex items-center gap-2"
          >
            {isSaving ? (
              <>
                <Loader2 className="h-5 w-5 animate-spin" />
                Saving...
              </>
            ) : (
              <>
                <Save className="h-5 w-5" />
                Save Usernames
              </>
            )}
          </Button>
        </div>
      </div>
    </div>
  );
}
