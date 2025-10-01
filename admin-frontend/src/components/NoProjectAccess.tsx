import React from 'react';
import { Empty, Button, Card, Typography } from 'antd';
import { TeamOutlined, MailOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useAuthStore } from '../stores/authStore';

const { Title, Paragraph } = Typography;

interface NoProjectAccessProps {
  onContactAdmin?: () => void;
}

const NoProjectAccess: React.FC<NoProjectAccessProps> = ({ onContactAdmin }) => {
  const { t } = useTranslation();
  const { user } = useAuthStore();

  return (
    <Card>
      <div style={{ textAlign: 'center', padding: '60px 20px' }}>
        <Empty
          image={<TeamOutlined style={{ fontSize: 64, color: '#d9d9d9' }} />}
          description={
            <div>
              <Title level={4} style={{ color: '#666', marginTop: 16 }}>
                {t('translation.noProjectAccess.title')}
              </Title>
              <Paragraph style={{ color: '#999', fontSize: 14 }}>
                {t('translation.noProjectAccess.description')}
              </Paragraph>
            </div>
          }
        >
          <div style={{ marginTop: 24 }}>
            <Paragraph style={{ color: '#666', marginBottom: 16 }}>
              {t('translation.noProjectAccess.instructions')}
            </Paragraph>
            
            <div style={{ marginBottom: 16 }}>
              <Paragraph style={{ fontSize: 12, color: '#999' }}>
                <strong>{t('translation.noProjectAccess.currentUser')}:</strong> {user?.username} ({user?.email})
              </Paragraph>
              <Paragraph style={{ fontSize: 12, color: '#999' }}>
                <strong>{t('translation.noProjectAccess.currentRole')}:</strong> {user?.role && t(`userManagement.roles.${user.role}`)}
              </Paragraph>
            </div>

            {onContactAdmin && (
              <Button 
                type="primary" 
                icon={<MailOutlined />}
                onClick={onContactAdmin}
              >
                {t('translation.noProjectAccess.contactAdmin')}
              </Button>
            )}
          </div>
        </Empty>
      </div>
    </Card>
  );
};

export default NoProjectAccess;
