CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Tenants
CREATE TABLE tenants (
    tenant_id   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(255) NOT NULL,
    domain      VARCHAR(255),
    status      VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW()
);

-- Roles
CREATE TABLE roles (
    role_id   INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    name      VARCHAR(50) NOT NULL,
    UNIQUE (tenant_id, name)
);

CREATE TABLE permissions (
    permission_id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE role_permissions (
    role_id INT REFERENCES roles(role_id) ON DELETE CASCADE,
    permission_id INT REFERENCES permissions(permission_id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE org_units (
    org_unit_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    parent_org_unit_id UUID REFERENCES org_units(org_unit_id),
    UNIQUE (tenant_id, name)
);

-- Users
CREATE TABLE users (
    user_id      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id    UUID REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    org_unit_id UUID REFERENCES org_units(org_unit_id),
    email        VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name    VARCHAR(255),
    status       VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at   TIMESTAMP DEFAULT NOW(),
    updated_at   TIMESTAMP DEFAULT NOW(),
    last_login_at TIMESTAMP,
    UNIQUE (tenant_id, email)
);

CREATE TABLE user_roles (
    user_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
    role_id INT REFERENCES roles(role_id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

-- Invitations
CREATE TABLE invitations (
    invitation_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id     UUID NOT NULL REFERENCES tenants(tenant_id) ON DELETE CASCADE,
    email         VARCHAR(255) NOT NULL,
    token         VARCHAR(255) NOT NULL UNIQUE,
    status        VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    expires_at    TIMESTAMP,
    created_at    TIMESTAMP DEFAULT NOW()
);

CREATE TABLE invitation_roles (
    invitation_id UUID REFERENCES invitations(invitation_id) ON DELETE CASCADE,
    role_id INT REFERENCES roles(role_id) ON DELETE CASCADE,
    PRIMARY KEY (invitation_id, role_id)
);

-- Password resets
CREATE TABLE password_resets (
    reset_id   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    token      VARCHAR(255) NOT NULL UNIQUE,
    status     VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Refresh tokens
CREATE TABLE refresh_tokens (
    token_id    UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    token_hash  VARCHAR(255) NOT NULL,
    revoked     BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at  TIMESTAMP NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW()
);

-- Audit Logs
CREATE TABLE audit_logs (
    log_id      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID REFERENCES users(user_id) ON DELETE SET NULL,
    tenant_id   UUID REFERENCES tenants(tenant_id) ON DELETE SET NULL,
    action      VARCHAR(100) NOT NULL,
    ip_address  VARCHAR(45),
    user_agent  VARCHAR(255),
    created_at  TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX idx_invitations_email ON invitations(email);
CREATE INDEX idx_password_resets_user_id ON password_resets(user_id);