# Contributing to now.ink

**Thank you for considering contributing to now.ink!** We're building the anti-feed, the quiet rebellion against algorithmic manipulation. Every line of code you write helps keep the truth open-source.

---

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing Requirements](#testing-requirements)
- [Documentation](#documentation)
- [Commit Message Format](#commit-message-format)
- [Pull Request Process](#pull-request-process)

---

## Code of Conduct

### Our Values

1. **Anti-Algorithm:** No recommendation engines, no shadow-banning, no manipulated feeds
2. **Anti-AI Generated Content:** Only authentic, human-captured moments
3. **Anti-Corporate:** No ads, no data mining, transparent monetization
4. **Pro-Truth:** Open-source by default, auditable code, honest communication
5. **Pro-User:** Users own their content, their data, their experience

### Expected Behavior

- Be respectful and constructive in code reviews
- Focus on technical merit, not personal preferences
- Assume good intent; disagree without being disagreeable
- Prioritize user privacy and data ownership
- Document decisions that impact the user experience

### Unacceptable Behavior

- Implementing algorithmic content ranking without consensus
- Adding tracking/analytics that violate user privacy
- Integrating AI content generation features
- Copying code without proper attribution
- Harassment or discrimination of any kind

---

## Getting Started

### Prerequisites

- **Node.js** 20+ (for Nuxt 4 & React Native)
- **Go** 1.21+ (for backend API)
- **Solana CLI** (for blockchain integration)
- **PostgreSQL** 15+ or **MongoDB** 6+ (database)
- **Git** for version control

See **[SETUP.md](SETUP.md)** for detailed installation instructions.

### Fork & Clone

```bash
# Fork the repository on GitHub
# Then clone your fork
git clone https://github.com/YOUR_USERNAME/now.ink.git
cd now.ink

# Add upstream remote
git remote add upstream https://github.com/originalrepo/now.ink.git
```

### Create a Branch

```bash
# Fetch latest changes
git fetch upstream
git checkout main
git merge upstream/main

# Create feature branch
git checkout -b feature/your-feature-name
# or
git checkout -b fix/bug-description
```

---

## Development Workflow

### 1. Pick an Issue

- Check [GitHub Issues](https://github.com/now.ink/issues)
- Look for `good first issue` or `help wanted` labels
- Comment on the issue to claim it

### 2. Write Code

- Follow [Coding Standards](#coding-standards)
- Write tests for new features
- Update documentation if needed

### 3. Test Locally

```bash
# Run tests for each component
cd web && npm test
cd mobile && npm test
cd backend && go test ./...
```

### 4. Commit Changes

```bash
# Stage changes
git add .

# Commit with emoji (see format below)
git commit -m "✨ Add time-slider component to map view"
```

### 5. Push & Create PR

```bash
# Push to your fork
git push origin feature/your-feature-name

# Create Pull Request on GitHub
# Fill out the PR template
```

---

## Coding Standards

### TypeScript (Web & Mobile)

**Style Guide:** Follow [Airbnb TypeScript Style Guide](https://github.com/airbnb/javascript)

**Key Rules:**
- Use TypeScript strict mode (`strict: true`)
- Prefer `const` over `let`, never use `var`
- Use arrow functions for callbacks
- Destructure props and imports
- Use absolute imports (configured in tsconfig)

**Example:**

```typescript
// Good ✅
import { useWallet } from '@/composables/useWallet';

export const WalletConnect = () => {
  const { connect, disconnect, address } = useWallet();
  
  const handleConnect = async () => {
    try {
      await connect();
    } catch (error) {
      console.error('Connection failed:', error);
    }
  };
  
  return (
    <button onClick={handleConnect}>
      {address ? `Connected: ${address.slice(0, 6)}...` : 'Connect Wallet'}
    </button>
  );
};

// Bad ❌
var WalletConnect = function() {
  let wallet = useWallet();
  // ... imperative spaghetti code
}
```

### Go (Backend)

**Style Guide:** Follow [Effective Go](https://go.dev/doc/effective_go) and [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

**Key Rules:**
- Use `gofmt` to format code
- Run `golangci-lint` before committing
- Keep functions small and focused (<50 lines)
- Use meaningful variable names (no single-letter vars except loops)
- Always handle errors explicitly

**Example:**

```go
// Good ✅
func (s *StreamService) StartStream(ctx context.Context, userID uuid.UUID, req *StartStreamRequest) (*Stream, error) {
    if err := req.Validate(); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }
    
    stream := &Stream{
        ID:        uuid.New(),
        UserID:    userID,
        Title:     req.Title,
        Latitude:  req.Latitude,
        Longitude: req.Longitude,
        IsLive:    true,
        StartedAt: time.Now(),
    }
    
    if err := s.repo.Create(ctx, stream); err != nil {
        return nil, fmt.Errorf("failed to create stream: %w", err)
    }
    
    return stream, nil
}

// Bad ❌
func StartStream(u string, t string, lat float64, lon float64) {
    // ... no error handling, unclear params
}
```

### Naming Conventions

- **Variables:** `camelCase` (TS) or `camelCase` (Go)
- **Functions:** `camelCase` (TS) or `PascalCase` (Go exported)
- **Components:** `PascalCase` (TS React/Vue)
- **Files:** `kebab-case.ts` (TS) or `snake_case.go` (Go)
- **Constants:** `SCREAMING_SNAKE_CASE`

---

## Testing Requirements

### Unit Tests

**Required for:**
- All new features
- Bug fixes
- Utility functions
- API handlers

**Coverage Target:** 80%+ for critical paths

**Example (Go):**

```go
func TestStreamService_StartStream(t *testing.T) {
    ctx := context.Background()
    mockRepo := &MockStreamRepository{}
    service := NewStreamService(mockRepo)
    
    req := &StartStreamRequest{
        Title:     "Test Stream",
        Latitude:  40.7128,
        Longitude: -74.0060,
    }
    
    stream, err := service.StartStream(ctx, uuid.New(), req)
    
    assert.NoError(t, err)
    assert.NotNil(t, stream)
    assert.Equal(t, "Test Stream", stream.Title)
    assert.True(t, stream.IsLive)
}
```

### Integration Tests

**Required for:**
- API endpoints
- Database interactions
- Blockchain operations (use devnet)

### E2E Tests

**Recommended for:**
- Critical user flows (signup, minting, playback)
- Use Playwright (web) or Detox (mobile)

---

## Documentation

### When to Update Docs

- Adding a new API endpoint → update `docs/API.md`
- Changing architecture → update `docs/ARCHITECTURE.md`
- New environment variable → update `.env.sample` and `docs/SETUP.md`
- User-facing feature → update `docs/USER_GUIDE.md`

### Documentation Style

- Use Markdown
- Include code examples
- Keep it concise (under 100 lines per section)
- Use diagrams (ASCII art or Mermaid) where helpful

---

## Commit Message Format

We use **emoji-prefixed commit messages** grouped by feature.

### Format

```
<emoji> <type>: <short description>

<optional body>

<optional footer>
```

### Emoji Guide

| Emoji | Type | Description |
|-------|------|-------------|
| ✨ | feat | New feature |
| 🐛 | fix | Bug fix |
| 📝 | docs | Documentation only |
| 💄 | style | UI/styling changes |
| ♻️ | refactor | Code refactor |
| ⚡ | perf | Performance improvement |
| ✅ | test | Adding/updating tests |
| 🔧 | chore | Build/config changes |
| 🚀 | deploy | Deployment-related |
| 🔒 | security | Security fix |

### Examples

```bash
✨ feat: add time-slider component to map view

Implements the time-range filter allowing users to filter
NFTs by date. Uses a dual-handle slider for start/end dates.

Closes #42

---

🐛 fix: prevent duplicate NFT mints

Check if stream already has an NFT before minting.
Adds unique constraint on streams.nft_mint_address.

---

📝 docs: update API.md with new geo endpoints

Added documentation for /geo/nearby and /geo/bounds
```

### Commit by Feature

Group related changes into a single commit:

```bash
# Good ✅
git add src/components/TimeSlider.tsx src/composables/useTimeFilter.ts
git commit -m "✨ feat: implement time-slider for map filtering"

# Bad ❌
git commit -m "update file"
git commit -m "fix typo"
git commit -m "add another thing"
```

---

## Pull Request Process

### Before Submitting

1. **Sync with upstream:**
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Run linters:**
   ```bash
   npm run lint        # TypeScript
   golangci-lint run   # Go
   ```

3. **Run tests:**
   ```bash
   npm test && go test ./...
   ```

4. **Update CHANGELOG** (if applicable)

### PR Template

```markdown
## Description
Brief summary of changes.

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Tested locally

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-reviewed code
- [ ] Commented hard-to-understand areas
- [ ] Updated documentation
- [ ] No new warnings
- [ ] Added tests

## Related Issues
Closes #123
```

### Review Process

1. **Automated Checks:** CI/CD runs tests and linters
2. **Code Review:** At least one maintainer reviews
3. **Approval:** Maintainer approves or requests changes
4. **Merge:** Squash and merge (keep history clean)

### After Merge

- **Delete your branch:**
  ```bash
  git branch -d feature/your-feature-name
  git push origin --delete feature/your-feature-name
  ```

- **Celebrate!** 🎉 You've contributed to the anti-algorithm movement.

---

## Questions?

- **Discord:** [TBD]
- **GitHub Discussions:** [TBD]
- **Email:** contribute@now.ink

---

**Remember:** We're not just building an app. We're building a movement. Every commit is a vote for transparency, authenticity, and user sovereignty. Thank you for being part of it.

---

**Build it. Ship it. Watch people finally look up.**
