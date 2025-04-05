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

  // Get project list
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
      console.error("Get project list failed:", error);
      setError("Get project list failed");
    } finally {
      setLoading(false);
    }
  }, [page, pageSize, keyword]);

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
        await axios.put(`/api/projects/update/${currentProject.id}`, values);
        message.success("Project updated successfully");
      } else {
        // Create project
        await axios.post("/api/projects", values);
        message.success("Project created successfully");
      }

      setModalVisible(false);
      form.resetFields();
      fetchProjects();
    } catch (error) {
      console.error("Save project failed:", error);
      message.error("Save project failed");
    }
  };

  // Delete project
  const handleDeleteProject = async (id: number) => {
    try {
      await axios.delete(`/api/projects/delete/${id}`);
      message.success("Project deleted successfully");
      fetchProjects();
    } catch (error) {
      console.error("Delete project failed:", error);
      message.error("Delete project failed");
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
      const response = await axios.get(`/api/projects/detail/${project.id}`);
      const projectData = response.data;

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
      message.error("Get project details failed");
    }
  };

  // View project details
  const showProjectDetail = async (project: Project) => {
    try {
      setDetailLoading(true);
      setDetailVisible(true);

      // Get project details
      const response = await axios.get(`/api/projects/detail/${project.id}`);
      setCurrentProject(response.data);
    } catch (error) {
      console.error("Get project details failed:", error);
      message.error("Get project details failed");
    } finally {
      setDetailLoading(false);
    }
  };

  // Table column definition
  const columns = [
    {
      title: "ID",
      dataIndex: "id",
      key: "id",
      width: 80,
    },
    {
      title: "Project Name",
      dataIndex: "name",
      key: "name",
    },
    {
      title: "Project Slug",
      dataIndex: "slug",
      key: "slug",
    },
    {
      title: "Description",
      dataIndex: "description",
      key: "description",
      ellipsis: true,
    },
    {
      title: "Status",
      dataIndex: "status",
      key: "status",
      render: (status: string) => (
        <span style={{ color: status === "active" ? "green" : "orange" }}>
          {status === "active" ? "Active" : "Archived"}
        </span>
      ),
    },
    {
      title: "Created At",
      dataIndex: "created_at",
      key: "created_at",
    },
    {
      title: "Action",
      key: "action",
      width: 200,
      render: (_: unknown, record: Project) => (
        <Space size="middle">
          <Button
            type="text"
            icon={<EyeOutlined />}
            onClick={() => showProjectDetail(record)}
          >
            Details
          </Button>
          <Button
            type="text"
            icon={<EditOutlined />}
            onClick={() => showEditModal(record)}
          >
            Edit
          </Button>
          <Popconfirm
            title="Are you sure you want to delete this project?"
            description="This action cannot be undone!"
            onConfirm={() => handleDeleteProject(record.id)}
            okText="Yes"
            cancelText="Cancel"
          >
            <Button type="text" danger icon={<DeleteOutlined />}>
              Delete
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div className="bg-white p-6 rounded-lg shadow">
      <div className="mb-6 flex justify-between items-center">
        <Title level={3}>Project Management</Title>
        <Space size="middle">
          <Input
            placeholder="Search projects"
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
            Create Project
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
          showTotal: (total) => `Total ${total} records`,
          onChange: (page, pageSize) => {
            setPage(page);
            setPageSize(pageSize);
          },
        }}
      />

      {/* Create/Edit project modal */}
      <Modal
        title={currentProject ? "Edit Project" : "Create Project"}
        open={modalVisible}
        onOk={handleSaveProject}
        onCancel={() => setModalVisible(false)}
        okText={currentProject ? "Update" : "Create"}
        cancelText="Cancel"
        maskClosable={false}
      >
        <Form form={form} layout="vertical" name="projectForm">
          <Form.Item
            name="name"
            label="Project Name"
            rules={[{ required: true, message: "Please enter project name" }]}
          >
            <Input placeholder="Please enter project name" />
          </Form.Item>

          <Form.Item name="description" label="Project Description">
            <Input.TextArea
              rows={4}
              placeholder="Please enter project description"
            />
          </Form.Item>

          <Form.Item
            name="slug"
            label="Project Slug"
            help="For URL-friendly identifier, leave blank to generate automatically based on project name"
          >
            <Input placeholder="For example: my-project" />
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
