#!/bin/bash
# Production Database Setup for now.ink
# Run this on your production server after initial deployment

set -e

echo "🗄️  now.ink Production Database Setup"
echo "======================================"
echo ""

# Check if running as root or with sudo
if [ "$EUID" -ne 0 ]; then 
   echo "❌ Please run as root or with sudo"
   exit 1
fi

# Prompt for database password
read -sp "Enter database password for nowink_user: " DB_PASSWORD
echo ""

if [ -z "$DB_PASSWORD" ]; then
    echo "❌ Password cannot be empty"
    exit 1
fi

echo "📦 Installing PostgreSQL 16 with PostGIS..."
apt update
apt install -y postgresql-16 postgresql-16-postgis-3

echo "✅ PostgreSQL installed"
echo ""

echo "🔧 Configuring PostgreSQL..."

# Configure PostgreSQL to listen on all interfaces (for Docker)
sed -i "s/#listen_addresses = 'localhost'/listen_addresses = '*'/g" /etc/postgresql/16/main/postgresql.conf

# Add pg_hba.conf entry for Docker network
echo "host    all             all             172.16.0.0/12           md5" >> /etc/postgresql/16/main/pg_hba.conf

# Restart PostgreSQL
systemctl restart postgresql

echo "✅ PostgreSQL configured"
echo ""

echo "👤 Creating database user and database..."

# Create user and database
sudo -u postgres psql << EOF
-- Create user
CREATE USER nowink_user WITH PASSWORD '$DB_PASSWORD';

-- Create database
CREATE DATABASE nowink OWNER nowink_user;

-- Connect to database
\c nowink

-- Enable PostGIS extension
CREATE EXTENSION IF NOT EXISTS postgis;

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE nowink TO nowink_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO nowink_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO nowink_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO nowink_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO nowink_user;

EOF

echo "✅ Database user and database created"
echo ""

echo "📊 Running migrations..."

# Run migrations
sudo -u postgres psql -U nowink_user -d nowink -f /home/nowink/now.ink/backend/internal/db/migrations/001_initial_schema.sql

echo "✅ Migrations complete"
echo ""

echo "🔒 Securing database..."

# Set proper file permissions
chmod 640 /etc/postgresql/16/main/postgresql.conf
chmod 640 /etc/postgresql/16/main/pg_hba.conf

# Enable and start PostgreSQL
systemctl enable postgresql
systemctl restart postgresql

echo "✅ Database secured"
echo ""

echo "🧪 Testing connection..."

# Test connection
if sudo -u postgres psql -U nowink_user -d nowink -c "SELECT PostGIS_Version();" > /dev/null 2>&1; then
    echo "✅ Database connection successful!"
else
    echo "❌ Database connection failed"
    exit 1
fi

echo ""
echo "================================"
echo "✅ Database Setup Complete!"
echo "================================"
echo ""
echo "Database Details:"
echo "  Host: localhost (or postgres in Docker)"
echo "  Port: 5432"
echo "  Database: nowink"
echo "  User: nowink_user"
echo "  Password: [saved]"
echo ""
echo "Next steps:"
echo "1. Update .env.production with DB_PASSWORD"
echo "2. Start Docker services: docker-compose up -d"
echo "3. Verify: docker-compose logs -f api"
echo ""
