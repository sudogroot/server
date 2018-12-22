# Joyread

Joyread is a self-hosted ebook reader written in Go.

### Easy installation
You can simply run `sudo snap install joyread` (on Linux platforms) or ship Joyread with Docker.

### Share your ebooks
Share ebooks with your family and friends: Being a multi-user product, Joyread can be used for sharing ebooks with all users on the platform or just with selected users. You can also keep your ebooks private.

### Categorization and search
Categorize your ebooks by tags and search for it. You can also search by metadata (title and author).

### Nextcloud integration
You might already have your ebook collection on your Nextcloud: Why not just use Nextcloud sync feature in Joyread to grab all of your ebooks and read them?

### Source folder sync
Do you find it cumbersome to upload your massive collection of ebooks via HTTP file upload? In Joyread, each user has a separate source folder where the respective user ebooks will be stored. You can SSH or FTP (or however else you like), upload all of your ebooks to your source folder, and sync it via the web interface.

# Setup
Joyread is under development. It is not ready for production use.

### Prerequisites
 - PostgreSQL 10
 - Create a new role and database in PostgreSQL. Example shown below
   ```
   CREATE ROLE joyreaduser WITH LOGIN PASSWORD 'jellyfish' VALID UNTIL 'infinity';
 
   CREATE DATABASE joyreaddb WITH ENCODING='UTF8' OWNER=joyreaduser CONNECTION LIMIT=-1;
   ```
 ### Development
  - Clone the repo and put it in an appropriate `GOPATH`. For eg: `$GOPATH/src/github.com/joyread/server`
  - Configure the values in `config/app.yaml`
  - Run `go get -d ./...` or `dep ensure` inside the project folder
  - Then `go run ./cmd/joyread/main.go`. This will run the Joyread server on the port mentioned in the `app.yaml` configuration. Default port is `8080`
  - In order to compile SCSS, you can do `gulp` inside the project folder
