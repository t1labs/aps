# APS

[![CircleCI](https://circleci.com/gh/t1labs/aps.svg?style=svg)](https://circleci.com/gh/t1labs/aps)

This attempt at an artificial pancreas system attempts to keep everything as simple as possible in a highly complex environment.

> This application is in no way complete and provides no promises as to up-time, accuracy, or reliability. **You should not use the data from this application for treatment decisions.**

## Example

To get started all you need is Docker.

```
docker run -e DEXCOM_SHARE_USERNAME=<your-username> -e DEXCOM_SHARE_PASSWORD=<your-password> t1labs/aps
```

> This should start a Docker container, and output your blood sugar level every minute.

## Contributing

Contributions are welcome and most definitely encouraged!

If you want to add support for a new CGM or pump, please start by opening an issue describing the type of support you wish to make, and a general plan for how it should work.

## Testing

To run all of the tests:

```
$ make test
```
