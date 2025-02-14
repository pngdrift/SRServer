# SRServer

[![Discord](https://img.shields.io/discord/1340027267304390769?label=Discord&style=plastic&link=https%3A%2F%2Fdiscord.gg%2F3gRGmyUJvx)](https://discord.gg/3gRGmyUJvx)

Open-source server for the android game "Drag Racing: Streets" aka "Street racing".


## Supported game versions
|APK Version|Client version|Obfuscated|Link|
|-|-|-|-|
|v1.8.1|35|No|[4PDA (Login required)](https://4pda.to/forum/dl/post/11818859/street-race.apk)|
|v1.8.2|35|Yes (weakly)|[Apk4Fun](https://www.apk4fun.com/go.php?id=mobi.square.sr.android&p=221002&s=ONmtIa09PDW7w&l=https%3A%2F%2Ff0.apk4fun.com%2Fget.php%3Fp%3D221002%26i%3Dmobi.square.sr.android%26v%3D1.8.2)|
|v1.8.5|36|Yes|[Web Archive](https://web.archive.org/web/20240724184522/https://f0.apk4fun.com/get.php?p=240006&i=mobi.square.sr.android&v=1.8.5&token=92b7d8506459889b89d1f9e4938d1dd61721852121)|
|v1.8.6|37|Yes|[Web Archive](https://web.archive.org/web/20240724184436/https://dl.apkhome.net/2018/1/drag-racing-v1.8.6-full.apk)|

[Changelog](https://www.apk4fun.com/history/185939/8/)

## Editing APK
The APK needs be modified, for e.g., using [MT Manager](https://mt2.cn).
* In the SRConfig class to replace the content server address from "85.25.237.169" to yours.
* If needed, clone the APK.
* Initially, the login works only through [Odnoklassniki](https://ok.ru) and [Facebook](https://www.facebook.com), in order for [VKontakte](https://vk.com) to work
you need to change in file resources.arsc `vk_api_version` from `5.60` to `5.90` or higher.
Unfortunately, login via Google Play Games Service is not available in any way.

## Building proto files
The game uses [protobuf](https://protobuf.dev).
```sh
mkdir "./proto/out"
protoc --proto_path=proto/src --go_out=./proto/out --go_opt=paths=source_relative proto/src/*.proto
```

## Setup config
Copy `conf/Conf.go_example` to `conf/Conf.go` and fill in the required fields for the server ip and mysql database credentials.

## Running
```sh
go run main.go
```

## Warning
The resources and items are not fully consistent with game version 1.8.x.

## Screenshot
![](https://i.imgur.com/qmRFEJp.png)

## Using Git LFS
Please note that this repository uses [Git LFS](https://git-lfs.com/).