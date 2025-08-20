import { create } from 'zustand';
import axios from 'axios';
import { message } from 'antd';
import api from '../utils/api';
import { ApiResponse } from '../types/api';

interface User {
  id: number;
  username: string;
}

interface AuthState {
  user: User | null;
  loading: boolean;
  isAuthenticated: boolean;
  setUser: (user: User | null) => void;
  setLoading: (loading: boolean) => void;
  checkAuthStatus: () => Promise<void>;
  login: (username: string, password: string) => Promise<void>;
  logout: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  loading: true,
  isAuthenticated: false,

  setUser: (user) => set({ user, isAuthenticated: !!user }),
  setLoading: (loading) => set({ loading }),

  checkAuthStatus: async () => {
    try {
      set({ loading: true });
      const token = localStorage.getItem('authToken');

      if (!token) {
        set({ loading: false });
        return;
      }

      axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
      const response = await axios.get('/api/user/info');
      set({ user: response.data, isAuthenticated: true });
    } catch (error) {
      console.error('Validate user session failed:', error);
      localStorage.removeItem('authToken');
      delete axios.defaults.headers.common['Authorization'];
      set({ user: null, isAuthenticated: false });
    } finally {
      set({ loading: false });
    }
  },

  login: async (username: string, password: string) => {
    try {
      const response: ApiResponse<{ token: string; user: User }> = await api.post('/api/login', {
        username,
        password,
      });

      const { token, user: userData } = response.data;

      localStorage.setItem('authToken', token);
      axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
      set({ user: userData, isAuthenticated: true });
    } catch (error) {
      console.error('Login failed:', error);
      throw error;
    }
  },

  logout: () => {
    localStorage.removeItem('authToken');
    delete axios.defaults.headers.common['Authorization'];
    message.success('Logout successfully');
    set({ user: null, isAuthenticated: false });
  },
})); 