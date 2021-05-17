# Table reservation app

## Prerequisites

Create an AWS Account in Amazon
https://aws.amazon.com/premiumsupport/knowledge-center/create-and-activate-aws-account/

Create IAM User and get Access and Secret key
https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/setting-up.html

Full Tutorial for Go
https://medium.com/wesionary-team/aws-sdk-for-go-and-uploading-a-file-using-s3-bucket-df7425317a40

Create .env file in root directory and add following values:
```dotenv
POSTGRES_URI=postgres://<username>:<password>@127.0.0.1:5432/<db_name>
SESSION_SECRET=<any string>
AWS_REGION=eu-north-1
AWS_ACCESS_KEY_ID=<your access key>
AWS_SECRET_ACCESS_KEY=<your secret key>
BUCKET_NAME=ez-booking-bucket
```
this bucket is public, you should be able to access it

Manually give admin role to certain user that is admin
```postgresql
    update users set role_id = 1 where id = ?
```
