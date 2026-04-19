import { Alert, Button, Card, Form, Input, Typography, message } from 'antd';
import { useMutation } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { login } from '@/services/auth';
import { useAuthStore } from '@/features/auth/store';

type LoginForm = {
  username: string;
  password: string;
};

export const LoginPage = () => {
  const navigate = useNavigate();
  const setSession = useAuthStore((state) => state.setSession);

  const loginMutation = useMutation({
    mutationFn: login,
    onSuccess: (data) => {
      setSession({
        accessToken: data.accessToken,
        refreshToken: data.refreshToken,
        user: data.user
      });
      void navigate('/');
    },
    onError: (error) => {
      const detail =
        error instanceof Error ? error.message : '未知错误';
      message.error(`登录失败：${detail}`);
    }
  });

  const onFinish = (values: LoginForm) => {
    loginMutation.mutate(values);
  };

  return (
    <div
      style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: 'linear-gradient(120deg, #f6f9fc 0%, #eaf2ff 100%)',
        padding: 16
      }}
    >
      <Card style={{ width: 360 }}>
        <Typography.Title level={3}>kbManage 登录</Typography.Title>
        <Alert
          type="info"
          showIcon
          style={{ marginBottom: 16 }}
          message="当前平台支持本地账号与外部身份并存，默认登录方式可在身份治理中心切换。"
        />
        <Form<LoginForm> layout="vertical" onFinish={onFinish}>
          <Form.Item
            label="用户名"
            name="username"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input autoComplete="username" />
          </Form.Item>
          <Form.Item
            label="密码"
            name="password"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password autoComplete="current-password" />
          </Form.Item>
          <Button
            type="primary"
            htmlType="submit"
            loading={loginMutation.isPending}
            block
          >
            登录
          </Button>
        </Form>
      </Card>
    </div>
  );
};
