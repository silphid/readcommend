# Readcommend

Readcommend is a highly-scalable, feature-rich book recommendation web app written in Go and React for the true book
aficionados and disavowed human-sized bookworms.

# Project structure

- [/app](app): the web application, written in Node.js/TypeScript/React.js
- [/dev](dev/README.md): development-time scripts and docker-compose.
- [/doc/open-api.yaml](doc/open-api.yaml): OpenAPI doc/specs.
- [/dist](dist): files intended to be included as-is in final docker image.
- [/src](src): the back-end service, written in Go, with Echo server and PostgreSQL DB.

# Development environment

This project has been developed and tested with the following tools and versions:

- macOS Big Sur 11.2.1
- Docker Desktop for Mac 3.0.4
  - Engine: 20.10.2
  - Compose: 1.27.4
- Node v10.18.1
- Go 1.16

# Running server and database

In project's root dir, run:

```bash
$ docker-compose up
```

Then point your browser to: http://localhost:5000

# Initial requirements

- Write a REST endpoints that interacts with a collection of books from a data source containing the following properties:
  - Author
  - Genre
  - Number of pages
  - Year of publication
  - Overall rating (from 1 to 5 stars)
- The program will provide an interface for the user to select his preferences, in ranking order:
  - by specific author
  - by specific genre
  - by number of pages
  - by classic or modern literature
- The user may enter multiple criteria if he so wishes.
- Basic Requirements:
  - Written in GO.
  - Uses REST to pass data between the server and the front-end.
  - Is a web application.
  - Deployed to GitHub.
  - Contains at least 1 unit test.
  - Contains some documentation.

# Wishlist

The following features and aspects have not been implemented yet, but are considered essential if the application was to be deployed in a real production environment:

- UI improvements:
  - Results paging.
  - Dynamically querying list of authors as user types (would be required for very large datasets).
- Integration/e2e tests, covering all REST endpoints.
- Health monitoring endpoint (it is listed in OpenAPI spec, but not yet implemented).
- Serving OpenAPI UI/doc via an endpoint.
- Infrastructure as Code to deploy PostgreSQL database.
- Kubernetes chart/helmfile to deploy the service.

# Notes

- Current test coverage for the data access layer is above 80% (actually 92% for the most crucial book querying package). That coverage basically amounts to 100% of everything that is not error handling (which cases are much more involved to cover, because they depend on database state and connectivity).
- Currently, test feed dataset is loaded into database as one of the migrations, which would not be the case in a real project (it would rather be loaded independently via a SQL script during development or integration/e2e tests).
- Some authors in the test dataset do not have associated books.
- If you're looking for an author in the test dataset with many associated books, try "Lynne Danticat" who has 8.
