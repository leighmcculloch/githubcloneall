# githubcloneall

Clone all github repositories for a user.

## Install

### Binary (Linux; macOS; Windows)

Download and install the binary from the [releases](https://github.com/leighmcculloch/githubcloneall/releases) page.

### Brew (macOS)

```
brew install 4d63/githubcloneall/githubcloneall
```

### From Source

```
go get 4d63.com/githubcloneall
```

## Usage

```
Usage: githubcloneall -u username -d dir -token TOKEN -type orgs

  -d string
        Output directory
  -h    Print help
  -token string
        Github personal access token or oauth token
  -type string
        github type (users, orgs) (default "users")
  -u string
        GitHub username
```
