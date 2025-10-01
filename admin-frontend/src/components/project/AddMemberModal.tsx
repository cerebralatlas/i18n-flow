import React, { useState, useEffect } from 'react';
import { Modal, Form, Select, message } from 'antd';
import { useTranslation } from 'react-i18next';
import { projectMemberService, ProjectMemberInfo, AddProjectMemberRequest } from '../../services/projectMemberService';
import { userService, User } from '../../services/userService';

const { Option } = Select;

interface AddMemberModalProps {
  visible: boolean;
  projectId: number;
  existingMembers: ProjectMemberInfo[];
  onCancel: () => void;
  onSuccess: () => void;
}

const AddMemberModal: React.FC<AddMemberModalProps> = ({
  visible,
  projectId,
  existingMembers,
  onCancel,
  onSuccess,
}) => {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [users, setUsers] = useState<User[]>([]);
  const [fetchingUsers, setFetchingUsers] = useState(false);

  // 获取可添加的用户列表
  const fetchUsers = async () => {
    try {
      setFetchingUsers(true);
      const response = await userService.getUsers(1, 100); // 获取前100个用户
      
      // 过滤掉已经是项目成员的用户
      const existingUserIds = existingMembers.map(member => member.user_id);
      const availableUsers = response.data.filter(user => 
        !existingUserIds.includes(user.id) && user.status === 'active'
      );
      
      setUsers(availableUsers);
    } catch (error) {
      message.error(t('projectMember.messages.fetchUsersError'));
    } finally {
      setFetchingUsers(false);
    }
  };

  useEffect(() => {
    if (visible) {
      fetchUsers();
    }
  }, [visible, existingMembers]);

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      setLoading(true);
      
      const memberData: AddProjectMemberRequest = {
        user_id: values.user_id,
        role: values.role,
      };

      await projectMemberService.addMember(projectId, memberData);
      message.success(t('projectMember.messages.addSuccess'));
      form.resetFields();
      onSuccess();
    } catch (error: any) {
      if (error.response?.data?.message) {
        message.error(error.response.data.message);
      } else {
        message.error(t('projectMember.messages.addError'));
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
      title={t('projectMember.addMember')}
      open={visible}
      onOk={handleSubmit}
      onCancel={handleCancel}
      confirmLoading={loading}
      okText={t('common.buttons.add')}
      cancelText={t('common.buttons.cancel')}
      destroyOnClose
    >
      <Form
        form={form}
        layout="vertical"
        requiredMark={false}
      >
        <Form.Item
          name="user_id"
          label={t('projectMember.form.user')}
          rules={[
            { required: true, message: t('projectMember.form.userRequired') },
          ]}
        >
          <Select
            placeholder={t('projectMember.form.userPlaceholder')}
            loading={fetchingUsers}
            showSearch
            filterOption={(input, option) =>
              option?.children?.toString().toLowerCase().includes(input.toLowerCase()) ||
              option?.value?.toString().includes(input)
            }
          >
            {users.map(user => (
              <Option key={user.id} value={user.id}>
                {user.username} ({user.email})
              </Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          name="role"
          label={t('projectMember.form.role')}
          rules={[
            { required: true, message: t('projectMember.form.roleRequired') },
          ]}
          initialValue="viewer"
        >
          <Select placeholder={t('projectMember.form.rolePlaceholder')}>
            <Option value="owner">{t('projectMember.roles.owner')}</Option>
            <Option value="editor">{t('projectMember.roles.editor')}</Option>
            <Option value="viewer">{t('projectMember.roles.viewer')}</Option>
          </Select>
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default AddMemberModal;
