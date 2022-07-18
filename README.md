# SWR
## Building
After cloning the repository, `make all` will download dependencies and build the server.

## Running
Once built, the executable `server` will be in the `bin` directory. Simply run it from the root
of the repo with `./bin/server` on POSIX or `bin\server.exe` on Windows.

## History
Growing up I used to play muds. I loved them. There was a mud called SWR based on SMAUG (which in turn was a merc/diku derivative)
that recreated the Star Wars universe in text based form. It was pretty good and other muds formed by forking the source and adding
their contributions.

SWR is a reimagining of that codebase without the merc/diku/smaug legacy. Instead it's a pure go mud with some javascript for scripting
that is more flexible and robust than the C mud engines of old.


## Features
- TELNET echo/no-echo
- ANSI Colors
- Color Codes in MERC/DIKU/SMAUG style (ex: &WThis is White&d) 'This is White' in ⬜ text `#ffffff`
- Fuzzy Matching. Commands in the mud are fuzzy matched. Which means the command `LOOK` can be called using just `l` or `look`. Similarly the command `SCORE` can be called with just `sc`. This makes it quick to execute commands with shortcuts.
- YAML based areas. Easy to edit.
- Progressive Language system with alphabet support.
- Scheduled Function calling.
- Scheduled Backups using `tar` shell command.
- Multi-threaded using *go* routines.
- Abstract command system makes it easy to add commands.

## Planned
- [ ] Area Resets / Mob Resets
- [ ] Browser based area editor
- [ ] Spaceships
- [ ] Crafting
- [ ] Combat
- [ ] Starting Newbie Area
- [ ] Training Skills
- [ ] Skill Tree
- [ ] Cargo/Trade
- [ ] Guilds/Clans
- [ ] Factions
- [ ] Space Combat
- [ ] Dungeons
- [ ] Raids
- [ ] Quests