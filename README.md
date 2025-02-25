# xfconf-profile

xfconf-profile is a command-line tool to manage Xfce settings using `xfconf-query`. Profiles are simple JSON documents containing only properties that differ from the default values.

## Example usage

Given a simple `profile.json` file:
```json
{
  "xsettings": {
    "/Net/ThemeName": "Chicago95"
  }
}
```
Applying and reverting the changes from the profile:
```bash
$ xfconf-profile apply profile.json
Setting xsettings::/Net/ThemeName => Chicago95

$ xfconf-profile revert profile.json
Resetting xsettings::/Net/ThemeName
```

## Design and Goals

This projects goals are:
  1. Define a schema for xfconf profiles that can be declaratively defined and provided by a distribution
  2. Provide a tool as a static binary that can modify a user's xfconf settings based on this config

This project is a core component of [winblues/vauxite](https://github.com/winblues/vauxite).

Thanks to https://github.com/jamescherti/watch-xfce-xfconf for the inspiration. 
