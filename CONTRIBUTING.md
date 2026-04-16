# Contributing to setupx

First off, thank you for considering contributing to `setupx`! It's people like you that make `setupx` such a great tool.

## How Can I Contribute?

### Reporting Bugs
* Check the [issue tracker](https://github.com/sumant1122/setupx/issues) to see if the bug has already been reported.
* If you find a new bug, please open a new issue using the **Bug Report** template.

### Suggesting Enhancements
* Open a new issue using the **Feature Request** template.
* Explain why this enhancement would be useful to most users.

### Pull Requests
1. Fork the repository.
2. Create a new branch for your feature or fix (`git checkout -b feature/amazing-feature`).
3. Commit your changes with clear, descriptive messages.
4. Push to the branch (`git push origin feature/amazing-feature`).
5. Open a Pull Request against the `main` branch.

## Development Setup

### Prerequisites
- Go 1.25 or higher

### Running Tests
Always ensure the tests pass before submitting a PR:
```bash
go test ./...
```

### Building
```bash
go build -o setupx main.go
```

## Coding Standards
- Follow standard Go formatting (`go fmt`).
- Ensure all new features are accompanied by unit tests.
- Keep the CLI output clean and user-friendly.

## License
By contributing to `setupx`, you agree that your contributions will be licensed under the MIT License.
