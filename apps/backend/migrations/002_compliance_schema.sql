-- Create compliance standards table
CREATE TABLE compliance_standards (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    description TEXT,
    version VARCHAR(20) DEFAULT '1.0',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tenant_id, code)
);

-- Create compliance requirements table
CREATE TABLE compliance_requirements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    standard_id UUID NOT NULL REFERENCES compliance_standards(id) ON DELETE CASCADE,
    code VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(50),
    priority VARCHAR(20) DEFAULT 'medium',
    is_mandatory BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create compliance assessments table
CREATE TABLE compliance_assessments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    requirement_id UUID NOT NULL REFERENCES compliance_requirements(id) ON DELETE CASCADE,
    status VARCHAR(20) DEFAULT 'pending',
    evidence TEXT,
    assessor_id UUID REFERENCES users(id),
    assessed_at TIMESTAMP,
    next_review_date TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create compliance gaps table
CREATE TABLE compliance_gaps (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    assessment_id UUID NOT NULL REFERENCES compliance_assessments(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    severity VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'open',
    remediation_plan TEXT,
    target_date DATE,
    responsible_id UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_compliance_standards_tenant_id ON compliance_standards(tenant_id);
CREATE INDEX idx_compliance_requirements_standard_id ON compliance_requirements(standard_id);
CREATE INDEX idx_compliance_assessments_tenant_id ON compliance_assessments(tenant_id);
CREATE INDEX idx_compliance_assessments_requirement_id ON compliance_assessments(requirement_id);
CREATE INDEX idx_compliance_gaps_assessment_id ON compliance_gaps(assessment_id);

-- Insert default compliance standards
INSERT INTO compliance_standards (id, tenant_id, name, code, description) VALUES 
('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'ISO 27001', 'ISO27001', 'Информационная безопасность - Системы управления'),
('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', 'GDPR', 'GDPR', 'Общий регламент по защите данных'),
('00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000001', 'НПА РК', 'NPA_KZ', 'Национальное право Республики Казахстан');

-- Insert sample ISO 27001 requirements
INSERT INTO compliance_requirements (standard_id, code, title, description, category, priority, is_mandatory) VALUES 
('00000000-0000-0000-0000-000000000001', 'A.5.1.1', 'Политики информационной безопасности', 'Политики для информационной безопасности должны быть определены, одобрены руководством, опубликованы и доведены до всех сотрудников', 'Политики', 'high', true),
('00000000-0000-0000-0000-000000000001', 'A.5.1.2', 'Периодический пересмотр политик', 'Политики информационной безопасности должны периодически пересматриваться', 'Политики', 'high', true),
('00000000-0000-0000-0000-000000000001', 'A.6.1.1', 'Роли и обязанности', 'Все обязанности в области информационной безопасности должны быть четко определены и назначены', 'Организация', 'high', true),
('00000000-0000-0000-0000-000000000001', 'A.7.1.1', 'Проверка при приеме на работу', 'Справочные проверки должны проводиться для всех кандидатов на работу', 'Управление персоналом', 'medium', true),
('00000000-0000-0000-0000-000000000001', 'A.8.1.1', 'Классификация активов', 'Активы должны быть классифицированы и защищены в соответствии с их важностью', 'Управление активами', 'high', true);

-- Insert sample GDPR requirements
INSERT INTO compliance_requirements (standard_id, code, title, description, category, priority, is_mandatory) VALUES 
('00000000-0000-0000-0000-000000000002', 'Art.5', 'Принципы обработки данных', 'Персональные данные должны обрабатываться законно, справедливо и прозрачно', 'Принципы', 'high', true),
('00000000-0000-0000-0000-000000000002', 'Art.25', 'Защита данных по умолчанию', 'Технические и организационные меры должны обеспечивать защиту данных по умолчанию', 'Технические меры', 'high', true),
('00000000-0000-0000-0000-000000000002', 'Art.32', 'Безопасность обработки', 'Организация должна обеспечить безопасность обработки персональных данных', 'Безопасность', 'high', true),
('00000000-0000-0000-0000-000000000002', 'Art.33', 'Уведомление о нарушении', 'О нарушениях персональных данных должно быть уведомлено в течение 72 часов', 'Инциденты', 'high', true);
