# ScaleChat 

This was a fun learning project, the main goal was building stateful services that scale horizontally,
along that, I also experimented with a few technologies I wanted to try out, including ScyllaDB and NATS.


### The Services

The project is split into 3 microservices, user, socket and chat.

##### User

Handles user registeration and authentication, it issues JWT tokens with a secret shared between services.
It would also handle token refresh, but that hasn't been implemented yet.
Since it's a toy project, it would also hold the user's profile data, if that gets added.

##### Socket

The core of the project, it handles the websocket connections to clients, for each unique chat room a user connects with,
the service subscribes to NATS with that chat room's ID. messages recieved from the client are published with the chatroom ID,
messages recieved from NATS are sent to clients connected on the topic. It also handes retrieving and saving messages from and to
ScyllaDB/CassandraDB.

##### Chat

Handles the CRUD operations for chatrooms, in the future it could also publish messages for events like chat room deletion, bans, etc.


### Things to improve

There's a lot to fix before this is production quality, main things to fix:

- [] Implement JWT refresh tokens
- [] Refactor repeated code between services into a shared lib
- [] Fix error handling and logging
- [] Add integration and stress testing



### How to run

This project runs on K8s only for now, you can run it on just docker with a few changes though.

- run `build.sh` to build the docker images
- run the nats and scylla deployments
- connect to the scylla deployment and add the keyspace `scalechatmessages`
- run the other deployments in whatever sequence