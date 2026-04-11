import '@testing-library/jest-dom/vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { useAuthStore } from '@/features/auth/store';
import { LoginPage } from '@/features/auth/pages/LoginPage';
import { login } from '@/services/auth';

vi.mock('@/services/auth', () => ({
  login: vi.fn()
}));

const installAntdDomShims = () => {
  const nativeGetComputedStyle = window.getComputedStyle.bind(window);
  Object.defineProperty(window, 'getComputedStyle', {
    writable: true,
    value: ((elt: Element) => nativeGetComputedStyle(elt)) as typeof window.getComputedStyle
  });

  if (!window.matchMedia) {
    Object.defineProperty(window, 'matchMedia', {
      writable: true,
      value: vi.fn().mockImplementation((query: string) => ({
        matches: false,
        media: query,
        onchange: null,
        addListener: vi.fn(),
        removeListener: vi.fn(),
        addEventListener: vi.fn(),
        removeEventListener: vi.fn(),
        dispatchEvent: vi.fn()
      }))
    });
  }

  if (!window.scrollTo) {
    Object.defineProperty(window, 'scrollTo', {
      writable: true,
      value: vi.fn()
    });
  }

  if (!globalThis.ResizeObserver) {
    class ResizeObserverMock {
      observe() {}

      unobserve() {}

      disconnect() {}
    }
    vi.stubGlobal('ResizeObserver', ResizeObserverMock);
  }

  if (!window.requestAnimationFrame) {
    Object.defineProperty(window, 'requestAnimationFrame', {
      writable: true,
      value: (callback: FrameRequestCallback) => window.setTimeout(() => callback(0), 16)
    });
  }

  if (!window.cancelAnimationFrame) {
    Object.defineProperty(window, 'cancelAnimationFrame', {
      writable: true,
      value: (id: number) => window.clearTimeout(id)
    });
  }
};

const renderPage = () => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false }
    }
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={['/login']}>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/" element={<div>HOME_PAGE</div>} />
        </Routes>
      </MemoryRouter>
    </QueryClientProvider>
  );
};

describe('LoginPage', () => {
  beforeAll(() => {
    installAntdDomShims();
  });

  beforeEach(() => {
    window.sessionStorage.clear();
    useAuthStore.setState({
      accessToken: null,
      refreshToken: null,
      user: null,
      isAuthenticated: false
    });
    vi.clearAllMocks();
  });

  it('submits credentials and stores session on success', async () => {
    vi.mocked(login).mockResolvedValue({
      accessToken: 'access-token',
      refreshToken: 'refresh-token',
      expiresIn: 3600,
      user: {
        id: 'u-1',
        username: 'alice',
        displayName: 'Alice'
      }
    });

    renderPage();

    fireEvent.change(screen.getByLabelText('用户名'), {
      target: { value: 'alice' }
    });
    fireEvent.change(screen.getByLabelText('密码'), {
      target: { value: 'password-123' }
    });
    fireEvent.click(screen.getByRole('button', { name: /登\s*录/ }));

    await waitFor(() => {
      expect(login).toHaveBeenCalledWith(
        {
          username: 'alice',
          password: 'password-123'
        },
        expect.anything()
      );
    });

    await waitFor(() => {
      expect(screen.getByText('HOME_PAGE')).toBeInTheDocument();
    });

    expect(useAuthStore.getState().accessToken).toBe('access-token');
    expect(useAuthStore.getState().refreshToken).toBe('refresh-token');
    expect(useAuthStore.getState().isAuthenticated).toBe(true);
    expect(useAuthStore.getState().user?.username).toBe('alice');
  });
});
