plocks is yet another status bar generator and was inspired by
[Goblocks](https://github.com/Stargarth/Goblocks). plocks allows for
individual update intervals for each block and triggering updates of
selected blocks via user signals. It thus uses less resources than many
other status bar generators and allows for immediate updates when e.g.
the volume is changed.

# Installation
```shell
git clone git@github.com:codesoap/plocks.git
cd plocks
go install
# The binary is now at ~/go/bin/plocks.
```

# Usage
Configure your blocks by modifying `config.go`. Provide
`updateOnSIGUSR1: true` or `updateOnSIGUSR2: true` for each block that
shall be updated when plocks receives the respective signal. You can
send these signals by using something like `pkill -USR1 plocks`. If no
`interval` is provided for a block, the block will never be updated
automatically.

To use plocks for dwm, execute something like this, once X11 is started:
```
plocks | while read line; do xsetroot -name " $line"; done &
```
