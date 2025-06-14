# oauth2
A lightweight OAuth2 authentication system powered by Vue 3 (frontend) and Golang (backend). Built with Clean Architecture, it supports access &amp; refresh token management and is designed for easy integration and scalability.

version: '3.8'  
services:
  oauth2:
    container_name: postgres
    image: postgres:14-alpine3.21
    environment:
      - POSTGRES_USER: ${DB_USER}
      - POSTGRES_PASSWORD: ${DB_PASSWORD}
      - POSTGRES_DB: ${DB_NAME}
    ports:
      -${DB_PORT}:5432
    volumes:
      - ./tmp/postgres:/var/lib/postgres/data