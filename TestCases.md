rutuja@rutuja:grpc-client-serve$ touch go.sum
docker-compose up -d
kv-service is up-to-date
api-service is up-to-date
rutuja@rutuja:~/grpc-client-serve$ docker-compose ps
   Name          Command         State                          Ports                    
-----------------------------------------------------------------------------------------
api-service   ./api-service   Up (healthy)   0.0.0.0:8081->8080/tcp,:::8081->8080/tcp    
kv-service    ./kv-service    Up (healthy)   0.0.0.0:50051->50051/tcp,:::50051->50051/tcp
rutuja@rutuja:~/grpc-client-serve$ curl http://localhost:8080/api/v1/health

Caution: You are using the Snap version of curl.
Due to Snap's sandbox nature, this version has some limitations.
For example, it may not be able to access hidden folders in your home directory
or other restricted areas of the os.
This means you may encounter errors when using snap curl to download files.
For those case, you might want to use the native curl package.
For details, see: https://github.com/boukendesho/curl-snap/issues/1

To stop seeing this message, run the following command:
curl.snap-acked

curl: (7) Failed to connect to localhost port 8080 after 0 ms: Could not connect to server
rutuja@rutuja:~/grpc-client-serve$ docker-compose ps
   Name          Command         State                          Ports                    
-----------------------------------------------------------------------------------------
api-service   ./api-service   Up (healthy)   0.0.0.0:8081->8080/tcp,:::8081->8080/tcp    
kv-service    ./kv-service    Up (healthy)   0.0.0.0:50051->50051/tcp,:::50051->50051/tcp
rutuja@rutuja:~/grpc-client-serve$ curl http://localhost:8081/api/v1/health

Caution: You are using the Snap version of curl.
Due to Snap's sandbox nature, this version has some limitations.
For example, it may not be able to access hidden folders in your home directory
or other restricted areas of the os.
This means you may encounter errors when using snap curl to download files.
For those case, you might want to use the native curl package.
For details, see: https://github.com/boukendesho/curl-snap/issues/1

To stop seeing this message, run the following command:
curl.snap-acked

{"service":"api-service","status":"healthy","timestamp":1760234310}
rutuja@rutuja:~/grpc-client-serve$ curl -X POST http://localhost:8081/api/v1/kv \
  -H "Content-Type: application/json" \
  -d '{"key": "test", "value": "hello"}'

Caution: You are using the Snap version of curl.
Due to Snap's sandbox nature, this version has some limitations.
For example, it may not be able to access hidden folders in your home directory
or other restricted areas of the os.
This means you may encounter errors when using snap curl to download files.
For those case, you might want to use the native curl package.
For details, see: https://github.com/boukendesho/curl-snap/issues/1

To stop seeing this message, run the following command:
curl.snap-acked

{"success":true,"message":"Successfully stored key 'test'"}
rutuja@rutuja:~/grpc-client-serve$ ^[[200~curl http://localhost:8081/api/v1/kv/test~^C
rutuja@rutuja:~/grpc-client-serve$ curl http://localhost:8081/api/v1/kv/test

Caution: You are using the Snap version of curl.
Due to Snap's sandbox nature, this version has some limitations.
For example, it may not be able to access hidden folders in your home directory
or other restricted areas of the os.
This means you may encounter errors when using snap curl to download files.
For those case, you might want to use the native curl package.
For details, see: https://github.com/boukendesho/curl-snap/issues/1

To stop seeing this message, run the following command:
curl.snap-acked

{"found":true,"value":"hello","message":"Successfully retrieved key 'test'"}
rutuja@rutuja:~/grpc-client-serve$ curl -X DELETE http://localhost:8081/api/v1/kv/test

Caution: You are using the Snap version of curl.
Due to Snap's sandbox nature, this version has some limitations.
For example, it may not be able to access hidden folders in your home directory
or other restricted areas of the os.
This means you may encounter errors when using snap curl to download files.
For those case, you might want to use the native curl package.
For details, see: https://github.com/boukendesho/curl-snap/issues/1

To stop seeing this message, run the following command:
curl.snap-acked

{"success":true,"message":"Successfully deleted key 'test'"}
rutuja@rutuja:~/grpc-client-serve$ docker-compose logs -f
Attaching to api-service, kv-service
kv-service     | 2025/10/12 01:56:06 Key-Value gRPC server starting on port :50051
api-service    | 2025/10/12 01:56:12 API server starting on port :8080
api-service    | 2025/10/12 01:56:12 Connecting to KV service at kv-service:50051
api-service    | 2025/10/12 01:56:17 GET /api/v1/health 34.84µs
api-service    | 2025/10/12 01:56:47 GET /api/v1/health 28.042µs
api-service    | 2025/10/12 01:57:17 GET /api/v1/health 40.184µs
api-service    | 2025/10/12 01:57:47 GET /api/v1/health 38.605µs
api-service    | 2025/10/12 01:58:17 GET /api/v1/health 29.218µs
api-service    | 2025/10/12 01:58:30 GET /api/v1/health 34.632µs
api-service    | 2025/10/12 01:58:44 POST /api/v1/kv 860.457µs
api-service    | 2025/10/12 01:58:47 GET /api/v1/health 27.852µs
api-service    | 2025/10/12 01:58:57 GET /api/v1/kv/test 484.725µs
api-service    | 2025/10/12 01:59:13 DELETE /api/v1/kv/test 508.546µs
api-service    | 2025/10/12 01:59:17 GET /api/v1/health 28.119µs
^CERROR: Aborting.
rutuja@rutuja:~/grpc-client-serve$ docker-compose down
Stopping api-service ... done
Stopping kv-service  ... done
Removing api-service ... done
Removing kv-service  ... done
Removing network grpc-client-serve_kvstore-network
