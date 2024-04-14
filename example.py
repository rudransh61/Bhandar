import requests

# Function to set a key-value pair in the cache server
def set_key_value(key, value, ttl):
    url = "http://127.0.0.1:8080/set"
    data = {"key": key, "value": value, "ttl": ttl}
    response = requests.post(url, data=data)
    print(response.text)

# Function to get the value associated with a key from the cache server
def get_value(key):
    url = f"http://127.0.0.1:8080/get?key={key}"
    response = requests.get(url)
    print(response.text)

# Example usage
set_key_value("name", "Alice", "60s")  # Set a key-value pair with TTL of 60 seconds
get_value("name")  # Get the value associated with the key "name"
