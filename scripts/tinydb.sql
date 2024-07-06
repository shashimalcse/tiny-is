CREATE TABLE application (
    id VARCHAR(255) PRIMARY KEY,
    client_id VARCHAR(255) NOT NULL,
    client_secret VARCHAR(255) NOT NULL,
    redirect_uri VARCHAR(255) NOT NULL,
    grant_types VARCHAR(255) NOT NULL
);

CREATE TABLE org_user (
    id VARCHAR(255) PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE authorization_server (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
);

CREATE TABLE scope (
    id VARCHAR(255) PRIMARY KEY,
    authorization_server_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    FOREIGN KEY (authorization_server_id) REFERENCES authorization_server(id)
);

CREATE TABLE policy (
    id VARCHAR(255) PRIMARY KEY,
    authorization_server_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    FOREIGN KEY (authorization_server_id) REFERENCES authorization_server(id)
);

CREATE TABLE policy_application (
    id VARCHAR(255) PRIMARY KEY,
    policy_id VARCHAR(255) NOT NULL,
    application_id VARCHAR(255) NOT NULL,
    FOREIGN KEY (policy_id) REFERENCES policy(id),
    FOREIGN KEY (application_id) REFERENCES application(id)
);

CREATE TABLE rule (
    id VARCHAR(255) PRIMARY KEY,
    order_id INTEGER NOT NULL,
    policy_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    access_token_lifetime INTEGER NOT NULL,
    FOREIGN KEY (policy_id) REFERENCES policy(id)
);

CREATE TABLE rule_grant_type (
    rule_id VARCHAR(255),
    grant_type VARCHAR(255),
    PRIMARY KEY (rule_id, grant_type),
    FOREIGN KEY (rule_id) REFERENCES rules(id)
);

CREATE TABLE rule_user (
    rule_id VARCHAR(255),
    user_id VARCHAR(255),
    PRIMARY KEY (rule_id, user_id),
    FOREIGN KEY (rule_id) REFERENCES rules(id),
    FOREIGN KEY (user_id) REFERENCES org_user(id)
);

CREATE TABLE rule_scope (
    rule_id VARCHAR(255),
    scope_id VARCHAR(255),
    PRIMARY KEY (rule_id, scope_id),
    FOREIGN KEY (rule_id) REFERENCES rules(id),
    FOREIGN KEY (scope_id) REFERENCES scope(id)
);

