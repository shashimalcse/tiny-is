CREATE TABLE organization (
    id UUID PRIMARY KEY,
    name VARCHAR(255) UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
); 

CREATE TABLE application (
    id UUID PRIMARY KEY,
    organization_id UUID REFERENCES organization(id),
    client_id VARCHAR(255) NOT NULL,
    client_secret VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    redirect_uris TEXT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (organization_id, client_id),
    UNIQUE (organization_id, name)
);

CREATE TABLE grant_type (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

INSERT INTO grant_type (name) VALUES 
    ('authorization_code'),
    ('client_credentials'),
    ('refresh_token');  

CREATE TABLE client_grant_type (
    application_id UUID REFERENCES application(id),
    grant_type_id INTEGER REFERENCES grant_type(id),
    PRIMARY KEY (application_id, grant_type_id)
);

CREATE TABLE org_user (
    id UUID PRIMARY KEY,
    organization_id UUID REFERENCES organization(id),
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (organization_id, username),
    UNIQUE (organization_id, email)
);


CREATE TABLE access_token (
    id UUID PRIMARY KEY,
    application_id UUID REFERENCES application(id),
    user_id UUID REFERENCES org_user(id),
    organization_id UUID REFERENCES organization(id),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE refresh_token (
    id UUID PRIMARY KEY,
    application_id UUID REFERENCES application(id),
    user_id UUID REFERENCES org_user(id),
    organization_id UUID REFERENCES organization(id),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE role (
    id UUID PRIMARY KEY,
    organization_id UUID REFERENCES organization(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (organization_id, name)
);

CREATE TABLE user_role (
    user_id UUID REFERENCES org_user(id),
    role_id UUID REFERENCES role(id),
    PRIMARY KEY (user_id, role_id)
);


