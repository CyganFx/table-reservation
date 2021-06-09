# Table reservation app

## Prerequisites
Create .env file in root directory and add following values:
```dotenv
POSTGRES_PASSWORD=<password>
SESSION_SECRET=<any string>
AWS_SECRET_ACCESS_KEY=<confidential>
```
This bucket is public, you should be able to access it

Manually give admin role to certain user that is admin
```postgresql
    update users set role_id = 1 where id = ?
```
## Useful links
Full Tutorial for AWS S3 in Go
https://medium.com/wesionary-team/aws-sdk-for-go-and-uploading-a-file-using-s3-bucket-df7425317a40


