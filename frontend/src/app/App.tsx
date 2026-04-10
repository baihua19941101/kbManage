import { Button, Layout, Typography } from 'antd';
import { Outlet, useNavigate } from 'react-router-dom';
import { AuthorizedMenu } from '@/app/AuthorizedMenu';
import { useAuthStore } from '@/features/auth/store';

const { Header, Content, Sider } = Layout;

const HomePage = () => (
  <div style={{ padding: 24, background: '#fff', borderRadius: 8 }}>
    <Typography.Title level={4} style={{ marginTop: 0 }}>
      平台首页
    </Typography.Title>
    <Typography.Paragraph style={{ marginBottom: 0 }}>
      前端基础骨架已就绪，可继续接入集群、工作空间和资源模块。
    </Typography.Paragraph>
  </div>
);

export const AppLayout = () => {
  const navigate = useNavigate();
  const user = useAuthStore((state) => state.user);
  const clearSession = useAuthStore((state) => state.clearSession);

  const onLogout = () => {
    clearSession();
    void navigate('/login');
  };

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          background: '#0b1f33'
        }}
      >
        <Typography.Title level={4} style={{ color: '#fff', margin: 0 }}>
          kbManage
        </Typography.Title>
        <div style={{ display: 'flex', gap: 12, alignItems: 'center' }}>
          <Typography.Text style={{ color: '#d6e4ff' }}>
            {user?.displayName ?? user?.username}
          </Typography.Text>
          <Button onClick={onLogout}>退出登录</Button>
        </div>
      </Header>
      <Layout>
        <Sider width={220} theme="light">
          <AuthorizedMenu />
        </Sider>
        <Content style={{ padding: 24 }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
};

export const Home = HomePage;
