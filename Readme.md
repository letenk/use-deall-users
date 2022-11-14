# Credential
**Admin:**
```
usename: admin
password: password
```

**User:**
```
usename: user
password: password
```

# Documentation
[Postman Documentation](https://documenter.getpostman.com/view/12132212/2s8YevnUpD)

[Swagger](https://app.swaggerhub.com/apis-docs/DARMAWANRIZKY43/use_deall_rest_api_users/1.0.0#/)


# Architecture Diagram
## Local Development
![architecture diagram local development](/assets/use-deall-architecture-diagram-local-development.png)

# Todo
- [x] Documentation
    - [x] Swagger
    - [x] Postman
- [x] Architecture diagram flow CRUD and Login
- [x] Attach credential in the Readme.md
- [ ] CRUD
    - [x] Create (admin only)
    - [x] Get all
    - [x] Get one (admin only)
    - [x] Update (admin only)
    - [x] Delete (admin only)
    - [x] Login
        - [x] Generate Token JWT
        - [ ] Refresh Token
- [x] Role `user` only can access `read` (Endpoint `Get all`)
- [x] Containerization
- [ ] Kubernetes 
- [x] Upload into Github
