# service-user-admin

### to do service-user-admin

- Admin req

- [x] Fix config: log and cors

- [x] Setup Middleware Auth user(unix_id)
- [x] Setup Middleware Auth admin(unix_id with os.Getenv("ADMIN_ID"))
- [-] Setup Expire for JWT

- [x] CORS CONFIG

### Admin

- Root admin by env ADMIN_ID
- [x] POST deactive admin
  - [x] with midlleware auth admin
- [x] POST active admin

  - [x] with midlleware auth admin

- [x] ~GET Log service
  - Auth middleware admin, Uri(unix_id on env)
- [x] GET Service status
  - Auth middleware admin, Uri(unix_id on env)
- [x] PUT Update User profile admin
  - Update data can be empty content
  - update with name, email, phone
  - Auth middleware admin, Uri(unix_id on env)
- [x] POST get all user by admin
  - Auth middleware admin, Uri(unix_id on env)
- [x] PUT update password by admin
  - Auth middleware admin, Uri(unix_id on env)

- [x] POST Register
  - [x] POST Check email
  - [x] POST Check phone
- [x] POST Login

- Dashboard (Middleware Auth)
- [x] PUT Update User profile 
  - Update data can be empty content
  - update with name, email, phone
- [x] GET User Profile
  - Auth middleware 
- [x] PUT update password
  - Auth middleware

  - with middleware on token

- [ ] POST Update_avatar

- [x] POST Logout
  - Delete token

### Investor

- [x] GET getAdminID/:unix_id
  - not use middleware auth

### CI/CD Github Actions

- [x] ~CI/CD Github Actions
- [] Push to Docker Hub
- [] Push to GCP registry

### Note

- For update user this is server must use middleware auth admin same env ADMIN_ID
- Admin can update user profile
- Admin can deactive user
- Admin can active user
- Admin can get user profile
- Remember to set env ADMIN_ID same with unix_id on user 

