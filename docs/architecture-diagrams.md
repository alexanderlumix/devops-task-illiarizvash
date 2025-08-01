# Architecture Documentation

## System Overview

This project implements a microservices architecture with MongoDB replica set, featuring:
- **MongoDB Replica Set** (3 nodes) for high availability
- **HAProxy** for load balancing
- **Go Application** for reading products
- **Node.js Application** for creating products
- **Docker Compose** for orchestration

## Architecture Diagram

```mermaid
graph TB
    subgraph "Client Layer"
        Client[Client Applications]
    end
    
    subgraph "Load Balancer"
        HAProxy[HAProxy Load Balancer<br/>Port: 27034]
    end
    
    subgraph "Application Layer"
        GoApp[Go Application<br/>Port: 8080]
        NodeApp[Node.js Application<br/>Port: 3000]
    end
    
    subgraph "Database Layer"
        Mongo0[MongoDB Node 0<br/>Port: 27030]
        Mongo1[MongoDB Node 1<br/>Port: 27031]
        Mongo2[MongoDB Node 2<br/>Port: 27032]
    end
    
    subgraph "Infrastructure"
        Docker[Docker Compose<br/>Orchestration]
        CI[GitHub Actions<br/>CI/CD Pipeline]
    end
    
    Client --> HAProxy
    HAProxy --> Mongo0
    HAProxy --> Mongo1
    HAProxy --> Mongo2
    
    GoApp --> HAProxy
    NodeApp --> HAProxy
    
    Docker --> Mongo0
    Docker --> Mongo1
    Docker --> Mongo2
    Docker --> HAProxy
    Docker --> GoApp
    Docker --> NodeApp
    
    CI --> Docker
```

## Component Details

### 1. MongoDB Replica Set

**Purpose**: High availability database cluster
**Components**:
- 3 MongoDB nodes (mongo-0, mongo-1, mongo-2)
- Replica set name: `rs0`
- Authentication enabled with keyfile
- Admin user for management
- Application user for data access

**Configuration**:
```yaml
# Example MongoDB node configuration
mongod --replSet rs0 --bind_ip_all --keyFile /etc/mongo-keyfile
```

### 2. HAProxy Load Balancer

**Purpose**: Distribute database connections across replica set
**Configuration**:
- Port: 27034
- Round-robin load balancing
- Health checks for MongoDB nodes
- Connection pooling

### 3. Go Application

**Purpose**: Read products from database
**Features**:
- Structured logging with zap
- Health check endpoint (/health)
- Environment variable configuration
- Error handling and retry logic

**API Endpoints**:
- `GET /health` - Health check

### 4. Node.js Application

**Purpose**: Create products in database
**Features**:
- Structured logging with winston
- Input validation with express-validator
- Rate limiting with express-rate-limit
- Health check endpoint (/health)
- REST API for product creation

**API Endpoints**:
- `GET /health` - Health check
- `POST /products` - Create product (with validation)

## Data Flow

### Product Creation Flow

```mermaid
sequenceDiagram
    participant Client
    participant NodeApp
    participant HAProxy
    participant MongoDB
    
    Client->>NodeApp: POST /products
    NodeApp->>NodeApp: Validate input
    NodeApp->>HAProxy: Connect to MongoDB
    HAProxy->>MongoDB: Route to primary
    MongoDB->>MongoDB: Write to database
    MongoDB->>HAProxy: Confirm write
    HAProxy->>NodeApp: Return result
    NodeApp->>Client: Return product data
```

### Product Reading Flow

```mermaid
sequenceDiagram
    participant Client
    participant GoApp
    participant HAProxy
    participant MongoDB
    
    GoApp->>HAProxy: Connect to MongoDB
    HAProxy->>MongoDB: Route to any node
    MongoDB->>MongoDB: Read from database
    MongoDB->>HAProxy: Return data
    HAProxy->>GoApp: Return results
    GoApp->>GoApp: Format and display
```

## Security Architecture

### Authentication & Authorization

