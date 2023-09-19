# bitbucket-archiver

This is a tool to automatically backup git repos hosted on Bitbucket locally. 
The primary use case for this tool is to reduce clutter on Bitbucket. 

## Modes of Operation 
The archiver has two options to specify the repos to archive.

1. Archive all repos that are marked as archived on Bitbucket
2. Provide a list of bitbucket projects to archive (All repos in the projects will be archived)


## Usage

### Only create a backup
1. Copy the `.env.example` file to `.env` and fill in the values
2. Choose the mode of operation
    - Default: Archive all repos that are marked as archived on Bitbucket 
        ```bash
        docker run --env-file .env ghcr.io/ls1admin/bitbucket-archiver:latest
        ```

    - Project based archival: 
        - Create a file with the name `projects.txt` in the root of the project
        - Add the project names to the file, one per line
        ```bash    
        docker run --env-file .env -v $(pwd)/projects.txt:/app/projects.txt ghcr.io/ls1admin/bitbucket-archiver:latest --project-file projects.txt
        ```
        > Explanation  
        > We need to explicitly provide the env file to the docker container.
        > We also need to mount the `projects.txt` file to the container so that it can be read by the script.

### Backup and Delete Repositories 
1. Copy the `.env.example` file to `.env` and fill in the values
2. Choose the mode of operation
    - Default: Archive all repos that are marked as archived on Bitbucket 
        ```bash
        docker run --env-file .env ghcr.io/ls1admin/bitbucket-archiver:latest --execute-delete
        ```

    - Project based archival: 
        - Create a file with the name `projects.txt` in the root of the project
        - Add the project names to the file, one per line
        ```bash    
        docker run --env-file .env -v $(pwd)/projects.txt:/app/projects.txt ghcr.io/ls1admin/bitbucket-archiver:latest --project-file projects.txt --execute-delete
        ```