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
import api from "../utils/api";
import ErrorAlert from "../components/ErrorAlert";
import ProjectDetailModal from "../components/ProjectDetailModal";
import { Project } from "../types/project";
import { ApiResponse, PaginatedResponse } from "../types/api";
import { useTranslation } from "react-i18next";

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
  const { t } = useTranslation();

  // Get project list
  const fetchProjects = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const response: PaginatedResponse<Project> = await api.get("/api/projects", {
        params: {
          page,
          page_size: pageSize,
          keyword,
        },
      });

      // API interceptor returns response.data, so response contains {success, data, meta}
      setProjects(response.data);
      setTotal(response.meta?.total_count || 0);
    } catch (error) {
      console.error("Get project list failed:", error);
      setError(t("project.message.getListFailed"));
    } finally {
      setLoading(false);
    }
  }, [page, pageSize, keyword, t]);

  // Get data when initial load and condition changes
  useEffect(() => {
    fetchProjects();
  }, [page, pageSize, keyword, fetchProjects]);

  // Create or update project
  const handleSaveProject = async () => {
    try {
      const values = await form.validateFields();

      if (currentProject) {
        // Update project
        await api.put(`/api/projects/update/${currentProject.id}`, values);
        message.success(t("project.message.updateSuccess"));
      } else {
        // Create project
        await api.post("/api/projects", values);
        message.success(t("project.message.createSuccess"));
      }

      setModalVisible(false);
      form.resetFields();
      fetchProjects();
    } catch (error) {
      console.error("Save project failed:", error);
      message.error(t("project.message.saveFailed"));
    }
  };

  // Delete project
  const handleDeleteProject = async (id: number) => {
    try {
      await api.delete(`/api/projects/delete/${id}`);
      message.success(t("project.message.deleteSuccess"));
      fetchProjects();
    } catch (error) {
      console.error("Delete project failed:", error);
      message.error(t("project.message.deleteFailed"));
    }
  };

  // Open create project modal
  const showCreateModal = () => {
    setCurrentProject(null);
    form.resetFields();
    setModalVisible(true);
  };

  // Open edit project modal
  const showEditModal = async (project: Project) => {
    try {
      // Get project details
      const response: ApiResponse<Project> = await api.get(`/api/projects/detail/${project.id}`);
      const projectData = response.data; // response contains {success, data, meta}, so response.data is the project

      setCurrentProject(projectData);
      form.setFieldsValue({
        name: projectData.name,
        description: projectData.description,
        slug: projectData.slug,
      });
      setModalVisible(true);
      setDetailVisible(false); // Close detail modal if it is open
    } catch (error) {
      console.error("Get project details failed:", error);
      message.error(t("project.message.getDetailsFailed"));
    }
  };

  // View project details
  const showProjectDetail = async (project: Project) => {
    try {
      setDetailLoading(true);
      setDetailVisible(true);

      // Get project details
      const response: ApiResponse<Project> = await api.get(`/api/projects/detail/${project.id}`);
      setCurrentProject(response.data); // response contains {success, data, meta}, so response.data is the project
    } catch (error) {
      console.error("Get project details failed:", error);
      message.error(t("project.message.getDetailsFailed"));
    } finally {
      setDetailLoading(false);
    }
  };

  // Table column definition
  const columns = [
    {
      title: t("project.table.id"),
      dataIndex: "id",
      key: "id",
      width: 80,
    },
    {
      title: t("project.table.name"),
      dataIndex: "name",
      key: "name",
    },
    {
      title: t("project.table.slug"),
      dataIndex: "slug",
      key: "slug",
    },
    {
      title: t("project.table.description"),
      dataIndex: "description",
      key: "description",
      ellipsis: true,
    },
    {
      title: t("project.table.statusTitle"),
      dataIndex: "status",
      key: "status",
      render: (status: string) => (
        <span style={{ color: status === "active" ? "green" : "orange" }}>
          {status === "active"
            ? t("project.table.status.active")
            : t("project.table.status.archived")}
        </span>
      ),
    },
    {
      title: t("project.table.createdAt"),
      dataIndex: "created_at",
      key: "created_at",
      render: (createdAt: string) => {
        return new Date(createdAt).toLocaleString();
      },
    },
    {
      title: t("project.table.action"),
      key: "action",
      width: 180,
      render: (_: unknown, record: Project) => (
        <Space size="small">
          <Button
            type="text"
            icon={<EyeOutlined />}
            size="small"
            onClick={() => showProjectDetail(record)}
          />
          <Button
            type="text"
            icon={<EditOutlined />}
            size="small"
            onClick={() => showEditModal(record)}
          />
          <Popconfirm
            title={t("project.delete.title")}
            description={t("project.delete.description")}
            onConfirm={() => handleDeleteProject(record.id)}
            okText={t("project.delete.confirm")}
            cancelText={t("project.delete.cancel")}
          >
            <Button type="text" danger icon={<DeleteOutlined />} size="small" />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div className="bg-white p-6 rounded-lg shadow">
      <div className="mb-6 flex justify-between items-center">
        <Title level={3}>{t("project.title")}</Title>
        <Space size="middle">
          <Input
            placeholder={t("project.search")}
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
            {t("project.create")}
          </Button>
        </Space>
      </div>

      {error && (
        <div className="mb-4">
          <ErrorAlert
            message={t("project.message.getListFailed")}
            onRetry={fetchProjects}
          />
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
          showTotal: (total) => t("project.table.total", { count: total }),
          onChange: (page, pageSize) => {
            setPage(page);
            setPageSize(pageSize);
          },
        }}
      />

      {/* Create/Edit project modal */}
      <Modal
        title={currentProject ? t("project.edit") : t("project.create")}
        open={modalVisible}
        onOk={handleSaveProject}
        onCancel={() => setModalVisible(false)}
        okText={currentProject ? t("project.edit") : t("project.create")}
        cancelText={t("project.delete.cancel")}
        maskClosable={false}
      >
        <Form form={form} layout="vertical" name="projectForm">
          <Form.Item
            name="name"
            label={t("project.form.name.label")}
            rules={[
              { required: true, message: t("project.form.name.required") },
            ]}
          >
            <Input placeholder={t("project.form.name.placeholder")} />
          </Form.Item>

          <Form.Item
            name="description"
            label={t("project.form.description.label")}
          >
            <Input.TextArea
              rows={4}
              placeholder={t("project.form.description.placeholder")}
            />
          </Form.Item>

          <Form.Item
            name="slug"
            label={t("project.form.slug.label")}
            help={t("project.form.slug.help")}
          >
            <Input placeholder={t("project.form.slug.placeholder")} />
          </Form.Item>
        </Form>
      </Modal>

      {/* Project details modal */}
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
