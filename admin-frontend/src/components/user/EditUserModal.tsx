import React, { useEffect } from 'react';
import { Modal, Form, Input, Select, message } from 'antd';
import { useTranslation } from 'react-i18next';
import { userService, User, UpdateUserRequest } from '../../services/userService';

const { Option } = Select;

interface EditUserModalProps {
  visible: boolean;
  user: User | null;
  onCancel: () => void;
  onSuccess: () => void;
}

const EditUserModal: React.FC<EditUserModalProps> = ({
  visible,
  user,
  onCancel,
  onSuccess,
}) => {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const [loading, setLoading] = React.useState(false);

  useEffect(() => {
    if (visible && user) {
      form.setFieldsValue({
        username: user.username,
        email: user.email,
        role: user.role,
        status: user.status,
      });
    }
  }, [visible, user, form]);

  const handleSubmit = async () => {
    if (!user) return;

    try {
      const values = await form.validateFields();
      setLoading(true);
      
      const userData: UpdateUserRequest = {
        username: values.username,
        email: values.email,
        role: values.role,
        status: values.status,
      };

      await userService.updateUser(user.id, userData);
      message.success(t('userManagement.messages.updateSuccess'));
      onSuccess();
    } catch (error: any) {
      if (error.response?.data?.message) {
        message.error(error.response.data.message);
      } else {
        message.error(t('userManagement.messages.updateError'));
      }
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = () => {
    form.resetFields();
    onCancel();
  };

  return (
    <Modal
      title={t('userManagement.editUser')}
      open={visible}
      onOk={handleSubmit}
      onCancel={handleCancel}
      confirmLoading={loading}
      okText={t('common.buttons.save')}
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
          name="role"
          label={t('userManagement.form.role')}
          rules={[
            { required: true, message: t('userManagement.form.roleRequired') },
          ]}
        >
          <Select placeholder={t('userManagement.form.rolePlaceholder')}>
            <Option value="admin">{t('userManagement.roles.admin')}</Option>
            <Option value="member">{t('userManagement.roles.member')}</Option>
            <Option value="viewer">{t('userManagement.roles.viewer')}</Option>
          </Select>
        </Form.Item>

        <Form.Item
          name="status"
          label={t('userManagement.form.status')}
          rules={[
            { required: true, message: t('userManagement.form.statusRequired') },
          ]}
        >
          <Select placeholder={t('userManagement.form.statusPlaceholder')}>
            <Option value="active">{t('userManagement.status.active')}</Option>
            <Option value="disabled">{t('userManagement.status.disabled')}</Option>
          </Select>
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default EditUserModal;
