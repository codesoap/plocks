package main

const separator = " | "

var blocks = []block{
	{
		command:         `sndioctl -n output.level | awk '{print "🔊 " int($0 * 100) "%"}'`,
		updateOnSIGUSR1: true,
	},
	{
		command:  `apm -l | awk '{print "⚡ " $0 "%"}'`,
		interval: 15,
	},
	{
		command:  `curl 'wttr.in/Berlin?format=1'`,
		interval: 30 * 60,
	},
	{
		command:  `ifconfig trunk0 | awk '/inet/ {print $2}'`,
		interval: 10,
	},
	{
		command:  `date +'%Y-%m-%d %H:%M'`,
		interval: 20,
	},
}
