# api-order-processor

Este projeto demonstra a implementação de um sistema de processamento de pedidos assíncrono utilizando uma arquitetura baseada em Microservices, MongoDB e RabbitMQ.

## Serviços

- **api_server**: API HTTP para criação e consulta de pedidos, publica mensagens no RabbitMQ.
- **worker_service**: Consome mensagens do RabbitMQ e processa os pedidos, atualizando o MongoDB.
- **mongodb**: Banco de dados para persistência dos pedidos.
- **rabbitmq**: Fila de mensagens para comunicação assíncrona entre os serviços.

## Como rodar com Docker Compose

1. Certifique-se de ter o [Docker](https://docs.docker.com/get-docker/) e o [Docker Compose](https://docs.docker.com/compose/install/) instalados.

2. No diretório do projeto, execute:

	 ```bash
	 docker compose up --build
	 ```

	 Isso irá subir todos os serviços necessários.

3. A API estará disponível em: `http://localhost:8081`

## Exemplo de chamada curl

### Criar um novo pedido

```bash
curl --location 'http://localhost:8081/orders' \
--header 'Content-Type: application/json' \
--data '{
    "product": "Produto Exemplo",
    "quantity": 2
}'
```

### Consultar pedidos

```bash
curl http://localhost:8081/orders
```

## Observações

- Os serviços `api_server` e `worker_service` só iniciam após o MongoDB e o RabbitMQ estarem prontos.
- As credenciais e URIs de conexão estão configuradas no `docker-compose.yml`.
