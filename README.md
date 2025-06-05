# oauth2
A lightweight OAuth2 authentication system powered by Vue 3 (frontend) and Golang (backend). Built with Clean Architecture, it supports access &amp; refresh token management and is designed for easy integration and scalability.

https://accounts.google.com/o/oauth2/v2/auth?
client_id=487291772648-olmn5125kgmujjerkru6ihjh3nrbefa4.apps.googleusercontent.com&
redirect_uri=http://localhost:8080/v1/auth/callback&
response_type=code&
scope=openid%20email%20profile&
prompt=consent&
access_type=offline

https://accounts.google.com/o/oauth2/auth?client_id=487291772648-olmn5125kgmujjerkru6ihjh3nrbefa4.apps.googleusercontent.com&redirect_uri=http://localhost:8080/v1/auth/callback&response_type=code&scope=openid%20email%20profile&state=test123