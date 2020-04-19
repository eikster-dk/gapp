# gapp

```
gapp is a small Github utility cli that will fill the gap of features that I miss when working github

Usage:
  gapp [command]

Available Commands:
  actions     Commands for github actions
  completion  Generate shell completion scripts
  help        Help about any command
  login       Login initiates gapp by allowing you to provide a personal access token
  version     Print the gapp version number
```

## Installation of gapp

### MacOS

### Linux
### Docker
### Downloading a Release from GitHub
It's possible to download a release directly from Github

Example with curl:
```
curl -OL https://github.com/eikc/gapp/releases/download/v<version>/gapp_<version>_Darwin_amd64.tar.gz
```

Extract the binary:
```
tar xf gapp_<version>_Darwin_amd64.tar.gz
```

Move the binary into your path: 
```
sudo mv gapp /usr/local/bin
```


Note: Windows build is a zip file, which can be extracted and placed somewhere in your [path](https://docs.microsoft.com/en-us/previous-versions/office/developer/sharepoint-2010/ee537574(v=office.14)?redirectedfrom=MSDN)

## Actions

### Secrets file

