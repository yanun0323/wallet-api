Wallet-Api
===

RESTful API for money transfers between wallet.

Requirements:
---
Go 1.17.6 or higher

Frameworks
---
* REST Service
    * ++[echo](https://echo.labstack.com/guide/)++
* Database
    * ++[GORM](https://gorm.io/docs/index.html)++ ( use SQLite ) 
* Test
    * ++[testify](https://github.com/stretchr/testify/blob/master/assert/assertions_test.go)++



Available Services
---
* starts a server on localhost port 1323 (by default). 
http://localhost:1323

|  Method  | Path                |    Usage                         |
| -------- | ------------------- | -------------------------------- |
| GET      | /wallet             | get all wallets                  |
| GET      | /wallet/{walletId}  | get wallet by walletId           |
| POST     | /wallet             | create a new wallet              |
| PUT      | /wallet/{walletId}  | deposite money in wallet         |
| PUT      | /wallet             | transfer money between 2 wallets |
| DELETE   | /wallet/{walletId}  | delete wallet by walletId        |


Http Status
---
* 200 OK: The request has succeeded
* 201 Created: The resourse has created
* 400 Bad Request: The request could not be understood by the server
* 404 Not Found: The requested resource cannot be found
* 500 Internal Server Error: The server encountered an unexpected condition

Sample JSON for Account
---
* ### GET : get all wallets 
    `/wallet` 
    ```
    none
    ```
    
* ### GET : get wallet by walletId
    `/wallet/{walletId}` 
    ```
    none
    ```

* ### POST : create a new wallet
    `/wallet `
    ```
    {
      "walletId": "4123567891234567",
      "balance": 100.00
    }
    ```
* ### PUT : deposite money in wallet
    `/wallet/{walletId} ` 
    ```
    {
      "amount": 100.00
    }
    ```
* ### PUT : transfer money between 2 wallets
    `/wallet ` 
    ```
    {
      "walletFromId": "4123567891234567",
      "walletToId": "4123567891234333",
      "amount": 100.00
    }
    ```
        
* ### DELETE : delete wallet by walletId
    `/wallet/{walletId} ` 
    ```
    none
    ```
