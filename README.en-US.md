[![mundipagg maturity](http://maturityapp.herokuapp.com/maturity.php?project=mundipagg/boleto-api&command=badge_image)](http://maturityapp.herokuapp.com/index.html?project=mundipagg/boleto-api)

[![GoDoc](https://godoc.org/github.com/mundipagg/boleto-api?status.svg)](https://godoc.org/github.com/mundipagg/boleto-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/mundipagg/boleto-api)](https://goreportcard.com/report/github.com/mundipagg/boleto-api)

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/b85953cc9fa84b56822e7e5d91203e91)](https://www.codacy.com/app/mundipagg/boleto-api?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=mundipagg/boleto-api&amp;utm_campaign=Badge_Grade)
[![Maintainability](https://api.codeclimate.com/v1/badges/b9ad683e9d8f87034339/maintainability)](https://codeclimate.com/github/mundipagg/boleto-api/maintainability)

# Boleto API

It is an API for online registration of boleto at banks and creation of boleto for payments.

Currently, we support the following banks:
* Banco do Brasil
* Caixa
* Citibank
* Santander 
* BradescoShopFacil
* BradescoNetEmpresas
* Itau

The integration order will follow the list above but we may have changes considering our clients demands.

## Prerequisite
Below are all the prerequisites for you to be able to run the application:
* [GO](https://golang.org/dl/)
* Container [Docker](https://docs.docker.com/desktop/)

____________________________________________________________

## Running the application

Before cloning the Project, you should create the file path inside $GOPATH
```
	% mkdir -p "$GOPATH/src/github.com/mundipagg"
	% cd $GOPATH/src/github.com/mundipagg 
	% git clone https://github.com/mundipagg/boleto-api
```

In the project we have two ways to start the application, which is through Docker or manually using a command line.

### Using Docker

Run this command to download the dependencies and start the application by running the command
```
 cd .\boleto-api\devops\
 docker-compose up -d
```

To identify the application container just locate the name `devops_boleto-api` using the command below:
```
 docker ps -a
```

To view the application logs, simply execute the following command
```
 docker logs -f ID_CONTAINER
```

### Utilizando linha de comando

Open a command prompt and type the command below to generate the application build
``` 
 go build
``` 

After running the above command successfully, it will be possible to identify the application build with the name `boleto-api` that was generated inside the root folder (folder where the project was cloned).

If you are using Linux (* NIX), run the file generated in the build
```
 % ./boleto-api -airplane-mode
```

If you are using Windows, run the exe file generated in the build
```
 % .\boletoStone-publisher.exe -airplane-mode
```

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

```
    % go build && mv boleto-api /usr/local/bin
``` 
------------------

## Example of using the API

You can use [Postman](https://www.postman.com/downloads/) to request the API's services or even the curl
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

## For more information

See [FAQ](./FAQ.md)

## Contributing

To contribute, see [CONTRIBUTING](CONTRIBUTING.md)
