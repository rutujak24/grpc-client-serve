**Diagram:**
```
Main Goroutine              Signal Goroutine
     |                            |
grpcServer.Serve()          Wait for signal
     |                            |
Handle requests             [Ctrl+C pressed]
     |                            |
     |                      GracefulStop()
     |                            |
     └──────── shutdown ──────────┘
```


## **🏗️ CLIENT-SERVER ARCHITECTURE**

### **What is Client-Server?**
```
Client                          Server
(Requester)                     (Responder)
   |                               |
   |──── Send Request ────────────>|
   |                               |
   |                          Process request
   |                               |
   |<─── Send Response ────────────|
   |                               |
```

### **In Our System - TWO Client-Server Relationships:**

#### **1. User → API Service (HTTP Client-Server)**
```
curl (Client)                API Service (Server)
      |                            |
      |── POST /api/v1/kv ────────>|
      |   {"key": "x"}             |
      |                       Parse JSON
      |                            |
      |<─── 200 OK ────────────────|
      |   {"success": true}        |
```

**Client Side:** User's browser, curl, mobile app
**Server Side:** API Service (HTTP server on port 8080)


#### **2. API Service → KV Service (gRPC Client-Server)**
```
API Service (Client)       KV Service (Server)
      |                            |
      |── Store(key, value) ──────>|
      |                       Save to map
      |                            |
      |<─── StoreResponse ─────────|
      |   {success: true}          |
```

**Client Side:** API Service uses `KeyValueServiceClient`
**Server Side:** KV Service implements `KeyValueServiceServer`



### **Full Request Flow:**
```
1. User's curl command
   ↓ HTTP POST /api/v1/kv
2. API Service (Server for HTTP, Client for gRPC)
   ├─ Receives HTTP request
   ├─ Parses JSON
   ├─ Creates gRPC request
   ↓ gRPC Store()
3. KV Service (Server for gRPC)
   ├─ Receives gRPC request
   ├─ Stores in map
   ├─ Returns gRPC response
   ↑ gRPC StoreResponse
4. API Service
   ├─ Receives gRPC response
   ├─ Converts to JSON
   ↑ HTTP 200 OK
5. User sees JSON response
```
