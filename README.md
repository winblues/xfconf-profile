# xfconf-profile

xfconf-profile is a command-line tool to manage Xfce settings using `xfconf-query`. Profiles are simple JSON documents containing only properties that differ from the default values.


## Example usage

Given a simple `profile.json` file:
```json
{
  "xsettings": {
    "/Net/ThemeName": "Chicago95",
    "/Net/IconThemeName": "Chicago95"
  },
  "xfwm4": {
    "/general/theme": "Chicago95"
  }
}

```
Applying and reverting the changes from the profile:
```bash
$ xfconf-profile apply profile.json
• Setting xsettings::/Net/ThemeName ➔ Chicago95
• Setting xsettings::/Net/IconThemeName ➔ Chicago95
• Setting xfwm4::/general/theme ➔ Chicago95

$ xfconf-profile revert profile.json
• Resetting xfwm4::/general/theme
• Resetting xsettings::/Net/IconThemeName
• Resetting xsettings::/Net/ThemeName
```

## Installation
```bash
mkdir ~/.local/bin
curl -L -o ~/.local/bin/xfconf-profile \
  https://github.com/winblues/xfconf-profile/releases/latest/download/xfconf-profile-linux-amd64
chmod +x ~/.local/bin/xfconf-profile
```

## Design and Goals

This project is a core component of [winblues/vauxite](https://github.com/winblues/vauxite) and partly inspired by [jamescherti/watch-xfce-xfconf](https://github.com/jamescherti/watch-xfce-xfconf).

Our goals are:
  1. Define a schema for xfconf profiles that can be declaratively defined and provided by a distribution
  2. Provide a tool as a static binary that can modify a user's xfconf settings based on this config

## Thanks

This project includes derivative work from the following third-party projects:

- [jamescherti/watch-xfce-xfconf](https://github.com/jamescherti/watch-xfce-xfconf) - MIT License - Copyright (C) 2021-2025 James Cherti
