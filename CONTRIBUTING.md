# Contributing to SkillFlow

Thank you for your interest in contributing to SkillFlow! This document provides guidelines for contributing to the project.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment.

## How to Contribute

### Reporting Bugs

1. Check if the bug has already been reported in Issues
2. If not, create a new issue with:
   - Clear title and description
   - Steps to reproduce
   - Expected vs actual behavior
   - Environment details (OS, Go version, etc.)
   - Screenshots if applicable

### Suggesting Features

1. Check existing feature requests in Issues
2. Create a new issue with:
   - Clear description of the feature
   - Use cases and benefits
   - Possible implementation approach

### Pull Requests

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (`make test`)
6. Run linters (`make lint`)
7. Commit with descriptive messages
8. Push to your fork
9. Open a Pull Request

### Commit Message Guidelines

Follow conventional commits:

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Example:
```
feat(api): add user search endpoint

Implement full-text search for users using Elasticsearch.
Includes pagination and filtering options.

Closes #123
```

## Development Setup

1. Clone the repository:
```bash
git clone https://github.com/vern/skillflow.git
cd skillflow
```

2. Install dependencies:
```bash
go mod download
```

3. Start development environment:
```bash
make dev-up
```

4. Run tests:
```bash
make test
```

## Code Style

- Follow Go best practices and idioms
- Use `gofmt` for formatting
- Follow the project's existing code structure
- Write clear, self-documenting code
- Add comments for complex logic

## Testing

- Write unit tests for new code
- Maintain test coverage above 80%
- Include integration tests where appropriate
- Test edge cases and error conditions

## Documentation

- Update README.md if needed
- Add/update API documentation for new endpoints
- Include code comments for exported functions
- Update deployment docs for infrastructure changes

## Review Process

1. Maintainers will review your PR
2. Address feedback and update your PR
3. Once approved, maintainers will merge

## Questions?

Feel free to open an issue for questions or join our community discussions.

Thank you for contributing to SkillFlow!
