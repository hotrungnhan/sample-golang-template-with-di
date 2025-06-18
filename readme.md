# Guide
Require:
* Devbox

## Setup And Run project

For testing purpose

Step 0: Please copy `.env.sample` into `.env`
Step 1: Go to devbox shell `devbox shell`
Step 2: Start infra `task up`
Step 3: Setup New DB `task db:setup` (it not have persistent volume)
Step 4: Generate code `task gen`
Step 5: Start service`task run -- http`
Step 6: Remove infra `task down`

## Check list
* [x] POST /api/shortlinks
* [x] GET /api/shortlinks/{id}
* [x] GET /shortlinks/{id}
* [x] ShortCode 8 alphanumeric characters
* [x] Handle duplicate original URLs by returning the existing short code
* [x] Validate Parameter (original URL)
* [x] Core business Unit test.
* [x] 3 tier architecture && CQRS pattern
* [ ] auto generate swagger docs

## Folder Structure Overview

| Folder Path        | Description                                                                                 |
| ------------------ | ------------------------------------------------------------------------------------------- |
| /cmds              | Houses main application entry points or command-line tools.                                 |
| /controllers       | Contains logic for handling incoming requests and returning responses in web applications.  |
| /generated         | Holds auto-generated code or files (should not be manually edited).                         |
| /migrations        | Contains database migration scripts for managing schema changes.                            |
| /models            | Defines data models or structures, often representing database tables or business entities. |
| /repositories      | Manages data access logic, such as database queries and persistence.                        |
| /serializers       | Handles data serialization and deserialization (e.g., JSON ↔ Go structs).                   |
| /services          | Contains business logic and service layer code.                                             |
| /utils/helpers     | Contains additional helper functions for specific or isolated tasks.                        |
| /utils/injects     | Likely used for dependency injection configurations or related code.                        |
| /utils/middlewares | Contains middleware logic, such as request/response processing layers.                      |
| /utils/types       | Defines custom types, interfaces, or type aliases used across the project.                  |
