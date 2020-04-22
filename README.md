# gapp
Gapp is a small utitily cli that makes life easier around Github actions.

* Easily dispatch repository_dispatch commands to a repository
* Manage secrets across multiple [repositories](https://github.com/eikc/gapp/tree/master/examples) with the help of one or multiple yaml files

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
It can easily be installed with brew:
```
brew install eikc/gapp/gapp
```

### Downloading a Release from GitHub
It's possible to download a release directly from Github

Example with curl:
```
curl -OL https://github.com/eikc/gapp/releases/download/v<version>/gapp_<version>_Darwin_x86_64.tar.gz
```

Extract the binary:
```
tar xf gapp_<version>_Darwin_x86_64.tar.gz
```

Move the binary into your path: 
```
sudo mv gapp /usr/local/bin
```


Note: Windows build is a zip file, which can be extracted and placed somewhere in your [path](https://docs.microsoft.com/en-us/previous-versions/office/developer/sharepoint-2010/ee537574(v=office.14)?redirectedfrom=MSDN)

