# Bhojpur Wallet - Data Processing Engine
The Bhojpur Wallet is a platform-as-a-service used as a Service Engine based
on the Bhojpur.NET Platform. It leverages PostgreSQL database for storage.

#### Bounded Contexts
The application uses the following bounded contexts:

1. Admin
2. Agent
3. Merchant
4. Subscriber
5. Transaction
6. Account
7. Statement
8. Tariff
9. Auth
10. Customer

##### 1. Admin Context
The application needs some form of administration by a super user in-charge
with the responsibility of running and maintaining the application to ensure
reliability and stability. This user is an `admin` and is given their own
bounded context. Some responsibilities/actions of this user are:

1. Can login to system or register.
2. Can assign float to a Super Agent.
3. Can configure tariff
4. Can suspend/change status of a customer account
5. Can view/edit/delete customer accounts

As the application grows and scales the administrator context would have more
responsibilities.

1. The application would need more than one administrator and more so more
than one category of administrators. Some examples of administrators with
their roles include:

    i. Customer Care - is a part admin who would assist customers with
    information about the system and troubleshoot problems.
    
    ii. Finance - is a part admin whose responsibilities would be financial
    and accounting aspect in the system.
    
    iii. IT - an admin whose responsible for the infrastructure that the
    system runs on.

##### 2. Agent Context
We have developed a wallet and money transfer service for a Telco and we have
been given the go ahead by the Central Bank to deploy the application and get
some customers to use our system. There are however some initial steps the
business has to perform to start onboarding new customers. The system should
be able to have some level of autonomy when it comes to the flow of money.
That is where agents come in.

> The initial obstacle in the pilot is gaining the Agent’s trust and
> encouraging them to process cash withdrawals and agent training.
>

Our first initial steps before we can roll out a Payment Wallet service

1. Acquire a business entity licensed to hold public money (i.e. a Bank).
2. Create a Super Agent(s) whose task would be depositing money to our bank
account.
3. Once a Super Agent deposits to our account, we assign them with an
equivalent amount of float they can sell to other agents.
4. When we onboard an ordinary Agent, they will have a balance of zero, and
they will approach the Super Agent to get float.

The Agents are important customers to the system. They can also have various
categories depending on the business use case. For example, we have two types
of agents:

1. Super Agent
2. Ordinary Agent

##### 3. Merchant Context
Our system has two kinds of Merchants.

1. A merchant to provides utility services to their customers
2. A merchant that sells goods and services to end customers

Both the merchants have unique ways of how customers pay for their good and
services. However, both the merchants have an account number.

1.  Pay bill number as an Account number for a merchant - Customer provides
the `pay bill number`, `a customer account number` and the `amount`.
2.  Till number as an account number for a merchant - Customer provides the
`till number` and `amount`.

`Pay bill number` is usually given to utility companies that need to identify
from whom the payment is coming from by the `customer account number`.

`Till number` is usually given to small scale traders that want to accept
payment via our system from their customers.

For example, we stick to one general Merchant that accepts payments.

##### 4. Subscriber Context
A subscriber does not have much going on. They can authenticate and perform a
transaction.

##### 5. Transaction Context
Contains all business logic in regard to transactions happening in the system.
It enforces the transaction rules and business policy.

Business Policies:

1. A transaction cannot happen between identical customers i.e. a customer
cannot transact with themselves
2. A deposit cannot be done by customer none other than an Agent
3. A customer cannot perform a withdrawal with no other customer than an Agent
4. A super agent is however only allowed to do deposits for other agents only
5. Customers are not allowed to deposit, withdraw or transfer money below the
minimum amount allowed
6. Apply transaction fee as per the tariff configured

##### 6. Account Context
The main responsibility of this context is managing customer accounts/wallets.
Responsibilities:

1. Updating account balances, credit/debit accounts
2. Updating system ledger after changing account balances
3. 

##### 7. Statement Context
The main responsibility of this context is managing the system ledger. If we
scale the system, we can view this ledger as the statements/transactions event
store. Borrowing from `event sourcing` design, our statement context is a record
of every event with customer transactions.

