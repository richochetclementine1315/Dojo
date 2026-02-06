import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { User } from '../types';
import { authService } from '../services/authService';

interface AuthState {
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  setAuth: (user: User, accessToken: string, refreshToken: string) => void;
  logout: () => void;
  updateUser: (user: User) => void;
  ensureFreshToken: () => Promise<string | null>;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,

      setAuth: (user, accessToken, refreshToken) => {
        localStorage.setItem('access_token', accessToken);
        localStorage.setItem('refresh_token', refreshToken);
        set({ 
          user, 
          accessToken, 
          refreshToken, 
          isAuthenticated: true 
        });
      },

      logout: () => {
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        set({ 
          user: null, 
          accessToken: null, 
          refreshToken: null, 
          isAuthenticated: false 
        });
      },

      updateUser: (user) => {
        set({ user });
      },

      ensureFreshToken: async () => {
        const state = get();
        if (!state.refreshToken) {
          return null;
        }

        try {
          // Try to refresh the token
          const response = await authService.refreshToken(state.refreshToken);
          
          // Update the store with new tokens
          localStorage.setItem('access_token', response.access_token);
          localStorage.setItem('refresh_token', response.refresh_token);
          set({ 
            accessToken: response.access_token, 
            refreshToken: response.refresh_token 
          });

          return response.access_token;
        } catch (error) {
          // If refresh fails, logout
          get().logout();
          return null;
        }
      },
    }),
    {
      name: 'dojo-auth',
    }
  )
);