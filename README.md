# Ventrata Octo Test Project

You can find the project description here.
https://paper.dropbox.com/doc/Golang-Task-SRHVn4MGthOvIUZS8R1By

## Getting Started

### Prerequisites
- Go 1.19 or higher
- Docker and Docker Compose

### Installing

#### 1. Clone the repository
```
git clone https://github.com/SDSTA/vent-octo.git
cd vent-octo
```

#### 2. Build the Docker containers
```
make docker-build
```

#### 3. Start the services
```
make docker-up
```

#### 4. Apply database migrations
```
make migrate-up
```

### Runing Tests
```
make test
```

## Built With
- Go (Golang)
- Docker
- golang-migrate/migrate
- CurrencyAPIss
- Data-dog
- sql-mock
- Swagger

## Developer

### Andrei Florea - andreiflorea19981014@gmail.com