import { fetchJSON } from '@/services/api/client';

export type LoginRequest = {
  username: string;
  password: string;
};

export type AuthUser = {
  id: string;
  username: string;
  displayName?: string;
  platformRoles?: string[];
};

export type LoginResponse = {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
  user: AuthUser;
};

export const login = async (payload: LoginRequest): Promise<LoginResponse> =>
  fetchJSON<LoginResponse>('/auth/login', {
    method: 'POST',
    body: JSON.stringify(payload)
  });

export const refreshSession = async (
  refreshToken: string
): Promise<LoginResponse> =>
  fetchJSON<LoginResponse>('/auth/refresh', {
    method: 'POST',
    body: JSON.stringify({ refreshToken })
  });
