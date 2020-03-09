# Routnies-go

REST based API server using GO

Build and run using docker

```
    make build
    docker-compose up --build

    docker-compose down --remove-orphans
```

Build and run on on local

```
    make build

    make run
```

Endpoints:
 
    - _create: Create a new routine
            ```
                curl --location --request GET 'http://localhost:8080/v1/_create?start=<start>&step=<step>'
            ```

    - _check: Check the current state of a routine
            ```
                curl --location --request GET 'http://localhost:8080/v1/_check?id=<id>'
            ```


    - _render: Returns an HTML page with info on all the routines and there current state 
        ```
            curl --location --request GET 'http://localhost:8080/v1/_render'
        ```

    - _clear: Clear the timer and cleanly exit the routine by setting stopped status
            ```
                curl --location --request PUT 'localhost:8080/v1/_pause?id=<>'
            ```
    - _pause: Pause the routine if exists, else error. Also sets modified time
            ```
                curl --location --request PUT 'http://localhost:8080/v1/_clear?id=<>'
            ```
    

    