##### 8. Tariff Context
This context has a responsibility of configuring and maintaining the tariff used
in various transactions.

##### 9. Auth Context
The system has four different types of users, `admin`, `agent`, `subscriber` and
`merchant`. The auth context is responsible for authenticating and authorizing
these users into the system.

##### 10. Customer Context
This context is mainly an aggregator of the `agent`, `merchant` and `subscriber`
contexts. It exposes common functionality for, which can be used in other core
contexts.


## Installation

To begin with, the application uses PostgreSQL database as the backend database.

### Database Installation

The application uses PostgreSQL as the database server. To setup PostgreSQL on
local your machine using Docker.

Get the official PostgreSQL docker image.
```bash
$ docker pull postgres
``` 

then create a container from the image with the following variables
```bash
$ docker create \
--name mpesa-db \
-e POSTGRES_USER=bhojpur \
-e POSTGRES_PASSWORD=bhojpur \
-p 5432:5432 \
postgres
```

Run the following command to start the container
```bash
$ docker start bhojpur-db
```

### Application Installation

Running the application is very simple. First, we need to copy and create our
configuration.

#### Lets begin. Cloning ...

```bash
$ git clone https://github.com/bhojpur/wallet.git
```

#### Configuring
```bash
$ cd wallet
$ cp config.yml.example config.yml
```

This configuration file looks something like this

```yaml
database:
  host: "127.0.0.1"
  port: "5432"
  user: "bhojpur"
  password: "bhojpur"
  dbname: "bhojpur"

app_secret_key: "eQig7GS4cHO2su"
```

You can change the config variables depending on your database setup. I have
followed the default setup shown at database installation step.

#### Building and running

##### Using the Binary
```bash
$ mkdir bin
$ go build -o bin/wallet-server 
$ ./bin/wallet-server
```

It will install all dependencies required and produce a binary for your platform.

##### Using the Dockerfile
Make sure you have docker installed and working properly.

```bash
$ docker build -t bhojpur-wallet:latest .
$ docker container create --network=host --name wallet-server --restart unless-stopped wallet-server
$ docker container start wallet-server
```

The server will start at port `6700`.

Enjoy.

## API Usage

A description of the RESTful Web Services APIs.

### Endpoints

All the routes exposed in the application are defined in this function
```go
func apiRouteGroup(api fiber.Router, domain *registry.Domain, config app.Config) {

	api.Post("/login/:user_type", user_handlers.Authenticate(domain, config))
	api.Post("/user/:user_type", user_handlers.Register(domain))

	// create group at /api/admin
	admin := api.Group("/admin", middleware.AuthByBearerToken(config.Secret))
	admin.Post("/assign-float", user_handlers.AssignFloat(domain.Admin))
	admin.Post("/update-charge", user_handlers.UpdateCharge(domain.Tariff))
	admin.Get("/get-tariff", user_handlers.GetTariff(domain.Tariff))
	admin.Put("/super-agent-status", user_handlers.UpdateSuperAgentStatus(domain.Agent))

	// create group at /api/account
	account := api.Group("/account", middleware.AuthByBearerToken(config.Secret))
	account.Get("/balance", account_handlers.BalanceEnquiry(domain.Account))
	account.Get("/statement", account_handlers.MiniStatement(domain.Statement))

	// create group at /api/transaction
	transaction := api.Group("/transaction", middleware.AuthByBearerToken(config.Secret))
	transaction.Post("/deposit", transaction_handlers.Deposit(domain.Transactor))
	transaction.Post("/transfer", transaction_handlers.Transfer(domain.Transactor))
	transaction.Post("/withdraw", transaction_handlers.Withdraw(domain.Transactor))
}
```

The routes are mounted on the prefix `/api` so your requests should point to
```
POST /api/login/<user_type>                     <-- user_type can be either of agent, admininistrator, merchant, subscriber
POST /api/user/<user_type> # for registration   <-- user_type can be either of agent, admininistrator, merchant, subscriber
POST /api/admin/assign-float
POST /api/admin/update-charge
GET /api/admin/get-tariff
PUT /api/admin/super-agent-status
GET /api/account/balance
POST /api/account/statement
POST /api/transaction/deposit
POST /api/transaction/transfer
POST /api/transaction/withdraw
```

