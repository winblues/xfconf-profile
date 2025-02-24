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
