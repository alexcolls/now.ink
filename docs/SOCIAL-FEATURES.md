# 👥 Social Features API Documentation

Complete API documentation for social features in now.ink.

---

## 📋 Overview

The social features API enables users to:
- Follow and unfollow other users
- View chronological feeds from followed users
- Search for users by username or wallet
- View user profiles with statistics
- Check follow relationships

All social endpoints require authentication (JWT token) unless specified otherwise.

---

## 🔑 Authentication

Social endpoints use JWT authentication. Include the token in the Authorization header:

```bash
Authorization: Bearer <your_jwt_token>
```

Get a JWT token by signing in with your wallet through the `/auth` endpoints.

---

## 📍 API Endpoints

### Follow User

**POST** `/api/v1/social/follow/:user_id`

Follow another user.

**Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Parameters:**
- `user_id` (path, required): UUID of the user to follow

**Response:**
```json
{
  "message": "followed successfully",
  "following": true
}
```

**Status Codes:**
- `200` - Success
- `400` - Invalid user_id or trying to follow yourself
- `401` - Unauthorized
- `500` - Server error

**Example:**
```bash
curl -X POST https://api.now.ink/api/v1/social/follow/123e4567-e89b-12d3-a456-426614174000 \
  -H "Authorization: Bearer <token>"
```

---

### Unfollow User

**DELETE** `/api/v1/social/follow/:user_id`

Unfollow a user.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Parameters:**
- `user_id` (path, required): UUID of the user to unfollow

**Response:**
```json
{
  "message": "unfollowed successfully",
  "following": false
}
```

**Status Codes:**
- `200` - Success
- `400` - Invalid user_id
- `401` - Unauthorized
- `500` - Server error or not following user

---

### Check Following Status

**GET** `/api/v1/social/following/:user_id/check`

Check if current user follows target user.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Parameters:**
- `user_id` (path, required): UUID of the user to check

**Response:**
```json
{
  "following": true
}
```

**Example:**
```bash
curl https://api.now.ink/api/v1/social/following/123e4567-e89b-12d3-a456-426614174000/check \
  -H "Authorization: Bearer <token>"
```

---

### Get Feed

**GET** `/api/v1/social/feed`

Get chronological feed of NFTs from users you follow.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Query Parameters:**
- `limit` (optional, default: 20, max: 50): Number of items to return
- `offset` (optional, default: 0): Pagination offset

**Response:**
```json
{
  "feed": [
    {
      "mint_address": "B1aKv...xyz",
      "metadata_uri": "ar://abc123...",
      "name": "Sunset at the beach",
      "creator": "6jK9z...abc",
      "creator_username": "alice_now",
      "creator_avatar": "ar://avatar123...",
      "latitude": 40.7128,
      "longitude": -74.0060,
      "timestamp": "2025-11-05T14:30:00Z",
      "duration_seconds": 45,
      "video_url": "ar://video456...",
      "thumbnail_url": "ar://thumb789...",
      "views": 156,
      "created_at": "2025-11-05T14:31:00Z"
    }
  ],
  "count": 1
}
```

**Example:**
```bash
curl "https://api.now.ink/api/v1/social/feed?limit=20&offset=0" \
  -H "Authorization: Bearer <token>"
```

**Notes:**
- Feed is sorted chronologically (newest first)
- Only shows NFTs from users you follow
- Empty array returned if not following anyone

---

### Search Users

**GET** `/api/v1/users/search`

Search for users by username or wallet address.

**Query Parameters:**
- `q` (required): Search query string
- `limit` (optional, default: 20, max: 50): Number of results
- `offset` (optional, default: 0): Pagination offset

**Response:**
```json
{
  "users": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "wallet_address": "6jK9z...abc",
      "username": "alice_now",
      "bio": "Capturing moments around the world 🌍",
      "avatar_url": "ar://avatar123...",
      "is_premium": false,
      "created_at": "2025-10-01T10:00:00Z",
      "updated_at": "2025-11-05T14:00:00Z"
    }
  ],
  "count": 1,
  "query": "alice"
}
```

**Example:**
```bash
curl "https://api.now.ink/api/v1/users/search?q=alice&limit=10" \
  -H "Authorization: Bearer <token>"
```

**Notes:**
- Search is case-insensitive
- Matches both username and wallet address
- Results are ranked: exact match > starts with > contains

---

### Get User Profile

**GET** `/api/v1/users/:user_id`

Get detailed user profile with statistics.

**Parameters:**
- `user_id` (path, required): UUID of the user

