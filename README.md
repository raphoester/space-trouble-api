# Space trouble API 

This is a simple backend that I have created to solve a technical test for a job application. 
Its goal is to allow users to book tickets for space trips.

The original assignment can be found [here](assignment.md).

This sample application can be used as a boilerplate for a more complex system.

## Tech stack 

- Presentation layer: gRPC with buf for code generation 
- CQRS pattern and tactical DDD
- Persistence layer: PostgreSQL + tests with ory/dockertest
- Runs on Docker Compose

## How to run

Make sure you have both Docker and Docker Compose installed.

Clone the repository and run the following commands:

```console
make dbuild
docker-compose up
```
