# 🌡️ Desafio Clima por CEP com OpenTelemetry e Zipkin

Este projeto é composto por dois microserviços escritos em Go que trabalham juntos para retornar a temperatura atual de uma cidade com base em um CEP informado. A comunicação entre os serviços é rastreada com OpenTelemetry e Zipkin.

---

## 📦 Visão Geral dos Serviços

### 🔹 Serviço A – Input CEP
- Expõe `POST /cep`
- Recebe JSON `{ "cep": "29902555" }`
- Valida se o CEP tem 8 dígitos e é uma string
- Encaminha o CEP para o Serviço B via HTTP
- Propaga o trace OpenTelemetry

### 🔹 Serviço B – Orquestração e Clima
- Expõe `GET /weather/{cep}`
- Valida o CEP
- Usa ViaCEP para buscar a cidade
- Usa WeatherAPI para buscar a temperatura
- Converte para Celsius, Fahrenheit e Kelvin
- Retorna:
```json
{
  "city": "São Paulo",
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}
```

---

## 🚀 Como executar localmente

### 🔧 Pré-requisitos:
- Docker
- Docker Compose

### 🔁 Passos:

```bash
git clone https://github.com/WENDELLDELIMA/desafio-clima-go2.git
cd desafio-clima-go2
docker-compose up --build
```

Acesse os serviços:
- Serviço A: `http://localhost:8081/cep`
- Serviço B: `http://localhost:8080/weather/{cep}`
- Zipkin (Tracing): [http://localhost:9411](http://localhost:9411)

> ⚠️ Substitua `your_weather_api_key` no `docker-compose.yml` com sua chave da [WeatherAPI](https://www.weatherapi.com/)

---

## 📩 Exemplo de Requisição

### POST para Serviço A
```
POST /cep
Content-Type: application/json

{
  "cep": "01001000"
}
```

### Respostas:
- ✅ **200 OK** com cidade e temperaturas
- ❌ **422 Unprocessable Entity** – se CEP inválido
- ❌ **404 Not Found** – se CEP não encontrado

---

## 📊 OpenTelemetry + Zipkin

- Cada requisição cria um trace que atravessa os dois serviços.
- Spans são criados para:
  - Validação de CEP
  - Chamada à API ViaCEP
  - Chamada à WeatherAPI
- Use Zipkin para visualizar latências e dependências entre os serviços.

---

## 🗂️ Estrutura do Projeto

```
.
├── docker-compose.yml
├── service-a/
│   └── main.go
├── service-b/
│   └── main.go
```

---

## 👨‍💻 Desenvolvido por
[Wendell Lima](https://github.com/WENDELLDELIMA)

---

## 📜 Licença
MCCCCasdasç≈çç   
# desafio-clima-go2
