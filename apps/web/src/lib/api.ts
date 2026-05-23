const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3001';

interface ApiResponse<T = unknown> {
  code: number;
  message: string;
  data: T;
}

interface PaginatedData<T> {
  items: T[];
  total: number;
  page: number;
  page_size: number;
}

let accessToken: string | null = null;
let refreshPromise: Promise<string | null> | null = null;

export function setAccessToken(token: string | null) {
  accessToken = token;
}

export function getAccessToken() {
  return accessToken;
}

async function refreshAccessToken(): Promise<string | null> {
  try {
    const res = await fetch(`${API_BASE}/api/v1/auth/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
    });
    if (!res.ok) return null;
    const json = await res.json();
    const newToken = json.data?.access_token;
    if (newToken) {
      accessToken = newToken;
      return newToken;
    }
    return null;
  } catch {
    return null;
  }
}

async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const headers: Record<string, string> = {
    ...(options.headers as Record<string, string>),
  };

  if (accessToken) {
    headers['Authorization'] = `Bearer ${accessToken}`;
  }

  if (!(options.body instanceof FormData)) {
    headers['Content-Type'] = 'application/json';
  }

  let res = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers,
  });

  if (res.status === 401 && accessToken) {
    if (!refreshPromise) {
      refreshPromise = refreshAccessToken().finally(() => {
        refreshPromise = null;
      });
    }
    const newToken = await refreshPromise;
    if (newToken) {
      headers['Authorization'] = `Bearer ${newToken}`;
      res = await fetch(`${API_BASE}${endpoint}`, {
        ...options,
        headers,
      });
    } else {
      accessToken = null;
      throw new AuthError('Session expired');
    }
  }

  const json: ApiResponse<T> = await res.json();

  if (json.code !== 0) {
    throw new ApiError(json.code, json.message);
  }

  return json.data;
}

export class ApiError extends Error {
  code: number;
  constructor(code: number, message: string) {
    super(message);
    this.code = code;
  }
}

export class AuthError extends Error {
  constructor(message: string) {
    super(message);
  }
}

// Auth API
export const authApi = {
  login: (provider: string, code: string) =>
    request<{ access_token: string; refresh_token: string; expires_in: number; user: unknown }>(
      '/api/v1/auth/login',
      { method: 'POST', body: JSON.stringify({ provider, code }) }
    ),

  refresh: () =>
    request<{ access_token: string; expires_in: number }>(
      '/api/v1/auth/refresh',
      { method: 'POST' }
    ),

  me: () => request<{ id: string; email: string; name: string; avatar_url: string }>(
    '/api/v1/auth/me'
  ),

  devLogin: (email: string, name: string) =>
    request<{ access_token: string; refresh_token: string; expires_in: number; user: unknown }>(
      '/api/v1/auth/dev-login',
      { method: 'POST', body: JSON.stringify({ email, name }) }
    ),
};

// Space API
export const spaceApi = {
  create: (data: { title: string; description?: string; cover_image?: string; theme?: string }) =>
    request<{ id: string; title: string; description: string; cover_image: string; theme: string; memory_count: number; created_at: string }>(
      '/api/v1/spaces',
      { method: 'POST', body: JSON.stringify(data) }
    ),

  list: (page = 1, pageSize = 20) =>
    request<PaginatedData<{ id: string; title: string; description: string; cover_image: string; theme: string; memory_count: number; created_at: string }>>(
      `/api/v1/spaces?page=${page}&page_size=${pageSize}`
    ),

  get: (id: string) =>
    request<{ id: string; title: string; description: string; cover_image: string; theme: string; created_at: string }>(
      `/api/v1/spaces/${id}`
    ),

  update: (id: string, data: { title?: string; description?: string; cover_image?: string; theme?: string }) =>
    request(`/api/v1/spaces/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),

  delete: (id: string) =>
    request(`/api/v1/spaces/${id}`, { method: 'DELETE' }),
};

// Memory API
export const memoryApi = {
  create: (formData: FormData) =>
    request<{ id: string; title: string; image_url: string; ai_content: string; memory_date: string }>(
      '/api/v1/memories',
      { method: 'POST', body: formData }
    ),

  list: (spaceId: string, page = 1, pageSize = 20) =>
    request<PaginatedData<{ id: string; title: string; image_url: string; thumbnail_url: string; ai_content: string; location: string; memory_date: string; created_at: string }>>(
      `/api/v1/memories?space_id=${spaceId}&page=${page}&page_size=${pageSize}`
    ),

  get: (id: string) =>
    request<{ id: string; title: string; image_url: string; ai_content: string; location: string; memory_date: string }>(
      `/api/v1/memories/${id}`
    ),

  delete: (id: string) =>
    request(`/api/v1/memories/${id}`, { method: 'DELETE' }),
};

// AI API
export const aiApi = {
  generate: (text: string, style?: string) =>
    request<{ content: string; model: string; tokens_used: number }>(
      '/api/v1/ai/generate-memory',
      { method: 'POST', body: JSON.stringify({ text, style }) }
    ),

  regenerate: (memoryId: string, style?: string) =>
    request<{ content: string; model: string; tokens_used: number }>(
      '/api/v1/ai/regenerate',
      { method: 'POST', body: JSON.stringify({ memory_id: memoryId, style }) }
    ),
};