#### To Register

The APIs can be used to register four types of users: `admin`, `agent`, `merchant`
and `subscriber`

#### Admin Registration
An admin can be registered to the API with the following `POST` parameters

`firstName`, `lastName`, `email`, `password`

Curl request example
```bash
curl --request POST \
  --url http://localhost:6700/api/user/administrator \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data firstName=Admin \
  --data lastName=Batua \
  --data email=admin_wallet@bhojpur.net \
  --data password=welcome
```

Response example
```json
{
  "status": "success",
  "message": "user created",
  "data": {
    "userID": "bc8f933b-b0aa-4049-9caf-dfe73008bc24",
    "userType": "administrator"
  }
}
```

##### Agent Registration
At minimum, you need to create two agents, one of which will become a
`super agent`. An agent can be registered to the APIs with the following
`POST` parameters

`firstName`, `lastName`, `email`,  `phoneNumber`, `password`

Curl request example
```bash
curl --request POST \
  --url http://localhost:6700/api/user/agent \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data firstName=Agent \
  --data lastName=Batua \
  --data email=agent_wallet@bhojpur.net \
  --data phoneNumber=16282004199 \
  --data password=welcome
```

Response example
```json
{
  "status": "success",
  "message": "user created",
  "data": {
    "userID": "cca4d238-74ae-4d47-aae8-b0ab952aac28",
    "userType": "agent"
  }
}
```

##### Merchant Registration
A merchant can be registered to the api with the following `POST` parameters

`firstName`, `lastName`, `email`,  `phoneNumber`, `password`

Curl request example
```bash
curl --request POST \
  --url http://localhost:6700/api/user/merchant \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data firstName=Merchant \
  --data lastName=Batua \
  --data email=merchant_wallet@bhojpur.net \
  --data phoneNumber=16282004199 \
  --data password=welcome
```

Response example
```json
{
  "status": "success",
  "message": "user created",
  "data": {
    "userID": "c3a71820-ef66-74d9-adc8-f365a234ed5c",
    "userType": "merchant"
  }
}
```

##### Subscriber Registration
A subscriber can be registered to the api with the following `POST` parameters

`firstName`, `lastName`, `email`,  `phoneNumber`, `password`

Curl request example
```bash
curl --request POST \
  --url http://localhost:6700/api/user/subscriber \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data firstName=Subscriber \
  --data lastName=Batua \
  --data email=subscriber_wallet@bhojpur.net \
  --data phoneNumber=16282004199 \
  --data password=welcome
```

Response example
```json
{
  "status": "success",
  "message": "user created",
  "data": {
    "userID": "cf8d7f25-367e-4ac7-8b5f-eaa7608e6c3f",
    "userType": "subscriber"
  }
}
```

#### To Login
You can use the following `POST` parameters for login with any of the four users
`email`, `password`

Curl request example for subscriber login
```bash
curl --request POST \
  --url http://localhost:6700/api/login/subscriber \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data email=subscriber_wallet@bhojpur.net \
  --data password=welcome
```

Response example

```json
{
  "userId": "cf8d7f25-367e-4ac7-8b5f-eaa7608e6c3f",
  "userType": "subscriber",
  "token": "eyJhbGciOiJIUzI1NiIsInR4cCI6IkpXVCJ8.eyJ1c2VyIjp7InVzZXJJZCI6ImNmOWQ5ZjI4LTM1N2UtNGFjNy05YjVmLWVhYTg2MDllNmMyZiIsInVzZXJUeXBlIjoic3Vic2NyaWJlciJ9LCJleHAiOjE2MDU0MTYwNjAsImlhdCI6MTYwNTM5NDQ2MH0.lAJ4WpF2Mnfg52iuTOoPV8nvbHV3JrMQOC-5xXrQ5EE"
}
```

**NOTE**: The remaining endpoints require the token acquired above for authentication


#### Initial Steps Before Transacting
There are some initial setups that need to be done before you can begin doing transactions.

