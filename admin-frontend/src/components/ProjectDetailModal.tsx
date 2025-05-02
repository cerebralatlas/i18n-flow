import React from "react";
import { Modal, Descriptions, Button, Spin } from "antd";
import { Project } from "../types/project";
import { useTranslation } from "react-i18next";

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
  const { t } = useTranslation();

  if (!project) return null;

  return (
    <Modal
      title={t("project.detail.title")}
      open={visible}
      onCancel={onClose}
      footer={[
        <Button key="close" onClick={onClose}>
          {t("project.detail.close")}
        </Button>,
        <Button key="edit" type="primary" onClick={() => onEdit(project)}>
          {t("project.detail.editButton")}
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
          <Descriptions.Item label={t("project.detail.fields.id")}>
            {project.id}
          </Descriptions.Item>
          <Descriptions.Item label={t("project.detail.fields.name")}>
            {project.name}
          </Descriptions.Item>
          <Descriptions.Item label={t("project.detail.fields.slug")}>
            {project.slug}
          </Descriptions.Item>
          <Descriptions.Item label={t("project.detail.fields.description")}>
            {project.description || t("project.detail.noDescription")}
          </Descriptions.Item>
          <Descriptions.Item label={t("project.detail.fields.status")}>
            <span
              style={{
                color: project.status === "active" ? "green" : "orange",
              }}
            >
              {project.status === "active"
                ? t("project.table.status.active")
                : t("project.table.status.archived")}
            </span>
          </Descriptions.Item>
          <Descriptions.Item label={t("project.detail.fields.createdAt")}>
            {project.created_at}
          </Descriptions.Item>
          <Descriptions.Item label={t("project.detail.fields.updatedAt")}>
            {project.updated_at}
          </Descriptions.Item>
        </Descriptions>
      )}
    </Modal>
  );
};

export default ProjectDetailModal;
