# System Architecture

## Overview

This project implements a microservices-based architecture with MongoDB replica set for high availability and data consistency.

## Architecture Components

### 1. MongoDB Replica Set
- **Purpose**: Provides high availability and data redundancy
- **Components**: 3 MongoDB instances (Primary, Secondary, Secondary)
- **Ports**: 27030, 27031, 27032
- **Load Balancer**: HAProxy on port 27034

### 2. Applications
- **Node.js App**: Product creation service
- **Go App**: Product reading service
- **Communication**: Both apps connect to MongoDB via HAProxy

### 3. Infrastructure
- **Containerization**: Docker for all services
- **Orchestration**: Docker Compose for local development
- **Load Balancing**: HAProxy for MongoDB connections

## Data Flow

```
Node.js App → HAProxy → MongoDB Replica Set
Go App     → HAProxy → MongoDB Replica Set
```

## Security Model

- **Authentication**: MongoDB users with role-based access
- **Network**: Isolated Docker networks
- **Secrets**: Environment variables and keyfiles

## Scalability Considerations

- **Horizontal Scaling**: Additional MongoDB nodes can be added
- **Vertical Scaling**: Resource limits configurable in Docker Compose
- **Load Distribution**: HAProxy provides connection pooling

## Monitoring Points

- **Database**: Replica set status, connection health
- **Applications**: Response times, error rates
- **Infrastructure**: Container health, resource usage 