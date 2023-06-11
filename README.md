# service-user-admin


### to do service-user-admin


- Admin req

- [x] Fix config: log and cors

- [x] Setup Middleware Auth user(unix_id)
- [x] Setup Middleware Auth admin(unix_id with os.Getenv("ADMIN_ID"))

- [x]CORS CONFIG

- [ ] POST deactive admin
    - [ ] with midlleware auth admin
- [ ] POST active admin
    - [ ] with midlleware auth admin

- [x] ~GET Log service
    - Auth middleware admin, Uri(unix_id on env), 
- [x] GET Service status
- [x] POST Register
    - [x] POST Check email
    - [x] POST Check phone
- [x] POST Login

- Dashboard (Middleware Auth)
- [x] PUT Update User profile admin
    - Update data can be empty content
    - update with name, email, phone
- [x] GET User Profile
    - with middleware on token

- [ ] POST Update_avatar


- [x] POST Logout
    - Delete token 


### CI/CD Github Actions

- [x] ~CI/CD Github Actions
- [] Push to Docker Hub
- [] Push to GCP registry


