# Auth
## Login/Signup
1. User sends request to /auth/login/google, redirect_uri /auth/callback/google
2. User redirected to /auth/callback/google after login, with authorization code in query params
3. Handler extracts authorization code, calls Google endpoint for user info
4. Populate database if needed, retrieve user info from database
5. Create session token using crypto, hash it using SHA256 to generate session ID
6. Store session ID in database, create cookie using session token

## ValidateToken (called when client gets user information)
1. Handler gets session token from cookie
2. Hash it to get session ID
3. Lookup session ID in database to get user id
4. If session expiry time is less than 15 days, then create a new session token + session id with new expiration of 30 days
5. Set cookies
6. Invalidate previous session ID



