# Simple-Blockchain-Client v1

한국어 버전 : [Korean README](README.ko.md)

## ▶️ Demo & Presentation Video (Please watch it first.)

[YouTube](https://www.youtube.com/watch?v=h-vbAZ6Hku4)

---

## 📖 Project Overview

This project is a simple Ethereum-inspired blockchain client implemented **from scratch** in Go.  
It runs entirely on a local environment and includes:

<u>**P2P networking, transaction and block validation, propagation, storage, and execution.**</u>

---

## 📖 Project Goals

- **Learn Go in depth**  
  I wanted to explore Go’s features (goroutines, channels, etc.) by using them in a real project instead of toy examples.

- **Deepen understanding of blockchain core concepts**  
  By implementing core blockchain logic myself, I aimed to translate theory into concrete code and feel why each component is necessary.  
  For example, I initially serialized data using JSON. Through this, I experienced vague typing and inefficient size, which helped me appreciate why encodings like RLP are more suitable and efficient in real blockchain systems.

- **Improve skills in distributed systems and object-oriented design**  
  As the project grew, I realized that maintaining and extending the codebase requires structured design.  
  This led me to apply object-oriented design principles and modularize the system so the client remains maintainable and extensible as it scales.

---

## ⚙️ Package-Based Architecture & Data Flow

![image](https://github.com/user-attachments/assets/3c368d0d-a882-4c3a-a203-50349dedb728)

### Core Packages

| Name           | Description                                                                                                        |
| -------------- | ------------------------------------------------------------------------------------------------------------------ |
| **P2P**        | Runs UDP/TCP servers and broadcasts messages to peers                                                              |
| **Blockchain** | Blockchain core. Includes the `BlockProcessor` for handling txs/blocks, a periodic `BlockCreator`, and a `Mempool` |
| **Rpc-server** | Exposes JSON-RPC endpoints for external communication                                                              |
| **LevelDB**    | Persists blocks and account state                                                                                  |

### Other Packages

| Name          | Description                                                            |
| ------------- | ---------------------------------------------------------------------- |
| **Account**   | Validates, creates, stores, and updates accounts                       |
| **Mediator**  | Mediates data exchange between packages                                |
| **Bootnode**  | Bootstrap node for node discovery and peer sharing                     |
| **Constants** | Manages configuration such as block interval, tx-per-block limit, etc. |
| **Utils**     | Global utility functions (e.g. Keccak256)                              |

Network protocols such as `Node Discovery` and `P2P` are described in more detail in  
[p2p/README.md](https://github.com/DreamBoysYJ/simple-blockchain-client/tree/main/p2p) (Cmd/Ctrl + click).

---

## 🛠 Usage Guide

### Prerequisites

Please install and prepare the following tools before running the project:

1. **Go**  
   Download & install: [Go official site](https://go.dev/dl/)

2. **Postman Desktop Agent**  
   Download: [Postman Desktop Agent](https://www.postman.com/downloads/)  
   You’ll need the Desktop Agent to send HTTP requests to `localhost` from Postman.

3. **Postman API Documentation**  
   Docs: [Postman API Docs](https://documenter.getpostman.com/view/25348775/2sAYQWLZZ9)  
   You can use these to send transactions, query blocks, and interact with the node via Postman.

4. **Notion API Documentation (Fallback)**  
   Docs: [Notion API Docs](https://ivory-gerbera-298.notion.site/Simple-Blockchain-client-v1-API-17c57a963d328091b9c1fbdc405345a0)  
   If the Postman docs get stuck on infinite loading multiple times, please refer to the Notion page instead.

5. **Firewall (Windows)**  
   On Windows, the program may fail to run due to firewall restrictions.  
   Temporarily disable the firewall (or allow the app/ports) if you run into connectivity issues.

6. **Restarting on Errors**  
   If you encounter unexpected errors during execution, try the following:
   - Delete the `db` folder in the project to reset the database.
   - Or restart from scratch: first run the bootstrap node, then start the full nodes again in order.

---

## 🚀 Installation & Run

1. **Clone the project and move into the directory:**

   ```bash
   git clone https://github.com/DreamBoysYJ/simple-blockchain-client.git
   cd simple-blockchain-client
   ```

2. **(Optional) Set up a global command:**:  
    Run the following so you can execute the main commands globally:
   `bash
make all
`
   This builds the project and installs the binary into /usr/local/bin, so you can use it system-wide.
   (Even if you set up the global command, it’s still recommended to run commands from the simple-blockchain-client directory when you want to be sure the DB is fully deleted and re-initialized.)

3. **`bootnode(bootstrap node)` RUN**:

   ```bash
   go run . -nodeID=boot -mode=bootnode

   # If you completed step 2:

   simple-blockchain-client -nodeID=boot -mode=bootnode
   ```

   The bootnode runs a UDP server and is used to collect fullnode addresses when nodes start up and connect.
   I designed a simple protocol inspired by Node Discovery, and the bootnode receives the addresses of all nodes that communicate with it.

4. **Run at least 3 `Fullnodes` in different terminals:**:

   ```bash
   go run . -nodeID=node1 -mode=fullnode -port=30301 -rpcport=8081
   go run . -nodeID=node2 -mode=fullnode -port=30302 -rpcport=8082
   go run . -nodeID=node3 -mode=fullnode -port=30303 -rpcport=8083

   ### If you completed step 2:

   simple-blockchain-client -nodeID=node1 -mode=fullnode -port=30301 -rpcport=8081
   simple-blockchain-client -nodeID=node2 -mode=fullnode -port=30302 -rpcport=8082
   simple-blockchain-client -nodeID=node3 -mode=fullnode -port=30303 -rpcport=8083
   ```

---

### Flag Description

| Flag      | Description                                                                                               | Default    |
| --------- | --------------------------------------------------------------------------------------------------------- | ---------- |
| `port`    | UDP server port for `Node Discovery` and TCP server port for P2P communication                            | 30303      |
| `rpcport` | JSON-RPC server port used to communicate with external tools (browser, DApp, Postman, etc.)               | 8080       |
| `nodeID`  | Identifier used to distinguish each node. Also used to set the database path (`dbPath`) for local testing | `default`  |
| `mode`    | Node role: `bootnode` or `fullnode`. If omitted, it defaults to fullnode                                  | `fullnode` |

![image](https://github.com/user-attachments/assets/5157266f-d262-4353-aa5c-ed9f64853e53)
As shown above, the client:

Creates accounts for each node,

Creates the genesis block,

Connects nodes together and builds the P2P network.

---

### 5. Testing Transactions with Postman

1. **Check the current state**:  
   First, call `getLastBlock`, `getBlockNumber`, and `getAccountInfo` to verify the current state.

2. **Create transactions**:  
   Use `tx1`–`tx10` and `SendTransaction` to send signed transactions.
   - `SendRawTransaction` is preconfigured with the `signature`.
   - The signed message is constructed by concatenating `from`, `to`, `value`, and `nonce`.
   - Currently, only a specific `from` address is allowed. Since MetaMask or similar tools are not used, the private key for this address is hardcoded, and in the genesis block this address is set as the miner with an initial balance of `10000`.
   - You can execute `tx1`–`tx10` in any order because the mempool sorts transactions by account and `nonce`.
   - However, at the moment only transactions from a single address are stored in the mempool, so you must send at least `tx1`–`tx5` for a block to be created. (When building a block, the node pulls transactions from the mempool in round-robin order by address and nonce.)

![image](https://github.com/user-attachments/assets/02d3b079-d030-4886-8f23-867848fb830b)
As shown above, each transaction is validated and then propagated to peers.  
 To prevent infinite propagation, any duplicate transaction that is already in the mempool is dropped.

---

### 6. Checking Block Creation in the Terminal

1. **Block creation**:  
   Each node periodically checks the mempool, and when a certain number of transactions is available, it attempts to create a block and broadcast it.

2. **Block validation**:  
   Nodes that receive a block verify the previous block reference, the Merkle tree, and all transactions inside the block.

3. **State changes**:  
   After validation, the node stores the block, executes the transactions, and adds `1000` as a block reward to the miner’s address.

4. **Block propagation**:  
   Finally, the node propagates the block to its peers.

![image](https://github.com/user-attachments/assets/6ed740de-9805-4b4c-839c-9bb2f8708163)
As shown above, after the validation process, the node stores the block and updates the DB.  
If any of the block’s transactions are still present in the node’s own mempool, they are removed.  
This prevents the same transaction from being included in multiple blocks when a node creates and broadcasts a block first.

---

### 7. Verifying the Updated State with Postman

Call `getLastBlock`, `getBlockNumber`, and `getAccountInfo` again to check the updated state.  
You should see that:

- The block producer’s address has received the block reward.
- The genesis miner, who is also the transaction sender (`from`), now has an updated `nonce` and balance.

<img width="1040" alt="image" src="https://github.com/user-attachments/assets/eec76975-7f7a-412d-b4d1-6de7c6181acf" />
