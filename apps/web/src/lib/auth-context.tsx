'use client';

import { createContext, useContext, useState, useEffect, useCallback, type ReactNode } from 'react';
import { authApi, setAccessToken, getAccessToken } from './api';

interface User {
  id: string;
  email: string;
  name: string;
  avatar_url: string;
}

interface AuthState {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (provider: string) => void;
  logout: () => void;
  handleCallback: (provider: string, code: string) => Promise<void>;
}

const AuthContext = createContext<AuthState>({
  user: null,
  isLoading: true,
  isAuthenticated: false,
  login: () => {},
  logout: () => {},
  handleCallback: async () => {},
});

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const token = getAccessToken();
    if (token) {
      authApi.me()
        .then((u) => setUser(u))
        .catch(() => setAccessToken(null))
        .finally(() => setIsLoading(false));
    } else {
      setIsLoading(false);
    }
  }, []);

  const login = useCallback((provider: string) => {
    const redirectUri = `${window.location.origin}/auth/callback?provider=${provider}`;
    const apiBase = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3001';

    switch (provider) {
      case 'google':
        window.location.href = `${apiBase}/api/v1/auth/google?redirect_uri=${encodeURIComponent(redirectUri)}`;
        break;
      case 'github':
        window.location.href = `${apiBase}/api/v1/auth/github?redirect_uri=${encodeURIComponent(redirectUri)}`;
        break;
    }
  }, []);

  const handleCallback = useCallback(async (provider: string, code: string) => {
    const result = await authApi.login(provider, code);
    setAccessToken(result.access_token);
    const u = await authApi.me();
    setUser(u);
  }, []);

  const logout = useCallback(() => {
    setAccessToken(null);
    setUser(null);
  }, []);

  return (
    <AuthContext.Provider
      value={{
        user,
        isLoading,
        isAuthenticated: !!user,
        login,
        logout,
        handleCallback,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}
