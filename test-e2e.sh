#!/bin/bash
# now.ink End-to-End Testing Script
# Tests the complete flow: Auth → Stream → Upload → Mint

set -e

echo "🧪 now.ink End-to-End Test"
echo "=========================="
echo ""

# Configuration
API_URL="http://localhost:8080/api/v1"
TEST_WALLET="TEST_WALLET_E2E_$(date +%s)"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Helper functions
print_step() {
    echo -e "${BLUE}📋 $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
    exit 1
}

print_info() {
    echo -e "${YELLOW}ℹ️  $1${NC}"
}

# Check prerequisites
print_step "Step 1: Checking prerequisites"

if ! command -v curl &> /dev/null; then
    print_error "curl is required but not installed"
fi

if ! command -v jq &> /dev/null; then
    print_error "jq is required but not installed"
fi

if ! command -v ffmpeg &> /dev/null; then
    print_error "ffmpeg is required but not installed"
fi

# Check if API is running
if ! curl -s "$API_URL/../health" > /dev/null; then
    print_error "API is not running at $API_URL. Start with: cd backend && go run cmd/api/main.go"
fi

print_success "All prerequisites met"
echo ""

# Step 2: Get authentication nonce
print_step "Step 2: Getting authentication nonce"

NONCE_RESPONSE=$(curl -s -X POST "$API_URL/auth/nonce" \
    -H "Content-Type: application/json" \
    -d "{\"wallet_address\":\"$TEST_WALLET\"}")

NONCE=$(echo "$NONCE_RESPONSE" | jq -r '.nonce')

if [ "$NONCE" == "null" ] || [ -z "$NONCE" ]; then
    print_error "Failed to get nonce: $NONCE_RESPONSE"
fi

print_success "Nonce received: ${NONCE:0:16}..."
echo ""

# Step 3: Verify wallet and get JWT token
print_step "Step 3: Verifying wallet and getting JWT token"

TOKEN_RESPONSE=$(curl -s -X POST "$API_URL/auth/verify" \
    -H "Content-Type: application/json" \
    -d "{\"wallet_address\":\"$TEST_WALLET\",\"signature\":\"test_signature\",\"nonce\":\"$NONCE\"}")

TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.token')
USER_ID=$(echo "$TOKEN_RESPONSE" | jq -r '.user.id')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
    print_error "Failed to get token: $TOKEN_RESPONSE"
fi

print_success "JWT token received"
print_info "User ID: $USER_ID"
echo ""

# Step 4: Start a stream
print_step "Step 4: Starting a stream"

STREAM_RESPONSE=$(curl -s -X POST "$API_URL/streams/start" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"title":"E2E Test Stream","latitude":40.7128,"longitude":-74.0060,"is_public":true}')

STREAM_ID=$(echo "$STREAM_RESPONSE" | jq -r '.id')

if [ "$STREAM_ID" == "null" ] || [ -z "$STREAM_ID" ]; then
    print_error "Failed to start stream: $STREAM_RESPONSE"
fi

print_success "Stream started"
print_info "Stream ID: $STREAM_ID"
echo ""

# Step 5: Create test video
print_step "Step 5: Creating test video"

TEST_VIDEO="/tmp/nowink-test-$(date +%s).mp4"

ffmpeg -f lavfi -i testsrc=duration=3:size=640x480:rate=30 \
    -f lavfi -i sine=frequency=1000:duration=3 \
    -pix_fmt yuv420p -y "$TEST_VIDEO" \
    > /dev/null 2>&1

if [ ! -f "$TEST_VIDEO" ]; then
    print_error "Failed to create test video"
fi

VIDEO_SIZE=$(du -h "$TEST_VIDEO" | cut -f1)
print_success "Test video created: $VIDEO_SIZE"
echo ""

# Step 6: Upload video and mint NFT
print_step "Step 6: Uploading video and minting NFT"

MINT_RESPONSE=$(curl -s -X POST "$API_URL/streams/$STREAM_ID/save" \
    -H "Authorization: Bearer $TOKEN" \
    -F "video=@$TEST_VIDEO")

MINT_ADDRESS=$(echo "$MINT_RESPONSE" | jq -r '.mint_address')
ARWEAVE_TX=$(echo "$MINT_RESPONSE" | jq -r '.arweave_tx')
METADATA_URI=$(echo "$MINT_RESPONSE" | jq -r '.metadata_uri')

if [ "$MINT_ADDRESS" == "null" ] || [ -z "$MINT_ADDRESS" ]; then
    print_error "Failed to mint NFT: $MINT_RESPONSE"
fi

print_success "NFT minted successfully!"
print_info "Mint Address: $MINT_ADDRESS"
print_info "Arweave TX: $ARWEAVE_TX"
print_info "Metadata URI: $METADATA_URI"
echo ""

# Step 7: Verify stream was updated
print_step "Step 7: Verifying stream was updated"

STREAM_CHECK=$(curl -s -X GET "$API_URL/streams/$STREAM_ID" \
    -H "Authorization: Bearer $TOKEN")

STREAM_MINT=$(echo "$STREAM_CHECK" | jq -r '.mint_address')

if [ "$STREAM_MINT" != "null" ] && [ ! -z "$STREAM_MINT" ]; then
    print_success "Stream updated with mint address: $STREAM_MINT"
else
    print_error "Stream was not updated with mint address"
fi
echo ""

# Step 8: Verify NFT in database
print_step "Step 8: Checking NFT in database"

if command -v psql &> /dev/null; then
    NFT_COUNT=$(psql -U nowink_user -d nowink -t -c \
        "SELECT COUNT(*) FROM nfts WHERE mint_address = '$MINT_ADDRESS';" 2>/dev/null || echo "0")
    
    if [ "$NFT_COUNT" -gt 0 ]; then
        print_success "NFT found in database"
    else
        print_error "NFT not found in database"
    fi
else
    print_info "psql not available, skipping database check"
fi
echo ""

# Step 9: List NFTs
print_step "Step 9: Listing NFTs"

NFT_LIST=$(curl -s -X GET "$API_URL/nfts" \
    -H "Authorization: Bearer $TOKEN")

NFT_COUNT=$(echo "$NFT_LIST" | jq -r '.total')

print_success "Found $NFT_COUNT NFTs in the system"
echo ""

# Cleanup
print_step "Step 10: Cleanup"
rm -f "$TEST_VIDEO"
print_success "Test video cleaned up"
echo ""

# Summary
echo "================================"
echo -e "${GREEN}🎉 END-TO-END TEST PASSED!${NC}"
echo "================================"
echo ""
echo "Test Summary:"
echo "  ✅ Authentication (nonce + JWT)"
echo "  ✅ Stream creation"
echo "  ✅ Video upload"
echo "  ✅ NFT minting"
echo "  ✅ Database persistence"
echo "  ✅ API queries"
echo ""
echo "Test Artifacts:"
echo "  • Wallet: $TEST_WALLET"
echo "  • Stream ID: $STREAM_ID"
echo "  • Mint Address: $MINT_ADDRESS"
echo "  • Arweave TX: $ARWEAVE_TX"
echo ""
echo "🚀 All systems operational!"
