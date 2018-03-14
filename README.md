## Sample Chaincode for a Liquidity Optimization Mechanism (LOM)

This chaincode - for Hyperledger Fabric 1.01 written in GoLang - is a first rough rewrite of IBM's example02 in the hyperledger Fabric samples. In this sample two accounts are created and via `invoke` and `query` transfers and postings of the accounts can be executed. In this version the user can create several accounts and list them.

In this variant I have added the feature of creating accounts via `add` and accounts can be listed by issuing a `list` command to the peer. Moreover I have added a small batch script `update` for updating the chaincode so DLT network does not have to be restarted for each try-out of new functionality in the chaincode. Also I have changed the internal name of the `invoke` to `transfer`which I think better reflects what is happening. Finally I created a `history` call which takes one argument - an account id - and returns the history of states for that account as JSON. 

The files are organised in three folders:

 - Scripts: Holds the main scriptfile which sets up channel, peers, certificates and chaincode
 - Workdir: Holds script files when inside the CLI and sending requests to PEERS
 - Chaincode/lom_simple: Hold the chaincode for this example emulating a very simple account system

In the Chaincode I have added a simple Struct reflecting a simple account with accountID, owner, balance and limit. The owner and limit currently do nothing as they are not implemented in the business logics.

One extension could be to ensure that only users who have created and account - can send money from the account. Also there could be limits to how many tokens will be on an account when created. 

#### Observation on Hyperledger and couchDB
When I run first-network with the option `-s couchDB` the internal DB is based on couchDB. With couchDB the user can view current world state (of the accounts). However, the user can also change data - that is world state - in couchDB. I believe this is somewhat of a flaw - as this is done outside of the chaincode - outside of the smartcontract  - and thus outside of the common agreement which everybody follows. 

I would suspect that the integrity of the blocks would somehow hinder this - but for reasons I do not understand Hyperledger Fabric allow changes directly to world state via couchDB. Hyperledger must somehow produce blocks underneath. I have not been able to verify this, as I have not yet added and Fabric Blockchain Explorer to the project. I have checked if invokes/changes to the blockchain via the smartcontract would override the couchDB changes if run after - but that is not the case. One can simply change world state without a smartcontract.