##### 1. Creating a Super Agent
Before you can start transacting, you need to login as an administrator and
create a Super Agent by changing the status of an existing agent. When registering an
Agent, you ought to have created at minimum two agents. It is now that we need make
one of those agents a Super Agent.

The following endpoint is used to update the `super agent status` of an Agent.

`PUT /api/admin/super-agent-status` requires the following post parameters: `email`

Curl request example
```bash
curl --request PUT \
  --url http://localhost:6700/api/admin/super-agent-status \
  --header 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJJZCI6Ijc2YmM0YWEzLTAyNWQtNGQ1YS1hNWZiLWY1NDk1NTdmNjM0YSIsInVzZXJUeXBlIjoiYWRtaW5pc3RyYXRvciJ9LCJleHAiOjE2MDU0NTE4MDUsImlhdCI6MTYwNTQzMDIwNX0.8lTWl9hGr9GTST7WpEpzKdm_gqhMkf4qUellLx4o5bw' \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data email=agent_wallet@bhojpur.net
```

Response example

```json
{
  "status": "success",
  "message": "Super Agent Status updated"
}
```

##### 2. Assigning Float
Login as an Administrator, you need to assign float to your `super-agent` using
the following endpoint

`POST /api/admin/assign-float`

Curl request example
```bash
curl --request POST \
  --url http://localhost:6700/api/admin/assign-float \
  --header 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJJZCI6ImE3OGVjNjNhLTA0ZWItNDAzNC1iZmVkLTBhNmMwMjU3ZTJlNCIsInVzZXJUeXBlIjoiYWRtaW5pc3RyYXRvciJ9LCJleHAiOjE2MDUzMjExNzUsImlhdCI6MTYwNTI5OTU3NX0.4fGlMJQB-eylKwOAwa4d16nVQQt3uYgwPbUjYt7j9zA' \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data accountNo=agent_wallet@bhojpur.net \
  --data amount=100000
```

Response example
```json
{
  "status": "success",
  "message": "Float has been assigned.",
  "data": {
    "balance": 100000
  }
}
```

##### 3. Transfer Float to Agents
The `super-agent` is limited to depositing to Agents only. You will need to
transfer the acquired float to other Agents you have registered.

##### 4. Configure Tariff
The default tariff in the system is set to zero amount for all chargeable
transactions. You could begin testing transactions using the default tariff and
later choose to configure your own tariff. Choose your poison :-).

You can configure a tariff by updating the available charges. The system doesn't
allow you to add any other charge band.

`GET /api/admin/get-tariff` - use this endpoint to get the available configured
transaction charges

Response example
```json
{
  "status": "success",
  "message": "Tariff retrieved",
  "data": [
    {
      "id": "acf3e6bf-c9de-45b4-a8b6-bf97f92b783a",
      "txnOperation": "WITHDRAW",
      "srcUserType": "subscriber",
      "destUserType": "agent",
      "fee": 0
    },
    {
      "id": "0e5a4aaa-135a-4464-96c9-d021f769bdb7",
      "txnOperation": "WITHDRAW",
      "srcUserType": "merchant",
      "destUserType": "agent",
      "fee": 0
    },
    {
      "id": "243e7ecc-c2dd-41bb-9953-1278050bfb64",
      "txnOperation": "WITHDRAW",
      "srcUserType": "agent",
      "destUserType": "agent",
      "fee": 0
    },
    {
      "id": "f8835176-316c-49de-b001-687e2c4a338d",
      "txnOperation": "TRANSFER",
      "srcUserType": "agent",
      "destUserType": "agent",
      "fee": 0
    },
    {
      "id": "4edeb6d0-37cd-4c67-997a-0b3fa93b722d",
      "txnOperation": "TRANSFER",
      "srcUserType": "subscriber",
      "destUserType": "subscriber",
      "fee": 0
    },
    {
      "id": "94c0ae8b-a131-41b9-b5af-5235b8926fa4",
      "txnOperation": "TRANSFER",
      "srcUserType": "merchant",
      "destUserType": "subscriber",
      "fee": 0
    },
    {
      "id": "450e4baa-58c3-41b3-abe5-a55555492e0c",
      "txnOperation": "TRANSFER",
      "srcUserType": "subscriber",
      "destUserType": "merchant",
      "fee": 0
    },
    {
      "id": "3623a89f-c496-41c8-b6c9-73429cc4ef9d",
      "txnOperation": "TRANSFER",
      "srcUserType": "agent",
      "destUserType": "merchant",
      "fee": 0
    }
  ]
}
```

