# ![RealWorld Example App](logo.png)

> ### Golang codebase containing real world examples (CRUD, auth, advanced patterns, etc) that adheres to the [RealWorld](https://github.com/gothinkster/realworld) spec and API.


### [Demo](https://demo.realworld.io/)&nbsp;&nbsp;&nbsp;&nbsp;[RealWorld](https://github.com/gothinkster/realworld)


This codebase was created to demonstrate a fully fledged fullstack application built with Golang including CRUD operations, authentication, routing, pagination, and more.

We've gone to great lengths to adhere to the Golang community styleguides & best practices.

For more information on how to this works with other frontends/backends, head over to the [RealWorld](https://github.com/gothinkster/realworld) repo.


# How it works

My system leverages the principles of hexagonal architecture to achieve a robust and flexible design. At its core, hexagonal architecture emphasizes a clear separation between the core business logic and external dependencies. This clear separation allows developers to focus on specific concerns without becoming tightly coupled to other parts of the system.

# Getting started

## Require
- [Docker](https://www.docker.com/)
- [Go](https://go.dev/)

##Development
1. Create PostgreSQL database.
2. Migration, using the following command (edit path to migrations to full path of real-world-api/pkg/databases/migrations).
>$ migrate -source file://path/to/migrations -database 'postgres://admin:123456@localhost:5432/realworld-db?sslmode=disable' -verbose up
3. Run the app (Development is using configuration in .env.dev).
>$ air -c .air.dev.toml
##Build
1. Run 
>$ docker compose up -d
2. Migration, using the following command (edit path to migrations to full path of real-world-api/pkg/databases/migrations).
>$ migrate -source file://path/to/migrations -database 'postgres://admin:123456@localhost:5432/realworld-db?sslmode=disable' -verbose up
3. Ensure realworld-app is running, you can send request to the app as config in .evn.prod