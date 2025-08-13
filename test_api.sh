#!/bin/bash

# Go Fintech Backend API Test Script
BASE_URL="http://localhost:8080"
JWT_TOKEN=""

echo "🚀 Starting API Tests..."
echo "================================"

# Test 1: Health Check
echo "1. Testing Health Check..."
curl -s -X GET "$BASE_URL/" | jq '.' || echo "Health check failed"

# Test 2: Register User
echo -e "\n2. Testing User Registration..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "role": "user"
  }')

echo "Register Response:"
echo "$REGISTER_RESPONSE" | jq '.' || echo "$REGISTER_RESPONSE"

# Test 3: Login User
echo -e "\n3. Testing User Login..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }')

echo "Login Response:"
echo "$LOGIN_RESPONSE" | jq '.' || echo "$LOGIN_RESPONSE"

# Extract JWT token
JWT_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token' 2>/dev/null)
if [ "$JWT_TOKEN" = "null" ] || [ -z "$JWT_TOKEN" ]; then
    echo "❌ Failed to get JWT token"
    exit 1
fi

echo "✅ JWT Token obtained: ${JWT_TOKEN:0:20}..."

# Test 4: Credit Transaction
echo -e "\n4. Testing Credit Transaction..."
CREDIT_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/transactions/credit" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "amount": 100.50,
    "description": "Test credit"
  }')

echo "Credit Response:"
echo "$CREDIT_RESPONSE" | jq '.' || echo "$CREDIT_RESPONSE"

# Test 5: Current Balance
echo -e "\n5. Testing Current Balance..."
BALANCE_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/balances/current" \
  -H "Authorization: Bearer $JWT_TOKEN")

echo "Balance Response:"
echo "$BALANCE_RESPONSE" | jq '.' || echo "$BALANCE_RESPONSE"

# Test 6: Transaction History
echo -e "\n6. Testing Transaction History..."
HISTORY_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/transactions/history" \
  -H "Authorization: Bearer $JWT_TOKEN")

echo "History Response:"
echo "$HISTORY_RESPONSE" | jq '.' || echo "$HISTORY_RESPONSE"

# Test 7: Prometheus Metrics
echo -e "\n7. Testing Prometheus Metrics..."
METRICS_RESPONSE=$(curl -s -X GET "$BASE_URL/metrics")
echo "Metrics endpoint accessible: $(echo "$METRICS_RESPONSE" | wc -l) lines"

# Test 8: List Users (Admin only)
echo -e "\n8. Testing List Users (Admin only)..."
USERS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/users" \
  -H "Authorization: Bearer $JWT_TOKEN")

echo "Users Response:"
echo "$USERS_RESPONSE" | jq '.' || echo "$USERS_RESPONSE"

# Test 9: Get User by ID
echo -e "\n9. Testing Get User by ID..."
USER_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/users/1" \
  -H "Authorization: Bearer $JWT_TOKEN")

echo "User Response:"
echo "$USER_RESPONSE" | jq '.' || echo "$USER_RESPONSE"

# Test 10: Update User
echo -e "\n10. Testing Update User..."
UPDATE_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/v1/users/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "username": "updateduser",
    "email": "updated@example.com",
    "role": "admin"
  }')

echo "Update Response:"
echo "$UPDATE_RESPONSE" | jq '.' || echo "$UPDATE_RESPONSE"

echo -e "\n✅ API Tests Completed!"
echo "================================"
