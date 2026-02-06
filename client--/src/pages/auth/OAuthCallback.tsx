import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useAuthStore } from '@/store/authStore';
import api from '@/lib/api';
import type { User, ApiResponse } from '@/types';
import { Loader2 } from 'lucide-react';

export default function OAuthCallback() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { setAuth } = useAuthStore();
  const [error, setError] = useState('');

  useEffect(() => {
    const handleCallback = async () => {
      try {
        // Get tokens from URL params
        const accessToken = searchParams.get('access_token');
        const refreshToken = searchParams.get('refresh_token');
        const error = searchParams.get('error');

        if (error) {
          setError(error);
          setTimeout(() => navigate('/login'), 3000);
          return;
        }

        if (accessToken && refreshToken) {
          try {
            // Fetch user data with the access token directly
            const response = await api.get<ApiResponse<User>>('/users/profile', {
              headers: {
                Authorization: `Bearer ${accessToken}`
              }
            });
            
            const user = response.data.data;
            
            if (!user) {
              throw new Error('No user data received from server');
            }
            
            // Set auth in store (this will also update localStorage)
            setAuth(user, accessToken, refreshToken);
            
            // Redirect to dashboard
            setTimeout(() => {
              navigate('/dashboard');
            }, 1000);
            return;
          } catch (apiError: any) {
            setError(`Failed to fetch user data: ${apiError.response?.data?.message || apiError.message}`);
            setTimeout(() => navigate('/login'), 3000);
            return;
          }
        }

        // If no tokens found, show error
        setError('No authentication tokens received');
        setTimeout(() => navigate('/login'), 3000);

      } catch (err: any) {
        setError(err.response?.data?.message || err.message || 'Authentication failed');
        setTimeout(() => navigate('/login'), 3000);
      }
    };

    handleCallback();
  }, [navigate, searchParams, setAuth]);

  return (
    <div className="min-h-screen bg-dojo-gradient flex items-center justify-center">
      <div className="text-center">
        {error ? (
          <div className="space-y-4">
            <div className="text-red-400 text-lg">{error}</div>
            <p className="text-gray-400">Redirecting to login...</p>
          </div>
        ) : (
          <div className="space-y-4">
            <Loader2 className="h-12 w-12 animate-spin text-dojo-red-500 mx-auto" />
            <p className="text-white text-lg">Completing authentication...</p>
            <p className="text-gray-400 text-sm">Please wait while we log you in</p>
          </div>
        )}
      </div>
    </div>
  );
}
