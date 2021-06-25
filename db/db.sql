USE policy;

CREATE TABLE resource (
    ID VARCHAR(50) NOT NULL PRIMARY KEY,
    rego TEXT,
    jsonSchema JSON
);

CREATE TABLE policy (
    ID VARCHAR(50) NOT NULL PRIMARY KEY
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
VIEW USER_RESOURCE(user_id, user_resource, resource_content)
AS 
SELECT
	usr.id,
	rsc.id,
	plc_rsc.content
FROM USER usr
INNER JOIN USER_POLICY usr_plc
	ON usr_plc.userID = usr.id
INNER JOIN POLICY plc
	ON plc.id = usr_plc.policyID
INNER JOIN POLICY_RESOURCE plc_rsc
	ON plc_rsc.policyID = plc.id
INNER JOIN RESOURCE rsc
	ON rsc.id = plc_rsc.resourceID;