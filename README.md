# 🚀 Rate Limiter

Rate Limiter API is a simple **rate-limiting service** using **Go, Redis, and Redigo**.  
It limits API requests **per IP** and **per API key** using Redis as a storage backend.

## 📦 Features

✅ Rate limiting by **IP address**  
✅ Rate limiting by **API key**  
✅ Configurable **request limits**  
✅ Uses **Redis** for efficient storage  
✅ Supports **Docker & Docker Compose**

---

## 🔧 Installation

### **1️⃣ Clone the Repository**

```sh
git clone https://github.com/yourusername/rate-limiter-api.git
cd rate-limiter-api
```

### **2️⃣ How to run**

**Tips: modify the .env to test others rate limiter**

```sh
docker-compose up -d
```

### **3️⃣ How to test**

Use the file in "docs/requests.http" with vscode plugin "REST Client" to test the API.
or use the curl command below:

```sh
curl -i http://localhost:8080/ -H "X-API-KEY: your-api-key
```

An other way to test is use an other project created by me for stress test the API, you can find it [here](https://github.com/isaacmirandacampos/stress-test-cli)
