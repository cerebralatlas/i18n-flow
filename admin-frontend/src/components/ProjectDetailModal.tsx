import React from "react";
import { Modal, Descriptions, Button, Spin } from "antd";
import { Project } from "../types/project";

interface ProjectDetailModalProps {
  project: Project | null;
  visible: boolean;
  loading: boolean;
  onClose: () => void;
  onEdit: (project: Project) => void;
}

const ProjectDetailModal: React.FC<ProjectDetailModalProps> = ({
  project,
  visible,
  loading,
  onClose,
  onEdit,
}) => {
  if (!project) return null;

  return (
    <Modal
      title="项目详情"
      open={visible}
      onCancel={onClose}
      footer={[
        <Button key="close" onClick={onClose}>
          关闭
        </Button>,
        <Button key="edit" type="primary" onClick={() => onEdit(project)}>
          编辑项目
        </Button>,
      ]}
      width={700}
    >
      {loading ? (
        <div className="flex justify-center items-center p-10">
          <Spin size="large" />
        </div>
      ) : (
        <Descriptions bordered column={1}>
          <Descriptions.Item label="项目ID">{project.id}</Descriptions.Item>
          <Descriptions.Item label="项目名称">{project.name}</Descriptions.Item>
          <Descriptions.Item label="项目标识">{project.slug}</Descriptions.Item>
          <Descriptions.Item label="项目描述">
            {project.description || "无"}
          </Descriptions.Item>
          <Descriptions.Item label="状态">
            <span
              style={{
                color: project.status === "active" ? "green" : "orange",
              }}
            >
              {project.status === "active" ? "活跃" : "归档"}
            </span>
          </Descriptions.Item>
          <Descriptions.Item label="创建时间">
            {project.created_at}
          </Descriptions.Item>
          <Descriptions.Item label="更新时间">
            {project.updated_at}
          </Descriptions.Item>
        </Descriptions>
      )}
    </Modal>
  );
};

export default ProjectDetailModal;
