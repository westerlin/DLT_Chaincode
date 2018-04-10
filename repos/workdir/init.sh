#!/bin/bash

chaincodes=($(peer chaincode list --installed | grep "Name" | cut -d " " -f2| cut -d "," -f1))

if [ -z "$POLICY" ] ; then
ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
CHANNEL_NAME=minidlt

alias transfer=". /wrk/transfer.sh"
alias invoke=". /wrk/invoke.sh"
alias query=". /wrk/query.sh"
alias add=". /wrk/add.sh"
alias list=". /wrk/list.sh"
alias delete=". /wrk/delete.sh"
alias update=". /wrk/update.sh"
alias history=". /wrk/history.sh"
alias instantiate=". /wrk/instantiate.sh"
alias menu=". /wrk/init.sh"
alias logging=". /wrk/logging.sh"
alias policy=". /wrk/policy.sh"
alias whoami="peer chaincode query -C ${CHANNEL_NAME} -n \${CHAINCODEID} -c '{\"Args\":[\"whoami\"]}'"
alias setpeer=". /wrk/setpeer.sh"
POLICY="OR	('Org1MSP.member','Org2MSP.member')"
cd /
#export PS1="\W \n>"
export PS1="\n> "
CHAINCODEID="None"
CHAINCODEID=${chaincodes[ ${#chaincodes[@]}-1 ]}
fi

#if [-z $CHANNEL_NAME] ; then
#fi

#peer chaincode list --installed -C $CHANNEL_NAME >> chains.txt
clear
echo "Welcome to MiniDLT network"
echo
echo "   DLT commands: "
echo "      transfer <sender> <receiver> <amount>                : Transfers amount fra sender to recevier. Both are account id's."
echo "      query <accountid>                                    : Lists the current balance of account"
echo "      add <accountid> <ownerid> <initial balance> <limit>  : Creates an account with initial balance and lower limit." 
echo "      delete <accountid>                                   : Deletes account." 
echo "      list                                                 : Lists all accounts and their status."
echo "      history <accountid>                                  : Lists history for a given account"	
echo "      whoami                                               : Provides certificate for MSP/Creator of transaction"	
echo
echo "   Special commands: "
echo "      menu                                                 : Returns to this menu"
echo "      policy                                               : Change the policy"
echo "      logging <logging level>                              : Sets logging level (CRITICAL | ERROR | WARNING | NOTICE | INFO | DEBUG)"
echo "      instantiate <version> <chaincode name>               : Installs and instantiates a new chaincode"
echo "      update <version>                                     : Updates a new chaincode"
echo "      peer chaincode list --instantiated                   : Lists alle instantiated chaincodes"
echo "      setpeer <peer index>                                 : Switch between peers 0-3"
echo "      peer chaincode list --installed                      : Lists alle installed chaincodes"
echo "      peer channel list                                    : Lists alle channels joined"
echo
echo "Current Policy is: $POLICY"
echo
echo "Current Chaincode is: $CHAINCODEID (installed: ${chaincodes[@]})"
echo
echo "Current Peer is: $CORE_PEER_ADDRESS"
echo
#RUN echo 'alias hi="echo hello"' >> ~/.bashrc
#RUN echo 'alias list=". /wrk/list.sh"' >> ~/.bashrc
