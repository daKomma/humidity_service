# Planing

## Requirements

- [ ] get Humidity
- [ ] get temperature
- [ ] live view of data
- [ ] diagram view of average data
  - [ ] daily
  - [ ] weekly
  - [ ] monthly
  - [ ] yearly
  
- [ ] get notification 
  - [ ] telegram
  - [ ] something to high
  - [ ] health check 

### Ideal Humidity
- winter 30-40
- summer 50-60

## Architecture

- 5 Microcontroller with API to provide sensor data
- 1 Dockercontainer with API to provide all data and handle logic
- 1 Database in a Dockercontainer to store the values
- Website to show data / get all data over Telegram...

## Routes

❗❗❗All dates in this UTC Format❗❗❗

### Microcontroller

#### healtcheck
```
{
    method: 'GET',
    url: '/health',
    response: {
        status: 'UP' || 'ERROR',
        message: STRING
    }
}
```

#### getData
```
{
    method: 'GET',
    url: '/data',
    response: {
        hum: NUMBER,
        temp: NUMBER
    }
}
```

### Server

#### healtcheck
```
{
    method: 'GET',
    url: '/health',
    response: {
        status: 'UP' || 'ERROR',
        message: STRING,
        timestamp: DATE
    }
}
```

#### register
```
{
    method: 'post',
    url: '/register',
    body: {
        url: STRING,
    },
    response: {
        health: STRING (URL),
        uuid: uuid of station
    }
}
```

#### getLiveData
for all or specific station
```
{
    method: 'get',
    url: '/data/live/:uuid',
    response: {
        timestamp: DATE,
        data: {
            id: string
            hum: NUMBER,
            temp: NUMBER
        }
    }
}
```

#### getData
```
{
    method: 'get',
    url: '/data',
    body: {
        start: DATE,
        end: DATE
    },
    response: [
        {
            data: {
                id: string,
                url: string,
                added: date,
                updated: date,
                hum: NUMBER,
                temp: NUMBER
            }
        }
    ]
}
```