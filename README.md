# imgur-cli
A cli tool for [Imgur](https://imgur.com) to upload and delete images.

# Features
1. Anonymous upload of images to imgur.
2. Delete anonymous image uploads.

# Getting the client id
This program requires a client id to work.
Get a client id from [https://api.imgur.com/oauth2/addclient](https://api.imgur.com/oauth2/addclient).

## If you are on windows
```
set IMGUR_CLIENT_ID=[CLIENT_ID]
```

## If you are on linux
```
$ export IMGUR_CLIENT_ID=[CLIENT_ID]
$ imgur-cli upload cat.png
```
Or
```
$ IMGUR_CLIENT_ID=[CLIENT_ID] imgur-cli upload cat.png
```
You can also place it in your shell profile, for example `~/.bashrc` or `~/.zshrc`

# Usage

## Anonymous upload of images to imgur
```
imgur-cli upload <path to image file>
```

## Delete uploaded images
```
imgur-cli delete <delete hash>
```