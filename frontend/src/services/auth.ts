export type LoginRequest = {
  username: string;
  password: string;
};

export type AuthUser = {
  id: string;
  username: string;
  displayName?: string;
};

export type LoginResponse = {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
  user: AuthUser;
};

class ApiError extends Error {
  status: number;

  constructor(status: number, message: string) {
    super(message);
    this.status = status;
  }
}

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1';

const fetchJSON = async <T>(path: string, init: RequestInit): Promise<T> => {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    headers: {
      'Content-Type': 'application/json',
      ...(init.headers || {})
    },
    ...init
  });

  if (!response.ok) {
    const text = await response.text();
    throw new ApiError(response.status, text || `Request failed with status ${response.status}`);
  }

  return (await response.json()) as T;
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
