CREATE TABLE organization (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
); 

CREATE TABLE application (
    id TEXT PRIMARY KEY,
    organization_id TEXT,
    client_id TEXT NOT NULL,
    client_secret TEXT NOT NULL,
    name TEXT NOT NULL,
    redirect_uris TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES organization(id) ON DELETE CASCADE,
    UNIQUE (organization_id, client_id),
    UNIQUE (organization_id, name)
);

CREATE TABLE grant_type (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE client_grant_type (
    application_id TEXT,
    grant_type_id INTEGER,
    PRIMARY KEY (application_id, grant_type_id),
    FOREIGN KEY (application_id) REFERENCES application(id) ON DELETE CASCADE,
    FOREIGN KEY (grant_type_id) REFERENCES grant_type(id) ON DELETE CASCADE
);

CREATE TABLE org_user (
    id TEXT PRIMARY KEY,
    organization_id TEXT,
    username TEXT NOT NULL,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES organization(id) ON DELETE CASCADE,
    UNIQUE (organization_id, username),
    UNIQUE (organization_id, email)
);

CREATE TABLE attribute (
    id TEXT PRIMARY KEY,
    organization_id TEXT,
    name TEXT NOT NULL,
    FOREIGN KEY (organization_id) REFERENCES organization(id) ON DELETE CASCADE,
    UNIQUE (organization_id, name)
);

CREATE TABLE user_attribute (
    user_id TEXT,
    attribute_id TEXT,
    value TEXT,
    PRIMARY KEY (user_id, attribute_id),
    FOREIGN KEY (user_id) REFERENCES org_user(id) ON DELETE CASCADE,
    FOREIGN KEY (attribute_id) REFERENCES attribute(id)
);

CREATE TABLE token (
    id TEXT PRIMARY KEY,
    client_id TEXT NOT NULL,
    entry_id TEXT NOT NULL,
    organization_id TEXT,
    created_at BIGINT NOT NULL,
    expires_at BIGINT NOT NULL,
    FOREIGN KEY (organization_id) REFERENCES organization(id) ON DELETE CASCADE,
    FOREIGN KEY (client_id, organization_id) REFERENCES application(client_id, organization_id)
);

CREATE TABLE role (
    id TEXT PRIMARY KEY,
    organization_id TEXT,
    name TEXT NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES organization(id),
    UNIQUE (organization_id, name)
);

CREATE TABLE user_role (
    user_id TEXT,
    role_id TEXT,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES org_user(id),
    FOREIGN KEY (role_id) REFERENCES role(id)
);
