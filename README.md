# Verisart Exercise
A REST API that allows users to create and exchanges certificates.

## design
The application is built according to the following principles:
- certificates can be created, edited and exchanged by existing users only.
- no authentication is enforced apart from the above cases.
- transactions are stored in a chronological order in the in memory store.
Not required but nice to have in case we need to retrieve the transaction history of a certificate.
- The application uses email addresses as user identifiers. This is ok as we know that email addresses are unique. But it's not great that we are also using the same is to generate URLs.
This should not be allowed in a production enviroment but it is accepted for demo purposes.
- at present transaction can only be accepted. But the system is desgined to allow for rejection too.

## Build and Run the app
If using Makefiles is not an option building and running the app be achieved by

### Without Docker
Requirements:
Golang
Make

the app can be build the app with
```
make build
```

and run
```
./build/verisart
```

If using Makefiles is not an option building and running the app be achieved by
```
go run main.go
```

The above will start a new server port 9091.

### With Docker
Requirements:
Docker > 17
Make

the app can be build and run with a single command
```
make docker_run
```

If using Makefiles is not an option building and running the container can be achieved by
```
docker build --tag verisart . && docker run --rm -d --name verisart_app -p 9091:9091 verisart latest
```

The above will generate a docker image and container available at http://0.0.0.0:9091

## Quick Start
For simplicity the example above use `curl`. The same can be achieved using similar tools.

1- The first thing to do in order to consume the api is to generate users.
```
# user 1
curl -X POST -d '{"email": "user1@email.com", "name": "joe"}' http://0.0.0.0:9091/users

# user 2
curl -X POST -d '{"email": "user2@email.com", "name": "mary"}' http://0.0.0.0:9091/users
```

Both example will return the user that was created or an error message. If all went well the output should look like
```json
{
  "id": "2b8ed671-8f1e-4246-83c3-c7b61425b291",
  "email": "user2@email.com",
  "name": "mary"
}
```

2 - Now let's create one certificate for each user.
One for Joe
```
curl -H "X-User-Email: user1@email.com" -X POST -d '{"title": "cert1", "year": 1998, "note": "some notes"}' http://0.0.0.0:9091/certificates
```

and one for Mary
```
curl -H "X-User-Email: user2@email.com" -X POST -d '{"title": "cert2", "year": 2018, "note": "some other notes"}' http://0.0.0.0:9091/certificates
```

Both commands should return the certificate that was generated or an error.
If all went well the output should look something like
```json
{
  "id": "7b96e24c-330f-4629-b736-d780432d9cf3",
  "title": "cert1",
  "createdAt": "2018-11-22T12:21:38.5902426Z",
  "ownerId": "user1@email.com",
  "year": 1998,
  "note": "some notes",
  "transfer": null
}
```

Please note the user of the `X-User-Email` header, which is required for authenticating the user. Not including the header will prodice an error.

3 - We can see ou users' certificates with
```
curl http://0.0.0.0:9091/users/<the-user-email-address>/certificates
```
where `<the-user-email-address>` should be replaced by either `user1@email.com` or `user2@email.com`,

4 - A certificate transaction from Joe to Mary can be created with
```
curl -X POST -d '{"email": "user2@email.com"}' http://0.0.0.0:9091/certificates/<certificate-id>/transfers
```

Where the `<certificate-id>` should be replaced with one of the certificate Ids that we saw in step 2.

5 - And finally to complete the transaction
```
curl -X PATCH -d '{"email": "user2@email.com", "status": "accepted"}' http://0.0.0.0:9091/certificates/<certificate-id>/transfers
```
Where the `<certificate-id>` should be replaced with one of the certificate Ids that we saw in step 2.

If no error was returned we can verify that the certificate was successfully transfered by repeting step 2. Joe should now have 0 certificates while Mary should have 2.

### Updating certificates

Existing certificates can be updated by specifying the fields that require updating.
Note that attempting to update a transaction object will result in an error as transaction can only updated via the a certificate transfer.

Method: PATCH
Endpoint: certificates/<the-certificate-id>

An example of updating a certificate could look like:
```
curl -X PATCH -d '{"title" : "my new shiny title", "year": 2018, "notes": "new notes" }' http://0.0.0.0:9091/certificates/<the-certificate-id>

```

### Deleting certificates
Existing certificates can be also deleted. Once deleted certificates cannot be recovered.

Method: DELETE
Endpoint: certificates/<the-certificate-id>

An example of updating a certificate could look like:
```
curl -X DELETE http://0.0.0.0:9091/certificates/<the-certificate-id>

```

## Test the app
To run unit testa and thir code coverage run
```
make test
```

Running the above command will also generate code coverage, accessible as an HTML file in the /artefacts folder.

## TODO/Nice to have
- user authentication
- better error handling
- CI for automated builds
- logging and monitoring
- A cli tool for allowing for runtime configuration