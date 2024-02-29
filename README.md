
# Project Sadeem-RestAPI

RestAPI written in Go using the Echo framework

## Getting Started

1. Clone this repo locally.

2. Make sure you have PostgreSQL set up and running.

3. Create a new database for this project.

4. Define your `$DATABASE_URL`, `$JWT_SIGNING_TOKEN`, `$DEFAULT_PROFILE_PICTURE` and `$PICTURE_DIR` environment variables.

5. run `make up` to apply up migrations.

5. run `make build` to buld the application. (If you aren't using make, please make sure to manually follow the same commands in the make script)

6. run the binary found in bin/ .

7. after creating a user, you can make them an admin with this query

```sql

INSERT INTO admin_users (user_id) VALUES("Your User's ID")

```

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

Respnses are available in english and in arabic
To choose the language simply change the Accept-Languae header to en or ar.

. /api/users
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
./api/user-categories  Activates or deactivates categories for a particular user

```json
{
    "userName" : "nonAdminUser",
    "categories" : ["Example1", "Category1", "Category2"],
    "activated" : true // set it to false to deactivate the categories
}
```

## GET

./api/users/:name/profile-picture  Get the profile picture of a particular user

./api/users:name  Get user info 

./api/categories?page=1&size=1&  Get all activated categories with pagination


## PUT

./api/users/:name/  updates user info
```json
{
    "userName" : "newUserName",
    "email" : "new@email.com",
    "password" : "currentPassword"
}
```
./api/users/:name/profile-picture  Updates the profile picture with the one attached in the body

## DELETE

./api/users/:name  deletes a user

./api/users/:name/profile-pictures deletes the user's profile pipcture, reseting it back to the default one

./api/categories/:name deletes a category
