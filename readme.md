**This project is currently under active development. It is experimental and not ready for production.**

# ⛊ kalkan ⛊
A simple, experimental, configurable, standalone reverse proxy service with built-in basic rate limiting to control and protect the backend API by filtering and throttling incoming requests.

# Configuration Guide

This guide provides an overview of the configuration options available in the `config.yaml` (located in the `internal/cfg` folder) file used by the reverse proxy service with rate limiting. Below, you'll find descriptions for each section and the parameters you can adjust.

## Configuration Structure

The `config.yaml` file is divided into three main sections:

1. **Server Configuration**
2. **Origin Configuration**
3. **Rate Limiter Configuration**

### 1. Server Configuration

This section defines the configuration settings for the server that will run the reverse proxy.

```yaml
server:
  port: 8080
```
- `port`: Specifies the port number on which the server will listen for incoming requests. Adjust this value to change the listening port.


### 2. Origin Configuration
The origins section is where you define the backend services that the proxy will route requests to. Each origin is defined by its name, an edge, and a URL.
```yaml
origins:
  - name: Server1
    edge: /server1
    url: "https://sample.com"
  ...
```
- `name`: A friendly name for the backend server. This is used for identification and logging purposes.
- `edge`: The endpoint path that the proxy listens to. When a request matches this path, it will be forwarded to the corresponding url.
- `url`: The actual URL of the backend server where requests will be forwarded.

### 3. Rate Limiter Configuration
This section allows you to configure the rate-limiting functionality, which helps protect your backend servers by controlling the rate of incoming requests.
```yaml
rl:
  enabled: true
  allowed: 5
  timeframe: 10
  block: 20
```
- `enabled`: A boolean flag (true or false) that determines whether the rate limiter is active. If set to true, the rate limiter will be enforced.
- `allowed`: The number of requests allowed per client within the specified timeframe.
- `timeframe`: The duration (in seconds) within which the number of requests specified in allowed can be made. After this timeframe, the request count is reset.
- `block`: The duration (in seconds) for which the client is blocked from making further requests if they exceed the allowed limit.