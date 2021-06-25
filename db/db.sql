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