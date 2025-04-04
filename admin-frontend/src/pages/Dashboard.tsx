import React from "react";
import { Card, Row, Col, Statistic } from "antd";

const DashboardHome: React.FC = () => {
  return (
    <div>
      <div className="bg-white rounded-lg shadow p-6 mb-6">
        <h2 className="text-xl font-semibold mb-4">仪表板</h2>
        <p>欢迎使用 i18n-flow 管理系统！</p>
      </div>

      <Row gutter={16}>
        <Col span={6}>
          <Card>
            <Statistic title="项目总数" value={0} loading={true} />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="翻译总数" value={0} loading={true} />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="语言总数" value={20} />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="用户数" value={1} suffix="人" />
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default DashboardHome;
