# 🦫 CapybaraDB in Go  🦫

## Project Goal
The goal of this project is to create a high-performance, scalable, and versatile database system written in Go. This database aims to support both **ACID compliance** and **clustering**, offering robust features for modern applications. It is designed to handle **column-oriented** and **row-oriented** storage, providing flexibility for both transactional and analytical workloads.

---

## Key Features
1. **ACID Compliance**
    - Ensures data integrity with support for transactions, locking mechanisms, and durability through Write-Ahead Logging (WAL).

2. **Clustering**
    - Built-in support for distributed storage, replication, and sharding.
    - Strong consistency protocols (Raft/Paxos) to synchronize data across nodes.

3. **Dual Storage Models**
    - **Column-Oriented Storage**: Optimized for analytical queries, enabling high-speed aggregation and column-based operations.
    - **Row-Oriented Storage**: Ideal for transactional workloads with rapid row-level access.

4. **Indexing and Optimization**
    - Advanced indexing structures (B-Tree, LSM-Tree) for efficient data retrieval.
    - Support for compression algorithms like Snappy or Zstandard to minimize storage footprint.

5. **Concurrency and Performance**
    - Leverages Go’s goroutines for handling multiple operations simultaneously.
    - Built-in caching and Bloom filters to accelerate query performance.

---

## Why Go?
Go is chosen as the primary language due to its:
- Efficient concurrency model with goroutines.
- Simplicity and performance for building low-level, high-performance systems.
- Extensive support for networking and distributed systems.

---

## Why Capybara DB?
Because capybaras are awesome! 🦫


---

## Architecture Overview
1. **Storage Engine**
    - Manages columnar and row-based data layouts.
    - Handles on-disk and in-memory storage seamlessly.

2. **Transaction Manager**
    - Implements isolation levels, locking, and commit/rollback functionality.

3. **Distributed Consensus**
    - Uses Raft for leader election and data synchronization between nodes.

4. **Query Engine**
    - Parses and executes queries, optimizing based on storage type.

5. **APIs**
    - Provides REST and gRPC endpoints for interacting with the database.

---

## Getting Started

### Prerequisites
- Go 1.23 or higher.
- Docker (optional, for clustering and testing).

### Installation
TODO...
