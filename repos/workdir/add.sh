#!/bin/bash

ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
CORE_PEER_TLS_ENABLED=true
#CHANNEL_NAME=minidlt

#peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n $CHAINCODEID -c '{"Args":["add","'$1'","'$2'","'$3'","0"]}' 
peer chaincode invoke -C $CHANNEL_NAME -n $CHAINCODEID --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA  -c '{"Args":["add","'$1'","'$2'","'$3'","0"]}' 
