# 🧠 NameRegChain — Educational Blockchain in Go

**NameRegChain** is a minimal, educational blockchain implementation written in Go, focused on understanding how account-based state, transaction validation, signing, consensus, and blocks work together.

This project implements:
- ✅ Account-based state (like Ethereum)
- ✅ ECDSA signature verification
- ✅ Nonce enforcement to prevent replay
- ✅ Simple transaction logic (name → address registration)
- ✅ A basic block structure and chain
- ✅ JSON HTTP API for interaction
- ✅ Basic test suite for key edge cases

---

## 🚀 Project Goals

This is not meant to be production-ready. It’s a hands-on educational tool to:
- Understand blockchain data structures
- Learn how signed transactions are validated
- Explore nonce-based transaction ordering
- Simulate a basic chain with name registration as state

---

## 📁 Features

| Feature                         | Description                                  |
|----------------------------------|----------------------------------------------|
| Account-based state              | Each account is identified by an address (public key hash) |
| ECDSA transaction signing        | Transactions must be signed by account owner |
| Nonce system                     | Ensures transaction order and prevents replay |
| REST API                         | POST `/register`, GET `/lookup`, `/nonce`, `/chain` |
| Simplified consensus             | 1 transaction per block, no forks or peers   |
| Error handling                   | Detects duplicate names, invalid signatures, nonce reuse |

---

## 🧪 API Endpoints

| Method | Endpoint        | Description                          |
|--------|------------------|--------------------------------------|
| POST   | `/register`      | Submit signed tx to register a name  |
| GET    | `/lookup?name=`  | Fetch registered address for a name  |
| GET    | `/nonce?address=`| Get current nonce for an address     |
| GET    | `/chain`         | View entire block history            |

---

## 🧪 Running the Project

```bash
go run ./cmd/namereg
```

Use curl or Postman to interact with the API. Sample registration flow:

```bash
curl -X POST localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{ "name": "alice", "address": "0x...", "nonce": 0, "pubKey": "...", "signature": "..." }'

curl localhost:8080/lookup?name=alice
```

---

## 🧪 Running Tests

```bash
go test ./tests
```

Covers:
- Successful registration
- Duplicate name error
- Invalid nonce (replay or future)
- Signature mismatch

---

## 📚 Educational Value

This project is ideal for:
- Students learning how blockchains work internally
- Engineers transitioning into Web3 from distributed systems
- Anyone wanting to build a minimal smart chain from scratch

---

## 🛠 Future Additions

- Mempool and batch block building
- libp2p peer discovery and gossiping
- Merkle roots for block state integrity
- Cosmos-SDK-style module system
- zkProof integration (eventually)

---

## 🧑‍💻 Author

Educational project maintained by a distributed systems & Web3 engineer.  
Meant to serve as a hands-on learning journey into the internals of blockchain design.

---

## 📜 License

MIT — use freely, modify, and build on top.