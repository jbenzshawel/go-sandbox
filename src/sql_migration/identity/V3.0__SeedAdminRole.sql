INSERT INTO identity.permissions(permission_id, name) VALUES (1, 'View Users');
INSERT INTO identity.permissions(permission_id, name) VALUES (2, 'Edit Users');
INSERT INTO identity.permissions(permission_id, name) VALUES (3, 'View Roles');
INSERT INTO identity.permissions(permission_id, name) VALUES (4, 'Edit Roles');

INSERT INTO identity.roles(role_id, name) VALUES (1, 'Admin');

INSERT INTO identity.role_permissions(role_id, permission_id)
    SELECT 1, permission_id FROM identity.permissions;