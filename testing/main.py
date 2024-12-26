import random
import requests
import json
import argparse
import logging
import time
from concurrent.futures import ThreadPoolExecutor, as_completed

# Config for the frequency of deletions (e.g., 20% chance of deletion)
DELETE_PROBABILITY = 0.2
API_URL = "http://localhost:8000/order"
ORDERBOOK_URL = "http://localhost:8000/orderbook"
TRADES_URL = "http://localhost:8000/trades/1"

# Random values range
MIN_PRICE = 50
MAX_PRICE = 55
MIN_QUANTITY = 100
MAX_QUANTITY = 1000
USER_RANGE = 100

# Global counter for contract numbers
global_counter = 0

# Dictionary to store contract IDs for each user
user_contracts = {}

# Configure logging
logging.basicConfig(
    filename='trading_simulation.log',
    level=logging.INFO,
    format='%(asctime)s - %(message)s',
    datefmt='%Y-%m-%d %H:%M:%S'
)


# Function to send a POST (add order) request
def send_post_request():
    global global_counter
    user_id = random.randint(1, USER_RANGE)
    order_type = random.choice(["buy", "sell", "limit_buy", "limit_sell"])
    if order_type in ["limit_buy", "buy"]:
        price = random.randint(MIN_PRICE, MAX_PRICE)
    else:
        price = random.randint(MAX_PRICE - 3, MAX_PRICE + 5)
    quantity = random.randint(MIN_QUANTITY, MAX_QUANTITY)

    order = {
        "user_id": str(user_id),
        "contract_id": str(global_counter),
        "request_type": "add",
        "order_type": order_type,
        "price": price,
        "quantity": quantity
    }

    response = requests.post(API_URL, json=order)

    if response.status_code in [200, 201]:
        log_message = f"POST order: {json.dumps(order)} - Status: {response.status_code}"
        print(log_message)
        logging.info(log_message)
        # Store the contract ID for the user
        if user_id not in user_contracts:
            user_contracts[user_id] = []
        user_contracts[user_id].append(global_counter)
        global_counter += 1  # Increment the global counter after a successful order
    else:
        log_message = f"POST order failed: {json.dumps(order)} - Status: {response.status_code}"
        print(log_message)
        logging.error(log_message)


# Function to send a DELETE (cancel order) request
def send_delete_request():
    user_id = random.randint(1, USER_RANGE)
    if user_id in user_contracts and user_contracts[user_id]:
        contract_id = random.choice(user_contracts[user_id])  # Randomly select a contract ID for the user

        order = {
            "user_id": str(user_id),
            "contract_id": str(contract_id),
            "request_type": "delete"
        }

        response = requests.delete(API_URL, json=order)

        if response.status_code == 200:
            log_message = f"DELETE order: {json.dumps(order)} - Status: {response.status_code}"
            print(log_message)
            logging.info(log_message)
            # Remove the contract ID from the user's list
            user_contracts[user_id].remove(contract_id)
        else:
            log_message = f"DELETE order failed: {json.dumps(order)} - Status: {response.status_code}"
            print(log_message)
            logging.error(log_message)
    else:
        log_message = f"DELETE order failed: No contract ID found for user {user_id}"
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


# Function to get the last traded price
def get_last_traded_price():
    response = requests.get(TRADES_URL)
    if response.status_code == 200:
        last_traded_price = response.json()
        log_message = f"Last traded price: {last_traded_price}"
        print(log_message)
        logging.info(log_message)
    else:
        log_message = f"Failed to get the last traded price. Status code: {response.status_code}"
        print(log_message)
        logging.error(log_message)


# Main function to manage concurrency
def simulate_trading(num_requests):
    with ThreadPoolExecutor(max_workers=10) as executor:
        futures = []
        order_count = 0

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

            # Random delay to mimic real-world trading
            time.sleep(random.uniform(0.1, 0.5))

        # Wait for all the futures to complete
        for future in as_completed(futures):
            future.result()

        # Get the last traded price after all requests are completed
        get_last_traded_price()


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Simulate trading by sending POST and DELETE requests.")
    parser.add_argument("num_requests", type=int, help="Number of requests to send")
    args = parser.parse_args()

    simulate_trading(args.num_requests)
