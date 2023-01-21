CREATE TABLE IF NOT EXISTS identity.roles (
     role_id SERIAL NOT NULL,
     role_uuid UUID NOT NULL,
     name varchar(250) NOT NULL,
     description varchar(500) NULL,
     created_at timestamp NOT NULL,
     last_updated_at timestamp NOT NULL,
     PRIMARY KEY (role_id)
);

CREATE TABLE IF NOT EXISTS identity.permissions (
    permission_id SERIAL NOT NULL,
    permission_uuid UUID NOT NULL,
    name varchar(250) NOT NULL,
    description varchar(500) NULL,
    created_at timestamp NOT NULL,
    last_updated_at timestamp NOT NULL,
    PRIMARY KEY (permission_id)
);

CREATE TABLE IF NOT EXISTS identity.role_permissions (
    role_id INT NOT NULL REFERENCES identity.roles(role_id),
    permission_id INT NOT NULL REFERENCES identity.permissions(permission_id),
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE IF NOT EXISTS identity.user_roles (
    user_id INT NOT NULL REFERENCES identity.users(user_id),
    role_id INT NOT NULL REFERENCES identity.roles(role_id),
    PRIMARY KEY (user_id, role_id)
);