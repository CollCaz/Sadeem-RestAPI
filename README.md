
# Project Sadeem-RestAPI

RestAPI written in Go using the Echo framework

## Getting Started

1. Clone this repo locally.

2. Make sure you have PostgreSQL set up and running.

3. Create a new database for this project.

4. Define your `$DATABASE_URL`, `$JWT_SIGNING_TOKEN`, `$DEFAULT_PROFILE_PICTURE` and `$PICTURE_DIR` environment variables.

5. run `make up` to apply up migrations.

5. run `make build` to buld the application.

6. run the binary found in bin/ .

## MakeFile

print all make options and their description
```bash
make help
```

build the application
```bash
make build
```

run the application
```bash
make run
```

clean up binary from the last build
```bash
make clean
```
	
apply down migrations
```bash
make down
```
apply up migrations
```bash
make up
```

## End Points

### POST

. /api/register
```json
{
    "userName" : "MyName",
    "emial"    : "test@email.com",
    "password" : "12345678",
}
```

. /api/login

``` json
{
    "email" : "test@email.com",
    "password" : "12345678"
}

```

./api/categories

```json
{
    "name" : "Example"
}

```
./api/user-categories  Activates or deactivates a catagory for a particular user

```json
{
    "userName" : "nonAdminUser",
    "categories" : ["Example1", "Category1", "Category2"],
    "activated" : true // set it to false to deactivate the categories
}
```

