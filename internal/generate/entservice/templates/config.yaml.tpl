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
    driver: sqlite3
    source: ':memory:'
  redis:
    addrs:
      - 127.0.0.1:6379