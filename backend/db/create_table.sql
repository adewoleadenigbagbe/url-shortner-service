BEGIN;

-- USERS
CREATE TABLE IF NOT EXISTS userRoles(
   Id CHAR(36) NOT NULL PRIMARY KEY,
   Role INTEGER NOT NULL,
   CreatedOn DATETIME NOT NULL,
   IsDeprecated BOOLEAN NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_userroles_role on userRoles (Role);

CREATE TABLE IF NOT EXISTS users(
   Id CHAR(36) NOT NULL PRIMARY KEY,
   Name VARCHAR(255) NOT NULL,
   Email VARCHAR(50) NOT NULL,
   Password VARCHAR(50) NOT NULL,
   CreatedOn DATETIME NOT NULL,
   ModifiedOn DATETIME NOT NULL,
   LastLogin DATETIME NOT NULL,
   OrganizationId CHAR(36) NOT NULL,
   ReferralUserId CHAR(36) NULL,
   RoleId CHAR(36) NOT NULL,
   IsDeprecated BOOLEAN NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_users_roleId on users (RoleId);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email on users (Email);

CREATE TABLE IF NOT EXISTS userkeys(
   Id CHAR(36)  NOT NULL PRIMARY KEY,
   Apikey VARCHAR(20) NOT NULL,
   CreatedOn DATETIME NOT NULL,
   ModifiedOn DATETIME NOT NULL,
   ExpirationDate DATETIME NOT NULL,
   UserId CHAR(36) NOT NULL,
   OrganizationId CHAR(36) NOT NULL,
   IsActive BOOLEAN NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_userkeys_userId ON userkeys (UserId);

CREATE TABLE IF NOT EXISTS organizations(
   Id CHAR(36) NOT NULL PRIMARY KEY,
   Name VARCHAR(255) NOT NULL,
   PhoneNumber VARCHAR(20) NOT NULL,
   TimeZone  VARCHAR(50) NULL,
   OwnerId CHAR(36) NOT NULL,
   CreatedOn DATETIME NOT NULL,
   ModifiedOn DATETIME NOT NULL,
   IsDeprecated BOOLEAN NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_organizations_Name on organizations (Name);
CREATE UNIQUE INDEX IF NOT EXISTS idx_organizations_ownerId on organizations (OwnerId);

CREATE TABLE IF NOT EXISTS invites(
   Id CHAR(36) NOT NULL PRIMARY KEY,
   Email VARCHAR(50) NOT NULL,
   ReferralUserId CHAR(36) NOT NULL,
   RoleId CHAR(36) NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_invites_email on invites(Email);

-- TEAMS
CREATE TABLE IF NOT EXISTS teams(
   Id CHAR(36) NOT NULL PRIMARY KEY,
   Name VARCHAR(255) NOT NULL,
   OrganizationId CHAR(36) NOT NULL,
   IsDeprecated BOOLEAN NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_teams_name on teams (Name);
CREATE INDEX IF NOT EXISTS idx_teams_organizationId on teams (OrganizationId);


CREATE TABLE IF NOT EXISTS teamusers(
   Id CHAR(36) NOT NULL PRIMARY KEY,
   TeamId CHAR(36) NOT NULL,
   UserId CHAR(36) NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_teamusers_teamId on teamusers (TeamId);
CREATE INDEX IF NOT EXISTS idx_teamusers_userId on teamusers (UserId);

--DOMAINS
CREATE TABLE IF NOT EXISTS domains(
   Id CHAR(36) NOT NULL PRIMARY KEY,
   Name VARCHAR(255) NOT NULL,
   IsCustom BOOLEAN NOT NULL,
   OrganizationId CHAR(36) NOT NULL,
   CreatedOn DATETIME NOT NULL,
   ModifiedOn DATETIME NOT NULL,
   CreatedById CHAR(36) NOT NULL,
   IsDeprecated BOOLEAN NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_domains_name on domains (Name);
CREATE INDEX IF NOT EXISTS idx_domains_createdById on domains (CreatedById);

-- SHORTLINKS
CREATE TABLE IF NOT EXISTS shortlinks(
   Id CHAR(36) NOT NULL PRIMARY KEY,
   Hash VARCHAR(8) NOT NULL,
   OriginalUrl VARCHAR(255) NOT NULL,
   DomainId CHAR(36) NOT NULL,
   Alias VARCHAR(10) NULL,
   CreatedOn DATETIME NOT NULL,
   ModifiedOn DATETIME NOT NULL,
   ExpirationDate DATETIME NOT NULL,
   OrganizationId CHAR(36) NOT NULL,
   CreatedById CHAR(36) NOT NULL,
   Cloaking BOOLEAN NOT NULL,
   IsDeprecated BOOLEAN NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_shortlinks_organizationId ON shortlinks (OrganizationId);
CREATE INDEX IF NOT EXISTS idx_shortlinks_domainId ON shortlinks (DomainId);


CREATE TABLE IF NOT EXISTS unusedshortlinks(
   Hash VARCHAR(8) PRIMARY KEY NOT NULL,
   CreatedOn DATETIME NOT NULL,
   ModifiedOn DATETIME NOT NULL,
   ExpirationDate DATETIME NOT NULL,
   Used BOOLEAN NOT NULL
);

-- CLICK LOGS
CREATE TABLE IF NOT EXISTS accesslogs(
   Id CHAR(36)  NOT NULL PRIMARY KEY,
   ShortId CHAR(36) NOT NULL,
   Country VARCHAR(50) NULL,
   TimeZone  VARCHAR(50) NULL,
   City VARCHAR(50) NULL,
   Os VARCHAR(50) NULL,
   Browser  VARCHAR(50) NULL,
   UserAgent  VARCHAR(50) NULL,
   Platform  VARCHAR(50) NULL,
   IpAddress  VARCHAR(255) NULL,
   Method VARCHAR(50) NULL,
   Status INTEGER NULL,
   OrganizationId CHAR(36) NOT NULL,
   CreatedOn DATETIME NOT NULL,
   IsDeprecated BOOLEAN NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_accesslogs_shortId ON accesslogs (ShortId);
CREATE INDEX IF NOT EXISTS idx_accesslogs_createdon ON accesslogs (CreatedOn);
CREATE INDEX IF NOT EXISTS idx_accesslogs_organizationId ON accesslogs (OrganizationId);

CREATE TABLE IF NOT EXISTS tags(
   Id CHAR(36)  NOT NULL PRIMARY KEY,
   Name VARCHAR(255) NOT NULL,
   CreatedOn DATETIME NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_tags_name on tags (Name);

CREATE TABLE IF NOT EXISTS shortlinktags(
   Id CHAR(36)  NOT NULL PRIMARY KEY,
   ShortId CHAR(36) NOT NULL,
   TagId CHAR(36) NOT NULL,
   CreatedOn DATETIME NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_shortlinktags_shortId ON shortlinktags (ShortId);
CREATE INDEX IF NOT EXISTS idx_shortlinktags_tagId ON shortlinktags (TagId);

-- PLANS & PAYMENTS
CREATE TABLE IF NOT EXISTS payplans(
   Id CHAR(36) NOT NULL PRIMARY KEY,
   Type INTEGER NOT NULL,
   Amount DOUBLE NOT NULL,
   CreatedOn DATETIME NOT NULL,
   ModifiedOn DATETIME NOT NULL,
   IsLatest BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS organizationpayplans(
   Id CHAR(36)  NOT NULL PRIMARY KEY,
   PayCycle INTEGER NOT NULL,
   PayPlanId CHAR(36) NOT NULL,
   OrganizationId CHAR(36) NOT NULL,
   CreatedOn DATETIME NOT NULL,
   ModifiedOn DATETIME NOT NULL,
   Status INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_organizationpayplans_PayPlanId ON organizationpayplans (PayPlanId);
CREATE INDEX IF NOT EXISTS idx_organizationpayplans_OrganizationId ON organizationpayplans (OrganizationId);

CREATE TABLE IF NOT EXISTS payschedules(
   Id CHAR(36) NOT NULL PRIMARY KEY,
   EffectiveDate DATETIME NOT NULL,
   EndDate DATETIME NOT NULL,
   CreatedOn DATETIME NOT NULL,
   ModifiedOn DATETIME NOT NULL,
   OrganizationId CHAR(36) NOT NULL,
   OrganizationPayPlanId CHAR(36) NOT NULL,
   IsNext BOOLEAN NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_payschedules_OrganizationPayPlanId ON payschedules (OrganizationPayPlanId);
CREATE INDEX IF NOT EXISTS idx_payschedules_OrganizationId ON payschedules (OrganizationId);


CREATE TABLE IF NOT EXISTS revenues(
   Id CHAR(36) NOT NULL PRIMARY KEY,
   Amount DOUBLE NOT NULL,
   StartDate DATETIME NOT NULL,
   EndDate DATETIME NOT NULL,
   PayScheduleId CHAR(36) NOT NULL,
   OrganizationId CHAR(36) NOT NULL,
   CreatedOn DATETIME NOT NULL,
   ModifiedOn DATETIME NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_revenues_OrganizationPayPlanId ON revenues (PayScheduleId);
CREATE INDEX IF NOT EXISTS idx_revenues_OrganizationId ON revenues (OrganizationId);

COMMIT;