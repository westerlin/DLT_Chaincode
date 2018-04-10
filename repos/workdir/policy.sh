#!/bin/bash
if [ "$POLICY" = "OR	('Org1MSP.member','Org2MSP.member')" ] ; then
#	POLICY="AND	('Org1MSP.member','Org2MSP.member')"
#	POLICY="AND	('Org2MSP.member')"
	POLICY="AND	('Org1MSP.member','Org2MSP.member')"
#	POLICY="AND	('Org2MSP.admin')"
#	POLICY="AND	('Org2MSP.peer1')"
#	POLICY="AND	('peer0.org2.example.com','peer1.org2.example.com')"
else
	POLICY="OR	('Org1MSP.member','Org2MSP.member')"
fi

echo "Current policy is set to:"
echo "    (+) $POLICY"
