# stress

```bash
NAME:
   stress - A tool for stress testing commands

USAGE:
   stress [global options] command [command options]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --cmd value                 Command to run for stress testing
   --runs value, -r value      Number of times to run the command (default: 100)
   --parallel value, -p value  Number of parallel executions (default: 4)
   --help, -h                  show help
```

## Install 

### Mac

```bash
sudo make install_stress_mac
```

## Usage

```bash
stress --cmd "npm test --workspaces" --runs 100 --parallel 10
```