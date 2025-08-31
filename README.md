# Real-Time Analytics Service "Beacon"

Простой, но мощный сервис для сбора и анализа веб-аналитики в реальном времени. Проект создан с целью изучения и демонстрации принципов построения распределенных систем на Go с использованием современного стека технологий.

## Архитектура

Система построена на микросервисной архитектуре, где каждый сервис выполняет свою специализированную задачу. Компоненты взаимодействуют друг с другом асинхронно через брокер сообщений NATS.

```mermaid
graph TD
    subgraph "Внешний мир"
        JS_Script[JS-скрипт на сайте клиента]
        Owner_Browser[Браузер владельца сайта]
    end

    subgraph "Система 'Beacon'"
        subgraph "Сервисы"
            Collector[Collector Service (Go)]
            Archiver[Archiver Service (Go)]
            Aggregator[Aggregator Service (Go)]
            API[API Service (Go)]
        end
        
        subgraph "Инфраструктура"
            NATS[NATS Message Bus]
            MongoDB[(MongoDB)]
            PostgreSQL[(PostgreSQL)]
            Redis[(Redis)]
        end

        Collector -- "1. Publish raw event" --> NATS
        NATS -- "2a. Consume raw event" --> Archiver
        NATS -- "2b. Consume raw event" --> Aggregator
        
        Archiver -- "3a. Write raw event" --> MongoDB
        Aggregator -- "3b. Write aggregated data" --> PostgreSQL
        Aggregator -- "3c. Update real-time cache" --> Redis
        
        API -- "4b. Read aggregated data" --> PostgreSQL
        API -- "4c. Read real-time data/cache" --> Redis
    end

    JS_Script -- "1. HTTP POST /track" --> Collector
    Owner_Browser -- "4a. HTTP GET /api/stats" --> API
    API -- "4d. JSON Response" --> Owner_Browser

    classDef services fill:#D6E8D5,stroke:#333,stroke-width:2px;
    class Collector,Archiver,Aggregator,API services;
```

## Технологический стек

*   **Язык:** Go
*   **Базы данных:**
    *   PostgreSQL (для структурированных данных: пользователи, сайты, агрегированная статистика)
    *   MongoDB (для хранения сырых событий)
    *   Redis (для кеширования и хранения real-time данных, например, счетчика онлайн-пользователей)
*   **Инфраструктура:**
    *   NATS (брокер сообщений для асинхронного взаимодействия сервисов)
    *   Docker & Docker Compose (для оркестрации и запуска окружения)
*   **Ключевые библиотеки Go:** Gin (для HTTP-сервисов), NATS.go, официальные драйверы для баз данных.

## Локальный запуск проекта

Для запуска проекта на локальной машине необходимо иметь установленный Docker.

1.  **Клонируйте репозиторий:**
    ```bash
    git clone [адрес_вашего_репозитория]
    cd beacon
    ```

2.  **Запустите инфраструктуру:**
    Эта команда поднимет все необходимые базы данных и брокер сообщений в Docker-контейнерах.
    ```bash
    docker-compose up -d
    ```

3.  **Запустите сервисы:**
    Откройте несколько терминалов и запустите каждый сервис в отдельном терминале.

    *   Терминал 1: Collector
        ```bash
        go run ./cmd/collector/main.go
        ```
    *   Терминал 2: Archiver
        ```bash
        go run ./cmd/archiver/main.go
        ```
    *   *... и так далее для других сервисов, когда они будут готовы.*

## Статус проекта

🚧 Проект находится в активной разработке.