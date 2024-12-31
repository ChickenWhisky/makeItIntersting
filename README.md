## Implementation

The implementation will be divided into stages focusing on key components, starting with backend services, followed by frontend integration, and finally, the admin portal.

### 1. User Service Implementation

The user service is responsible for user account management, balance tracking, and transaction history.

**Key Features:**
- User registration, login, and authentication.
- Fund management (simulated for now).
- APIs to query user details and balances.

**Endpoints:**
- `POST /users/register`: Register a new user.
- `POST /users/login`: Login a user.
- `GET /users/{id}/balance`: Fetch current balance.
- `POST /users/{id}/add-funds`: Simulate adding funds.

**Data Flow:**
- When a user registers, a new entry is created in the `users` table.
- The user’s balance will be initialized to zero and can be incremented via the `add-funds` endpoint.

### 2. Order Book and Contract Service

This service manages event contracts, order creation, matching, and cancellation. It interacts with the **Matching Engine** to handle real-time matching.

**Key Features:**
- Creating new contracts for events.
- Managing the order book, with buy/sell orders.
- Matching engine for real-time matching of contracts.
- Order cancellation for unmatched contracts.

**Endpoints:**
- `POST /contracts`: Create a new contract (Issuer starts a contract and pools a portion of the $1).
- `POST /contracts/{id}/orders`: Place a buy/sell order for a contract.
- `DELETE /contracts/{id}/orders/{orderId}`: Cancel an open order.
- `GET /contracts/{id}/orders`: List current orders for a contract.

**Data Flow:**
- A user creates a contract (e.g., for event A) using the `/contracts` endpoint.
- The contract appears in the order book, and users can match it using the `/orders` endpoint. Partial matching is allowed.
- The matching engine handles real-time contract matching, logging all successful matches for future price history tracking.

### 3. Event Service and Admin Portal

The Event Service handles the lifecycle of an event, including tracking, expiration, and resolution through the admin portal.

**Key Features:**
- Create events with contract expiration dates.
- Manage event outcomes through the admin portal.
- Update contract statuses and user balances based on event outcomes.

**Endpoints:**
- `POST /events`: Admin creates a new event.
- `POST /events/{id}/resolve`: Admin resolves the event with a specified outcome.
- `GET /events/{id}`: Get details of an event.

**Data Flow:**
- The admin creates an event through the admin portal. When the event expires, the outcome is manually resolved via the `/resolve` endpoint.
- Once an event is resolved, the system will automatically calculate the resulting payouts and update user balances.

### 4. Matching Engine Implementation

The **Matching Engine** operates in real-time, continuously matching buy/sell orders on the order book based on partial or full fills. It processes contracts in real-time and checks the order book for possible matches. When a match is found, the order is executed, and funds are locked for the event duration.

**Flow:**
- When an order is placed or updated, the matching engine is triggered.
- The engine continuously checks for potential matches.
- If an order is partially matched, the remaining amount is kept open for future matches.
- Upon event resolution, the winning side gets the full $1 per contract.

**Matching Algorithm:**
- First, check the order book for opposite orders (e.g., a buy order will be matched with the closest sell order).
- If a match is found, allocate the corresponding amount (either partial or full).
- If partial, leave the remaining order open until fully matched or expired.

### 5. Historical Price Tracking for Events

To provide users with a way to visualize the historical matched prices of contracts for each event, we’ll maintain a **Price History** log. This log will store all match events for contracts within an event, including the fractions contributed by each user.

**Steps:**
1. **Database Update**: Add the `PriceHistory` table to track matched prices at the event level.
2. **Update Matching Engine**: Modify the matching engine to log each match into the event's price history.
3. **API Endpoint**: Implement the `GET /events/{id}/history` endpoint to provide the frontend with historical data for graphing.
4. **Graphing on WebUI**: Integrate the historical price data into the frontend, allowing users to visualize trends for specific events.

**Endpoints:**
- `GET /events/{id}/history`: Fetch historical price matches for a specific event.

**Data Flow:**
- Each time a contract is matched, either partially or fully, the **matched price** (e.g., $0.6 pooled by buyer, $0.4 by seller) is recorded with the timestamp.
- The frontend can use this data to graph trends for each event, showing how much users were willing to pool for or against an event over time.

---

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

## Base URL
```
http://localhost:8080
```

## Endpoints

### Create Order
Creates a new order in the order book.

**Endpoint:** `POST /order`  
**Content-Type:** `application/json`

**Request Body:**
```json
{
    "user_id": "string",
    "order_type": "string",    // "buy", "sell", "limit_buy", "limit_sell"
    "price": float,           // required for limit orders
    "quantity": float,
    "request_type": "add"
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
Modifies an existing order in the order book.

**Endpoint:** `PUT /order`  
**Content-Type:** `application/json`

**Request Body:**
```json
{
    "contract_id": "string",
    "user_id": "string",
    "price": float,           // new price
    "quantity": float,        // new quantity
    "request_type": "modify"
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
Cancels an existing order in the order book.

**Endpoint:** `DELETE /order`  
**Content-Type:** `application/json`

**Request Body:**
```json
{
    "contract_id": "string",
    "user_id": "string",
    "request_type": "delete"
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
