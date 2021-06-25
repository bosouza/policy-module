USE policy;

INSERT INTO user(ID)
VALUES ('test-user-1'),('test-user-2'),('test-user-3');

INSERT INTO resource(ID)
VALUES ('R1');

INSERT INTO resource(ID)
VALUES ('R2');

INSERT INTO system_policy(ID)
VALUES ('1');
INSERT INTO system_policy(ID)
VALUES ('2');
INSERT INTO system_policy(ID)
VALUES ('3');

DELETE FROM system_policy
WHERE ID = '3';

UPDATE system_policy
SET ID = '3'
WHERE ID = '2';

INSERT INTO system_policy_resource(policyID, resourceID, content)
VALUES ('1','R1','Test 1');
INSERT INTO system_policy_resource(policyID, resourceID, content)
VALUES ('3','R2','Test 2');

DELETE FROM system_policy_resource
WHERE policyID = '3';

UPDATE system_policy_resource
SET policyID = '2'
WHERE policyID = '1';