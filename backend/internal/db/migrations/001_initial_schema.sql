-- now.ink Initial Database Schema
-- PostgreSQL with PostGIS extension

-- Enable PostGIS extension for geospatial queries
CREATE EXTENSION IF NOT EXISTS postgis;

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_address VARCHAR(44) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE,
    bio TEXT,
    avatar_url TEXT,
    is_premium BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Streams table (live or recorded)
CREATE TABLE IF NOT EXISTS streams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    title TEXT,
    is_live BOOLEAN DEFAULT TRUE,
    is_public BOOLEAN DEFAULT TRUE,
    started_at TIMESTAMP DEFAULT NOW(),
    ended_at TIMESTAMP,
    location GEOGRAPHY(POINT, 4326), -- PostGIS geospatial column
    nft_mint_address VARCHAR(44),
    arweave_tx_id VARCHAR(43),
    viewer_count INT DEFAULT 0,
    duration_seconds INT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- NFTs table (minted streams)
CREATE TABLE IF NOT EXISTS nfts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stream_id UUID REFERENCES streams(id) ON DELETE CASCADE,
    mint_address VARCHAR(44) UNIQUE NOT NULL,
    metadata_uri TEXT NOT NULL,
    creator_wallet VARCHAR(44) NOT NULL,
    title TEXT,
    latitude FLOAT,
    longitude FLOAT,
    timestamp TIMESTAMP NOT NULL,
    duration_seconds INT,
    thumbnail_url TEXT,
    video_url TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Follows table (social graph)
CREATE TABLE IF NOT EXISTS follows (
    follower_id UUID REFERENCES users(id) ON DELETE CASCADE,
    following_id UUID REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (follower_id, following_id),
    CONSTRAINT no_self_follow CHECK (follower_id != following_id)
);

-- Sessions table (for JWT nonces)
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_address VARCHAR(44) NOT NULL,
    nonce VARCHAR(64) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_users_wallet ON users(wallet_address);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

CREATE INDEX IF NOT EXISTS idx_streams_user ON streams(user_id);
CREATE INDEX IF NOT EXISTS idx_streams_is_live ON streams(is_live) WHERE is_live = TRUE;
CREATE INDEX IF NOT EXISTS idx_streams_location ON streams USING GIST(location);
CREATE INDEX IF NOT EXISTS idx_streams_started_at ON streams(started_at DESC);

CREATE INDEX IF NOT EXISTS idx_nfts_mint ON nfts(mint_address);
CREATE INDEX IF NOT EXISTS idx_nfts_creator ON nfts(creator_wallet);
CREATE INDEX IF NOT EXISTS idx_nfts_timestamp ON nfts(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_nfts_location ON nfts(latitude, longitude);

CREATE INDEX IF NOT EXISTS idx_follows_follower ON follows(follower_id);
CREATE INDEX IF NOT EXISTS idx_follows_following ON follows(following_id);

CREATE INDEX IF NOT EXISTS idx_sessions_nonce ON sessions(nonce);
CREATE INDEX IF NOT EXISTS idx_sessions_wallet ON sessions(wallet_address);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger for users table
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Function to clean expired sessions
CREATE OR REPLACE FUNCTION clean_expired_sessions()
RETURNS void AS $$
BEGIN
    DELETE FROM sessions WHERE expires_at < NOW();
END;
$$ LANGUAGE plpgsql;

-- Comments for documentation
COMMENT ON TABLE users IS 'User accounts identified by Solana wallet address';
COMMENT ON TABLE streams IS 'Live streams and recordings with geospatial data';
COMMENT ON TABLE nfts IS 'Minted NFTs on Solana blockchain';
COMMENT ON TABLE follows IS 'Social graph - chronological feed only, no algorithm';
COMMENT ON TABLE sessions IS 'Authentication nonces for wallet signature verification';

COMMENT ON COLUMN streams.location IS 'PostGIS GEOGRAPHY point for geospatial queries';
COMMENT ON COLUMN streams.is_public IS 'Public streams appear on map, private only visible to user';
COMMENT ON COLUMN users.is_premium IS 'Premium users can view all NFTs globally';
