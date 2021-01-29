# ghf
ghf is cli to manage file in GitHub repository.

```sh
$ ghf up skanehira images --clip
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

## Settings
At first, please set access token and email in config.yaml.

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

```sh
# upload file
# if upload from clipboard then use --clip
$ ghf up {owner} {repo} [filename] [--clip]

# list file
$ ghf ls {owner} {repo} {branch} [--f]

# delete file
$ ghf del {owner} {repo} {branch}
```

## TODO
- download file

## Author
skanehira
