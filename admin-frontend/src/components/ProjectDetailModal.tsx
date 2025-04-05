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
      title="Project Details"
      open={visible}
      onCancel={onClose}
      footer={[
        <Button key="close" onClick={onClose}>
          Close
        </Button>,
        <Button key="edit" type="primary" onClick={() => onEdit(project)}>
          Edit project
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
          <Descriptions.Item label="Project ID">{project.id}</Descriptions.Item>
          <Descriptions.Item label="Project Name">
            {project.name}
          </Descriptions.Item>
          <Descriptions.Item label="Project Slug">
            {project.slug}
          </Descriptions.Item>
          <Descriptions.Item label="Project Description">
            {project.description || "No description"}
          </Descriptions.Item>
          <Descriptions.Item label="Status">
            <span
              style={{
                color: project.status === "active" ? "green" : "orange",
              }}
            >
              {project.status === "active" ? "Active" : "Archived"}
            </span>
          </Descriptions.Item>
          <Descriptions.Item label="Created at">
            {project.created_at}
          </Descriptions.Item>
          <Descriptions.Item label="Updated at">
            {project.updated_at}
          </Descriptions.Item>
        </Descriptions>
      )}
    </Modal>
  );
};

export default ProjectDetailModal;
