### User management CRUD service using
> Golang
> Postgres
> gRPC

### User Table
```sql
CREATE TABLE IF NOT EXISTS "User" (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    email VARCHAR(50) NOT NULL UNIQUE,
    phone_number INT
);
```

### gRPC Methods
- `CreateNewUser` - Creates and stores a new user profile in database.
- `GetAllUsers` - Retrieves all the users from database.
- `GetUser` - Retrieves a user from database based on unique ID.
- `UpdateUser` - Updates a user based on unique ID.
- `DeleteUser` - Deleted the user from database.

### Compiling instructions
```cmd
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out= --go-grpc_opt=paths=source_relative proto/user.proto
```