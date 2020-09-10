plocks is yet another status bar generator and was inspired by
[Goblocks](https://github.com/Stargarth/Goblocks). plocks is simpler
than any status bar generator I've met before: It just executes shell
code to generate the output for each block; there are no built-in blocks
at all. This makes plocks very easy to use and applicable for many
different operating systems.

The big downside of this is, of course, worse performance than most
other status bar generators. This is, however, mitigated by plocks'
ability to update the blocks at independent intervals or even just when
receiving SIGUSR1 or SIGUSR2. If you have many blocks that need very
frequent updates, plocks is probably not the best choice for you.

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