```mermaid
graph LR
    subgraph "Authentication"
        Admin[Admin User<br/>mongo-1/mongo-1]
        AppUser[App User<br/>appuser/appuserpassword]
    end
    
    subgraph "Authorization"
        AdminRole[Admin Role<br/>userAdminAnyDatabase]
        AppRole[App Role<br/>readWrite on appdb]
    end
    
    Admin --> AdminRole
    AppUser --> AppRole
```

### Network Security

- **Internal Network**: All services communicate via Docker network
- **Port Exposure**: Only necessary ports exposed to host
- **Authentication**: MongoDB keyfile authentication
- **Environment Variables**: Sensitive data in environment variables

## Deployment Architecture

### Development Environment

```mermaid
graph TB
    subgraph "Local Development"
        DockerCompose[Docker Compose]
        LocalEnv[.env file]
        Scripts[Python Scripts]
    end
    
    DockerCompose --> LocalEnv
    Scripts --> DockerCompose
```

### Production Environment

```mermaid
graph TB
    subgraph "Production"
        K8s[Kubernetes]
        Secrets[Secret Management]
        Monitoring[Monitoring Stack]
    end
    
    K8s --> Secrets
    K8s --> Monitoring
```

## Monitoring & Observability

### Health Checks

- **MongoDB**: Ping command to verify connectivity
- **HAProxy**: Configuration validation
- **Go App**: HTTP health endpoint
- **Node.js App**: HTTP health endpoint

### Logging

- **Structured Logging**: JSON format with timestamps
- **Log Levels**: DEBUG, INFO, WARN, ERROR
- **Context**: Service name, version, environment
- **Security**: Sensitive data masked in logs

### Metrics

- **Application Metrics**: Request count, response time
- **Database Metrics**: Connection count, query performance
- **Infrastructure Metrics**: CPU, memory, disk usage

## Error Handling

### Application Errors

```mermaid
graph TD
    Error[Error Occurs]
    Error --> LogError[Log Error with Context]
    LogError --> Retry{Retryable?}
    Retry -->|Yes| RetryLogic[Retry with Backoff]
    Retry -->|No| FailGracefully[Fail Gracefully]
    RetryLogic --> LogError
    FailGracefully --> NotifyUser[Notify User]
```

### Database Errors

- **Connection Errors**: Retry with exponential backoff
- **Authentication Errors**: Log and fail fast
- **Query Errors**: Log details and continue
- **Replica Set Errors**: Switch to healthy nodes

## Performance Considerations

### Database Performance

- **Connection Pooling**: Reuse connections
- **Indexes**: Proper indexing on frequently queried fields
- **Read Preferences**: Read from secondary nodes when possible
- **Write Concerns**: Appropriate write concern levels

### Application Performance

- **Caching**: In-memory caching for frequently accessed data
- **Rate Limiting**: Prevent abuse and ensure fair usage
- **Async Operations**: Non-blocking I/O operations
- **Resource Limits**: Docker resource constraints

## Scalability

### Horizontal Scaling

- **MongoDB**: Add more replica set members
- **Applications**: Scale application instances
- **Load Balancer**: Add more HAProxy instances

### Vertical Scaling

- **Resource Allocation**: Increase CPU/memory limits
- **Database Optimization**: Tune MongoDB parameters
- **Application Optimization**: Profile and optimize code

## Disaster Recovery

### Backup Strategy

- **MongoDB**: Regular backups with mongodump
- **Configuration**: Version control for all configs
- **Data**: Point-in-time recovery capabilities

### Recovery Procedures

1. **Database Recovery**: Restore from backup
2. **Application Recovery**: Redeploy from container images
3. **Configuration Recovery**: Restore from version control
4. **Data Validation**: Verify data integrity after recovery

## Security Checklist

- [ ] All credentials in environment variables
- [ ] MongoDB authentication enabled
- [ ] Network isolation with Docker networks
- [ ] Input validation implemented
- [ ] Rate limiting configured
- [ ] Structured logging with sensitive data masking
- [ ] Health checks implemented
- [ ] Error handling with proper logging
- [ ] CI/CD pipeline with security scanning
- [ ] Regular security updates 