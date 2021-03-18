[![mundipagg maturity](http://maturityapp.herokuapp.com/maturity.php?project=mundipagg/boleto-api&command=badge_image)](http://maturityapp.herokuapp.com/index.html?project=mundipagg/boleto-api)

[![GoDoc](https://godoc.org/github.com/mundipagg/boleto-api?status.svg)](https://godoc.org/github.com/mundipagg/boleto-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/mundipagg/boleto-api)](https://goreportcard.com/report/github.com/mundipagg/boleto-api)

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/b85953cc9fa84b56822e7e5d91203e91)](https://www.codacy.com/app/mundipagg/boleto-api?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=mundipagg/boleto-api&amp;utm_campaign=Badge_Grade)
[![Maintainability](https://api.codeclimate.com/v1/badges/b9ad683e9d8f87034339/maintainability)](https://codeclimate.com/github/mundipagg/boleto-api/maintainability)

# Boleto API

É uma API para cadastro online de boleto em bancos e criação de boleto para pagamentos.

Atualmente, oferecemos suporte aos seguintes bancos:
* Banco do Brasil
* Caixa
* Citibank
* Santander 
* BradescoShopFacil
* BradescoNetEmpresas
* Itau

A ordem de integração seguirá a lista acima, mas podemos ter alterações considerando as demandas de nossos clientes.

## Pré-requisitos
Abaixo estão todos os pré requisitos para que consiga rodar a aplicação:
* [GO](https://golang.org/dl/)
* Container [Docker](https://docs.docker.com/desktop/)

____________________________________________________________
## Rodando a aplicação
A API foi desenvolvida em linguagem GO e por isso é necessário instalar as ferramentas de linguagem caso seja necessário compilar a aplicação a partir do código-fonte.

Antes de clonar o projeto, você deve criar o caminho do arquivo dentro $GOPATH
```
	% mkdir -p "$GOPATH/src/github.com/mundipagg"
	% cd $GOPATH/src/github.com/mundipagg 
	% git clone https://github.com/mundipagg/boleto-api
```

No projeto temos duas formas de inicializar a aplicação que é por meio do Docker ou manualmente por uma linha de comando.

### Utilizando Docker

Execute esse comando para baixar as dependências e iniciar a aplicação, executando o comando:
```
 cd .\boleto-api\devops\
 docker-compose up -d
```
Para identificar o container da aplicação basta localizar o nome `devops_boleto-api` através do comando abaixo:
```
 docker ps -a
```

Para visualizar os logs da aplicação basta executar o comando a seguir:
```
 docker logs -f ID_CONTAINER
```
With that you will have access to the logs of the application that is running through the container.

### Utilizando linha de comando

Abra um prompt de comando e digite o comando abaixo para gerar o build da aplicação:
``` 
 go build
``` 
Após rodar o comando acima com sucesso, será possível identificar o build da aplicação com o nome `boleto-api` que foi gerada dentro da pasta root (pasta onde foi clonado o projeto).

Se você estiver utilizando Linux (*NIX), execute o arquivo gerado no build
```
 % ./boleto-api -airplane-mode
```

Se você estiver utilizando Windows, execute o arquivo exe gerado no build
```
 % .\boletoStone-publisher.exe -airplane-mode
```

Se você deseja executar a API no modo dev, que carregará todas as variáveis ​​de ambiente padrão, você deve executar o aplicativo assim:
```
 % ./boleto-api -dev
```

Caso você queira executar o aplicativo no modo simulado usando um banco de dados de memória em vez de integração com o banco, você deve usar a opção simulada:
```
 % ./boleto-api -mock
```
Caso queira rodar o aplicativo com o log desligado, deve-se usar a opção -nolog:
```
 % ./boleto-api -nolog
```
Você pode combinar todas essas opções e, caso queira usá-las todas juntas, você pode simplesmente usar a opção -airplane-mode
```
 % ./boleto-api -airplane-mode
```

Por padrão, a api do boleto irá instalar e executar um servidor https, mas você pode executar no modo http com a seguinte opção
```
 % ./boleto-api -http-only
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

### Exemplo de utilização da API

Você pode utilizar o [Postman](https://www.postman.com/downloads/) para solicitar os serviços da API ou mesmo o curl
Veja os exemplos a seguir

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
No caso da Caixa, a impressão do boleto ficará a cargo da Caixa. Assim, a API retornará a URL do boleto da Caixa.


A resposta da API terá o seguinte padrão se ocorrer algum erro:
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

Layout da aplicação
---

A raiz do aplicativo contém apenas o arquivo main.go e alguns arquivos de configuração e documentação.

Na raiz, temos os seguintes pacotes:

* `api`: Rest Controllers
* `auth`: Autorização com os bancos
* `bank`: Interface que registra os boletos
* `bb`: Implementação para o Banco do Brasil
* `caixa`: Implementação para o Caixa
* `citibank`: Implementação para o Citibank
* `boleto`: Criação de um boleto
* `cache`: Base de dados (Chave/Valor) em memória, usado apenas quando a aplicação está rodando em modo mock
* `config`: Configuração da aplicação
* `db`: Base de dados de persistência
* `devops`: Contém os arquivos de upload/deploy da aplicação
* `validations`: Validação básia dos dados
* `log`: Logs da aplicação
* `models`: Modelo de dados da aplicação
* `parser`: XML parser
* `test`: Utilitário de testes
* `tmpl`: Utilitário de Template
* `util`: Utilitário de Miscellaneous
* `integrationTests`: Contém todos os testes de caixa preta
* `vendor`: Bibliotecas de tericeiros

## Para mais informações
Veja o [FAQ](./FAQ.md)

## Contribuir
Veja as regras para [contribuir](CONTRIBUTING.md)