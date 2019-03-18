# keepsake

Experiments in OAuth2 and JWT

---

## create access token

POST http://localhost:10101/api/v2/token/oauth/authorize
Content-Type: application/x-www-form-urlencoded

grant_type=client_credentials&client_id=test&client_secret=gremlin

---

## verify access token

GET http://localhost:10101/api/v2/token/oauth/verify?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTIzOTE2NDYsImlhdCI6MTU1MjM4ODA0NiwiaXNzIjoidGVzdCIsImp0aSI6IjhkNjA0YmMxLTZkMjktNGYxYy1iYWU0LTQ0MGM5NGM4MmQ1OCIsInN1YiI6InRlc3QifQ.p8Gk0d8WPiKJus57BO-Z4LsAuA6zAHNbW2ms8KlECdM

---

## sign JWT

POST http://localhost:10101/api/v2/token/jwt/sign HTTP/1.1
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTIzOTAwMDUsImlhdCI6MTU1MjM4NjQwNSwiaXNzIjoidGVzdCIsImp0aSI6Ijk4OTkzMWQ4LTBlN2QtNGU3Ni1iYmY2LTNiNDE1MmIyYzg4OCIsInN1YiI6InRlc3QifQ.OESutED2FrqNrrqGLfbssJA9sK3fbhJ1TqIoEY9HXSA

{
"iss": "test",
"sub": "a6d5c443-1f51-4783-ba1a-7686ffe3b54a",
"aud": [
"962fa4d8-bcbf-49a0-94b2-2de05ad274af"
],
"exp": 1510185728,
"iat": 1510185228,
"azp": "962fa4d8-bcbf-49a0-94b2-2de05ad274af",
"nonce": "fc5fdc6d-5dd6-47f4-b2c9-5d1216e9b771",
"name": "Ms Jane Marie Doe",
"given_name": "Jane",
"family_name": "Doe",
"middle_name": "Marie",
"picture": "https://example.org/jane.jpg",
"email": "jane@example.org"
}

---

## verify and decode JWT

Bearer is oath

GET http://localhost:10101/api/v2/token/jwt/verify?token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsiOTYyZmE0ZDgtYmNiZi00OWEwLTk0YjItMmRlMDVhZDI3NGFmIl0sImF6cCI6Ijk2MmZhNGQ4LWJjYmYtNDlhMC05NGIyLTJkZTA1YWQyNzRhZiIsImVtYWlsIjoiamFuZUBleGFtcGxlLm9yZyIsImV4cCI6MTU1MjM5MjUzOSwiZmFtaWx5X25hbWUiOiJEb2UiLCJnaXZlbl9uYW1lIjoiSmFuZSIsImlhdCI6MTUxMDE4NTIyOCwiaXNzIjoidGVzdCIsIm1pZGRsZV9uYW1lIjoiTWFyaWUiLCJuYW1lIjoiTXMgSmFuZSBNYXJpZSBEb2UiLCJub25jZSI6ImZjNWZkYzZkLTVkZDYtNDdmNC1iMmM5LTVkMTIxNmU5Yjc3MSIsInBpY3R1cmUiOiJodHRwczovL2V4YW1wbGUub3JnL2phbmUuanBnIiwic3ViIjoiYTZkNWM0NDMtMWY1MS00NzgzLWJhMWEtNzY4NmZmZTNiNTRhIn0.JCLiGu_QqjBmFviX887ioIffbkMDbo4KACzR-ss-fB8Z4lXQGW98IARhvtTPmelMqlU1ac7Xi-OzdfacbA9ng-4C3vhsAmpVdIbVfxeE96qYOLMq9TdRiIdOVXZ63YfV-CIvIJwZ7nfFXQ-LqgQvb2nDZrgUteCXLUX5WFV6ri9FokdM3snL6hTZLXcW-n10Sur9RHdEBL6XYX7Os-NeJADuEPhmPAp4lXKY6lzjTfFd8aD0EeLnFLvRWlBhqj7XIPVFX9Ou2PN9-QfEgJWkfjhxenC8hY17yOwsFsHmxkSK-nq3vetQzOSUklyt4glDQoPCkReIK_sRf--q8mTbRg
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTIzOTAwMDUsImlhdCI6MTU1MjM4NjQwNSwiaXNzIjoidGVzdCIsImp0aSI6Ijk4OTkzMWQ4LTBlN2QtNGU3Ni1iYmY2LTNiNDE1MmIyYzg4OCIsInN1YiI6InRlc3QifQ.OESutED2FrqNrrqGLfbssJA9sK3fbhJ1TqIoEY9HXSA
