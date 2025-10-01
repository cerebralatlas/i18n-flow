import React, { useState, useEffect } from 'react';
import {
  Table,
  Button,
  Space,
  Tag,
  Modal,
  message,
  Popconfirm,
  Card,
  Row,
  Col,
  Typography,
  Select,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  UserOutlined,
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { projectMemberService, ProjectMemberInfo } from '../../services/projectMemberService';
import { userService, User } from '../../services/userService';
import { useAuthStore } from '../../stores/authStore';
import AddMemberModal from './AddMemberModal';

const { Title } = Typography;
const { Option } = Select;

interface ProjectMemberManagementProps {
  projectId: number;
  projectName: string;
}

const ProjectMemberManagement: React.FC<ProjectMemberManagementProps> = ({
  projectId,
  projectName,
}) => {
  const { t } = useTranslation();
  const { user } = useAuthStore();
  const [members, setMembers] = useState<ProjectMemberInfo[]>([]);
  const [loading, setLoading] = useState(false);
  const [addModalVisible, setAddModalVisible] = useState(false);
  const [editingMember, setEditingMember] = useState<ProjectMemberInfo | null>(null);
  const [newRole, setNewRole] = useState<string>('');

  // 获取项目成员列表
  const fetchMembers = async () => {
    try {
      setLoading(true);
      const memberList = await projectMemberService.getProjectMembers(projectId);
      setMembers(memberList);
    } catch (error) {
      message.error(t('projectMember.messages.fetchError'));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchMembers();
  }, [projectId]);

  // 处理添加成员
  const handleAddMember = () => {
    setAddModalVisible(true);
  };

  // 处理编辑成员角色
  const handleEditRole = (member: ProjectMemberInfo) => {
    setEditingMember(member);
    setNewRole(member.role);
  };

  // 保存角色更改
  const handleSaveRole = async () => {
    if (!editingMember) return;

    try {
      await projectMemberService.updateMemberRole(projectId, editingMember.user_id, {
        role: newRole as 'owner' | 'editor' | 'viewer',
      });
      message.success(t('projectMember.messages.updateSuccess'));
      setEditingMember(null);
      fetchMembers();
    } catch (error) {
      message.error(t('projectMember.messages.updateError'));
    }
  };

  // 处理移除成员
  const handleRemoveMember = async (userId: number) => {
    try {
      await projectMemberService.removeMember(projectId, userId);
      message.success(t('projectMember.messages.removeSuccess'));
      fetchMembers();
    } catch (error: any) {
      if (error.response?.data?.message) {
        message.error(error.response.data.message);
      } else {
        message.error(t('projectMember.messages.removeError'));
      }
    }
  };

  // 获取角色标签颜色
  const getRoleTagColor = (role: string) => {
    switch (role) {
      case 'owner':
        return 'red';
      case 'editor':
        return 'blue';
      case 'viewer':
        return 'green';
      default:
        return 'default';
    }
  };

  // 检查是否可以编辑成员
  const canEditMember = (member: ProjectMemberInfo) => {
    if (user?.role === 'admin') return true;
    if (!user) return false;
    
    // 项目所有者可以编辑所有成员
    const currentUserMember = members.find(m => m.user_id === user.id);
    return currentUserMember?.role === 'owner';
  };

  // 检查是否可以移除成员
  const canRemoveMember = (member: ProjectMemberInfo) => {
    if (member.role === 'owner') return false; // 不能移除所有者
    return canEditMember(member);
  };

  const columns = [
    {
      title: t('projectMember.table.username'),
      dataIndex: 'username',
      key: 'username',
      render: (text: string) => (
        <Space>
          <UserOutlined />
          {text}
        </Space>
      ),
    },
    {
      title: t('projectMember.table.email'),
      dataIndex: 'email',
      key: 'email',
    },
    {
      title: t('projectMember.table.role'),
      dataIndex: 'role',
      key: 'role',
      render: (role: string, record: ProjectMemberInfo) => {
        if (editingMember?.id === record.id) {
          return (
            <Space>
              <Select
                value={newRole}
                onChange={setNewRole}
                style={{ width: 120 }}
                size="small"
              >
                <Option value="owner">{t('projectMember.roles.owner')}</Option>
                <Option value="editor">{t('projectMember.roles.editor')}</Option>
                <Option value="viewer">{t('projectMember.roles.viewer')}</Option>
              </Select>
              <Button type="link" size="small" onClick={handleSaveRole}>
                {t('common.buttons.save')}
              </Button>
              <Button type="link" size="small" onClick={() => setEditingMember(null)}>
                {t('common.buttons.cancel')}
              </Button>
            </Space>
          );
        }
        
        return (
          <Tag color={getRoleTagColor(role)}>
            {t(`projectMember.roles.${role}`)}
          </Tag>
        );
      },
    },
    {
      title: t('projectMember.table.actions'),
      key: 'actions',
      render: (_, record: ProjectMemberInfo) => (
        <Space>
          {canEditMember(record) && (
            <Button
              type="link"
              icon={<EditOutlined />}
              onClick={() => handleEditRole(record)}
            size="small"
          >
            {t('common.buttons.edit')}
          </Button>
          )}
          {canRemoveMember(record) && (
            <Popconfirm
              title={t('projectMember.confirmRemove')}
              onConfirm={() => handleRemoveMember(record.user_id)}
              okText={t('common.buttons.yes')}
              cancelText={t('common.buttons.no')}
            >
              <Button
                type="link"
                danger
                icon={<DeleteOutlined />}
                size="small"
              >
                {t('common.buttons.remove')}
              </Button>
            </Popconfirm>
          )}
        </Space>
      ),
    },
  ];

  return (
    <Card>
      <Row justify="space-between" align="middle" style={{ marginBottom: 16 }}>
        <Col>
          <Title level={4} style={{ margin: 0 }}>
            {t('projectMember.title', { projectName })}
          </Title>
        </Col>
        <Col>
          {(user?.role === 'admin' || members.find(m => m.user_id === user?.id)?.role === 'owner') && (
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={handleAddMember}
            >
              {t('projectMember.addMember')}
            </Button>
          )}
        </Col>
      </Row>

      <Table
        columns={columns}
        dataSource={members}
        rowKey="id"
        loading={loading}
        pagination={false}
        size="small"
      />

      {/* 添加成员模态框 */}
      <AddMemberModal
        visible={addModalVisible}
        projectId={projectId}
        existingMembers={members}
        onCancel={() => setAddModalVisible(false)}
        onSuccess={() => {
          setAddModalVisible(false);
          fetchMembers();
        }}
      />
    </Card>
  );
};

export default ProjectMemberManagement;
