#!/usr/bin/env bash
url="localhost:8080/api"
create_bodies=(
  '{"email":"sample@gmail.com","password":"password"}'
  '{"email":"sample1@gmail.com","password":"password", "expires_in_seconds":3600}'
)

createUsers() {
  for body in "${create_bodies[@]}"; do
    response=$(curl -X POST "$url/users" \
      -H "Content-Type: application/json" \
      -d "$body")

    echo "${response}"
  done
}

login_bodies=(
  '{"email":"sample@gmail.com","password":"password"}'
  '{"email":"sample1@gmail.com","password":"password", "expires_in_seconds":3600}'
)

loginUsers() {
  for body in "${login_bodies[@]}"; do
    response=$(curl -X POST "$url/login" \
      -H "Content-Type: application/json" \
      -d "$body")

    echo "${response}"
  done
}

createUsers
loginUsers
