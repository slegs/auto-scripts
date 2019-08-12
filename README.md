# Auto Scripts

Generate bash scripts for managing your remote servers
* creates folder in ~/.auto-scripts
* Assumes SSH available through keys (no passwords)
* config.yaml sets default values
* create templates which are stored in templates folder
* create scripts from templates which are stored in local folder

Command Types for use in Templates
1. Pull - pulls config or files from remote server
2. Push - pushes config or files to remote server
3. Local - local bash command execution
4. Remote - remote bash command execution 