**Response:**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "wallet_address": "6jK9z...abc",
  "username": "alice_now",
  "bio": "Capturing moments around the world 🌍",
  "avatar_url": "ar://avatar123...",
  "is_premium": false,
  "created_at": "2025-10-01T10:00:00Z",
  "updated_at": "2025-11-05T14:00:00Z",
  "followers_count": 245,
  "following_count": 78,
  "nft_count": 156
}
```

**Example:**
```bash
curl https://api.now.ink/api/v1/users/123e4567-e89b-12d3-a456-426614174000
```

**Status Codes:**
- `200` - Success
- `400` - Invalid user_id format
- `404` - User not found

---

### Get Followers

**GET** `/api/v1/users/:user_id/followers`

Get list of users following this user.

**Parameters:**
- `user_id` (path, required): UUID of the user

**Query Parameters:**
- `limit` (optional, default: 50, max: 100)
- `offset` (optional, default: 0)

**Response:**
```json
{
  "followers": [
    {
      "id": "234e5678-e89b-12d3-a456-426614174001",
      "wallet_address": "7kL0a...def",
      "username": "bob_now",
      "bio": "Street photographer",
      "avatar_url": "ar://avatar456...",
      "is_premium": true,
      "created_at": "2025-09-15T08:00:00Z",
      "updated_at": "2025-11-04T12:00:00Z"
    }
  ],
  "count": 1
}
```

**Example:**
```bash
curl "https://api.now.ink/api/v1/users/123e4567-e89b-12d3-a456-426614174000/followers?limit=20"
```

---

### Get Following

**GET** `/api/v1/users/:user_id/following`

Get list of users this user follows.

**Parameters:**
- `user_id` (path, required): UUID of the user

**Query Parameters:**
- `limit` (optional, default: 50, max: 100)
- `offset` (optional, default: 0)

**Response:**
```json
{
  "following": [
    {
      "id": "234e5678-e89b-12d3-a456-426614174001",
      "wallet_address": "7kL0a...def",
      "username": "bob_now",
      "bio": "Street photographer",
      "avatar_url": "ar://avatar456...",
      "is_premium": true,
      "created_at": "2025-09-15T08:00:00Z",
      "updated_at": "2025-11-04T12:00:00Z"
    }
  ],
  "count": 1
}
```

**Example:**
```bash
curl "https://api.now.ink/api/v1/users/123e4567-e89b-12d3-a456-426614174000/following?limit=20"
```

---

## 💻 Client Integration Examples

### JavaScript/TypeScript

```typescript
// Follow a user
async function followUser(userId: string, token: string) {
  const response = await fetch(
    `https://api.now.ink/api/v1/social/follow/${userId}`,
    {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    }
  );
  return response.json();
}

// Get feed
async function getFeed(token: string, limit = 20, offset = 0) {
  const response = await fetch(
    `https://api.now.ink/api/v1/social/feed?limit=${limit}&offset=${offset}`,
    {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    }
  );
  return response.json();
}

// Search users
async function searchUsers(query: string, limit = 10) {
  const response = await fetch(
    `https://api.now.ink/api/v1/users/search?q=${encodeURIComponent(query)}&limit=${limit}`
  );
  return response.json();
}

// Get user profile
async function getUserProfile(userId: string) {
  const response = await fetch(
    `https://api.now.ink/api/v1/users/${userId}`
  );
  return response.json();
}
```

### React Native

```typescript
import AsyncStorage from '@react-native-async-storage/async-storage';

async function followUser(userId: string) {
  const token = await AsyncStorage.getItem('jwt_token');
  
  const response = await fetch(
    `https://api.now.ink/api/v1/social/follow/${userId}`,
    {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    }
  );
  
  if (!response.ok) {
    throw new Error('Failed to follow user');
  }
  
  return response.json();
}

async function loadFeed(limit = 20, offset = 0) {
  const token = await AsyncStorage.getItem('jwt_token');
  
  const response = await fetch(
    `https://api.now.ink/api/v1/social/feed?limit=${limit}&offset=${offset}`,
    {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    }
  );
  
  const data = await response.json();
  return data.feed;
}
```

---

## 🗄️ Database Schema

### follows table

```sql
CREATE TABLE IF NOT EXISTS follows (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    follower_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(follower_id, following_id),
    CHECK (follower_id != following_id)
);

CREATE INDEX idx_follows_follower ON follows(follower_id);
CREATE INDEX idx_follows_following ON follows(following_id);
CREATE INDEX idx_follows_created ON follows(created_at DESC);
```

---

## 🔍 Feed Algorithm

The feed follows these principles:

1. **Chronological Only**: No algorithmic ranking or recommendations
2. **Following-Based**: Only shows content from users you follow
3. **Newest First**: Ordered by NFT creation time (DESC)
4. **Paginated**: Default 20 items, max 50 per page
5. **Includes Creator Info**: Username and avatar for each NFT

**SQL Query (simplified):**
```sql
SELECT n.*, u.username, u.avatar_url
FROM nfts n
INNER JOIN users u ON n.creator_wallet = u.wallet_address
INNER JOIN follows f ON u.id = f.following_id
WHERE f.follower_id = $1
ORDER BY n.created_at DESC
LIMIT $2 OFFSET $3
```

---

## 🔐 Security Notes

1. **Self-Follow Prevention**: Users cannot follow themselves
2. **Duplicate Follow Protection**: `UNIQUE` constraint on `(follower_id, following_id)`
3. **Cascade Deletion**: Follows are deleted when users are deleted
4. **Auth Required**: Follow/unfollow operations require JWT authentication
5. **Rate Limiting**: Consider implementing rate limits on follow/unfollow actions

---

## 📈 Performance Considerations

### Indexes
- `idx_follows_follower`: Fast lookups for "who am I following?"
- `idx_follows_following`: Fast lookups for "who follows me?"
- `idx_follows_created`: Chronological feed ordering

### Caching Strategy
- Cache user profile stats (followers/following counts)
- Cache feed results for 30-60 seconds
- Invalidate on follow/unfollow actions

### Query Optimization
- Use `LIMIT` on all list endpoints
- Implement cursor-based pagination for large feeds
- Add database query timeouts

---

## 🚀 Future Enhancements

### Planned Features
- **Notifications**: Push notifications for new followers
- **Mutual Follows**: Badge/indicator for mutual follows
- **Block/Mute**: Block or mute users
- **Private Profiles**: Option to approve followers
- **Follow Requests**: For private accounts
- **Trending Feed**: Popular content discovery (optional)
- **Tags/Categories**: Filter feed by content type

---

## 📚 Related Documentation

- [API Overview](../README.md#api-endpoints)
- [Authentication](./AUTH.md)
- [Database Schema](./DATABASE-SCHEMA.md)
- [Mobile App Integration](./MOBILE-COMPLETE.md)

---

**Created**: November 5, 2025  
**Version**: 1.0.0  
**Status**: Production Ready
