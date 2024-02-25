BEGIN;
CREATE TABLE IF NOT EXISTS users(
   Id CHAR(36) NOT NULL PRIMARY KEY,
   Name VARCHAR(255) NOT NULL,
   Email VARCHAR(50) NOT NULL,
   CreatedOn DATETIME NOT NULL,
   ModifiedOn DATETIME NOT NULL,
   LastLogin DATETIME NOT NULL
);
CREATE UNIQUE INDEX idx_users_email on users (Email);

CREATE TABLE IF NOT EXISTS domains(
   Id CHAR(36) NOT NULL PRIMARY KEY,
   Name VARCHAR(255) NOT NULL,
   CreatedOn DATETIME NOT NULL,
   ModifiedOn DATETIME NOT NULL,
   IsDeprecated BOOLEAN NOT NULL,
   CreatedById CHAR(36) NOT NULL
);
CREATE UNIQUE INDEX idx_domains_name on domains (Name);
CREATE INDEX idx_domains_createdById on domains (CreatedById);


CREATE TABLE IF NOT EXISTS shortlinks(
   Hash VARCHAR(6) PRIMARY KEY NOT NULL,
   OriginalUrl VARCHAR(255) NOT NULL,
   DomainId CHAR(36) NOT NULL,
   Alias VARCHAR(10) NULL,
   Hits INTEGER NULL,
   CreatedOn DATETIME NOT NULL,
   ModifiedOn DATETIME NOT NULL,
   ExpirationDate DATETIME NOT NULL,
   IsDeprecated BOOLEAN NOT NULL,
   UserId CHAR(36) NOT NULL
);
CREATE UNIQUE INDEX idx_shortlinks_original_url ON shortlinks (OriginalUrl);
CREATE INDEX idx_shortlinks_userId ON shortlinks (UserId);
CREATE INDEX idx_shortlinks_domainId ON shortlinks (DomainId);


CREATE TABLE IF NOT EXISTS unusedshortlinks(
   Hash VARCHAR(6) PRIMARY KEY NOT NULL,
   CreatedOn DATETIME NOT NULL,
   ModifiedOn DATETIME NOT NULL,
   ExpirationDate DATETIME NOT NULL,
   Used BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS userkeys(
   Id CHAR(36)  NOT NULL PRIMARY KEY,
   ApiKey VARCHAR(20) NOT NULL,
   CreatedOn DATETIME NOT NULL,
   ModifiedOn DATETIME NOT NULL,
   ExpirationDate DATETIME NOT NULL,
   UserId CHAR(36) NOT NULL,
   IsActive BOOLEAN NOT NULL
);
CREATE INDEX idx_userkeys_userId ON userkeys (UserId);

COMMIT;