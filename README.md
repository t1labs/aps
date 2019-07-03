# APS

This attempt at an artificial pancreas system attempts to keep everything as simple as possible in a highly complex environment.

## Example

To get started all you need is Docker.

```
docker run -e DEXCOM_SHARE_USERNAME=<your-username> -e DEXCOM_SHARE_PASSWORD=<your-password> t1labs/aps
```

> This should start a docker container, and output your blood sugar level every minute.
