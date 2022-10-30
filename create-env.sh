#!/usr/bin/env /bin/sh

SECRET=$(openssl rand -base64 32)

echo -e SESSION_SECRET="$SECRET" > .env
{
  echo -e DB_USERNAME="root"
  echo -e DB_PASSWORD="password"
  echo -e DB_NAME="annotation"
  echo -e HOST="video_db"
} >> .env

echo "$SECRET"
