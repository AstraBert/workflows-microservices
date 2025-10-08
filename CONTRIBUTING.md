# Contributing to `workflows-microservices`

Thank you for your interest in contributing to this project! Please review these guidelines before getting started.

## Issue Reporting

### When to Report an Issue

- You've discovered bugs but lack the knowledge or time to fix them
- You have feature requests but cannot implement them yourself

> ⚠️ **Important:** Always search existing open and closed issues before submitting to avoid duplicates.

### How to Report an Issue

1. Open a new issue
2. Provide a clear, concise title that describes the problem or feature request
3. Include a detailed description of the issue or requested feature

## Code Contributions

### When to Contribute

- You've identified and fixed bugs
- You've optimized or improved existing code
- You've developed new features that would benefit the community

### How to Contribute

1. **Fork the repository and checkout a feature branch**

2. **Make your changes and ensure everything is working properly**

   For `main`, ensure the program builds:

   ```bash
   go build .
   ```

   And that all the templates have been correctly generated:

   ```bash
   templ generate
   ```

   For python microservices, ensure the package can be installed and the script run:

   ```bash
   uv pip install .
   process-{service-name} # replace service-name with the name of the service
   ```

   If you modify `schema.sql` or `query.sql` in any service, please run:

   ```bash
   sqlc generate
   ```

   And make any needed changes in the code.

3. **Commit your changes**


4. **Submit a pull request**
   Include a comprehensive description of your changes.

---

**Thank you for contributing!**
