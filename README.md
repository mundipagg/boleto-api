[![GoDoc](https://godoc.org/github.com/mundipagg/boleto-api?status.svg)](https://godoc.org/github.com/mundipagg/boleto-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/mundipagg/boleto-api)](https://goreportcard.com/report/github.com/mundipagg/boleto-api)

What is the Online Registered "Boleto" API?
--------------

BoletoOnline is an API for boleto's online register in banks and boleto's creation for payments.

Currently, we support the following banks:
* Banco do Brasil
* Caixa
* Citi
* Santander 
* Bradesco
* Itaú (coming soon)

The integration order will follow the list above but we may have changes considering our clients demands.

API Building
--------------

The API was developed using GO language and therefore it is necessary to install the language tools in case you need to compile the application from the source.

GO can be downloaded [here](https://golang.org/dl/)

Before cloning the Project, you should create the file path inside $GOPATH

	% mkdir -p "$GOPATH/src/github.com/mundipagg"
	% cd $GOPATH/src/github.com/mundipagg 
	% git clone https://github.com/mundipagg/boleto-api


Before compiling the application, you should install the [Glide](http://glide.sh/), which is the application dependency manager.

After installing GO, do:

	% cd devops
	% ./build

The build.sh script will download the application dependencies and install wkhtmltox, which is necessary to create boleto's in PDF format.

Running the application
-------------

Running the API with default configurations
Eg:

Linux (*NIX):

	% ./boleto-api
	
Windows:

	% boleto-api.exe

If you want to run the API in dev mode, which will load all standard environment variables, you should execute the application like this:

	% ./boleto-api -dev

In case you want to run the application in mock mode using in memory database instead of bank integration, you should use the mock option:

	% ./boleto-api -mock

In case you want to run the application with log turned off, you should use the option -nolog:

	% ./boleto-api -nolog

You can combine all these options and, in case you want to use them altogether, you can simply use the -airplane-mode option

	% ./boleto-api -airplane-mode

By default, boleto api will up and running a https server but you can run in http mode with the following option

	% ./boleto-api -http-only
	

Using Boleto API Online
------------------

You can use [Postman](https://chrome.google.com/webstore/detail/postman/fhbjgbiflinjbdggehcddcbncdddomop) to request the API's services or even the curl
See following examples
### Banco do Brasil
```curl
% curl -X POST \
  http://localhost:3000/v1/boleto/register \
  -d '{
    "Authentication" : {
        "Username":"user",
        "Password":"pass"
    },
    "Agreement":{
        "AgreementNumber":11111,
        "WalletVariation":19,
        "Wallet":17,
        "Agency":"123",
        "AgencyDigit":"2",
        "Account":"1231231",
        "AccountDigit":"3"
    },
    "Title":{
      "ExpireDate": "2017-05-25",
        "AmountInCents":200,
        "OurNumber":101405187,
        "Instructions":"Instruções"
    },
    "Buyer":{
        "Name":"BoletoOnlione",
        "Document": {
            "Type":"CNPJ",
            "Number":"73400584000166"
        },
        "Address":{
            "Street":"Rua Teste",
            "Number": "11",
            "Complement":"",
            "ZipCode":"12345678",
            "City":"Rio de Janeiro",
            "District":"Melhor bairro",
            "StateCode":"RJ"
        }
    },
    "Recipient":{
        "Name":"Nome do Recebedor",
        "Document": {
            "Type":"CNPJ",
            "Number":"12312312312366"
        },
        "Address":{
            "Street":"Rua do Recebedor",
            "Number": "322",
            "Complement":"2º Piso loja 404",
            "ZipCode":"112312342",
            "City":"Rio de Janeiro",
            "District":"Outro bairro",
            "StateCode":"RJ"
        }
    },
    "BankNumber":1
}
```
### Response Banco do Brasil
```
{
  "Url": "http://localhost:3000/boleto?fmt=html&id=g8HXWatft9oMLdTMAqzxbnPYFv3sqgV_KD0W7j8Cy9nkCLZMIK1WH2p9JwP1Jzz4ZtohmQ==",
  "DigitableLine": "00190000090101405100500066673179971340000010000",
  "BarCodeNumber": "00199713400000100000000001014051000006667317",
  "Links": [
    {
      "href": "http://localhost:3000/boleto?fmt=html&id=wOKZh6K_moLwXTW0Xr3oelh9YkYWXdl3VyURiQ-bu6TcuDzxdZI52BnQnuzNpGeh4TapUA==",
      "rel": "html",
      "method": "GET"
    },
    {
      "href": "http://localhost:3000/boleto?fmt=pdf&id=wOKZh6K_moLwXTW0Xr3oelh9YkYWXdl3VyURiQ-bu6TcuDzxdZI52BnQnuzNpGeh4TapUA==",
      "rel": "pdf",
      "method": "GET"
    }
  ]
}

```
### Caixa
```curl
curl -X POST \
  http://localhost:3000/v1/boleto/register \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -H 'postman-token: 1bc1dd5f-cc34-0716-3d56-f18798d3fb39' \
  -d '{
    "BankNumber": 104,
    "Agreement": {
        "AgreementNumber": 200656,
        "Agency":"1679"
    },
    "Title": {
        "ExpireDate": "2017-08-30",
        "AmountInCents": 1000,
        "OurNumber": 0,
        "Instructions": "Mensagem",
        "DocumentNumber": "NPC160517"
    },
    "Buyer": {
        "Name": "TESTE PAGADOR 001",
        "Document": {
            "Type": "CPF",
            "Number": "57962014849"
        },
        "Address": {
            "Street": "SAUS QUADRA 03",
            "Number": "",
            "Complement": "",
            "ZipCode": "20520051",
            "City": "Rio de Janeiro",
            "District": "Tijuca",
            "StateCode": "RJ"
        }
    },
    "Recipient": {
        "Document": {
            "Type": "CNPJ",
            "Number": "00732159000109"
        }
    }
}'
```
### Response
```
{
    "id": "e1EVv1KRwuGX6OXOo7PNGYR-ePD1VPtjv5iqya1LJiLiaIKozN11YMiePNk-WebdgP4eIA==",
    "digitableLine": "10492.00650 61000.100042 09922.269841 3 72670000001000",
    "barCodeNumber": "10493726700000010002006561000100040992226984",
    "ourNumber": "14000000099222698",
    "links": [
        {
            "href": "https://200.201.168.67:8010/ecobranca/SIGCB/imprimir/0200656/14000000099222698",
            "rel": "pdf",
            "method": "GET"
        }
    ]
}
```
In case of Caixa, the impress of boleto will be handled by Caixa. So API will return the Caixa's boleto URL.


The API's response will have the following pattern if any error occur:
```
{
  "Errors": [
    {
      "Code": "MPExpireDate",
      "Message": "Data de expiração não pode ser menor que a data de hoje"
    }
  ]
}
```

API Installation
-----------------

For installation of the executable it's only necessary environment variables configured and the compiled application.

Edit file $HOME/.bashrc.sh
```
    export API_PORT="3000"
    export API_VERSION="0.0.1"
    export ENVIRONMENT="Development"
    export SEQ_URL="http://example.mundipagg.com"
    export SEQ_API_KEY="API_KEY"
    export ENABLE_REQUEST_LOG="false"
    export ENABLE_PRINT_REQUEST="true"
    export URL_BB_REGISTER_BOLETO="https://cobranca.desenv.bb.com.br:7101/registrarBoleto"
    export URL_BB_TOKEN="https://oauth.desenv.bb.com.br:43000/oauth/token"
    export MONGODB_URL="10.0.2.15:27017"
    export APP_URL="http://localhost:8080/boleto"
```
    % go build && mv boleto-api /usr/local/bin

Then the application will be installed locally.

API Installation using Docker
-----------------

Before deploying, you should open the file [docker-compose](/devops/docker-compose.yml) and configure the information which is relevant to the environment. After setting up the docker-compose, you can install the application using the file deploy.sh

    % cd devops
    % ./build.sh
    % ./deploy.sh . local

The script will create the Docker's volume directories, compile the application, mount the API's and MongoDB's images and upload the containers.
For more information about docker-compose, see [doc](https://docs.docker.com/compose/). 

The parameters sent to the script show that the deploy will run locally and not via TFS. In case you don't send the argument "local", the script will use docker-compose.release.yml.

After being deployed, the application can be stoped or started.
    
    % cd devops/
    % ./stop.sh
    % ./start.sh
 
Backup & Restore
----------------- 

Run the following command to backup the database:

    % cd devops
    % ./doBackup.sh

The generated backups, by default, will be stored in the diretory `$HOME/backups` with the name `bck_boleto-api-YYYY-MM-DD.tar`.
To restore a backup:
    
    % cd devops
    % ./doRestore.sh

To do the restoration of the database it'll be asked for the restoration file date in the following pattern: `YYYY-MM-DD`.
    
Contributing
-----------------

To contribute, see [CONTRIBUTING](CONTRIBUTING.md)


Using Mongo in Dev
--------------------
You can use mongo in dev on Docker by using the following command

     % docker run --name mongo-boleto -p 27017:27017 -d mongo

Source Code Layout
---

The application root contains only the file main.go and some config and documentation files.

In the root, we have the following packages:

* `api`: Rest Controllers
* `auth`: Bank authentication
* `bank`: Boleto's register interface
* `bb`: Implementation of Banco do Brasil
* `caixa`: Implementation of Caixa
* `citibank`: Implementation of Citibank
* `boleto`: User boleto's creation
* `cache`: Database (key value) in-memory used only when the application is running in mock mode
* `config`: Application config
* `db`: Database persistency
* `devops`: Contains the upload, deploy, backup and restore files from the application
* `validations`: Basic data validations
* `log`: Application log
* `models`: Application's data structure
* `parser`: XML parser
* `test`: Tests utilitaries
* `tmpl`: Template utilitaries
* `util`: Miscellaneous utilitaries
* `integrationTests`: Contains all black box tests
* `vendor`: Thirdpart libraries

For more information
-----------------

See [FAQ](./FAQ.md)
