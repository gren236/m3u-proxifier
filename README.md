# m3u-proxifier

![Go](https://github.com/gren236/m3u-proxifier/workflows/Go/badge.svg?branch=master)

This small tool is used to retrieve remote m3u playlist and insert a custom proxy URL for each entry.

## Install

You have to have Golang >=1.13 version installed. To compile the program to a single binary executable:
```bash
go build -v ./cmd/m3u-proxifier/
```
This command will create a binary inside the current working directory.

## Usage

This tool receives a single command parameter `--config`. Pass a valid configuration JSON file to the executable like this:
```bash
m3u-proxifier --config=/some/dir/config.json
```
It is strictly recommended passing an absolute path to the config as there might be some unexpected behavior using relative paths. Also, it would be wise to use absolute paths inside the config JSON as well.

A config file example can be found at `configs` directory of this repo, but the overall structure is this:
```json
{
  "location": "http://m3u.example.com",
  "proxy": "192.168.1.42:7777",
  "old": "/test/file.m3u",
  "new": "/test/new_file.m3u"
}
```

* `location` - Path to the new updated playlist to be proxified. Location can be either HTTP URL or local file path.
* `proxy` - Address of a proxy to be added to every entry of a playlist.
* `old` - Path to the existing playlist. If no path specified, just a new proxified updated playlist is going to be saved. Otherwise, new proxfied entries are going to be merged to old playlist (added without deleting old entries).
* `new` - Name (path) of a new resulting file created.

## Bug Reporting

Feel free to report any found issues at GitHub page. Have fun :)
