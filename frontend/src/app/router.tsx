import { createBrowserRouter, Navigate } from 'react-router-dom';
import { AppLayout, Home } from '@/app/App';
import { ProtectedRoute } from '@/app/ProtectedRoute';
import { LoginPage } from '@/features/auth/pages/LoginPage';
import { ClusterOverviewPage } from '@/features/clusters/pages/ClusterOverviewPage';
import { ProjectPage } from '@/features/projects/pages/ProjectPage';
import { ResourceListPage } from '@/features/resources/pages/ResourceListPage';
import { WorkspacePage } from '@/features/workspaces/pages/WorkspacePage';
import { AuditEventPage } from '@/features/audit/pages/AuditEventPage';

export const router = createBrowserRouter([
  {
    path: '/login',
    element: <LoginPage />
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
            element: <ClusterOverviewPage />
          },
          {
            path: '/resources',
            element: <ResourceListPage />
          },
          {
            path: '/workspaces',
            element: <WorkspacePage />
          },
          {
            path: '/projects',
            element: <ProjectPage />
          },
          {
            path: '/audit-events',
            element: <AuditEventPage />
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
