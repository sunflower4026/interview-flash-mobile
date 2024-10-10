
# Golang Boiler Plate

## Getting started

To get started, simply run `docker-compose up --build`, and it will install and run. In the Docker Compose file, I have already included RabbitMQ, but I haven't used it in the project yet. Unfortunately, something came up that delayed my technical test.

If you don't want to run it using Docker Compose, you can simply build it with `go build -o main cmd/api/application.go` and then run it as usual. Don't forget to connect the backend app to your local PostgreSQL instance, and it will automatically run the migration to the latest schema.

## Folder Structure

```
1. cmd/
    - api/
        * application.go
    - job/
        * job_application.go
2. common/
    - constants/
    - hashing/
    - httpserver/
    - httpservice/
3. controller/
4. helper/
5. middleware/
6. model/
    - domain/
    - repository/
    - web/
7. routes/
8. service/
9. toolkit/

```
