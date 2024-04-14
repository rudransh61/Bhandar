
# Bhandar
<div align="center">
    <h1>Bhandar</h1>
    <img src="./logo.jpg" width="200">
    <p>A basic in-memory storage Real-Time Database in GoLang</p>
</div>

**Description:**
Bhandar is a lightweight, high-performance key-value cache implemented in Go. It provides a simple yet efficient caching solution for storing and retrieving data with support for TTL (time-to-live) expiration.

**Features:**
- In-memory key-value caching.
- Support for TTL-based expiration.
- Real-time command line interface for interacting with the cache.
- HTTP API for setting and getting key-value pairs.
- Built-in support for monitoring cache usage and performance.
- Customizable TTL for each key-value pair.
- Thread-safe operations with mutex locks.

**Installation:**
1. Clone the repository:
   ```
   git clone https://github.com/rudransh61/Bhandar.git
   ```

2. Navigate to the project directory:
   ```
   cd Bhandar
   ```

3. Build the project:
   ```
   go build
   ```

**Usage:**
1. Run the GoCache server:
   ```
   ./bhandar
   ```

2. Access the HTTP API endpoints:
   - Set a key-value pair:
     ```
     curl -X POST "http://localhost:8000/set?key=mykey&value=myvalue&ttl=10s"
     ```
   - Get the value associated with a key:
     ```
     curl "http://localhost:8000/get?key=mykey"
     ```

3. Interact with the real-time command line interface:
   ```
   127.0.0.1:8000> set mykey myvalue 10s
   127.0.0.1:8000> get mykey
   ```

**Contributing:**
Contributions are welcome! Feel free to fork the repository, make improvements, and submit pull requests.
<!-- 
**License:**
This project is licensed under the MIT License. See the LICENSE file for details. -->

**Contact:**
For any questions or feedback, please contact [Rudransh Bhardwaj] at [aarti19830@gmail.com].

**Acknowledgments:**
Special thanks to the Go community and contributors for their support and contributions to this project.