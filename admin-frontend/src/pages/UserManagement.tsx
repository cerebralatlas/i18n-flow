import React, { useState, useEffect } from 'react';
import {
  Table,
  Button,
  Space,
  Input,
  Tag,
  Modal,
  message,
  Popconfirm,
  Card,
  Row,
  Col,
  Typography,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  SearchOutlined,
  UserOutlined,
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { userService, User } from '../services/userService';
import { useAuthStore } from '../stores/authStore';
import CreateUserModal from '../components/user/CreateUserModal';
import EditUserModal from '../components/user/EditUserModal';

const { Title } = Typography;
const { Search } = Input;

const UserManagement: React.FC = () => {
  const { t } = useTranslation();
  const { user, loading: authLoading } = useAuthStore();
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [keyword, setKeyword] = useState('');
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [editModalVisible, setEditModalVisible] = useState(false);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);

  // 获取用户列表
  const fetchUsers = async () => {
    try {
      setLoading(true);
      const response = await userService.getUsers(currentPage, pageSize, keyword);
      setUsers(response.data);
      setTotal(response.meta.total_count);
    } catch (error) {
      message.error(t('userManagement.messages.fetchError'));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (!authLoading && user?.role === 'admin') {
      fetchUsers();
    }
  }, [currentPage, pageSize, keyword, authLoading, user?.role]);

  // 处理搜索
  const handleSearch = (value: string) => {
    setKeyword(value);
    setCurrentPage(1);
  };

  // 处理创建用户
  const handleCreateUser = () => {
    setCreateModalVisible(true);
  };

  // 处理编辑用户
  const handleEditUser = (user: User) => {
    setSelectedUser(user);
    setEditModalVisible(true);
  };

  // 处理删除用户
  const handleDeleteUser = async (userId: number) => {
    try {
      await userService.deleteUser(userId);
      message.success(t('userManagement.messages.deleteSuccess'));
      fetchUsers();
    } catch (error) {
      message.error(t('userManagement.messages.deleteError'));
    }
  };

  // 获取角色标签颜色
  const getRoleTagColor = (role: string) => {
    switch (role) {
      case 'admin':
        return 'red';
      case 'member':
        return 'blue';
      case 'viewer':
        return 'green';
      default:
        return 'default';
    }
  };

  // 获取状态标签颜色
  const getStatusTagColor = (status: string) => {
    return status === 'active' ? 'success' : 'error';
  };

  const columns = [
    {
      title: t('userManagement.table.id'),
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: t('userManagement.table.username'),
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
      title: t('userManagement.table.email'),
      dataIndex: 'email',
      key: 'email',
    },
    {
      title: t('userManagement.table.role'),
      dataIndex: 'role',
      key: 'role',
      render: (role: string) => (
        <Tag color={getRoleTagColor(role)}>
          {t(`userManagement.roles.${role}`)}
        </Tag>
      ),
    },
    {
      title: t('userManagement.table.status'),
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={getStatusTagColor(status)}>
          {t(`userManagement.status.${status}`)}
        </Tag>
      ),
    },
    {
      title: t('userManagement.table.createdAt'),
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => new Date(date).toLocaleDateString(),
    },
    {
      title: t('userManagement.table.actions'),
      key: 'actions',
      render: (_, record: User) => (
        <Space>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEditUser(record)}
            size="small"
          >
            {t('common.buttons.edit')}
          </Button>
          <Popconfirm
            title={t('userManagement.confirmDelete')}
            onConfirm={() => handleDeleteUser(record.id)}
            okText={t('common.buttons.yes')}
            cancelText={t('common.buttons.no')}
          >
            <Button
              type="link"
              danger
              icon={<DeleteOutlined />}
              size="small"
              disabled={record.role === 'admin'}
            >
              {t('common.buttons.delete')}
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  // 如果还在加载用户信息，显示加载状态
  if (authLoading) {
    return (
      <Card>
        <div style={{ textAlign: 'center', padding: '50px 0' }}>
          <div>Loading...</div>
        </div>
      </Card>
    );
  }

  // 如果不是管理员，显示无权限提示
  if (user?.role !== 'admin') {
    return (
      <Card>
        <div style={{ textAlign: 'center', padding: '50px 0' }}>
          <Title level={3}>{t('common.noPermission')}</Title>
          <p>{t('userManagement.adminRequired')}</p>
        </div>
      </Card>
    );
  }

  return (
    <div>
      <Card>
        <Row justify="space-between" align="middle" style={{ marginBottom: 16 }}>
          <Col>
            <Title level={2} style={{ margin: 0 }}>
              {t('userManagement.title')}
            </Title>
          </Col>
          <Col>
            <Space>
              <Search
                placeholder={t('userManagement.searchPlaceholder')}
                allowClear
                onSearch={handleSearch}
                style={{ width: 300 }}
                enterButton={<SearchOutlined />}
              />
              <Button
                type="primary"
                icon={<PlusOutlined />}
                onClick={handleCreateUser}
              >
                {t('userManagement.createUser')}
              </Button>
            </Space>
          </Col>
        </Row>

        <Table
          columns={columns}
          dataSource={users}
          rowKey="id"
          loading={loading}
          pagination={{
            current: currentPage,
            pageSize: pageSize,
            total: total,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total, range) =>
              t('common.pagination.total', { total, start: range[0], end: range[1] }),
            onChange: (page, size) => {
              setCurrentPage(page);
              setPageSize(size || 10);
            },
          }}
        />
      </Card>

      {/* 创建用户模态框 */}
      <CreateUserModal
        visible={createModalVisible}
        onCancel={() => setCreateModalVisible(false)}
        onSuccess={() => {
          setCreateModalVisible(false);
          fetchUsers();
        }}
      />

      {/* 编辑用户模态框 */}
      <EditUserModal
        visible={editModalVisible}
        user={selectedUser}
        onCancel={() => {
          setEditModalVisible(false);
          setSelectedUser(null);
        }}
        onSuccess={() => {
          setEditModalVisible(false);
          setSelectedUser(null);
          fetchUsers();
        }}
      />
    </div>
  );
};

export default UserManagement;
