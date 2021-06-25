USE policy;

CREATE TABLE resource (
    ID VARCHAR(30) NOT NULL PRIMARY KEY,
    rego TEXT,
    jsonSchema JSON
);

CREATE TABLE policy (
    ID VARCHAR(30) NOT NULL PRIMARY KEY
);

CREATE TABLE policy_resource (
    policyID VARCHAR(30),
    resourceID VARCHAR(30),
    FOREIGN KEY (policyID) REFERENCES policy(ID),
    FOREIGN KEY (resourceID) REFERENCES resource(ID)
);