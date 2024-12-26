import requests
import matplotlib.pyplot as plt

# API URL
API_URL = "http://localhost:8080/tr"

def fetch_order_book():
    response = requests.get(API_URL)
    if response.status_code == 200:
        return response.json()
    else:
        print(f"Failed to fetch order book: Status {response.status_code}")
        return None

def plot_last_50_matched_prices(prices):
    # Set the figure size dynamically based on the number of prices
    fig_width = max(10, len(prices) * 0.2)  # Minimum width of 10, increase by 0.2 per data point
    fig_height = 5  # Fixed height

    plt.figure(figsize=(fig_width, fig_height))
    plt.plot(prices, marker='o', linestyle='-', color='b')
    plt.title('Last 50 Matched Prices')
    plt.xlabel('Order Index')
    plt.ylabel('Price')
    plt.grid(True)
    plt.show()

if __name__ == "__main__":
    order_book = fetch_order_book()
    if order_book and "last_50_matched_prices" in order_book:
        last_50_matched_prices = order_book["last_50_matched_prices"]
        plot_last_50_matched_prices(last_50_matched_prices)
    else:
        print("No data available for last 50 matched prices.")