`POST /api/admin/update-charge` - use this endpoint to update a charge using
its `id`.

You need the following `POST` parameters

`amount`, `chargeId` - The amount should be in `cents`.

Curl request example
```bash
curl --request POST \
  --url http://localhost:6700/api/admin/update-charge \
  --header 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJJZCI6Ijc2YmM0YWEzLTAyNWQtNGQ1YS1hNWZiLWY1NDk1NTdmNjM0YSIsInVzZXJUeXBlIjoiYWRtaW5pc3RyYXRvciJ9LCJleHAiOjE2MDU0NTE4MDUsImlhdCI6MTYwNTQzMDIwNX0.8lTWl9hGr9GTST7WpEpzKdm_gqhMkf4qUellLx4o5bw' \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data amount=1050 \
  --data chargeId=acf3e6bf-c9de-45b4-a8b6-bf97f92b783a
```

Response example
```json
{
  "status": "success",
  "message": "charge configured"
}
```

#### Performing Transactions
While configuring a charge requires you to provide the amount in `paisas`,
performing transactions requires the amount to be in whole units i.e. `rupees`

Transacting also requires you to provide an `accountNo`, use the `email` of
the customer as the `accountNo`

`customerType` can be either of `agent`, `merchant` or `subscriber`

##### 1. To Deposit
A deposit is only done by an `agent`. You need an `agent` token to perform
this transaction.

You need the following `POST` parameters

`amount`, `accountNo` and `customerType`

Curl request example
```bash
curl --request POST \
  --url http://localhost:6700/api/transaction/deposit \
  --header 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJJZCI6ImNjYTdkMjI3LTc0YWUtNGQ0Ny1hYWU4LWEwYWI5NTJhYWMyOCIsInVzZXJUeXBlIjoiYWdlbnQifSwiZXhwIjoxNjA1MzIxODEyLCJpYXQiOjE2MDUzMDAyMTJ9.jFLfjScuvHaOV68n11sRticy2ntzQRhwbNq5E4sPmQI' \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data amount=400 \
  --data accountNo=subscriber_wallet@bhojpur.net \
  --data customerType=subscriber
```

Response example

```json
{
  "status": "success",
  "message": "Success",
  "data": {
    "message": "Transaction under processing. You will receive a message shortly."
  }
}
``` 

##### 2. To Withdraw
You need the following `POST` parameters

`amount`, `agentNumber`

Use agent email for `agentNumber`

Curl request example
```bash
curl --request POST \
  --url http://localhost:6700/api/transaction/withdraw \
  --header 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJJZCI6ImNmOWQ4ZjI4LTM1N2UtNGFjNy05YjVmLWVhYTg2MDllNmMyZiIsInVzZXJUeXBlIjoic3Vic2NyaWJlciJ9LCJleHAiOjE2MDU0MTYwNjAsImlhdCI6MTYwNTM5NDQ2MH0.lAJ4WpF2Mnfg52iuTOoPV8nvbHV3JrMQOC-5xXrQ5EE' \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data amount=40 \
  --data agentNumber=agent_wallet@bhojpur.net
```

Response example

```json
{
  "status": "success",
  "message": "Success",
  "data": {
    "message": "Transaction under processing. You will receive a message shortly."
  }
}
```

##### 3. To Transfer
You need the following `POST` parameters

`amount`, `accountNo` and `customerType`

Curl request example
```bash
curl --request POST \
  --url http://localhost:6700/api/transaction/transfer \
  --header 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJJZCI6ImNmOWQ4ZjI4LTM1N2UtNGFjNy05YjVmLWVhYTg2MDllNmMyZiIsInVzZXJUeXBlIjoic3Vic2NyaWJlciJ9LCJleHAiOjE2MDUzMjE5MDMsImlhdCI6MTYwNTMwMDMwM30.vLiHdNTr4onTVqUZbLbdpwgbH98VYzHJJU-JKtFOHVg' \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data amount=30 \
  --data accountNo=merch_wallet@bhojpur.net \
  --data customerType=merchant
```

