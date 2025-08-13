-- Initialize gofintech database
IF NOT EXISTS (SELECT name FROM sys.databases WHERE name = 'gofintech')
BEGIN
    CREATE DATABASE gofintech;
    PRINT 'Database gofintech created successfully!';
END
ELSE
BEGIN
    PRINT 'Database gofintech already exists!';
END
GO


