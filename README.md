# <img alt=logo src=rot.png width=40px> Rot

> Future proof secrets management

[![Integration](https://github.com/candiddev/rot/actions/workflows/integration.yaml/badge.svg?branch=main)](https://github.com/candiddev/rot/actions/workflows/integration.yaml)

[:book: Docs](https://rotx.dev/docs/)\
[:motorway: Roadmap](https://github.com/orgs/candiddev/projects/6/views/31)

Rot is an open source command line (CLI) tool for managing secrets.

Rot makes encrypting and decrypting secrets easy:

- Generate keys and values using current best encryption
- Easily rekey secrets to the latest encryption standards
- Share your secrets with other users and devices
- One-way encryption for production secrets
- Run commands and scripts with secrets injected via environment variables
- Store your secrets securely in git with human-readable diffs

Visit https://rotx.dev for more information.

## Why Rot?

Aside from the infamous [ROT13 Caeser cipher](https://en.wikipedia.org/wiki/ROT13), cryptographic keys have a tendency to "rot" rather quickly.  Frequent use of keys inevitably leads to leakage and/or compromise, and the underlying encryption algorithms may not be secure in the future.  Cryptographic material doesn't `age` well in general.

## License

The code in this repository is licensed under the [GNU AGPL](https://www.gnu.org/licenses/agpl-3.0.en.html).  Visit https://rotx.dev/pricing/ to purchase a license exemption.

## Development

Our development process is mostly trunk-based with a `main` branch that folks can contribute to using pull requests.  We tag releases as necessary using CalVer.

### Repository Layout

- `./github:` Reusable GitHub Actions
- `./go:` Rot code
- `./hugo:` Rot website
- `./shell:` Development tooling
- `./shared:` Shared libraries from https://github.com/candiddev/shared

Make sure you initialize the shared submodule:

```bash
git submodule update --init
```

### CI/CD

We use GitHub Actions to lint, test, build, release, and deploy the code.  You can view the pipelines in the `.github/workflows` directory.  You should be able to run most workflows locally and validate your code before opening a pull request.

### Tooling

Visit [shared/README.md](shared/README.md) for more information.
