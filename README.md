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

❗❗❗All dates in this [Format](https://pradeepl.com/blog/dates-in-apis/#:~:text=The%20pattern%20for%20this%20date,in%20your%20RESTful%20web%20APIs.)❗❗❗

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
        message: STRING
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
        connection: BOOLEAN,
        health: STRING (URL)
    }
}
```

#### getLiveData
```
{
    method: 'post',
    url: '/data/live',
    response: {
        timestamp: DATE,
        data: {
            hum: NUMBER,
            temp: NUMBER
        }
    }
}
```

#### getData
```
{
    method: 'post',
    url: '/data',
    body: {
        start: DATE,
        end: DATE
    },
    response: [
        {
            timestamp: DATE,
            data: {
                hum: NUMBER,
                temp: NUMBER
            }
        }
    ]
}
```