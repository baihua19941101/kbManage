import { useMutation } from '@tanstack/react-query';
import { Alert, Button, Card, Layout, Typography, message } from 'antd';
import { Outlet, useNavigate } from 'react-router-dom';
import { AuthorizedMenu } from '@/app/AuthorizedMenu';
import { normalizeErrorMessage } from '@/app/queryClient';
import { RoleBindingForm } from '@/features/auth/components/RoleBindingForm';
import { useAuthStore } from '@/features/auth/store';
import { createRoleBinding } from '@/services/roleBindings';

const { Header, Content, Sider } = Layout;

const HomePage = () => {
  const user = useAuthStore((state) => state.user);

  const canManageBindings = (user?.platformRoles || []).includes('platform-admin');

  const createBindingMutation = useMutation({
    mutationFn: createRoleBinding,
    onSuccess: () => {
      message.success('角色绑定创建成功');
    },
    onError: (error) => {
      message.error(normalizeErrorMessage(error, '角色绑定创建失败'));
    }
  });

  return (
    <div style={{ padding: 24, background: '#fff', borderRadius: 8 }}>
      <Typography.Title level={4} style={{ marginTop: 0 }}>
        平台首页
      </Typography.Title>
      <Typography.Paragraph>
        前端基础骨架已就绪，可继续接入集群、工作空间和资源模块。
      </Typography.Paragraph>
      {canManageBindings ? (
        <Card size="small" title="角色绑定管理">
          <RoleBindingForm
            loading={createBindingMutation.isPending}
            onSubmit={(payload) => {
              createBindingMutation.mutate(payload);
            }}
          />
        </Card>
      ) : (
        <Alert
          type="info"
          showIcon
          message="当前账号无角色绑定管理权限"
          description="仅 platform-admin 可在首页执行角色绑定操作。"
        />
      )}
    </div>
  );
};

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
