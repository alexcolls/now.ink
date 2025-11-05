# Contributing to now.ink

Thank you for your interest in contributing to now.ink! We welcome contributions from the community.

---

## 🌟 Philosophy

Before contributing, understand our core principles:

- **Anti-algorithm** - No recommendation engines, chronological only
- **Anti-AI** - Only real, live-captured content
- **User-owned** - NFTs on blockchain, permanent on Arweave
- **Open source** - Transparent, auditable, community-driven
- **No ads** - Revenue from platform commission only

If your contribution aligns with these values, we'd love to have you!

---

## 🚀 How to Contribute

### Types of Contributions We Accept

✅ **Bug fixes** - Found a bug? Fix it!  
✅ **Performance improvements** - Make it faster  
✅ **Documentation** - Clarify or expand docs  
✅ **Tests** - Improve test coverage  
✅ **UI/UX improvements** - Better user experience  
✅ **New features** - Discuss first in an issue  

❌ **Algorithm-based features** - Goes against our principles  
❌ **AI-generated content** - Not aligned with our mission  
❌ **Ad integrations** - We don't do ads  

---

## 📋 Before You Start

1. **Check existing issues** - Someone might already be working on it
2. **Create an issue** - Discuss new features before coding
3. **Read the docs** - Understand the architecture
4. **Fork the repo** - Work in your own fork
5. **Follow conventions** - Match existing code style

---

## 🔧 Development Setup

### Prerequisites
```bash
# Required
- Node.js 20+
- Go 1.22+
- PostgreSQL 16 with PostGIS
- Docker (optional)

# Optional but recommended
- Solana CLI (for blockchain testing)
- React Native dev tools
```

### Quick Setup
```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/now.ink.git
cd now.ink

# Backend
cd backend
go mod download
cp .env.sample .env
go run cmd/api/main.go

# Mobile (separate terminal)
cd mobile
npm install
npm start

# Run tests
./test-e2e.sh
```

See [MVP-COMPLETE.md](MVP-COMPLETE.md) for detailed setup.

---

## 📝 Contribution Workflow

### 1. Create a Branch
```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/bug-description
```

### 2. Make Your Changes
- Write clean, readable code
- Follow existing patterns
- Add tests if applicable
- Update documentation

### 3. Test Your Changes
```bash
# Backend tests
cd backend
go test ./...

# E2E test
./test-e2e.sh

# Mobile (manual testing)
npm start
```

### 4. Commit Your Changes
Use conventional commits with emojis (following project style):

```bash
git commit -m "✨ Add user search functionality"
git commit -m "🐛 Fix video upload timeout issue"
git commit -m "📚 Update API documentation"
git commit -m "♻️ Refactor map component"
```

Commit types:
- `✨` `:sparkles:` - New features
- `🐛` `:bug:` - Bug fixes
- `📚` `:books:` - Documentation
- `♻️` `:recycle:` - Refactoring
- `⚡` `:zap:` - Performance
- `🧪` `:test_tube:` - Tests
- `🔒` `:lock:` - Security
- `🎨` `:art:` - UI/UX

### 5. Push and Create PR
```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub with:
- Clear title and description
- Link to related issue
- Screenshots (if UI changes)
- Test results

---

## 💻 Code Style Guidelines

### Go (Backend)
```go
// Use standard Go formatting
gofmt -w .

// Follow Go conventions
- Use camelCase for private, PascalCase for public
- Keep functions small and focused
- Add comments for exported functions
- Use meaningful variable names
```

### TypeScript (Mobile)
```typescript
// Use Prettier formatting
npm run format

// Follow conventions
- Use const/let, not var
- Prefer arrow functions
- Add JSDoc comments for complex functions
- Use TypeScript types, not any
```

### General
- **80-120 characters** per line
- **2 spaces** indentation
- **No trailing whitespace**
- **Descriptive** variable/function names
- **Comments** for complex logic

---

## 🧪 Testing Guidelines

### What to Test
- ✅ New features
- ✅ Bug fixes
- ✅ API endpoints
- ✅ Database queries
- ✅ Blockchain interactions

### How to Test
```bash
# Backend unit tests
cd backend
go test ./internal/... -v

# E2E API test
./test-e2e.sh

# Mobile (manual for now)
npm start
# Test on device/simulator
```

### Test Coverage
- Aim for 70%+ coverage on backend
- Test happy paths and error cases
- Include edge cases

---

## 📚 Documentation

Update documentation when you:
- Add new features
- Change APIs
- Modify configuration
- Fix important bugs

Files to update:
- `README.md` - For major changes
- `CHANGELOG.md` - For all changes
- `API.md` - For API changes (if exists)
- Code comments - For complex logic

---

## 🔍 Code Review Process

### What We Look For
✅ **Code quality** - Clean, readable, maintainable  
✅ **Tests** - Adequate test coverage  
✅ **Documentation** - Clear comments and docs  
✅ **Performance** - No obvious bottlenecks  
✅ **Security** - No vulnerabilities  
✅ **Alignment** - Matches project philosophy  

### Review Timeline
- Small PRs: 1-2 days
- Medium PRs: 3-5 days
- Large PRs: 1 week+

Be patient! We review carefully.

---

## 🎯 Good First Issues

New to the project? Look for issues tagged:
- `good first issue`
- `help wanted`
- `documentation`

These are great starting points!

---

## 🤝 Community Guidelines

### Be Respectful
- Treat everyone with respect
- Be constructive in feedback
- Assume good intentions
- No harassment or discrimination

### Communication
- Use GitHub Issues for bugs/features
- Use Pull Requests for code
- Keep discussions on-topic
- English language preferred

### Getting Help
- Check documentation first
- Search existing issues
- Ask in GitHub Discussions
- Be specific in questions

---

## 🚫 What Not to Do

❌ Submit PRs without an issue  
❌ Make large changes without discussion  
❌ Ignore code review feedback  
❌ Copy code without attribution  
❌ Include breaking changes without notice  
❌ Add dependencies without justification  

---

## 📜 License Agreement

By contributing, you agree that:
- Your contributions will be licensed under the project's dual license
- You have the right to contribute the code
- Your contributions are your original work

See [LICENSE](LICENSE) for details.

---

## 🏆 Recognition

Contributors will be:
- Listed in CONTRIBUTORS.md
- Mentioned in release notes
- Credited in project README
- Given public thanks on social media

Top contributors may receive:
- Free commercial license
- Beta feature access
- Exclusive project swag

---

## 📞 Contact

Questions about contributing?
- GitHub Issues: https://github.com/alexcolls/now.ink/issues
- Create a discussion post
- Tag maintainers in PR comments

---

## 🌟 Thank You!

Every contribution, no matter how small, helps make now.ink better. We appreciate your time and effort!

**Build it. Ship it. Watch people finally look up.** ✨
