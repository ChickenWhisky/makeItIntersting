import random
import requests
import json
import logging

API_URL = "http://localhost:8080/order"
ORDERBOOK_URL = "http://localhost:8080/orderbook"

# Configure logging
logging.basicConfig(
    filename='trading_simulation.log',
    level=logging.INFO,
    format='%(asctime)s - %(message)s',
    datefmt='%Y-%m-%d %H:%M:%S'
)

# Function to send a POST (add order) request
def send_post_request(order):
    response = requests.post(API_URL, json=order)
    
    if response.status_code == 200 or response.status_code == 201:
        log_message = f"POST order: {json.dumps(order)} - Status: {response.status_code}"
        print(log_message)
        logging.info(log_message)
    else:
        log_message = f"POST order failed: {json.dumps(order)} - Status: {response.status_code}"
        print(log_message)
        logging.error(log_message)

def simulate_trading():
    # Situation Explanation:
    # - A sell order of 1100 units at a price of 100
    # - A buy order for 100 units at a price of 20
    # - A sell order of 1000 units at a price of 20
    # - A buy order for 50 units at a price of 100
    # - A sell order for 100 units at price 30
    # - A buy order for 20 units at price 100
    # - A buy order for 100 units at price 30

    # Order 1: Sell 1100 units at price 100
    sell_order_1 = {
        "user_id": "1",  # Example user ID
        "order_type": "sell",
        "price": 100,
        "quantity": 1100
    }
    send_post_request(sell_order_1)

    # Order 2: Buy 100 units at price 20
    buy_order_1 = {
        "user_id": "2",  # Another example user ID
        "order_type": "buy",
        "price": 20,
        "quantity": 100
    }
    send_post_request(buy_order_1)

    # Order 3: Sell 1000 units at price 20
    sell_order_2 = {
        "user_id": "3",  # Another example user ID
        "order_type": "sell",
        "price": 20,
        "quantity": 1000
    }
    send_post_request(sell_order_2)

    # Order 4: Buy 50 units at price 100
    buy_order_2 = {
        "user_id": "4",  # Another example user ID
        "order_type": "buy",
        "price": 100,
        "quantity": 50
    }
    send_post_request(buy_order_2)

    # New Orders:
    # Order 5: Sell 100 units at price 30
    sell_order_3 = {
        "user_id": "5",  # Another example user ID
        "order_type": "sell",
        "price": 30,
        "quantity": 100
    }
    send_post_request(sell_order_3)

    # Order 6: Buy 20 units at price 100
    buy_order_3 = {
        "user_id": "6",  # Another example user ID
        "order_type": "buy",
        "price": 100,
        "quantity": 20
    }
    send_post_request(buy_order_3)

    # Order 7: Buy 100 units at price 30
    buy_order_4 = {
        "user_id": "7",  # Another example user ID
        "order_type": "buy",
        "price": 30,
        "quantity": 100
    }
    send_post_request(buy_order_4)

if __name__ == "__main__":
    simulate_trading()
