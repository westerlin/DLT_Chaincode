#!/bin/bash

#ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
#CHANNEL_NAME=minidlt

#peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n mycc -c '{"Args":["list"]}' 
peer chaincode query -C $CHANNEL_NAME -n $CHAINCODEID -c '{"Args":["list","'$1'","'$2'"]}'
