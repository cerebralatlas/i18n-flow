import React from 'react';
import { Modal, Form, Input, Select, message, Button, Space, Tooltip } from 'antd';
import { CopyOutlined, ReloadOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { userService, CreateUserRequest } from '../../services/userService';

const { Option } = Select;

interface CreateUserModalProps {
  visible: boolean;
  onCancel: () => void;
  onSuccess: () => void;
}

const CreateUserModal: React.FC<CreateUserModalProps> = ({
  visible,
  onCancel,
  onSuccess,
}) => {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const [loading, setLoading] = React.useState(false);
  const [generatedPassword, setGeneratedPassword] = React.useState('');

  // 生成随机密码
  const generatePassword = () => {
    const length = 12;
    const charset = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*';
    let password = '';
    
    // 确保密码包含至少一个大写字母、小写字母、数字和特殊字符
    const lowercase = 'abcdefghijklmnopqrstuvwxyz';
    const uppercase = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';
    const numbers = '0123456789';
    const special = '!@#$%^&*';
    
    password += lowercase[Math.floor(Math.random() * lowercase.length)];
    password += uppercase[Math.floor(Math.random() * uppercase.length)];
    password += numbers[Math.floor(Math.random() * numbers.length)];
    password += special[Math.floor(Math.random() * special.length)];
    
    // 填充剩余长度
    for (let i = 4; i < length; i++) {
      password += charset[Math.floor(Math.random() * charset.length)];
    }
    
    // 打乱密码字符顺序
    return password.split('').sort(() => Math.random() - 0.5).join('');
  };

  // 处理生成密码
  const handleGeneratePassword = () => {
    const newPassword = generatePassword();
    setGeneratedPassword(newPassword);
    form.setFieldsValue({ password: newPassword });
    message.success(t('userManagement.messages.passwordGenerated'));
  };

  // 复制密码到剪贴板
  const handleCopyPassword = async () => {
    const password = form.getFieldValue('password');
    if (password) {
      try {
        await navigator.clipboard.writeText(password);
        message.success(t('userManagement.messages.passwordCopied'));
      } catch (error) {
        // 如果剪贴板API不可用，使用传统方法
        const textArea = document.createElement('textarea');
        textArea.value = password;
        document.body.appendChild(textArea);
        textArea.select();
        document.execCommand('copy');
        document.body.removeChild(textArea);
        message.success(t('userManagement.messages.passwordCopied'));
      }
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      setLoading(true);
      
      const userData: CreateUserRequest = {
        username: values.username,
        email: values.email,
        password: values.password,
        role: values.role,
      };

      await userService.createUser(userData);
      
      // 如果使用了生成的密码，显示密码信息供管理员复制
      if (generatedPassword) {
        Modal.info({
          title: t('userManagement.messages.userCreatedTitle'),
          content: (
            <div>
              <p>{t('userManagement.messages.userCreatedContent')}</p>
              <div style={{ 
                background: '#f5f5f5', 
                padding: '12px', 
                borderRadius: '4px', 
                marginTop: '12px',
                fontFamily: 'monospace'
              }}>
                <div><strong>{t('userManagement.form.username')}:</strong> {values.username}</div>
                <div><strong>{t('userManagement.form.password')}:</strong> {generatedPassword}</div>
                <div><strong>{t('userManagement.form.email')}:</strong> {values.email}</div>
              </div>
              <p style={{ marginTop: '12px', color: '#666' }}>
                {t('userManagement.messages.passwordReminder')}
              </p>
            </div>
          ),
          width: 500,
          okText: t('common.buttons.ok'),
        });
      } else {
        message.success(t('userManagement.messages.createSuccess'));
      }
      
      form.resetFields();
      setGeneratedPassword('');
      onSuccess();
    } catch (error: any) {
      if (error.response?.data?.message) {
        message.error(error.response.data.message);
      } else {
        message.error(t('userManagement.messages.createError'));
      }
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = () => {
    form.resetFields();
    setGeneratedPassword('');
    onCancel();
  };

  return (
    <Modal
      title={t('userManagement.createUser')}
      open={visible}
      onOk={handleSubmit}
      onCancel={handleCancel}
      confirmLoading={loading}
      okText={t('common.buttons.create')}
      cancelText={t('common.buttons.cancel')}
      destroyOnClose
    >
      <Form
        form={form}
        layout="vertical"
        requiredMark={false}
      >
        <Form.Item
          name="username"
          label={t('userManagement.form.username')}
          rules={[
            { required: true, message: t('userManagement.form.usernameRequired') },
            { min: 3, message: t('userManagement.form.usernameMinLength') },
            { max: 50, message: t('userManagement.form.usernameMaxLength') },
          ]}
        >
          <Input placeholder={t('userManagement.form.usernamePlaceholder')} />
        </Form.Item>

        <Form.Item
          name="email"
          label={t('userManagement.form.email')}
          rules={[
            { required: true, message: t('userManagement.form.emailRequired') },
            { type: 'email', message: t('userManagement.form.emailInvalid') },
          ]}
        >
          <Input placeholder={t('userManagement.form.emailPlaceholder')} />
        </Form.Item>

        <Form.Item
          name="password"
          label={
            <Space>
              {t('userManagement.form.password')}
              <Tooltip title={t('userManagement.form.passwordTooltip')}>
                <Button
                  type="link"
                  size="small"
                  icon={<ReloadOutlined />}
                  onClick={handleGeneratePassword}
                  style={{ padding: 0 }}
                >
                  {t('userManagement.form.generatePassword')}
                </Button>
              </Tooltip>
            </Space>
          }
          rules={[
            { required: true, message: t('userManagement.form.passwordRequired') },
            { min: 6, message: t('userManagement.form.passwordMinLength') },
          ]}
        >
          <Input.Password 
            placeholder={t('userManagement.form.passwordPlaceholder')}
            addonAfter={
              <Tooltip title={t('userManagement.form.copyPassword')}>
                <Button
                  type="text"
                  size="small"
                  icon={<CopyOutlined />}
                  onClick={handleCopyPassword}
                  style={{ border: 'none', padding: '0 4px' }}
                />
              </Tooltip>
            }
          />
        </Form.Item>

        <Form.Item
          name="role"
          label={t('userManagement.form.role')}
          rules={[
            { required: true, message: t('userManagement.form.roleRequired') },
          ]}
          initialValue="member"
        >
          <Select placeholder={t('userManagement.form.rolePlaceholder')}>
            <Option value="admin">{t('userManagement.roles.admin')}</Option>
            <Option value="member">{t('userManagement.roles.member')}</Option>
            <Option value="viewer">{t('userManagement.roles.viewer')}</Option>
          </Select>
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default CreateUserModal;
