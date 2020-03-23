# mqtt-mail-notifier

## Description
//TBD

## Try it out 
1. start environment
    ```shell script
    docker-compose -f example/docker-compose.yml up
    ```
1. wait a few seconds until all components has been started
1. publish mqtt message
    1. open locahost:18083
    1. login with username=admin password=public
    1. go to tools>websocket
    1. press connect
    1. fill message section with
        1. Topic= notify/mail/4711/correpationID
        2. Message= {"user-id":"fl1"}
    1. Press send
1. verify incoming mail 
    1. open localhost:8025
    2. enjoy :)

