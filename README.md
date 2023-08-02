# Humidity Service

## Description
In this project i developed a REST service to manage multiple sensor stations which provide values for humidity and temperature. 
I used the GO programming language because i wanted to learn this language with its features. In de developing process where lots of ups and downs but overall it was a great experience and i learned a lot of things.

## Improvement Options
Of course is this code not perfect. Here are some improvements that can be done in the next versions.
- Change Dockerfile to a Multistage build to increase security.
- Add unit test to the code.
- Add a pause mode to the stations so stations which are currently not available can be handled better.

## Usage
### Start
To start the service along side the database in this case a MySql database run:
```bash
docker-compose up
```

### Endpoints
A detailed documentation of the provided endpoints can then be found at [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html#/)

❗❗Note that all entpoints beside `/health` have the base path `/api/v1` this isn´t rendered by swagger :/  ❗❗

## Development
To develop the application the stations where mocked. To use the development enviroment start the `docker-compose-dev.yml`.
```bash
docker-compose -f docker-compose-dev.yml up
```
