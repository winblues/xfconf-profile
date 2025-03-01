# xfconf-profile

xfconf-profile is a command-line tool to manage Xfce settings using `xfconf-query`. Profiles are simple JSON documents containing only properties that differ from the default values.


## Example usage

Given a `profile.json` file such as:
```json
{
  "properties": {
    "xsettings": {
      "/Net/ThemeName": "Chicago95",
      "/Net/IconThemeName": "Chicago95"
    },
    "xfwm4": {
      "/general/theme": "Chicago95"
    }
  },
  "metadata": {
    "name": "winblues-blue95",
    "version": 1
  }
}
```
You can apply (and revert) the changes to properties from the profile:
```bash
$ xfconf-profile apply profile.json
• Setting xsettings/Net/ThemeName ➔ Chicago95
• Setting xsettings/Net/IconThemeName ➔ Chicago95
• Setting xfwm4/general/theme ➔ Chicago95

$ xfconf-profile revert profile.json
• Resetting xfwm4/general/theme
• Resetting xsettings/Net/IconThemeName
• Resetting xsettings/Net/ThemeName
```

## Installation
```bash
mkdir ~/.local/bin
curl -L -o ~/.local/bin/xfconf-profile \
  https://github.com/winblues/xfconf-profile/releases/latest/download/xfconf-profile-linux-amd64
chmod +x ~/.local/bin/xfconf-profile
```

## Configuration

xfconf-profile can be configured via a config file at `$XDG_CONFIG_HOME/xfconf-profile/config.yml`. This file is created on first use. The most up-to-date version of this file can be found at [src/assets/config.yml](https://github.com/winblues/xfconf-profile/blob/main/src/assets/config.yml).

```yaml
# Merge behavior for applying properties
# Options are:
#   - soft:  Only change properties if they are set to their default values.
#            This is the default behavior as it preserves changes to properties
#            that users have manually configured.
#   - hard:  Change all properties not explicitly excluded
#   - force: Change all propertiess
merge: "soft"

# List of properties to exclude when applying profiles. Each entry is a
# regular expression on the fully qualified name of the property
#
# Example:
#   exclude:
#     - "^xsettings/Net/ThemeName$" # Specific property
#     - "^xsettings/Net"            # Everything starting with /Net in the xsettings channel
#     - "^xfwm4"                    # Everything in the xfwm4 channel
#     - "mycustomname"              # Everything containing the string "mycustomname"
exclude: []

# Enable or disable the sync feature that winblues uses on login
sync:
  auto: true
```

## Design and Goals

This project is a core component of [winblues/vauxite](https://github.com/winblues/vauxite) and partly inspired by [jamescherti/watch-xfce-xfconf](https://github.com/jamescherti/watch-xfce-xfconf).

Our goals are:
  1. Define a schema for Xfce profiles that can be declaratively defined and provided by a distribution such as [blue95](https://github.com/winblues/blue95).
  2. Provide a tool as a static binary that can modify a user's xfconf settings based on such configs.

## Thanks

This project includes derivative work from the following third-party projects:

- [jamescherti/watch-xfce-xfconf](https://github.com/jamescherti/watch-xfce-xfconf) - MIT License - Copyright (C) 2021-2025 James Cherti
