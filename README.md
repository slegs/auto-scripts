# Auto Scripts

Generate and run bash scripts for managing your remote servers with controlled flags and defaults
* Written in Golang
* Creates configuration folder in ~/.auto-scripts
* Assumes host folder where auto-scripts is run is named as remote server FQDN
* Assumes SSH available through keys (no passwords)
* config.yaml sets default values
* create further templates which are stored in templates folder and reference in config.yaml
* Default scripts use rsync to synchronise directories and files
* Auto-backups are taken of changed files when getting or pushing

Default Command Types with precreated templates
1. Pull - pulls config or files from remote server
2. Push - pushes config or files to remote server
3. Local - local bash command execution
4. Remote - remote bash command execution
5. Manage - remote execution of systemctl

Uses
* Viper for reading yaml config file (github.com/spf13/viper)
* Urfave Cli for calling the bash scripts with flags and defaults (github.com/urfave/cli)
