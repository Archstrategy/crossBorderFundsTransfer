# Cross Border Funds Transfer (CBFT)

CBFT is an illustration of how blockchain and hyperledger can be used to model cross border payments involving banks and customers in different currencies. *This is not about any cryptocurrency or its usage for payments.* The proposed solution uses fiat currencies but payments flow through distributed ledgers thus reducing correspondent banks/additional hops in the existing payments networks. If you are interested in reading the *motivation* behind creating this app then  please visit [A case for banks to introduce B2B "real-time" Cross Border Funds Transfer (CBFT) Services for SMEs using Blockchain - Part 1](https://www.linkedin.com/pulse/case-banks-introduce-b2b-real-time-cross-border-funds-raj-m-shimpi/)

You can also access youtube video on howto setup and run this app at https://www.youtube.com/watch?v=w0Ybr5_Gh8A


For our payment network, we will setup  of 1 peer, one orderer service, one CA authority service for certificate issuance and verification and one channel on which our peers will communicate. The initial setup of the network is populated with:

1. Three representative banks - US_Bank, UK_Bank and Japan_Bank. Each one with the domicile country's currency reserves. These reserves are required because one of the competitive advantage we have considered for banks to provide liquidity and thus manage risk.

2. Each bank also has customer accounts. When customers make payments to other customers, it is assumed that the payments will originate in the domiciled currency. For e.g., US customer sending payments in USD to UK customer will have respective account debited, respective bank's reserves depleted and UK customer's account will be credited with GBP after conversion and UK bank's reserves will be increased since they are receiving GBPs.

3. During the setup, forex pairs are populated. This is for the demo purposes only since  forex will be ideally retrieved in real time and perhaps can be provided as an input to the smart contract by the clients or through a forex backend of the banks.

4. We are setting up only one channel and all the banks in the network will subscribe to it. In real scenarios, consider setting up separate channels among group of banks to protect privacy of transactions for competitive purposes.


## CBFT Features
Please download/clone the repo from github available at: https://github.com/Archstrategy/crossBorderFundsTransfer.git. At a high level, package has all the modules and sources to install components. Node.js scripts (clients) are provided for the following functionalities:

1. pay.js - script to execute and invoke a payment between customers at different or same banks

2. createBank.js - create a new bank

3. createCustomer.js - create a new customer associated with an existing bank or a new bank. Logically you should create a bank first, before adding a customer to a new bank.

4. ceateForex.js - create a new forex pair with average exchange rate (bid/ask are not supported)

5. queryBank.js, queryCustomer.js, queryForex.js - to search individual records

6. queryAll.js - to dump the entire ledger


## Directory Structure
```
├── banks-network           fabric network setup
├── blockchain-explorer			blockchain explorer setup
├── chaincode				        chaincode or smart contract for banks blockchain network
├── first-network	Basic fabric network setup
├── enrollAdmin.js					Enrolling admin with CA
├── pay.js									Invoke chaincode for transactions
├── createBank.js						create a new bank
├── createCustomer.js				create a new customer
├── createForex.js					create a new forex pair
├── queryBank.js						query a single bank from ledger
├── queryCustomer.js				query a single customer from ledger
├── queryForex.js						query a single forex pair from ledger
├── queryAll.js							query the entire ledger (dumps all entries from ledger)
├── registerUser.js		 		  register a user with CA
├── setting-up-fabric.pdf		setup fabric from scratch (a prereq for this app to work)
└── startFabric.sh					start this sample
```


## Requirements
Please resolve all issues and get the first network up and running before attempting to install and run this demo. You can find ubuntu cheat sheet to get basic fabric running by following setting-up-fabric.pdf document included in this repo and could be found at the root of the folder structure.

Following are the software dependencies required to install and run hyperledger explorer
* docker 17.06.2-ce [https://www.docker.com/community-edition]
* docker-compose 1.14.0 [https://docs.docker.com/compose/]
* nodejs 6.9.x (Note that v7.x is not yet supported)
* mysql 5.7 or greater
* GO programming language 1.7.x
* git, curl, and other binaries needed to run on windows or OS X.
* Optionals - atom/vscode for editing files, kitematic for docker view.

## Clone Repository

Clone this repository to get the latest using the following command.
1. `git clone https://github.com/Archstrategy/crossBorderFundsTransfer.git`
2. `cd crossBorderFundsTransfer`  

## Setup banks network
* Run the following command to kill any stale or active containers:
`docker rm -f $(docker ps -aq)`

* Clear any cached networks:
`docker network prune`

Run the following command to install the Fabric dependencies for the applications. We are concerned with fabric-ca-client which will allow our app(s) to communicate with the CA server and retrieve identity material, and with fabric-client which allows us to load the identity material and talk to the peers and ordering service.

`npm install`

Launch banks network using the startFabric.sh shell script. This command will spin up our various Fabric entities and launch a smart contract container for chaincode written in Golang:

`./startFabric.sh`

To stream  CA logs, split your terminal or open a new shell and issue the following:

`docker logs -f ca.example.com`

or you can also use kitematic docker viewer to look at all the docker images created and running. You can also look at the logs.

### Creating admin and registering users with CA authority
When we launched banks network, an admin user - admin - was registered with our Certificate Authority. Now we need to send an enroll call to the CA server and retrieve the enrollment certificate (eCert) for this user. Node.js SDK need this cert in order to form a user object for the admin. We will then use this admin object to subsequently register and enroll a new user. Send the admin enroll call to the CA server:

`node enrollAdmin.js`

This program will invoke a certificate signing request (CSR) and ultimately output an eCert and key material into a newly created folder - hfc-key-store - at the root of this project. Our node.js scripts will then look to this location when they need to create or load the identity objects for our various users.
`node registerUser.js`

### Initializing ledger
When banks network was started it installed chaincode on peers. It also initialized the ledger.

### Querying ledger for initial setup entities created

Run the following to query all records on the ledger

`node queryAll.js` This should return ledger in JSON format.
`Query has completed, checking results
Response is  [{"Key":"GBP:GBP", "Record":{"pair":"GBP:GBP","rate":1}}

,{"Key":"GBP:JPY", "Record":{"pair":"GBP:JPY","rate":155}}

,{"Key":"GBP:USD", "Record":{"pair":"GBP:USD","rate":1.35}}

,{"Key":"JPY:GBP", "Record":{"pair":"JPY:GBP","rate":0.0065}}

,{"Key":"JPY:JPY", "Record":{"pair":"JPY:JPY","rate":1}}

,{"Key":"JPY:USD", "Record":{"pair":"JPY:USD","rate":0.0088}}

,{"Key":"JPY_Alice_456", "Record":{"balance":1000000.0,"country":"Japan","currency":"JPY","custID":"456","customerBankID":"Japan_Bank","name":"JPY_Alice"}}

,{"Key":"JPY_John_Doe_123", "Record":{"balance":1000000.0,"country":"Japan","currency":"JPY","custID":"123","customerBankID":"Japan_Bank","name":"JPY_John_Doe"}}

,{"Key":"Japan_Bank", "Record":{"bankID":"Japan_Bank","country":"JAPAN","currency":"JPY","name":"Japan_Bank","reserves":10000000.0}}

,{"Key":"UK_Alice_456", "Record":{"balance":10000,"country":"UK","currency":"GBP","custID":"456","customerBankID":"UK_Bank","name":"UK_Alice"}}

,{"Key":"UK_Bank", "Record":{"bankID":"UK_Bank","country":"UK ","currency":"GBP","name":"UK_Bank","reserves":1000000.0}}

,{"Key":"UK_John_Doe_123", "Record":{"balance":10000,"country":"UK","currency":"GBP","custID":"123","customerBankID":"UK_Bank","name":"UK_John_Doe"}}

,{"Key":"USD:GBP", "Record":{"pair":"USD:GBP","rate":0.75}}

,{"Key":"USD:JPY", "Record":{"pair":"USD:JPY","rate":115}}

,{"Key":"USD:USD", "Record":{"pair":"USD:USD","rate":1}}

,{"Key":"US_Alice_456", "Record":{"balance":10000,"country":"US","currency":"USD","custID":"456","customerBankID":"US_Bank","name":"US_Alice"}}

,{"Key":"US_Bank", "Record":{"bankID":"US_Bank","country":"USA","currency":"USD","name":"US_Bank","reserves":1000000.0}}

,{"Key":"US_John_Doe_123", "Record":{"balance":10000,"country":"US","currency":"USD","custID":"123","customerBankID":"US_Bank","name":"US_John_Doe"}}

]
`

Query a single bank, customer or forex.

`node queryBank.js`
`Response is  {"bankID":"US_Bank","country":"USA","currency":"USD","name":"US_Bank","reserves":1000000.0}
	`

`Query Result: {"name":"Japan_Bank","bankID":"Japan_Bank","country":"JAPAN","currency":"JPY","reserves":1e+07}`


`node queryCustomer.js`		
`Response is  {"balance":1000000.0,"country":"Japan","currency":"JPY","custID":"123","customerBankID":"Japan_Bank","name":"JPY_John_Doe"}
`

`node queryForex.js	`
`Response is  {"pair":"USD:JPY","rate":115}`

Execute a payment between two customers

`node pay.js`

`Assigning transaction_id:  cf0f134a71e02026036ec4e8029a285c433177ea5c8b7143f33a7a3065061104
Transaction proposal was good
Successfully sent Proposal and received ProposalResponse: Status - 200, message - "OK"
info: [EventHub.js]: _connect - options {}
The transaction has been committed on peer localhost:7053
Send transaction promise and event listener promise have completed
Successfully sent transaction to the orderer.
Successfully committed the change to the ledger by the peer
`

Now query the bank and customer to check the credit and debit to respective accounts.
`,{"Key":"JPY_Alice_456", "Record":{"balance":999000,"country":"Japan","currency":"JPY","custID":"456","customerBankID":"Japan_Bank","name":"JPY_Alice"}}`

`,{"Key":"UK_Alice_456", "Record":{"balance":10006.5,"country":"UK","currency":"GBP","custID":"456","customerBankID":"UK_Bank","name":"UK_Alice"}}
`

`,{"Key":"UK_Bank", "Record":{"bankID":"UK_Bank","country":"UK ","currency":"GBP","name":"UK_Bank","reserves":1000006.5}}`

`,{"Key":"Japan_Bank", "Record":{"bankID":"Japan_Bank","country":"JAPAN","currency":"JPY","name":"Japan_Bank","reserves":9999000.0}}`

Insufficient funds if payment is made over the account balance

`error: [client-utils.js]: sendPeersProposal - Promise is rejected: Error: chaincode error (status: 500, message: Insufficent funds in customer: JPY_Alice Customer ID: 456)
`

Use the following to create entries on ledger.. the scripts include a sample format, you can change it to create additional records.

`node createBank.js`

`node creatCustomer.js`

`node createForex.js`

### Setup and start blockchain-explorer

#### setup mysql database
`cd blockchain-explorer`

 Run the database setup scripts located under `db/fabricexplorer.sql`
 `mysql -u<username> -p < db/fabricexplorer.sql`

 #### Running blockchain-explorer

 1. `cd blockchain-explorer`
 2. Modify config.json to update one of the channel
 	* mysql host, username, password details. Don't change anything else, the path to certs are relative to this directory.
 ```json
  "channel": "mychannel",
  "mysql":{
       "host":"127.0.0.1",
       "database":"fabricexplorer",
       "username":"root",
       "passwd":"pwd_for_mysql"
    }
 ```
 3. `npm install`
 4. `./start.sh`

	You can check log.log file `more log.log` and it should only have one line "Please open Internet explorer to access". If there are errors then you will see them in these files.  If error is shown then resolve it and kill the process using `ps -ef| grep main.js` to get the pid of the blockchain explorer process and then issuing `kill pid#` If all goes well and thre are no errors in log.log file, then continue to launch the URL http://localhost:8080 on a browser.

  ### stopping banks network and cleaning up
  After you are done with running the app, shut down the network and remove all docker images and other files creted.
  `cd ./banks-network`
  `./stop,sh`
  `./teardown.sh`

	kill the process running blockchain explorer
	`ps -ef| grep main.js`
	kill <pid#>
