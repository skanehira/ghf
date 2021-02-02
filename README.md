# ghf
ghf is cli to manage file in GitHub repository.

```sh
$ ghf up skanehira images main --clip
https://raw.githubusercontent.com/skanehira/images/main/20210129125731.png?token=AB4F5T5GEVU3VYFT5CWI2ILACOD6Q

$ ghf ls skanehira images main
20201112090031.png
20210129124442.png
20210129124829.png
20210129125731.png
```

## Requrements
- file(Linux only)
- xclip(Linux only)

## Installation
```
$ git clone https://github.com/skanehira/ghf.git
$ go install
```

## Settings
At first, please set access [token](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token) and email in config.yaml.

```yaml
email: "sho19921005@gmail.com"
token: "xxxxxxxxxxxxxxxxxx"
```

The config.yaml must be in the bellow place.

| OS         | place                                               |
|------------|-----------------------------------------------------|
| Window     | `%AppData%¥ghf¥config.yaml`                         |
| Mac        | `$HOME/Library/Application Support/ghf/config.yaml` |
| Linux/Unix | `$HOME/.config/ghf/config.yaml`                     |

## Usage

### upload

```sh
Usage:
  ghf up {owner} {repo} {branch} [file...] [flags]

Examples:
  $ ghf up skanehira images main sample1.png sample2.png
  $ ghf up skanehira images main --clip
  $ ghf up skanehira images main sample.png --dir gorilla
```

### list

```sh
Usage:
  ghf ls {owner} {repo} {branch} [--f]

Examples:
  $ ghf ls skanehira ghf master
  $ ghf ls skanehira ghf master --f
```

### del

```Sh
Usage:
  ghf del {owner} {repo} {branch}

Examples:
  $ ghf del skanehira ghf master
```

## Author
skanehira
