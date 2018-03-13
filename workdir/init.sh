#!/bin/bash
ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
CHANNEL_NAME=minidlt

alias transfer=". /wrk/transfer.sh"
alias invoke=". /wrk/invoke.sh"
alias query=". /wrk/query.sh"
alias account=". /wrk/account.sh"
alias list=". /wrk/list.sh"
alias update=". /wrk/update.sh"
#RUN echo 'alias hi="echo hello"' >> ~/.bashrc
#RUN echo 'alias list=". /wrk/list.sh"' >> ~/.bashrc
