# API Documentation

**now.ink Backend API Specification**

Base URL (Development): `http://localhost:8080/api/v1`  
Base URL (Production): `https://api.now.ink/v1`

---

## Table of Contents

- [Authentication](#authentication)
- [User Endpoints](#user-endpoints)
- [Stream Endpoints](#stream-endpoints)
- [NFT Endpoints](#nft-endpoints)
- [Geo Query Endpoints](#geo-query-endpoints)
- [Social Graph Endpoints](#social-graph-endpoints)
- [Premium Tier Endpoints](#premium-tier-endpoints)
- [WebSocket Protocols](#websocket-protocols)
- [Error Handling](#error-handling)

---

## Authentication

All authenticated endpoints require a JWT token in the `Authorization` header:

```
Authorization: Bearer <jwt_token>
```

### POST `/auth/nonce`

**Description:** Request a nonce for wallet signature  
**Auth Required:** No

**Request Body:**
```json
{
  "wallet_address": "7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU"
}
```

**Response:**
```json
{
  "nonce": "550e8400-e29b-41d4-a716-446655440000",
  "message": "Sign this message to authenticate with now.ink: 550e8400-e29b-41d4-a716-446655440000"
}
```

### POST `/auth/verify`

**Description:** Verify signed nonce and issue JWT  
**Auth Required:** No

**Request Body:**
```json
{
  "wallet_address": "7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU",
  "signature": "base64_encoded_signature",
  "nonce": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "uuid",
    "wallet_address": "7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU",
    "username": "satoshi",
    "is_premium": false
  }
}
```

---

## User Endpoints

### GET `/users/me`

**Description:** Get current user profile  
**Auth Required:** Yes

**Response:**
```json
{
  "id": "uuid",
  "wallet_address": "7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU",
  "username": "satoshi",
  "bio": "Living in the now",
  "avatar_url": "https://cdn.now.ink/avatars/...",
  "is_premium": false,
  "created_at": "2025-11-05T01:23:45Z"
}
```

### PUT `/users/me`

**Description:** Update user profile  
**Auth Required:** Yes

**Request Body:**
```json
{
  "username": "satoshi_new",
  "bio": "Updated bio",
  "avatar_url": "https://cdn.now.ink/avatars/new.jpg"
}
```

**Response:** Updated user object

### GET `/users/:wallet_address`

**Description:** Get public user profile  
**Auth Required:** No

**Response:** User object (public fields only)

---

## Stream Endpoints

### POST `/streams/start`

**Description:** Start a new live stream  
**Auth Required:** Yes

**Request Body:**
```json
{
  "title": "Sunset at the pier",
  "latitude": 40.7128,
  "longitude": -74.0060,
  "is_public": true
}
```

**Response:**
```json
{
  "stream_id": "uuid",
  "signaling_url": "wss://api.now.ink/ws/stream/uuid",
  "started_at": "2025-11-05T01:23:45Z"
}
```

### POST `/streams/:id/end`

**Description:** End a live stream  
**Auth Required:** Yes (must be stream owner)

**Response:**
```json
{
  "stream_id": "uuid",
  "ended_at": "2025-11-05T01:25:30Z",
  "duration_seconds": 105
}
```

### POST `/streams/:id/save`

**Description:** Save stream as NFT (triggers minting process)  
**Auth Required:** Yes (must be stream owner)

**Request Body:**
```json
{
  "video_data": "base64_encoded_video_or_upload_url"
}
```

**Response:**
```json
{
  "message": "Minting in progress",
  "mint_job_id": "uuid",
  "estimated_completion": "2025-11-05T01:26:00Z"
}
```

### GET `/streams/:id`

**Description:** Get stream details  
**Auth Required:** No (for public streams)

**Response:**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "creator": {
    "wallet_address": "...",
    "username": "satoshi"
  },
  "title": "Sunset at the pier",
  "is_live": false,
  "started_at": "2025-11-05T01:23:45Z",
  "ended_at": "2025-11-05T01:25:30Z",
  "location": {
    "latitude": 40.7128,
    "longitude": -74.0060
  },
  "nft_mint_address": "8x...",
  "arweave_tx_id": "abc123...",
  "viewer_count": 42
}
```

### GET `/streams/live`

**Description:** Get all currently live streams  
**Auth Required:** No

**Query Parameters:**
- `limit` (default: 50, max: 100)
- `offset` (default: 0)

**Response:**
```json
{
  "streams": [ /* array of stream objects */ ],
  "total": 123,
  "limit": 50,
  "offset": 0
}
```

---

## NFT Endpoints

### GET `/nfts`

**Description:** List all NFTs (with filters)  
**Auth Required:** No

**Query Parameters:**
- `latitude` (float, required with longitude)
- `longitude` (float, required with latitude)
- `radius_km` (float, default: 10)
- `start_date` (ISO 8601, e.g., `2025-01-01T00:00:00Z`)
- `end_date` (ISO 8601)
- `limit` (default: 50, max: 100)
- `offset` (default: 0)

**Response:**
```json
{
  "nfts": [
    {
      "id": "uuid",
      "mint_address": "8xKXtg...",
      "metadata_uri": "ar://abc123...",
      "creator_wallet": "7xKXtg...",
      "title": "Sunset at the pier",
      "latitude": 40.7128,
      "longitude": -74.0060,
      "created_at": "2025-11-05T01:25:30Z",
      "thumbnail_url": "ar://thumbnail123",
      "duration_seconds": 42
    }
  ],
  "total": 456,
  "limit": 50,
  "offset": 0
}
```

### GET `/nfts/:mint_address`

**Description:** Get specific NFT details  
**Auth Required:** No

**Response:**
```json
{
  "id": "uuid",
  "mint_address": "8xKXtg...",
  "metadata_uri": "ar://abc123...",
  "metadata": {
    "name": "now.ink Moment #12345",
    "description": "...",
    "attributes": [ /* ... */ ]
  },
  "creator": {
    "wallet_address": "7xKXtg...",
    "username": "satoshi"
  },
  "latitude": 40.7128,
  "longitude": -74.0060,
  "created_at": "2025-11-05T01:25:30Z",
  "video_url": "ar://video123",
  "can_play": true  // based on proximity or premium status
}
```

### GET `/nfts/:mint_address/playback`

**Description:** Get playback URL (requires proximity or premium)  
**Auth Required:** Yes (optional for free local playback)

**Response:**
```json
{
  "playback_url": "ar://video123",
  "expires_at": "2025-11-05T02:23:45Z",
  "access_reason": "proximity" | "premium" | "owner"
}
```

**Error Response (403):**
```json
{
  "error": "Access denied",
  "message": "You must be within 100m or have premium to view this NFT",
  "required_proximity_meters": 100,
  "current_distance_meters": 5420
}
```

### GET `/nfts/user/:wallet_address`

**Description:** Get all NFTs created by a user  
**Auth Required:** No

**Query Parameters:**
- `limit`, `offset` (same as `/nfts`)

**Response:** Same format as `/nfts`

---

## Geo Query Endpoints

### POST `/geo/nearby`

**Description:** Find NFTs near a location  
**Auth Required:** No

**Request Body:**
```json
{
  "latitude": 40.7128,
  "longitude": -74.0060,
  "radius_km": 5,
  "start_date": "2025-01-01T00:00:00Z",
  "end_date": "2025-12-31T23:59:59Z"
}
```

**Response:** Same as `/nfts` response format

### POST `/geo/bounds`

**Description:** Find NFTs within a bounding box (for map view)  
**Auth Required:** No

**Request Body:**
```json
{
  "north": 40.8,
  "south": 40.6,
  "east": -73.9,
  "west": -74.1,
  "start_date": "2025-01-01T00:00:00Z",
  "end_date": "2025-12-31T23:59:59Z"
}
```

**Response:** Array of NFTs with locations

---

## Social Graph Endpoints

### POST `/social/follow`

**Description:** Follow a user  
**Auth Required:** Yes

**Request Body:**
```json
{
  "wallet_address": "7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU"
}
```

**Response:**
```json
{
  "message": "Successfully followed user",
  "following": true
}
```

### DELETE `/social/follow/:wallet_address`

**Description:** Unfollow a user  
**Auth Required:** Yes

**Response:**
```json
{
  "message": "Successfully unfollowed user",
  "following": false
}
```

### GET `/social/following`

**Description:** Get list of users you follow  
**Auth Required:** Yes

**Response:**
```json
{
  "following": [
    {
      "wallet_address": "...",
      "username": "satoshi",
      "avatar_url": "...",
      "followed_at": "2025-11-05T01:23:45Z"
    }
  ],
  "total": 42
}
```

### GET `/social/followers`

**Description:** Get list of your followers  
**Auth Required:** Yes

**Response:** Same format as `/social/following`

### GET `/social/feed`

**Description:** Get chronological feed of followed users' NFTs  
**Auth Required:** Yes

**Query Parameters:**
- `limit` (default: 50, max: 100)
- `offset` (default: 0)

**Response:**
```json
{
  "feed": [ /* array of NFT objects with creator info */ ],
  "total": 789,
  "limit": 50,
  "offset": 0
}
```

---

## Premium Tier Endpoints

### GET `/premium/status`

**Description:** Check premium status  
**Auth Required:** Yes

**Response:**
```json
{
  "is_premium": true,
  "expires_at": "2026-11-05T01:23:45Z",
  "auto_renew": true
}
```

### POST `/premium/subscribe`

**Description:** Subscribe to premium (payment flow TBD)  
**Auth Required:** Yes

**Request Body:**
```json
{
  "payment_method": "solana",
  "transaction_signature": "base64_signature"
}
```

**Response:**
```json
{
  "message": "Premium activated",
  "expires_at": "2026-11-05T01:23:45Z"
}
```

---

## WebSocket Protocols

### WebRTC Signaling: `wss://api.now.ink/ws/stream/:stream_id`

**Auth:** Include JWT in query param: `?token=<jwt>`

**Message Types:**

#### Client → Server

**1. Join as viewer:**
```json
{
  "type": "join",
  "role": "viewer"
}
```

**2. WebRTC Offer (broadcaster):**
```json
{
  "type": "offer",
  "sdp": "v=0\no=- ..."
}
```

**3. WebRTC Answer (viewer):**
```json
{
  "type": "answer",
  "sdp": "v=0\no=- ..."
}
```

**4. ICE Candidate:**
```json
{
  "type": "ice_candidate",
  "candidate": {
    "candidate": "...",
    "sdpMid": "...",
    "sdpMLineIndex": 0
  }
}
```

#### Server → Client

**1. Viewer joined:**
```json
{
  "type": "viewer_joined",
  "viewer_count": 5
}
```

**2. Forward offer/answer/ICE:**
```json
{
  "type": "offer" | "answer" | "ice_candidate",
  "from": "peer_id",
  ...
}
```

**3. Stream ended:**
```json
{
  "type": "stream_ended",
  "reason": "broadcaster_ended"
}
```

---

## Error Handling

All errors follow this format:

```json
{
  "error": "error_code",
  "message": "Human-readable error message",
  "details": { /* optional additional context */ }
}
```

### Common Error Codes

| Status | Code | Description |
|--------|------|-------------|
| 400 | `bad_request` | Invalid request body or parameters |
| 401 | `unauthorized` | Missing or invalid JWT token |
| 403 | `forbidden` | Insufficient permissions (e.g., not premium) |
| 404 | `not_found` | Resource doesn't exist |
| 409 | `conflict` | Resource already exists (e.g., username taken) |
| 429 | `rate_limited` | Too many requests |
| 500 | `internal_error` | Server error |

### Rate Limits

- **Free users:** 10 requests/second, 100 mints/day
- **Premium users:** 50 requests/second, unlimited mints
- **Anonymous:** 5 requests/second

Headers:
```
X-RateLimit-Limit: 10
X-RateLimit-Remaining: 7
X-RateLimit-Reset: 1699143600
```

---

## Versioning

The API uses URL versioning (`/api/v1`). Breaking changes will increment the version number. Non-breaking changes (new optional fields, new endpoints) are added to the current version.

---

## Pagination

All list endpoints support pagination:

**Query Parameters:**
- `limit` (default varies, max usually 100)
- `offset` (default 0)

**Response includes:**
```json
{
  "data": [ /* ... */ ],
  "total": 1234,
  "limit": 50,
  "offset": 0
}
```

---

## Timestamps

All timestamps are in ISO 8601 format with UTC timezone:  
`2025-11-05T01:23:45Z`

---

**For implementation details, see [ARCHITECTURE.md](ARCHITECTURE.md)**
