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
![architecture diagram local development](https://previews.dropbox.com/p/thumb/ABv1gh2Aqm_45ZdnA2LGpnUe3VgWm_vpB53TvVALCxdQqW2AO3XWPJMlR9XxxFXzBZtCr6JQCuneId-VBoCA1ABSnzAQ_Q9QPEdsR86JHrdCwy5fa1Ymqinrub1h1kzgXDHpIySvnpAweq035N8fRBISx3CN1ikJ8C9g-3NK8AejVIs09BgiR4hHiSeq8_eX3IycauBDJbtBXkPguPRALUBYhTxpvwoVM8EkoqQ6x43DPaO03CpQeHAaSMC-XG1J5nYtXm8XpZnTyXZo1EIUP0eNc4hTljBvrPEwL-e3AL86oJFPv5BXWs9h6eIdrNmRs4aiUFf0mSGdkwvrYpBHOXy8cMqjfM7l3uGDKgj3CTtK6rfv8Q02pOmvndxKyrIB_-A/p.png)

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
