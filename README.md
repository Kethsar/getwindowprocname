# getwindowprocname
~~Not actually usable yet haha~~ Usable but also does more than it was originally made for but >changing the name

Gets the name of the process that spawned the top-level window under the mouse cursor, for Windows.

Can be run as a 'server' on Windows and called with grpc. This will allow you to get the window's process from a separate machine.

## Why
This is only really meant to be used with [LookingGlass](https://github.com/gnif/LookingGlass). The idea is to take a screenshot of LookingGlass and then remotely call this program running on the VM so that you can properly label the screenshot, instead of always having it labeled as "looking-glass-client". Of course, this requires some screenshotting tool on Linux that will allow for this, such as by using a callback like [screenshotter](https://github.com/awused/screenshotter) does.

## Usage
Copy the example toml and remove "example" from its name.

On Windows, set the port and turn on Server mode (0 use for this on Windows as a normal client, I mean really now). Start the server.

On Linux, copy the toml config to ~/.config/getwindowprocname/ and set the address to your Windows machine (meant for local VMs but I'm not going to enforce that). Use in scripts as needed or something, I dunno. No one else should be using this.

## TODO
- [x] Actually get the process name, strip the extension, and remove white space.
- [x] Add grpc code.
  - [x] server
  - [x] client
- [x] Add TOML config
- [ ] Remove TODO list from README