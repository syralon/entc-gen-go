# trace:
#   endpoint:
server:
  http:
    network:
    addr: 0.0.0.0:8000
    timeout: 5s
  grpc:
    network:
    addr: 0.0.0.0:9000
    timeout: 5s
data:
  database:
    driver: mysql
    source: root:root@tcp(127.0.0.1:3306)/example
  redis:
    addrs:
      - 127.0.0.1:6379