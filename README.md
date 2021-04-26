[![Build Status](https://mundipagg.visualstudio.com/Processing%20and%20Reconciliation/_apis/build/status/Banking/Boleto/boleto-api%20-%20PRODUCTION?branchName=master)](https://mundipagg.visualstudio.com/Processing%20and%20Reconciliation/_build/latest?definitionId=319&branchName=master)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=mundipagg_boleto-api&metric=alert_status)](https://sonarcloud.io/dashboard?id=mundipagg_boleto-api)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=mundipagg_boleto-api&metric=coverage)](https://sonarcloud.io/dashboard?id=mundipagg_boleto-api)

# boleto-api [PT-BR](./README.md)/[ENG-US](./README-US-ENG.md)
API responsável pelo registro online de um boleto. Atualmente, oferecemos suporte aos seguintes bancos emissores:
* Banco do Brasil
* Caixa
* Citibank
* Santander 
* BradescoShopFacil
* BradescoNetEmpresas
* Itaú
* Pefisa

## Pré-requisitos
 * Estar logado na VPN através do programa [CiscoVPN](https://www.cisco.com/c/pt_br/products/security/anyconnect-secure-mobility-client/index.html)
 * [Golang](https://golang.org/dl/)
 * Container [Docker](https://docs.docker.com/desktop/)

____________________________________________________________
## Rodando a aplicação
Antes de clonar o projeto, você deve criar o caminho do arquivo dentro $GOPATH
```
    % mkdir -p "$GOPATH/src/github.com/mundipagg"
    % cd $GOPATH/src/github.com/mundipagg 
    % git clone https://github.com/mundipagg/boleto-api
```

No projeto temos duas formas de inicializar: via Docker ou via aplicação.

### Executando via Docker
Baixe as dependências e inicie a aplicação executando:
```
 cd .\boleto-api\devops\
 docker-compose up -d
```

Para identificar o container basta localizar o nome _boletostone-consumer através do comando:
```
 docker ps -a
```

Para visualizar os logs da aplicação basta executar o comando a seguir:
```
 docker logs -f ID_CONTAINER
```

### Executando via aplicação

Gere a build da aplicação através do comando:
``` 
 go build
``` 

Finalizado o comando acima será possível identificar o build da aplicação com o nome `boleto-api` que foi gerada dentro da pasta root do projeto.

#### Parâmetros para execução
Existem vários comandos que auxiliam na execução da app.
* *dev*: este comando carregará todas as variáveis de ambiente no modo padrão, conectando com BD, integrado com bancos emissores e etc;
* *mock*: caso seja necessário utilizar app sem integração com o banco emissor e com BD em memória;
* *nolog*: caso deseje execução sem log;
* *airplane-mode*: é a junção do _mock_ e _nolog_! Ou seja: não há dependências externas;


Para SO Linux (*NIX) execute o arquivo gerado no build:
```
 % ./boleto-api -airplane-mode
```

Para SO Windows execute o arquivo gerado no build:
```
 % .\boletoStone-publisher.exe -dev
```

Para instalação do executável são necessárias apenas as variáveis ​​de ambiente configuradas e a aplicação compilada.

Edite o arquivo $HOME/.bashrc.sh
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
____________________________________________________________    

### Utilizando

#### Banco do Brasil
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

#### Caixa
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

#### Response
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
No caso da Caixa, a impressão do boleto ficará a cargo da Caixa. Assim, a API retornará a URL do boleto da Caixa.


#### Erro
Caso ocorra erro esta é a resposta padrão:
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

## Para mais informações
Veja o [FAQ](./FAQ.md)

## Contribuir
Veja as regras para [contribuir](CONTRIBUTING.md)