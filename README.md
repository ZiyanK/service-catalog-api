# Service Catalog API

Service Catalog API is a REST API written in Golang that can be used as a storage of a collection of services along with it's respective versions.

## Tools/Technologies Used
1. [Gin](https://github.com/gin-gonic/gin) web framework
2. PostgreSQL (database)
3. [sqlx](https://github.com/jmoiron/sqlx) is used to interact with the database
4. [zap](https://github.com/uber-go/zap) for logging
5. [viper](https://github.com/spf13/viper) to get the env variables
6. [jwt](https://github.com/golang-jwt/jwt) for authentication
7. [validator](https://github.com/go-playground/validator/v10) to validate incoming requests
8. [goose](https://github.com/pressly/goose) to handler migrations
9. Docker and Kubernetes for easy deployment

## Database schema

### users
| column name | type                                |
|-------------|-------------------------------------|
| user_uuid   | UUID PRIMARY KEY                    |
| email       | VARCHAR(50) UNIQUE NOT NULL         |
| password    | VARCHAR(255) NOT NULL               |
| updated_at  | TIMESTAMP DEFAULT CURRENT_TIMESTAMP |
| created_at  | TIMESTAMP DEFAULT CURRENT_TIMESTAMP |

### services
| column name | type                                |
|-------------|-------------------------------------|
| service_id  | SERIAL PRIMARY KEY                  |
| name        | VARCHAR(255) NOT NULL               |
| description | TEXT                                |
| user_uuid   | UUID NOT NULL                       |
| updated_at  | TIMESTAMP DEFAULT CURRENT_TIMESTAMP |
| created_at  | TIMESTAMP DEFAULT CURRENT_TIMESTAMP |

### service_versions
| column name | type                                |
|-------------|-------------------------------------|
| sv_id       | SERIAL PRIMARY KEY                  |
| version     | VARCHAR(8)                          |
| changelog   | TEXT                                |
| service_id  | INTEGER NOT NULL                    |
| updated_at  | TIMESTAMP DEFAULT CURRENT_TIMESTAMP |
| created_at  | TIMESTAMP DEFAULT CURRENT_TIMESTAMP |

There is a foreign key for `user_uuid` in the `services` table and another foreign key for `service_id` in the `service_versions` table.

## To use
Prerequisites: Install [goose](https://github.com/pressly/goose)
* Create a `config.yaml` file and paste the content of the `config.sample.yaml` file (Change values as per usage)
* Create a postgres instance and run the migration file provided in the `database` directory. To make it easy, I have added command in the Makfile to do the same.
```bash
make db-up
make db-migrate
```
* Run the server
```bash
make run
```
* You can run the tests using the given command
```bash
make test
```
* If you want to build the application, you can use
```bash
make build
```

## Spec Open API
To view the OpenAPI documentation, you can run the following command in the terminal to start a swagger docker container
```bash
make spec-up
```
You can then view the documentation on [http://localhost:8080/](http://localhost:8080/)

To stop the swagger container, run
```bash
make spec-down
```

## Design considerations/Tradeoffs
* Use of preseverd JWT token that does not expire
* Writing all the tables and references in a single file. Done considering it is a simple CRUD API
* Pass of a context in each function to allow easy integration of tracing if required
* Separating of database operations as much as possible to increase code readability
* Addition of similar databse operations for transactional queries as done for the database operations in this [file](https://github.com/ZiyanK/service-catalog-api/app/db/sqlx.go) for easier code readability
* Haven't written tests for all cases of the handlers

## Assumptions
* Services can be updated but versions cannot be updated
* A version can only be created or deleted
* A single service cannot have multiple rows of the same version