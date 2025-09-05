# Security design

## Picture of current docker setup

## Picture of planned portal setup

- locally accessible only (both web and db)
- db password protected, disk encryption as rest of the server
    - password encrypted in frontend config
- frontend password protected - local account or AD
- action logging
    - rotated logs - limited size
    - redacted sensitive data
    - auditability
    - structured - json
- session based access - limited duration (1h)