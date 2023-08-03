# Humidity Service

## Description
In this project I developed a REST service to manage several sensor stations that provide values for humidity and temperature. I used the GO programming language because I wanted to learn about this language and its features. There were many ups and downs in the development process, but overall it was a great experience and I learned a lot.

## Improvement Options
Of course, this code is not perfect. Here are some improvements that can be made in the next releases:
- Change Dockerfile to a Multistage build to increase security.
- Adding unit tests to the code.
- Adding a pause mode for stations to better handle stations that are currently unavailable.

## Usage
### Start
To start the service along with the database, a MySql database is run in this case:
```bash
docker-compose up
```

### Endpoints
A detailed documentation of the provided endpoints can then be found at [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html#/)

❗❗Note that all endpoints except `/health` have the `/api/v1` base path, which is not reflected by Swagger :/❗❗

## Development
To develop the application, the stations were mocked. To use the development environment, start the `docker-compose-dev.yml`.
```bash
docker-compose -f docker-compose-dev.yml up
```
