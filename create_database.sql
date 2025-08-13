-- Create gofinch database
IF NOT EXISTS (SELECT name FROM sys.databases WHERE name = 'gofintech')
BEGIN
    CREATE DATABASE gofintech;
END
GO

PRINT 'Database gofinch created successfully!';
GO

