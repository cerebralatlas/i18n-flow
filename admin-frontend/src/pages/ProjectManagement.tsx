import React, { useState, useEffect, useCallback } from "react";
import {
  Table,
  Button,
  Input,
  Space,
  Modal,
  Form,
  message,
  Popconfirm,
  Typography,
} from "antd";
import {
  PlusOutlined,
  SearchOutlined,
  EditOutlined,
  DeleteOutlined,
  EyeOutlined,
} from "@ant-design/icons";
import axios from "axios";
import ErrorAlert from "../components/ErrorAlert";
import ProjectDetailModal from "../components/ProjectDetailModal";
import { Project } from "../types/project";

const { Title } = Typography;

const ProjectManagement: React.FC = () => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [total, setTotal] = useState<number>(0);
  const [page, setPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(10);
  const [keyword, setKeyword] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);
  const [modalVisible, setModalVisible] = useState<boolean>(false);
  const [detailVisible, setDetailVisible] = useState<boolean>(false);
  const [detailLoading, setDetailLoading] = useState<boolean>(false);
  const [currentProject, setCurrentProject] = useState<Project | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [form] = Form.useForm();

  // 获取项目列表
  const fetchProjects = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await axios.get("/api/projects", {
        params: {
          page,
          page_size: pageSize,
          keyword,
        },
      });

      setProjects(response.data.data);
      setTotal(response.data.total);
    } catch (error) {
      console.error("获取项目列表失败:", error);
      setError("获取项目列表失败");
    } finally {
      setLoading(false);
    }
  }, [page, pageSize, keyword]);

  // 初始加载和条件变化时获取数据
  useEffect(() => {
    fetchProjects();
  }, [page, pageSize, keyword, fetchProjects]);

  // 创建或更新项目
  const handleSaveProject = async () => {
    try {
      const values = await form.validateFields();

      if (currentProject) {
        // 更新项目
        await axios.put(`/api/projects/update/${currentProject.id}`, values);
        message.success("项目更新成功");
      } else {
        // 创建项目
        await axios.post("/api/projects", values);
        message.success("项目创建成功");
      }

      setModalVisible(false);
      form.resetFields();
      fetchProjects();
    } catch (error) {
      console.error("保存项目失败:", error);
      message.error("保存项目失败");
    }
  };

  // 删除项目
  const handleDeleteProject = async (id: number) => {
    try {
      await axios.delete(`/api/projects/delete/${id}`);
      message.success("项目删除成功");
      fetchProjects();
    } catch (error) {
      console.error("删除项目失败:", error);
      message.error("删除项目失败");
    }
  };

  // 打开创建项目弹窗
  const showCreateModal = () => {
    setCurrentProject(null);
    form.resetFields();
    setModalVisible(true);
  };

  // 打开编辑项目弹窗
  const showEditModal = async (project: Project) => {
    try {
      // 获取项目详情
      const response = await axios.get(`/api/projects/detail/${project.id}`);
      const projectData = response.data;

      setCurrentProject(projectData);
      form.setFieldsValue({
        name: projectData.name,
        description: projectData.description,
        slug: projectData.slug,
      });
      setModalVisible(true);
      setDetailVisible(false); // 关闭详情弹窗如果已打开
    } catch (error) {
      console.error("获取项目详情失败:", error);
      message.error("获取项目详情失败");
    }
  };

  // 查看项目详情
  const showProjectDetail = async (project: Project) => {
    try {
      setDetailLoading(true);
      setDetailVisible(true);

      // 获取项目详情
      const response = await axios.get(`/api/projects/detail/${project.id}`);
      setCurrentProject(response.data);
    } catch (error) {
      console.error("获取项目详情失败:", error);
      message.error("获取项目详情失败");
    } finally {
      setDetailLoading(false);
    }
  };

  // 表格列定义
  const columns = [
    {
      title: "ID",
      dataIndex: "id",
      key: "id",
      width: 80,
    },
    {
      title: "项目名称",
      dataIndex: "name",
      key: "name",
    },
    {
      title: "项目标识",
      dataIndex: "slug",
      key: "slug",
    },
    {
      title: "描述",
      dataIndex: "description",
      key: "description",
      ellipsis: true,
    },
    {
      title: "状态",
      dataIndex: "status",
      key: "status",
      render: (status: string) => (
        <span style={{ color: status === "active" ? "green" : "orange" }}>
          {status === "active" ? "活跃" : "归档"}
        </span>
      ),
    },
    {
      title: "创建时间",
      dataIndex: "created_at",
      key: "created_at",
    },
    {
      title: "操作",
      key: "action",
      width: 200,
      render: (_: unknown, record: Project) => (
        <Space size="middle">
          <Button
            type="text"
            icon={<EyeOutlined />}
            onClick={() => showProjectDetail(record)}
          >
            详情
          </Button>
          <Button
            type="text"
            icon={<EditOutlined />}
            onClick={() => showEditModal(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除该项目吗？"
            description="删除后将无法恢复！"
            onConfirm={() => handleDeleteProject(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="text" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div className="bg-white p-6 rounded-lg shadow">
      <div className="mb-6 flex justify-between items-center">
        <Title level={3}>项目管理</Title>
        <Space size="middle">
          <Input
            placeholder="搜索项目"
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            prefix={<SearchOutlined />}
            style={{ width: 250 }}
            allowClear
          />
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={showCreateModal}
          >
            创建项目
          </Button>
        </Space>
      </div>

      {error && (
        <div className="mb-4">
          <ErrorAlert message={error} onRetry={fetchProjects} />
        </div>
      )}

      <Table
        columns={columns}
        dataSource={projects}
        rowKey="id"
        loading={loading}
        pagination={{
          current: page,
          pageSize: pageSize,
          total: total,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total) => `共 ${total} 条记录`,
          onChange: (page, pageSize) => {
            setPage(page);
            setPageSize(pageSize);
          },
        }}
      />

      {/* 创建/编辑项目弹窗 */}
      <Modal
        title={currentProject ? "编辑项目" : "创建项目"}
        open={modalVisible}
        onOk={handleSaveProject}
        onCancel={() => setModalVisible(false)}
        okText={currentProject ? "更新" : "创建"}
        cancelText="取消"
        maskClosable={false}
      >
        <Form form={form} layout="vertical" name="projectForm">
          <Form.Item
            name="name"
            label="项目名称"
            rules={[{ required: true, message: "请输入项目名称" }]}
          >
            <Input placeholder="请输入项目名称" />
          </Form.Item>

          <Form.Item name="description" label="项目描述">
            <Input.TextArea rows={4} placeholder="请输入项目描述" />
          </Form.Item>

          <Form.Item
            name="slug"
            label="项目标识"
            help="用于URL友好的标识，留空将根据项目名自动生成"
          >
            <Input placeholder="例如：my-project" />
          </Form.Item>
        </Form>
      </Modal>

      {/* 项目详情弹窗 */}
      <ProjectDetailModal
        project={currentProject}
        visible={detailVisible}
        loading={detailLoading}
        onClose={() => setDetailVisible(false)}
        onEdit={(project) => {
          setDetailVisible(false);
          showEditModal(project);
        }}
      />
    </div>
  );
};

export default ProjectManagement;
