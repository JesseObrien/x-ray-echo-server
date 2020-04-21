# AWS X-Ray Echo Server

### What?

[AWS X-Ray](https://aws.amazon.com/xray/) is a tool for tracking segment information about requests through services. The typical way to run this locally is to run the [X-Ray daemon](https://docs.aws.amazon.com/xray/latest/devguide/xray-daemon.html) provided by Amazon. One flaw here is that the daemon requires _real_ credentials to connect to AWS and expects to send data back even locally.

_This service is for_

- If you've ever wanted to see what the X-Ray segment data looks like locally before it's sent to amazon's servers
- You just want to have a local daemon available that doesn't send data to AWS when you're doing testing or running local services

### Install

`go get -u github.com/jesseobrien/x-ray-echo-server`

### Running

`$ x-ray-echo-server`

### Docker
