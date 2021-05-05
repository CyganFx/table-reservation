# Table reservation app

## Prerequisites

Create .env file in root directory and add following values:
```dotenv
POSTGRES_URI=postgres://<username>:<password>@127.0.0.1:5432/<db_name>
SESSION_SECRET=<any secret>
```

Manually give admin role to certain user that is admin
```postgresql
    update users set role_id = 1 where id = ?
```
