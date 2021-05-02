# Usman's Bed and breakfast

This is the repository for my bed and breakfast project.

- Built in Go version 1.16.3
- Uses [chi](https://github.com/go-chi/chi/v5) router
- Uses [scs](https://github.com/alexedwards/scs/v2) for session management
- Uses [nosurf](https://github.com/justinas/nosurf) for CSRF 

## Testing
- `go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out`

## Deploying to server

1. Create a server instance (preferably Ubuntu)
1. SSH into the server as root
1. Run the commands:
    1. `apt update`
    1. `apt upgrade`
    1. `apt install postgres-12`
    1. `service postgres status` (to check if it is running)
    1. Install [Caddy server](https://caddyserver.com/docs/install)
    1. `apt install supervisor` (for starting services)
1. Restart instance
1. Make a new user as you shouldn't be root
    - `add user <username>`
1. Give user access to execute commands as sudo
    - `usermode -aG sudo <username>`
1. Verify that you can connect to instance via that user
    - `ssh <username>@<instance-ip>`
    - `sudo ls` (to verify that user can run commands as sudo)
### Install golang on the server
1. Go to their [website](https://golang.org/dl/) and copy the linux package
1. Run `wget <package>`
1. Follow the installation instructions from the website
1. Put the path in the `.profile` by
    1. `vi .profile`
    1. At the end of the file, put `export PATH=$PATH:/usr/local/go/bin
       `
1. Logout and log back in and verify that go is installed
1. Pull the source code from github in the home directory

### Setup database and migrations
1. Setup the remote database
    1. `cd /etc/postgresql/12/main`
    1. `sudo vi pg_hba.conf`
    1. At the bottom of the file, find comment `IPv4 local connections`
    1. Change the `md5` at the end of the line below the comment to `trust`
    1. Do the same for under the comment `IPv4 local connections`
    1. **Note:** Secure the db as necessary, as the above is not very secure
    1. Save changes and exit
    1. `sudo service postgresql stop`
    1. `sudo service postgresql start`
    1. Make sure it's running: `ps ax | grep postgr`
    
1. Create a database in postgres called 'bedandbreakfast'
1. cd into the codebase
1. `cp database.yml.example database.yml`
1. Update the database.yml as appropriate
1. To install migrations, soda needs to be installed
    - `go get github.com/gobuffalo/pop/...`
    - `vi .profile`
    - At the end of file add: `export PATH=$path:~/go/bin`, save and exit
    - `export PATH=$path:~/go/bin`
    - `which soda` (to verify soda is working)
1. Go back to project root directory and run
    - `soda migrate`
1. `go build -o bedandbreakfast cmd/web/*.go`

### Setup caddy configuration
1. `cd /etc/caddy`
1. `sudo mv Caddyfile Caddyfile.dist`
1. `sudo vi Caddyfile`
1. Paste in the following:
   ```json
    {
        email usmanzaheer1995@gmail.com
    }
    (static) {
        @static {
            file
            path *.ico *.css *.js *.gif *.jpg *.jpeg *.png *.svg *.woff *.json
        }
        header @static Cache-Control max-age=5184000
    }
    (security) {
        header {
            # enable HSTS
            Strict-Transport-Security max-age=31536000;
            # disable clients from sniffing media type
            X-Content-Type-Options nosniff
            # keep referrer data off of HTTP connections
            Referrer-Policy no-referrer-when-downgrade
        }
    }
    import conf.d/*.conf

1. Save and quit
1. `sudo mkdir conf.d` -> `cd conf.d`
1. `sudo vi bedandbreakfast.conf`
1. Paste in the following (and update stuff where appropriate)
   ```json
    <PUT ip/domain/subdomain HERE> {
        encode zstd gzip
        import static
        import security
    
        log {
                output file /var/www/<Project dir>/logs/caddy-access.log
                format single_field common_log
        }
        reverse_proxy http://localhost:<Put PORT on which app will run>
    }

1. `cd /var` -> `sudo mkdir www` -> `cd www`
1. Move the project files from home to /var/www
    - `sudo mv ~/bedandbreakfast bedandbreakfast` (must be the same dir as specified in caddy file)
1. Create logs folder
    - `mkdir logs`
    - `sudo chmod 777 logs` (make writable for anyone)
1. Restart caddy and make sure it's running
    1. `sudo service restart caddy`
    1. `sudo service caddy status`
    
### Setup supervisor
1. `cd /etc/supervisor/conf.d`
1. `sudo vi bedandbreakfast.conf`
1. Paste in the following and (the command takes command line parameters):
    ```json
    [program:bedandbreakfast]
    command=/var/www/bedandbreakfast/bedandbreakfast -dbname=bedandbreakfast -dbpass=Doranboots101$ -dbuser=postgres
    directory=/var/www/bedandbreakfast
    autorestart=true
    autostart=true
    stdout_logfile=/var/www/bedandbreakfast/logs/supervisord.log
    stderr_logfile=/var/www/bedandbreakfast/logs/supervisorerrod.log

1. Save and exit
1. Verify that supervisor can run our app
    - `sudo supervisorctl`
    - `update`
    - `status` (app should now be running)
1. Add supervisorctl to require no password for sudo (required for github actions)
    1. `sudo visudo`
    1. Paste in the following line at the bottom of the file
        - `<username> ALL= NOPASSWD: /bin/supervisorctl`
    
### Add an update script for the server
1. `cd /var/www/bedandbreakfast`
1. `vi update.sh`
1. Paste in the following commands:
```shell
#!/bin/bash

git pull

# need full path for github actions
~/go/bin/soda migrate
/usr/local/go/bin/go build -o bedandbreakfast cmd/web/*.go

sudo supervisorctl stop bedandbreakfast
sudo supervisorctl start bedandbreakfast
```
1. Save and exit
1. Make it executable
    - `chmod +x update.sh`
    
### Add github action
1. Make a new github action (found in the github repo)
1. Add in the following
    - **NOTE**: Remove the leading slashes from the dollar signs
    - Add secrets from `Settings -> Secrets -> New repository secrets`
    - `Host` is the server IP
```shell
# This is a basic workflow to help you get started with Actions

name: CI

# Controls when the action will run. 
on:
  # Triggers the workflow on push or pull request events but only for the main branch
  push:
    branches: [ main ]

# A workflow run is made up of one or more jobs that can run sequentially or in parallel
jobs:
  # This workflow contains a single job called "build"
  build:
    # The type of runner that the job will run on
    runs-on: ubuntu-latest

    # Steps represent a sequence of tasks that will be executed as part of the job
    steps:

      # Runs a set of commands using the runners shell
      - name: executing remote ssh commands using password
        uses: appleboy/ssh-action@master
        with:
          host: ${{ secrets.HOST }}
          username: ${{ secrets.USERNAME }}
          password: ${{ secrets.PASSWORD }}
          script: |
            cd /var/www/bedandbreakfast
            ./update.sh
```

### Restrict access to your app port so it cannot be accessed except from your domain
1. `sudo ufw status`
1. `sudo ufw deny <port>/tcp`
1. `sudo ufw enable`
1. `sudo ufw status`
