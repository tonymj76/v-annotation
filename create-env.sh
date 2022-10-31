#!/usr/bin/env /bin/sh

SECRET=$(openssl rand -base64 32)

echo  SESSION_SECRET="$SECRET" > .env
{
  echo  DB_USERNAME="root"
  echo  DB_PASSWORD="password"
  echo  DB_NAME="annotation"
  echo  HOST="video_db"
} >> .env

echo "$SECRET"
