-- Create tables in gofintech database
CREATE TABLE users (
    id INT IDENTITY(1,1) PRIMARY KEY,
    username NVARCHAR(50) NOT NULL UNIQUE,
    email NVARCHAR(100) NOT NULL UNIQUE,
    password_hash NVARCHAR(255) NOT NULL,
    role NVARCHAR(20) NOT NULL,
    created_at DATETIME2 NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME2 NOT NULL DEFAULT GETDATE()
);

CREATE TABLE transactions (
    id INT IDENTITY(1,1) PRIMARY KEY,
    from_user_id INT FOREIGN KEY REFERENCES users(id),
    to_user_id INT FOREIGN KEY REFERENCES users(id),
    amount DECIMAL(18,2) NOT NULL,
    type NVARCHAR(20) NOT NULL,
    status NVARCHAR(20) NOT NULL,
    created_at DATETIME2 NOT NULL DEFAULT GETDATE()
);

CREATE TABLE balances (
    user_id INT PRIMARY KEY FOREIGN KEY REFERENCES users(id),
    amount DECIMAL(18,2) NOT NULL,
    last_updated_at DATETIME2 NOT NULL DEFAULT GETDATE()
);

CREATE TABLE audit_logs (
    id INT IDENTITY(1,1) PRIMARY KEY,
    entity_type NVARCHAR(50) NOT NULL,
    entity_id INT NOT NULL,
    action NVARCHAR(50) NOT NULL,
    details NVARCHAR(MAX),
    created_at DATETIME2 NOT NULL DEFAULT GETDATE()
);

-- Performance indexes
CREATE INDEX IX_transactions_from_user_id ON transactions(from_user_id);
CREATE INDEX IX_transactions_to_user_id ON transactions(to_user_id);
CREATE INDEX IX_transactions_created_at ON transactions(created_at);
CREATE INDEX IX_transactions_type_status ON transactions(type, status);
CREATE INDEX IX_users_email ON users(email);
CREATE INDEX IX_users_role ON users(role);
CREATE INDEX IX_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX IX_audit_logs_created_at ON audit_logs(created_at);

PRINT 'Tables created successfully!';

