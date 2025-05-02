import React, { useEffect, useState } from "react";
import { Card, Row, Col, Statistic, message } from "antd";
import {
  getDashboardStats,
  DashboardStats,
} from "../services/dashboardService";
import { useTranslation } from "react-i18next";

const DashboardHome: React.FC = () => {
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const { t } = useTranslation();

  console.log(t("Welcome to React"));

  useEffect(() => {
    const fetchDashboardStats = async () => {
      try {
        const response = await getDashboardStats();
        setStats(response.data);
      } catch (error) {
        console.error("Get dashboard stats failed:", error);
        message.error("Get dashboard stats failed");
      } finally {
        setLoading(false);
      }
    };

    fetchDashboardStats();
  }, []);

  console.log(stats);

  return (
    <div>
      <div className="bg-white rounded-lg shadow p-6 mb-6">
        <h2 className="text-xl font-semibold mb-4">{t("dashboard.title")}</h2>
        <p>{t("dashboard.welcome")}</p>
      </div>

      <Row gutter={16}>
        <Col span={6}>
          <Card>
            <Statistic
              title={t("dashboard.stats.projectCount")}
              value={stats?.project_count || 0}
              loading={loading}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title={t("dashboard.stats.translationCount")}
              value={stats?.translation_count || 0}
              loading={loading}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title={t("dashboard.stats.languageCount")}
              value={stats?.language_count || 0}
              loading={loading}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title={t("dashboard.stats.userCount")}
              value={stats?.user_count || 0}
              loading={loading}
            />
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default DashboardHome;
