import { createBrowserRouter, Navigate } from 'react-router-dom';
import { AppLayout, Home } from '@/app/App';
import { ProtectedRoute } from '@/app/ProtectedRoute';

export const router = createBrowserRouter([
  {
    path: '/login',
    lazy: async () => {
      const { LoginPage } = await import('@/features/auth/pages/LoginPage');
      return { Component: LoginPage };
    }
  },
  {
    element: <ProtectedRoute />,
    children: [
      {
        path: '/',
        element: <AppLayout />,
        children: [
          {
            index: true,
            element: <Home />
          },
          {
            path: '/clusters',
            lazy: async () => {
              const { ClusterOverviewPage } = await import(
                '@/features/clusters/pages/ClusterOverviewPage'
              );
              return { Component: ClusterOverviewPage };
            }
          },
          {
            path: '/resources',
            lazy: async () => {
              const { ResourceListPage } = await import(
                '@/features/resources/pages/ResourceListPage'
              );
              return { Component: ResourceListPage };
            }
          },
          {
            path: '/workspaces',
            lazy: async () => {
              const { WorkspacePage } = await import('@/features/workspaces/pages/WorkspacePage');
              return { Component: WorkspacePage };
            }
          },
          {
            path: '/projects',
            lazy: async () => {
              const { ProjectPage } = await import('@/features/projects/pages/ProjectPage');
              return { Component: ProjectPage };
            }
          },
          {
            path: '/audit-events',
            lazy: async () => {
              const { AuditEventPage } = await import('@/features/audit/pages/AuditEventPage');
              return { Component: AuditEventPage };
            }
          }
        ]
      }
    ]
  },
  {
    path: '*',
    element: <Navigate to="/" replace />
  }
]);
