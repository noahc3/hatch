# hatch cli

A CLI tool for inspecting and executing Pterodactyl eggs without Pterodactyl or the Wings daemon.

## Usage

Usage: `hatch-cli <tag> [<operation>...]`
- `<tag>` is a pointer to the Egg to use. This can be a file path, web URL, or a tag specifying remote provider information (more information below).
- `<operation>` is a list of one or more valid operations:
  - `print`: Print the egg in JSON format.
  - `checkvars`: Load environment variables and check their values against the validation rules defiend by the Egg.
- Optional arguments:
  - `--quiet`, `-q`: Suppress output.
  - `--unattended`, `-y`: Run in unattended mode (no prompts).
  - `--help`: Display this help message.
- Supported remote providers (square brackets indicate an optional parameter):
  - GitHub - `@github:owner/repo[/ref]:path/to/file.json`
    - Fetch the egg from a GitHub repository. 'ref' is optional and defaults to 'master'.
    - Example: `@github:pelican-eggs/eggs:game_eggs/minecraft/java/paper/egg-paper.json`

Example: `hatch-cli @github:pelican-eggs/eggs:game_eggs/minecraft/java/paper/egg-paper.json checkvars`

## FAQ

- Why was this tool created?
  - I'm building a custom server control panel as a personal project, but I've already learned from my past projects
    that building an extensive set of software install scripts is absolutely mind-numbing. Pterodactyl already
    provides dozens of install scripts for tons of games and services, all that was left was to build a way to
    utilize them without actually needing Pterodactyl/Wings.

- Why doesn't this use Wings directly (either as a go module or the binary)
  - Several reasons, all of which can be summarized by me being inexperienced with Go:
    - Wings is a very complex piece of software with high dependency on the environment it's being run in. Information
      about the server being operated on is expected to be present for the installation logic, but obviously, this
      CLI doesn't install it with any server context information. Of course I could just make dummy data that
      satisfies the requirements, but my inexperience with Go and with the Pterodactyl/Wings codebase makes this 
      harder.
    - Along with the above, and given that I have little experience with Go and wanted to learn more, I thought this
      CLI would be a worthwhile project to write from scratch to get used to Go. On one hand that means I'm not 
      necessarily following good Go programming standards, but for me thats OK. That might also mean however that 
      this tool may not be stable or robust enough for your use-case, so use this program carefully!

- Does this support all Pterodactyl Egg scripts?
  - Maybe, but probably not. I have not searched through the Egg struct definitions extensively nor have I
   read or tested all of the available Eggs in the pelican-eggs repo. Only a subset of Minecraft Java scripts
   have been tested!


## Acknowledgements

This project obviously only exists thanks to the hard work by the Pterodactyl team. While this tool was created
to allow the use of Pterodactyl Eggs with other control panels, my hope is that myself or others who use this 
tool will be able to sponsor the Pterodactyl project in the future to promote it's development. Pterodactyl is 
an incredibly important project in the FOSS server management space with few alternatives, so it's imperative 
that the Pterodactyl project remains alive and healthy.

Also, thanks to the Pterodactyl team for commenting their codebase so extensively, it was a pleasant sight and
very handy!

## License

In good faith, this project is licensed under the same core MIT license as Pterodactyl.