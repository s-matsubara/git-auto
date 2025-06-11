# git-auto
![Go Version](https://img.shields.io/github/go-mod/go-version/s-matsubara/git-auto/main)
![Release](https://img.shields.io/github/v/release/s-matsubara/git-auto)
![CI](https://img.shields.io/github/actions/workflow/status/s-matsubara/git-auto/release.yml)
![Homebrew](https://img.shields.io/badge/homebrew-FBB040?logo=homebrew&logoColor=white)

`git-auto` is a CLI for automating common Git tasks.

## Installation

### Homebrew
```sh
brew tap s-matsubara/homebrew-tap
brew install --cask git-auto
```

Binaries are also available on the [releases page](https://github.com/s-matsubara/git-auto/releases).

## Usage

```
git-auto [command]
```

### tag
Increment semantic versions or specify a version.

```
git-auto tag [<version>] [major|minor|patch] [-p] [-m <message>]
```

- `--push` push the created tag to origin.
- `--message` add a tag message.

### delete-merged-branch
Remove local branches that have already been merged.

```
git-auto delete-merged-branch
```

### Git aliases
Set convenient aliases if desired:

```sh
git config --global alias.auto '!git-auto'
git config --global alias.mergedd '!git-auto delete-merged-branch'
git config --global alias.auto-tag '!git-auto tag'
```

## License
MIT

