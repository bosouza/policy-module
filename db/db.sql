USE policy;

CREATE TABLE resource (
    ID VARCHAR(50) NOT NULL PRIMARY KEY,
    rego TEXT,
    jsonSchema JSON
);

CREATE TABLE policy (
    ID VARCHAR(50) NOT NULL PRIMARY KEY,
    systemPolicy BOOLEAN
);

CREATE TABLE policy_resource (
    policyID VARCHAR(50),
    resourceID VARCHAR(50),
    content JSON,
    FOREIGN KEY (policyID) REFERENCES policy(ID),
    FOREIGN KEY (resourceID) REFERENCES resource(ID)
);

# we don't really own this table, just need it so we can assign policies to users
CREATE TABLE user (
    ID VARCHAR(50) NOT NULL PRIMARY KEY
);

CREATE TABLE user_policy (
    userID VARCHAR(50),
    policyID VARCHAR(50),
    FOREIGN KEY (userID) REFERENCES user(ID),
    FOREIGN KEY (policyID) REFERENCES policy(ID)
);

CREATE OR REPLACE
VIEW user_resource(userID, resourceID, content)
AS 
SELECT
	usr.ID,
	rsc.ID,
	plc_rsc.content
FROM user usr
INNER JOIN user_policy usr_plc
	ON usr_plc.userID = usr.id
INNER JOIN policy plc
	ON plc.id = usr_plc.policyID
INNER JOIN policy_resource plc_rsc
	ON plc_rsc.policyID = plc.id
INNER JOIN resource rsc
	ON rsc.id = plc_rsc.resourceID;

CREATE TABLE system_policy (
    ID VARCHAR(50) NOT NULL PRIMARY KEY
);

CREATE TABLE system_policy_resource (
    policyID VARCHAR(50),
    resourceID VARCHAR(50),
    content JSON,
    FOREIGN KEY (policyID) REFERENCES system_policy(ID),
    FOREIGN KEY (resourceID) REFERENCES resource(ID)
);

CREATE TRIGGER insert_policy
AFTER INSERT ON system_policy
FOR EACH ROW 
INSERT INTO policy (ID, systemPolicy)
VALUES (new.ID, 1);

CREATE TRIGGER delete_policy
AFTER DELETE ON system_policy
FOR EACH ROW 
DELETE FROM policy
WHERE policy.ID = old.ID;

CREATE TRIGGER update_policy
AFTER UPDATE ON system_policy
FOR EACH ROW 
UPDATE policy
SET policy.ID = new.ID
WHERE policy.ID = old.ID;

CREATE TRIGGER insert_policy_resource
AFTER INSERT ON system_policy_resource
FOR EACH ROW 
INSERT INTO policy_resource (policyID, resourceID, content)
VALUES (new.policyID, new.resourceID, new.content);

CREATE TRIGGER delete_policy_resource
AFTER DELETE ON system_policy_resource
FOR EACH ROW 
DELETE FROM policy_resource
WHERE policy_resource.policyID = old.policyID;

CREATE TRIGGER update_policy_resource
AFTER UPDATE ON system_policy_resource
FOR EACH ROW 
UPDATE policy_resource
SET policy_resource.policyID = new.policyID
WHERE policy_resource.policyID = old.policyID;