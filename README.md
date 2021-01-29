# ghf
ghf is cli to manage file in GitHub repository.

## Settings
At first, please set access token and email in config.yaml.

```yaml
user: "skanehira"
email: "sho19921005@gmail.com"
token: "xxxxxxxxxxxxxxxxxx"
```

config.yaml must be in the bellow place.

| OS         | place                                               |
|------------|-----------------------------------------------------|
| Window     | `%AppData%¥ghf¥config.yaml`                         |
| Mac        | `$HOME/Library/Application Support/ghf/config.yaml` |
| Linux/Unix | `$HOME/.config/ghf/config.yaml`                     |

## Usage

```sh
# upload file
# if upload from clipboard then use --clip
$ ghf up {owner} {repo} [filename] [--clip]

# list file
$ ghf ls {owner} {repo} {branch} [--f]
```

## TODO
- download file
- delete file

## Author
skanehira
