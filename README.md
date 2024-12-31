# Setup

To start the server either build the executable and run the executable

##### Build executable
```bash
go build main.go
./main.go
```

or we could simply run the following code
```bash
go run main.go
```
The service will then be available at `localhost:8000`
To change the above simply change the variable PORT from the .env file


# Endpoints
#### Create Order
Creates a new order in the order book.
**Endpoint:** `POST /order`  
**Content-Type:** `application/json`

**Request Body:**
```json
{
    "contract_id" : int64     // Unique Contract ID in the trading system  
    "user_id": "string",      // Unique User ID
    "order_type": "string",   // "buy", "sell", "limit_buy", "limit_sell"
    "price": float,           // required for limit orders
    "quantity": float,        // 
}
```

**Response:**
```json
{
    "contract_id": "string",
    "user_id": "string",
    "order_type": "string",
    "price": float,
    "quantity": float,
    "timestamp": int64
}
```

**Status Codes:**
- 200: Success
- 400: Invalid request
- 500: Server error

### Modify Order
Modifies an existing order in the order book if it hasnt been executed yet.

**Endpoint:** `PUT /order`  
**Content-Type:** `application/json`
**Notes :** User cannot change **Order Type**. Any attempt to add order type to request body is ignored
**Request Body:**
```json
{
    "contract_id": "string",  // Unique Contract ID in the trading system  
    "user_id": "string",      // Unique User 
    "price": float,           // new price
    "quantity": float,        // new quantity
}
```

**Response:**
```json
{
    "contract_id": "string",
    "user_id": "string",
    "order_type": "string",
    "price": float,
    "quantity": float,
    "timestamp": int64
}
```

**Status Codes:**
- 200: Success
- 400: Invalid request
- 404: Order not found
- 500: Server error

### Cancel Order
Cancels an existing order in the order book if it hasnt been executed yet.

**Endpoint:** `DELETE /order`  
**Content-Type:** `application/json`

**Request Body:**
```json
{
    "contract_id": "string",
    "user_id": "string",
}
```

**Response:**
```json
{
    "message": "Order successfully cancelled"
}
```

**Status Codes:**
- 200: Success
- 400: Invalid request
- 404: Order not found
- 500: Server error

### Get Last Trades
Retrieves the last N completed trades.

**Endpoint:** `GET /trades/:noOfContracts`

**Parameters:**
- `noOfContracts`: Number of recent trades to retrieve (integer)

**Response:**
```json
{
    "trades": [
        {
            "price": float,
            "quantity": float,
            "timestamp": int64,
            "buyer_id": "string",
            "seller_id": "string"
        }
    ]
}
```

**Status Codes:**
- 200: Success
- 400: Invalid request
- 500: Server error


## Implementation

The implementation will be divided into stages focusing on key components, starting with backend services, followed by frontend integration, and finally, the admin portal.




### Database Schema (Summary)
Here’s a summary of the key database tables involved in the implementation:

1. **Users**:
   - `id`: Unique ID of the user.
   - `username`: Username of the user.
   - `balance`: Current balance available to the user.
   - `transaction_history`: A JSON field to store transaction logs.

2. **Contracts**:
   - `contract_id`: Unique ID of the contract.
   - `commodity`: The commodity or event tied to this contract.
   - `issuer_id`: ID of the user who issued the contract.
   - `amount_issued`: Number of contracts issued by the user.
   - `fraction_pooled`: The amount of each contract pooled by the issuer (e.g., 0.6).
   - `remaining_fraction`: The portion left to be filled by other users (e.g., 0.4).
   - `expiration_date`: Expiry of the contract based on the event's resolution.
   - `status`: Status of the contract (Open, Matched, Expired).

3. **OrderBook**:
   - `order_id`: Unique ID of the order.
   - `contract_id`: Contract associated with this order.
   - `user_id`: User who placed the order.
   - `type`: Order type (Buy or Sell).
   - `amount`: Number of contracts issued or requested.
   - `fraction_pooled`: The fraction of each contract pooled by the user.
   - `remaining_fraction`: The fraction left to be pooled by other users.
   - `price`: Always set to $1 per contract.
   - `status`: Status of the order (Pending, Partially Matched, Completed, Cancelled).

4. **Events**:
   - `event_id`: Unique ID of the event.
   - `name`: Name of the event.
   - `end_date`: Date when the event ends.
   - `outcome`: Outcome of the event (Resolved by admin).

5. **PriceHistory**:
   - `history_id`: Unique ID for the price history entry.
   - `event_id`: The ID of the event being tracked.
   - `timestamp`: The time when the match occurred.
   - `buy_fraction`: Fraction pooled by the buyer (betting for the event).
   - `sell_fraction`: Fraction pooled by the seller (betting against the event).
   - `match_amount`: Number of contracts matched in this transaction.

---

### Key Features Implemented

- [ ] **User Management**: Registration, login, balance, and simulated fund transfers.
- [ ] **Order Book**: Users can place and cancel orders for contracts.
- [ ] **Real-Time Matching Engine**: Handles partial and full matching of contracts.
- [ ] **Price History Tracking**: Logs matched prices for each event and exposes this data via REST APIs.
- [ ] **Event Management**: Admins create events and resolve outcomes, with user balances updated accordingly.

# Order Book API Documentation

## Overview
This API provides endpoints to interact with an order matching engine. Users can create, modify, and cancel orders, as well as retrieve trade history. The system includes self-trade prevention mechanisms to prevent users from accidentally trading with themselves.

## Self-Trade Prevention
The system implements self-trade prevention to stop users from accidentally trading with themselves. The following modes are available:

- `cancel_newest`: Rejects the incoming order if it would match with an existing order from the same user
- `cancel_oldest`: Cancels existing orders and adds the new one
- `cancel_both`: Cancels both existing and new orders

## Error Responses
All error responses follow this format:
```json
{
    "error": "Error message describing what went wrong"
}
```

## Rate Limiting
Please note that API rate limits may apply. Contact the system administrator for specific limits.

## Example Usage

### Creating a Limit Buy Order
```bash
curl -X POST http://localhost:8080/order \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "0",
    "order_type": "buy",
    "price": 100,
    "quantity": 10,
    "request_type": "add"
  }'
```

### Cancelling an Order
```bash
curl -X DELETE http://localhost:8080/order \
  -H "Content-Type: application/json" \
  -d '{
    "contract_id": "abc123",
    "user_id": "user123",
    "request_type": "delete"
  }'
```

### Getting Last 5 Trades
```bash
curl http://localhost:8080/trades/5
```
