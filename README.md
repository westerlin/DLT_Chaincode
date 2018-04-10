## Sample Chaincode for a Liquidity Optimization Mechanism (LOM)

This chaincode - for Hyperledger Fabric 1.01 written in GoLang - is a first rough rewrite of IBM's example02 in the hyperledger Fabric samples. In this sample two accounts are created and via `invoke` and `query` transfers and postings of the accounts can be executed. In this version the user can create several accounts and list them.

You create the crypto material (certificates) by issuing this command:

`./byfn.sh -m generate -c minidlt`


NOTE! This is not nice - but I have hardcoded the cannelname __minidlt__ elsewhere in the code. Not pretty but haven't bothered fixing this.

You start the network by issuing the command:

`./byfn.sh -m up -c minidlt -s couchdb`

In this variant I have added the feature of creating accounts via `add` and accounts can be listed by issuing a `list` command to the peer. Moreover I have added a small batch script `update` for updating the chaincode so DLT network does not have to be restarted for each try-out of new functionality in the chaincode. Also I have changed the internal name of the `invoke` to `transfer`which I think better reflects what is happening. Finally I created a `history` call which takes one argument - an account id - and returns the history of states for that account as JSON. 

The files are organised in three folders:

 - Scripts: Holds the main scriptfile which sets up channel, peers, certificates and chaincode
 - Workdir: Holds script files when inside the CLI and sending requests to PEERS
 - Chaincode/lom_simple: Hold the chaincode for this example emulating a very simple account system

In the Chaincode I have added a simple Struct reflecting a simple account with accountID, owner, balance and limit. The owner do nothing. The limit represents a lower limit.

### How to interact ###
The docker-compose also fires up a container containing af HLF client name __CLI__. You can enter this client by writing: 

`docker exec -it cli bash --init-file /wrk/init.sh`

Thereby you enter the client terminal, which will present you with a menu (this is defined in the init.sh file) covering the most important commands. You can always revert to this menu by typing `menu` in the command prompt. 

```

Welcome to MiniDLT network

   DLT commands: 
      transfer <sender> <receiver> <amount>                : Transfers amount fra sender to recevier. Both are account id's.
      query <accountid>                                    : Lists the current balance of account
      add <accountid> <ownerid> <initial balance> <limit>  : Creates an account with initial balance and lower limit.
      list                                                 : Lists all accounts and their status.
      history <accountid>                                  : Lists history for a given account

   Special commands: 
      menu                                                 : Returns to this menu
      policy                                               : Change the policy
      logging <logging level>                              : Sets logging level (CRITICAL | ERROR | WARNING | NOTICE | INFO | DEBUG)
      instantiate <version> <chaincode name>               : Installs and instantiates a new chaincode
      update <version>                                     : Updates a new chaincode
      peer chaincode list --instantiated                   : Lists alle instantiated chaincodes
      peer chaincode list --installed                      : Lists alle installed chaincodes
      peer channel list                                    : Lists alle channels joined

Current Policy is: OR	('Org1MSP.member','Org2MSP.member')

Current Chaincode is: None
```

The client is set up to communicate with the peer0.org1.example.com. Thus, all transactions are originated from this part of the network initally.

I have changed the prompt so it does not show path, data, host etc. This is done to ensure more space for writing long commands (for some reason the CLI-container terminal is fixed at 80 characters across which is a hassle when you want to write long commands). 

After the HLF is started you have to install and instantiate the chaincode. This is done by this command:

`> instantiate 1.0 lomcode lom_simple`

NB! As default the policy is set to __"OR (Org1MSP.member, Org2MSP.member)"__ which basically states that one peer alone can endorse new TX's. By sending the command 

`> policy`

The policy is switched to __"AND (Org2MSP.member)"__ so the opposing organisation _Org2_ is responsible for endorsing what is originated from _Org1_. 


### Possible extensions ###
One extension could be to ensure that only users who have created and account - can send money from the account. Also there could be limits to how many tokens will be on an account when created. 

#### Observation on Hyperledger and couchDB
When I run first-network with the option `-s couchDB` the internal DB is based on couchDB. With couchDB the user can view current world state (of the accounts). However, the user can also change data - that is world state - in couchDB. I believe this is somewhat of a flaw - as this is done outside of the chaincode - outside of the smartcontract  - and thus outside of the common agreement which everybody follows.

I have checked if invokes/changes to the blockchain via the smartcontract would override the couchDB changes if run after - but that is not the case. One can simply change world state without a smartcontract.

It should be noted that this is only possible in the default policy mode - that is where just one peer can endorse a transaction. When another Peer (with it's untampered state database) is verifying an transaction which goes outside business logic - ie. the resulting balance from a transfer falls below limit - then the transaction is rejected. The state databases are however not update and can as such be out of sync.   

