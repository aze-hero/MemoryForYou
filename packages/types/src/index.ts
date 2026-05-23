export interface User {
  id: string;
  email: string;
  name: string;
  avatarUrl: string;
  provider: 'google' | 'github';
  createdAt: string;
}

export interface MemorySpace {
  id: string;
  userId: string;
  title: string;
  description: string;
  coverImage: string;
  theme: 'default' | 'warm' | 'ocean' | 'forest' | 'sunset';
  memoryCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface Memory {
  id: string;
  userId: string;
  spaceId: string;
  imageUrl: string;
  thumbnailUrl: string;
  title: string;
  content: string;
  aiContent: string;
  aiModel: string;
  location: string;
  memoryDate: string;
  sortOrder: number;
  createdAt: string;
  updatedAt: string;
}

export interface ApiResponse<T = unknown> {
  code: number;
  message: string;
  data: T;
}

export interface PaginatedData<T> {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
}

export interface AiGenerateRequest {
  text: string;
  style?: 'warm' | 'poetic' | 'cinematic';
  memoryId?: string;
}

export interface AiGenerateResponse {
  content: string;
  model: string;
  tokensUsed: number;
}
