### Sample Chaincode for and Liquidity Optimization Mechanism System (LOM)

This chaincode for Hyperledger Fabric 1.01 written in GoLang is a first rough rewrite of IBM's example02 in the hyperledger Fabric samples. In this sample two accounts are created and via `invoke` and `query` transfers and postings of the accounts can be executed. In this version the user can create several accounts and list them.

In this variant I have added the feature of creating accounts via `account` and accounts can be listed by issueing a `list` command to the peer. Moreover I have added a small batch script `update` for updating the chaincode so DLT network does not have to be restarted for each try-out of new functionality in the chaincode. Also I have changed the internal name of the `invoke` to `transfer`which I think better reflects what is happening.

The files are organised in three folders

 - Scripts: Holds the main scriptfile which sets up channel, peers, certificates and chaincode
 - Workdir: Holds script files when inside the CLI and sending requests to PEERS
 - Chaincode/lom_simple: Hold the chaincode for this example emulating a very simple account system

In the Chaincode I have added a simple Struct reflecting a simple account with accountID, owner, balance and limit. The owner and limit currently do nothing as they are not implemented in the business logics.

One extension could be to ensure that only users who have created and account - can send money from the account. Also there could be limits to how many tokens will be on an account when created. 

