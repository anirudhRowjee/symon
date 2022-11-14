# SYMon - Backend

This is the code for the Backend of `symon`, a system monitoring tool.

## Configurable Parameters
- Tick Rate

    It is possible to cofigure how often the backend polls the kernel for metrics. 

    On lines `17-18` of `main.go`, tweak the number at the end of the line to change the tick rate in Millisecond.

    ```go
    // Define the standard tick rate of 500 ms
    const TICK_RATE = time.Millisecond * 500
    ```


## Build And Run Instructions
Ensure you have Golang installed and your `$GOPATH` configured properly.

> :warn: you do not need to install dependencies this if you've already got it installed and changed code that doesn't introduce a dependency - for example, the Tick Rate

### Installing Dependencies


Navigate to the backend directory and run the following command to install the dependencies.

> This will be run only once, or whenever there are changes made to the `go.mod` file.  
```bash
$ go mod install
```

### Running the Backend

Once the dependencies have been installed **ensure there's nothing else running on port `1337`** .

Then type the following - 
```bash
$ go run main.go
```

## Polling for Metrics
The latest metrics per the tick rate will be available at
```
http://localhost:1337/metrics
```

Metrics are available as JSON, and will have a defined schema.

## Example API Response
Let's run a request from the CLI via `curl` in a Linux environment 

```
$ curl http://localhost:1337/metrics | jq
```

Here's a sample (but schematically accurate) response - 

```json
{
  "metrics_timestamp": "2006.01.02 15:04:05",
  "cpu_metrics": {
    "cpu_count": 4,
    "cpu_cores": 8,
    "cpu_usage_percentage": [
      5.050505050513054,
      6.060606060601885,
      5.05050505049693,
      11.224489795914343,
      6.930693069310386,
      8.000000000008185,
      4.9999999999954525,
      13.999999999998636
    ]
  },
  "memory_metrics": {
    "mem_usage_percentage": 15.873484733225402,
    "mem_total_available": 16486055936,
    "mem_total_used": 3286306816
  }
}
```
