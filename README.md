# Table reservation app

## Prerequisites
Create .env file in root directory and add following values:
```dotenv
POSTGRES_URI=postgres://<username>:<password>@127.0.0.1:5432/<db_name>
SESSION_SECRET=<any string>
AWS_REGION=<write me in dm to get it>
AWS_ACCESS_KEY_ID=<write me in dm to get it>
AWS_SECRET_ACCESS_KEY=<write me in dm to get it>
BUCKET_NAME=<write me in dm to get it>
```
This bucket is public, you should be able to access it

Manually give admin role to certain user that is admin
```postgresql
    update users set role_id = 1 where id = ?
```
## Useful links
Full Tutorial for AWS S3 in Go
https://medium.com/wesionary-team/aws-sdk-for-go-and-uploading-a-file-using-s3-bucket-df7425317a40


