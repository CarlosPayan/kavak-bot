# Kavak AI Commercial Agent Bot

Este repositorio contiene un bot en Go que utiliza LLMs (OpenAI) para simular el comportamiento de un agente comercial de Kavak. El bot está disponible vía HTTP (`/qa`) y puede integrarse con WhatsApp usando Twilio.

## Tabla de Contenidos

1. [Descripción](#descripción)  
2. [Requisitos](#requisitos)  
3. [Estructura del Proyecto](#estructura-del-proyecto)  
4. [Configuración](#configuración)  
   - [config.yaml](#configyaml)  
5. [Ejecución Local](#ejecución-local)  
   - [Prueba con curl](#prueba-con-curl)  
   - [Integración con WhatsApp (Twilio)](#integración-con-whatsapp-twilio)  
6. [Contenerización con Docker](#contenerización-con-docker)  
7. [Métricas y Monitoreo](#métricas-y-monitoreo)  
7. [Roadmap](#roadmap)  
9. [Enlaces Útiles](#enlaces-útiles)  

---

## Descripción

Este bot utiliza:
- **Go** como lenguaje de backend.  
- **OpenAI** para Chat Completions y embeddings (RAG).  
- **CSV** con catálogo de autos para realizar búsquedas semánticas.  
- **Twilio WhatsApp Sandbox** para la integración de chat en WhatsApp.  
- **Prometheus** para métricas de latencia (opcional).

Las principales funcionalidades son:
- Proporcionar información general sobre Kavak (propuesta de valor, sucursales, garantías).  
- Recomendar autos del catálogo, manteniendo un “Último auto recomendado” en contexto.  
- Calcular planes de financiamiento con tasa anual 10%, plazos de 3 a 6 años, y enganches variables.  
- Mantener un flujo conversacional natural y amistoso.

---

## Requisitos

- **Go 1.24+** instalado en tu sistema.  
- **Docker** (opcional, para contenedorización).  
- **Cuenta y API Key de OpenAI**.  
- **Cuenta de Twilio** con acceso a WhatsApp Sandbox (Account SID, Auth Token, número de sandbox).  
- (Opcional) **Redis** o similar si deseas persistencia de sesiones en lugar de memoria.

---

## Estructura del Proyecto

```
.
├── cmd
│   └── bot
│       └── main.go            # Punto de entrada de la aplicación
├── configs
│   └── config.yaml            # Archivo de configuración principal
├── data
│   └── catalog.csv            # CSV con el catálogo de autos
├── Dockerfile                 # Dockerfile para construir la imagen
├── internal
│   ├── catalog
│   │   ├── catalog.go         # Lógica de carga CSV y embeddings
│   │   └── catalog_test.go    # Tests del catálogo
│   ├── config
│   │   └── config.go          # Carga de config.yaml con Viper
│   ├── handlers
│   │   ├── rag.go             # Handler de /qa (RAG, recomendaciones, financiamiento)
│   │   └── whatsapp.go        # Handler de /whatsapp (Twilio webhook)
│   ├── llm
│   │   └── client.go          # Cliente envoltorio para OpenAI
│   ├── metrics
│   │   └── metrics.go         # Métricas Prometheus (histogramas)
│   ├── store
│   │   └── store.go           # Almacenamiento en memoria de sesiones
│   └── utils
│       └── fetch_kavak.go     # Scraper para información de Kavak
└── README.md                  # Este archivo
```

---

## Configuración

### config.yaml

Crea `configs/config.yaml` con el siguiente contenido (reemplazando los valores con tus credenciales):

```yaml
openai:
  api_key: "sk-XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

catalog:
  path: "data/catalog.csv"

kavakInfoURL: "https://www.kavak.com/mx/blog/sedes-de-kavak-en-mexico"

server:
  address: ":8080"


- **`openai.api_key`**: Tu clave de API de OpenAI.  
- **`catalog.path`**: Ruta al CSV con catálogo de autos.  
- **`kavakInfoURL`**: URL de donde se extrae la info general de Kavak.  
- **`server.address`**: Puerto en el que el servidor escuchará (ej. `:8080`).

---

## Ejecución Local

1. **Clona el repositorio**  

   ```bash
   git clone https://github.com/carlospayan/agent-comercial-ai.git
   cd agent-comercial-ai
   ```

2. **Instala dependencias y compila**  

   ```bash
   go mod tidy
   go build -o kavak-bot cmd/bot/main.go
   ```

3. **Ejecuta el binario**  

   ```bash
   ./kavak-bot
   ```

   Verás en consola:  
   ```
   Listening on :8080…
   ```

4. **Prueba con curl**  

   - `/qa` (sin sesión previa):  
     ```bash
     curl -i "http://localhost:8080/qa?q=Hola"
     ```
   - `/qa` (manteniendo cookie de `session_id`):
     ```bash
     curl -i -X GET "http://localhost:8080/qa?q=¿Qué+SUV+tienen?"           -b "session_id=<UUID_de_la_sesión>"
     ```

---

### Integración con WhatsApp (Twilio)

1. **Levanta tu servidor local**  

   ```bash
   ./kavak-bot
   ```

2. **Expón localmente con ngrok** (u otro túnel):

   ```bash
   ngrok http 8080
   ```

   Copia la URL pública que te proporciona (ej. `https://abcd1234.ngrok.io`).

3. **Configura Twilio Sandbox**  

   - En tu consola de Twilio:  
     - Ve a **Messaging > Try it Out > WhatsApp Sandbox**.  
     - En “When a message comes in”, pega:
       ```
       https://abcd1234.ngrok.io/whatsapp
       ```
   - En tu WhatsApp envía:  
     ```
     join <sandbox_code>
     ```
     (Código que Twilio te asigna).

4. **Envía mensajes de prueba** desde tu WhatsApp:  
   - “Hola”  
   - “¿Qué SUV tienen?”  
   - “¿Cuánto cuesta ese auto?”  
   - “Si te doy 100000 de enganche, ¿cómo quedaría el financiamiento?”  
   - “Oye, ¿tienes Audis?”  

---

## Contenerización con Docker

### Dockerfile

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o kavak-bot cmd/bot/main.go

FROM alpine:latest
WORKDIR /root
COPY --from=builder /app/kavak-bot .
COPY --from=builder /app/configs/config.yaml configs/config.yaml
EXPOSE 8080
ENTRYPOINT ["./kavak-bot"]
```

### Construir la imagen

```bash
docker build -t kavak-bot:latest .
```

### Ejecutar el contenedor

```bash
docker run -d   -p 8080:8080 --name kavak-bot   kavak-bot:latest
```

---

## Métricas y Monitoreo

Se exponen métricas en `/metrics` usando **Prometheus**:

- `catalog_search_latency_ms` (histograma)  
- `openai_chat_latency_ms` (histograma)  
- `qa_request_latency_ms` (histograma)
- `whatsapp_request_latency_ms` (histograma)  


Para habilitar, visita:

```bash
curl http://localhost:8080/metrics
```

Configura tu Prometheus para que haga scrape a `http://<tu_servidor>:8080/metrics`.

---

## Roadmap
¿Cómo pondrías esto en producción?

    Contenerización con Docker, pipeline CI/CD para construir imágenes y testear automáticamente.

    Despliegue en un contenedor gestionado (AWS ECS), configurado con variables de entorno seguras.

    Persistencia de sesiones en Redis (en lugar de memoria).

    Logs y métricas exponen /metrics para Prometheus, con dashboards en Grafana y alertas.

    Seguridad: validación de firma Twilio, rate limiting, HTTPS obligatorio, secret management.

¿Cómo evaluarías el desempeño del agente?

    Métricas cuantitativas: latencia de RAG, latencia de OpenAI, latencia total, tasa de errores, consumo de tokens.

    Métricas cualitativas: encuestas en chat (👍/👎), revisiones manuales de conversaciones, tracking de “alucinaciones”.

    Dashboards que muestren p50/p90/p99 de latencia, errores por endpoint, ratio de satisfacción.

¿Cómo probarías que una nueva versión no tiene retroceso en su funcionalidad?

    Tests automáticos: suite de unitarios para lógica de catálogo; integración con httptest para /qa y /whatsapp usando mock LLM y catálogo de prueba.

    Tests end-to-end en staging: script de curl o Postman que valide flujos completos (Recomendación → Precio → Financiamiento).

    Deploy Canary/Blue-Green para enrutar tráfico gradualmente a la nueva versión y monitorear métricas antes de hacer swap total.

    Control de versiones de prompt: asegurarse de que, al cambiar SYSTEM_INSTRUCCIONS, no se introduzcan ambigüedades; añadir tests que verifiquen presencia de frases clave en el prompt.


## Enlaces Útiles

- [OpenAI Go SDK](https://github.com/sashabaranov/go-openai)  
- [Twilio WhatsApp Sandbox](https://www.twilio.com/console/sms/whatsapp/sandbox)  
- [Prometheus Go Client](https://github.com/prometheus/client_golang)  
- [Chi Router](https://github.com/go-chi/chi)  
- [Viper para configuración en Go](https://github.com/spf13/viper)  

