# 📱 SMS Gateway - Intelligent Queue Management System

## 🏗️ System Architecture

### High-Level Architecture

```mermaid
graph TB
    subgraph "Client Layer"
        API[REST API Endpoints]
        AUTH[Authentication & Authorization]
    end
    
    subgraph "Application Layer"
        SMS_SVC[SMS Service]
        USER_SVC[User Service]
        TRANS_SVC[Transaction Service]
    end
    
    subgraph "Queue Management Layer"
        QDS[Queue Distribution Strategy]
        MQC[Multi-Queue Consumer]
        
        subgraph "RabbitMQ Queues"
            MAIN_Q[sms-gateway<br/>Main Queue]
            PRIMARY_Q[sms-gateway-primary<br/>Primary Overflow<br/>90% weight]
            SECONDARY_Q[sms-gateway-secondary<br/>Secondary Overflow<br/>10% weight]
        end
    end
    
    subgraph "External Services"
        SMS_PROVIDER[SMS Provider]
        MOCK_PROVIDER[Mock Provider]
    end
    
    subgraph "Data Layer"
        MYSQL[(MySQL Database)]
        REDIS[(Redis Cache)]
    end
    
    API --> SMS_SVC
    API --> USER_SVC
    SMS_SVC --> QDS
    QDS --> MAIN_Q
    QDS --> PRIMARY_Q
    QDS --> SECONDARY_Q
    
    MQC --> MAIN_Q
    MQC --> PRIMARY_Q
    MQC --> SECONDARY_Q
    
    MQC --> SMS_PROVIDER
    MQC --> MOCK_PROVIDER
    
    SMS_SVC --> MYSQL
    USER_SVC --> MYSQL
    TRANS_SVC --> MYSQL
    
    style MAIN_Q fill:#e1f5fe
    style PRIMARY_Q fill:#f3e5f5
    style SECONDARY_Q fill:#fff3e0
    style QDS fill:#e8f5e8
    style MQC fill:#fff8e1
```

### SMS Processing Flow

```mermaid
sequenceDiagram
    participant Client
    participant API
    participant SMS_Service
    participant User_Service
    participant Queue_Strategy
    participant RabbitMQ
    participant Consumer
    participant Provider
    participant Database
    
    Client->>API: POST /send-sms
    API->>SMS_Service: SendSMS()
    SMS_Service->>User_Service: HasEnoughCredit()
    User_Service->>Database: Check balance
    Database-->>User_Service: Balance info
    User_Service-->>SMS_Service: Credit status
    
    alt Sufficient Credit
        SMS_Service->>Database: Create SMS record
        SMS_Service->>Database: Create transaction
        SMS_Service->>User_Service: DecreaseCredit()
        SMS_Service->>Queue_Strategy: PublishToQueue()
        
        Queue_Strategy->>RabbitMQ: Check main queue count
        RabbitMQ-->>Queue_Strategy: Queue count
        
        alt Main queue available
            Queue_Strategy->>RabbitMQ: Publish to main queue
        else Main queue full
            Queue_Strategy->>Queue_Strategy: Select overflow queue (weighted)
            Queue_Strategy->>RabbitMQ: Publish to overflow queue
        end
        
        Queue_Strategy-->>SMS_Service: Published
        SMS_Service-->>API: Success
        API-->>Client: SMS queued successfully
        
        Consumer->>RabbitMQ: Consume message
        Consumer->>Database: Get SMS details
        Consumer->>Provider: Send SMS
        Provider-->>Consumer: Delivery status
        Consumer->>Database: Update SMS & transaction status
        Consumer->>RabbitMQ: ACK message
    else Insufficient Credit
        SMS_Service-->>API: Error: Insufficient credit
        API-->>Client: Error response
    end
```

## 🧠 Intelligent Queue Distribution Strategy

### Queue Selection Algorithm

```mermaid
flowchart TD
    START([New SMS Message])
    CHECK_MAIN{Main Queue Count<br/>< Prefetch Threshold?}
    USE_MAIN[Route to Main Queue<br/>sms-gateway]
    OVERFLOW_LOGIC{Overflow Distribution}
    RANDOM[Generate Random<br/>Number 1-100]
    PRIMARY_CHECK{Random ≤<br/>Primary Weight<br/>90%?}
    USE_PRIMARY[Route to Primary<br/>sms-gateway-primary]
    USE_SECONDARY[Route to Secondary<br/>sms-gateway-secondary]
    PUBLISH[Publish Message]
    
    START --> CHECK_MAIN
    CHECK_MAIN -->|Yes| USE_MAIN
    CHECK_MAIN -->|No| OVERFLOW_LOGIC
    OVERFLOW_LOGIC --> RANDOM
    RANDOM --> PRIMARY_CHECK
    PRIMARY_CHECK -->|Yes| USE_PRIMARY
    PRIMARY_CHECK -->|No| USE_SECONDARY
    USE_MAIN --> PUBLISH
    USE_PRIMARY --> PUBLISH
    USE_SECONDARY --> PUBLISH
    
    style START fill:#e1f5fe
    style USE_MAIN fill:#c8e6c9
    style USE_PRIMARY fill:#f8bbd9
    style USE_SECONDARY fill:#ffcc02
    style OVERFLOW_LOGIC fill:#fff3e0
```

### Key Features

- **🎯 Adaptive Load Balancing**: Automatically switches to overflow queues when main queue reaches capacity
- **⚖️ Weighted Distribution**: 90% to primary overflow, 10% to secondary overflow
- **📊 Real-time Monitoring**: Continuously monitors queue depths
- **🔄 Graceful Degradation**: Maintains service availability under high load

## 🔄 Multi-Queue Consumer Architecture

Our multi-queue consumer pattern ensures efficient processing across all queues:

### Concurrent Consumer Pattern

```mermaid
graph TB
    subgraph "Multi-Queue Consumer"
        INIT[Initialize Queues]
        WG[Wait Group]
        
        subgraph "Concurrent Goroutines"
            C1[Consumer 1<br/>Main Queue]
            C2[Consumer 2<br/>Primary Queue]
            C3[Consumer 3<br/>Secondary Queue]
        end
        
        subgraph "Message Processing"
            PARSE[Parse SMS ID]
            FETCH[Fetch SMS Details]
            SEND[Send via Provider]
            UPDATE[Update Status]
            ACK[ACK/NACK Message]
        end
    end
    
    subgraph "RabbitMQ"
        Q1[sms-gateway]
        Q2[sms-gateway-primary]
        Q3[sms-gateway-secondary]
    end
    
    INIT --> WG
    WG --> C1
    WG --> C2
    WG --> C3
    
    C1 --> Q1
    C2 --> Q2
    C3 --> Q3
    
    Q1 --> PARSE
    Q2 --> PARSE
    Q3 --> PARSE
    
    PARSE --> FETCH
    FETCH --> SEND
    SEND --> UPDATE
    UPDATE --> ACK
    
    style C1 fill:#e1f5fe
    style C2 fill:#f3e5f5
    style C3 fill:#fff3e0
    style PARSE fill:#e8f5e8
```


### Queue Configuration Details

| Queue Name | Purpose | Weight | Use Case |
|------------|---------|---------|----------|
| `sms-gateway` | Main queue | N/A | Primary message processing under normal load |
| `sms-gateway-primary` | Primary overflow | 90% | High-priority overflow traffic |
| `sms-gateway-secondary` | Secondary overflow | 10% | Additional capacity for peak loads |

