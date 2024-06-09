# freepbx-s3-upload

**Config**

`/etc/asterisk/s3_go.conf`

`[aws]`

`access_key_id: AKIAX`

`secret_access_key: Ksa2D5aj`

`s3_bucket_name: freepbx`

`aws_region: ap-southeast-1`

**Run**

`go build`

`s3upload --file=recording_file_name`
