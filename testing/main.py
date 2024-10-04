import random
import requests
import json
import argparse
import logging
from concurrent.futures import ThreadPoolExecutor, as_completed

# Config for the frequency of deletions (e.g., 20% chance of deletion)
DELETE_PROBABILITY = 0.2
API_URL = "http://localhost:8080/order"
ORDERBOOK_URL = "http://localhost:8080/orderbook"

# Random values range
MIN_PRICE = 45
MAX_PRICE = 60
MIN_QUANTITY = 5
MAX_QUANTITY = 20
USER_RANGE = 10

# Configure logging
logging.basicConfig(
    filename='trading_simulation.log',
    level=logging.INFO,
    format='%(asctime)s - %(message)s',
    datefmt='%Y-%m-%d %H:%M:%S'
)

# Function to send a POST (add order) request
def send_post_request():
    user_id = random.randint(1, USER_RANGE)
    order_type = random.choice(["buy", "sell"])
    price = random.randint(MIN_PRICE, MAX_PRICE)
    quantity = random.randint(MIN_QUANTITY, MAX_QUANTITY)

    order = {
        "user_id": str(user_id),
        "order_type": order_type,
        "price": price,
        "quantity": quantity
    }

    response = requests.post(API_URL, json=order)
    
    if response.status_code == 200 or response.status_code == 201:
        log_message = f"POST order: {json.dumps(order)} - Status: {response.status_code}"
        print(log_message)
        logging.info(log_message)
    else:
        log_message = f"POST order failed: {json.dumps(order)} - Status: {response.status_code}"
        print(log_message)
        logging.error(log_message)

# Function to send a DELETE (cancel order) request
def send_delete_request():
    user_id = random.randint(1, USER_RANGE)
    order_type = random.choice(["buy", "sell"])
    price = random.randint(MIN_PRICE, MAX_PRICE)

    order = {
        "user_id": str(user_id),
        "order_type": order_type,
        "price": price
    }

    response = requests.delete(API_URL, json=order)
    
    if response.status_code == 200:
        log_message = f"DELETE order: {json.dumps(order)} - Status: {response.status_code}"
        print(log_message)
        logging.info(log_message)
    else:
        log_message = f"DELETE order failed: {json.dumps(order)} - Status: {response.status_code}"
        print(log_message)
        logging.error(log_message)

# Function to check the order book for negative values
def check_order_book():
    response = requests.get(ORDERBOOK_URL)
    if response.status_code == 200:
        order_book = response.json()
        if order_book is None:
            log_message = "ERROR: Order book is None."
            print(log_message)
            logging.error(log_message)
            return

        negative_asks = []
        negative_bids = []

        asks = order_book.get("asks")
        if asks is not None:
            for ask in asks:
                try:
                    price = ask['price']
                    quantity = ask['quantity']
                    if int(quantity) < 0:
                        negative_asks.append((price, quantity))
                except (ValueError, KeyError):
                    continue

        bids = order_book.get("bids")
        if bids is not None:
            for bid in bids:
                try:
                    price = bid['price']
                    quantity = bid['quantity']
                    if int(quantity) < 0:
                        negative_bids.append((price, quantity))
                except (ValueError, KeyError):
                    continue
        
        if negative_asks or negative_bids:
            log_message = f"ERROR: Negative values found in order book - Asks: {negative_asks}, Bids: {negative_bids}"
            print(log_message)
            logging.error(log_message)
        else:
            log_message = "No negative values found in order book."
            print(log_message)
            logging.info(log_message)
    else:
        log_message = f"Failed to fetch order book: Status {response.status_code}"
        print(log_message)
        logging.error(log_message)

# Main function to manage concurrency
def simulate_trading(num_requests):
    order_count = 0
    with ThreadPoolExecutor(max_workers=10) as executor:
        futures = []
        
        for _ in range(num_requests):
            # Randomly choose between sending a POST or DELETE request
            if random.random() < DELETE_PROBABILITY:
                futures.append(executor.submit(send_delete_request))
            else:
                futures.append(executor.submit(send_post_request))
            
            order_count += 1

            # Every 50 orders, check the order book
            if order_count % 50 == 0:
                futures.append(executor.submit(check_order_book))
            
        # Wait for all the futures to complete
        for future in as_completed(futures):
            future.result()

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Simulate trading by sending POST and DELETE requests.")
    parser.add_argument("num_requests", type=int, help="Number of requests to send")
    args = parser.parse_args()

    simulate_trading(args.num_requests)