Response example

```json
{
  "status": "success",
  "message": "Success",
  "data": {
    "message": "Transaction under processing. You will receive a message shortly."
  }
}
```


#### To Query Balance
This is just a `GET` request, no params

Curl request example
```bash
curl --request GET \
  --url http://localhost:6700/api/account/balance \
  --header 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJJZCI6ImNmOWQ4ZjI4LTM1N2UtNGFjNy05YjVmLWVhYTg2MDllNmMyZiIsInVzZXJUeXBlIjoic3Vic2NyaWJlciJ9LCJleHAiOjE2MDUzNjY3NTMsImlhdCI6MTYwNTM0NTE1M30.-Piib6bXzYqb0S8nLo76SBTyGmWi7UPUMExptIcqBZI'
```

Response example

```json
{
  "status": "success",
  "message": "Your current balance is 690",
  "data": {
    "userID": "cf8d7f25-367e-4ac7-8b5f-eaa7608e6c3f",
    "balance": 690
  }
}
```

#### To Get Mini Statement
This is just a `GET` request, no params

Curl request example
```bash
curl --request GET \
  --url http://localhost:6700/api/account/statement \
  --header 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJJZCI6Ijk4YmNmMmY1LWRiY2ItNDk1NS04NTU0LTc0OWYxMTVhZjU5OCIsImVtYWlsIjoiIn0sImV4cCI6MTYwNDA2OTE0MywiaWF0IjoxNjA0MDQ3NTQzfQ.IYyclrC66aweehs_A4Sigmc83a27udmPofM2yOeut9Q'
```

Response example

```json
{
  "status": "success",
  "message": "mini statement retrieved for the past 5 transactions",
  "data": {
    "message": "mini statement retrieved for the past 5 transactions",
    "userID": "cf8d7f25-367e-4ac7-8b5f-eaa7608e6c3f",
    "transactions": [
      {
        "transactionId": "97c3ff6d-72d5-479d-8838-85a5c32985a2",
        "transactionType": "DEPOSIT",
        "createdAt": "2020-11-14T01:59:28.613007+03:00",
        "creditedAmount": 400,
        "debitedAmount": 0,
        "userId": "cf8d7f25-367e-4ac7-8b5f-eaa7608e6c3f",
        "accountId": "63978e26-9c0d-40eb-a24b-d1ae51e21942"
      },
      {
        "transactionId": "4be4a008-b18e-4d4b-95d3-58b660d5b931",
        "transactionType": "TRANSFER",
        "createdAt": "2020-11-14T01:59:05.949066+03:00",
        "creditedAmount": 0,
        "debitedAmount": 30,
        "userId": "cf9d8f28-357e-4ac7-9b5f-eaa8609e6c2f",
        "accountId": "63978e26-9c0d-40eb-a24b-d1ae51e21942"
      },
      {
        "transactionId": "45da6c6a-03d8-4d58-849a-fd80bbfabbb4",
        "transactionType": "TRANSFER",
        "createdAt": "2020-11-14T01:57:04.622507+03:00",
        "creditedAmount": 0,
        "debitedAmount": 40,
        "userId": "cf9d8f28-357e-4ac7-9b5f-eaa8609e6c2f",
        "accountId": "63978e26-9c0d-40eb-a24b-d1ae51e21942"
      }
    ]
  }
}
```

## Testing

Tests have been written for the application. 

An approach to testing this application would be something in the following lines.

1. Test the code in the interactor files
2. Test the code in the repository files

The above files carry the bulk of the behaviour of the whole application, they
are the business logic of the application. The rest of the files are just
implementation details that could change rapidly and the tests written for them
would certainly fail after change.

e.g.
the HTTP handler functions in the application uses
[gofiber](https://github.com/gofiber/fiber), writing unit tests for them is good
but not desired, because [gofiber](https://github.com/gofiber/fiber) can be
replaced with [mux](https://github.com/gorilla/mux) easily and that would break
your tests.