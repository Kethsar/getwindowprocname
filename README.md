# getwindowprocname
Not actually usable yet haha

Gets the name of the process that spawned the top-level window under the mouse cursor, for Windows.

Can be run as a 'server' on Windows and called with grpc. This will allow you to get the window's process from a separate machine.

## Why
This is only really meant to be used with [LookingGlass](https://github.com/gnif/LookingGlass). The idea is to take a screenshot of LookingGlass and then remotely call this program running on the VM so that you can properly label the screenshot, instead of always having it labeled as "looking-glass-client". Of course, this requires some screenshotting tool on Linux that will allow for this, such as by using a callback like [screenshotter](https://github.com/awused/screenshotter) does.

## TODO
- [x] Actually get the process name, strip the extension, and remove white space.
- [x] Add grpc code.
  - [x] server
  - [x] client
- [x] Add TOML config