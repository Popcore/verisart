# Verisart Exercise
A REST API that allows users to create and exchanges certificates.

## Desing Notes
The application is built according to the following principles:
- certificates can be created, edited and exchanged by existing users only.
- no authentication is enforced apart from the above cases.
- transactions are stored in a chronological order in the in memory store.
Not required but nice to have in case we need to retrieve the transaction history of a certificate.
- the application uses email addresses as user identifiers. This is ok emails are guarnteed to be unique, but it has the disadvantage of using the same ids to generate URLs. This should not be allowed in a production enviroment but it is accepted for demo purposes.
- at present transaction can only be accepted. But the system is desgined to allow for rejection too.

## Build and Run the app
After this repositry has been cloned or downloaded `cd` in the `verisart` directory and follow the steps below.
The application can be built and run with and without docker.

### Without Docker
Requirements:
- [Golang](https://golang.org/dl/)
- [Make](https://www.gnu.org/software/make/manual/html_node/Introduction.html)

```
make build
```
Wll generate an executable in the /build folder.
The executable ca be run with
```
./build/verisart
```

If using Makefiles is not an option building and running the app can be achieved with
```
go run main.go
```
This command will start a new server on localhost at port `:9091`.

### With Docker
Requirements:
- [Docker > 17](https://docs.docker.com/v17.12/install/)
- [Make](https://www.gnu.org/software/make/manual/html_node/Introduction.html)

the app can be built and run with a single command
```
make docker_run
```

If using Makefiles is not an option building and running the container can be achieved by
```
docker build --tag verisart . && docker run --rm -d --name verisart_app -p 9091:9091 verisart latest
```

This commands will generate a docker image and start a container available at http://0.0.0.0:9091

## Quick Start
The examples below can be followed to quickly see and test how the app works and what functionalities it exposes.
For simplicity the code samples use `curl` for issuing HTTP requests, but the same result can be achieved using similar tools.

0 - Ensure the application is up and running.

1- The first thing to do in order to consume the API is to generate users
```
# user 1
curl -X POST -d '{"email": "user1@email.com", "name": "joe"}' http://0.0.0.0:9091/users

# user 2
curl -X POST -d '{"email": "user2@email.com", "name": "mary"}' http://0.0.0.0:9091/users
```

Both commands will return the user that was created or an error message explaining what went wrong. If all went well the output should look like
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
Please note the user of the `X-User-Email` header, which is required for authenticating the user. Not including the header will prodice an error.

Both commands should return the certificate that was generated or an error message explaining what went wrong.
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


3 - We can see our users' certificates with
```
curl http://0.0.0.0:9091/users/<the-user-email-address>/certificates
```
where `<the-user-email-address>` should be replaced by either `user1@email.com` or `user2@email.com`,


4 - A new certificate transaction from Joe to Mary can be created with
```
curl -X POST -d '{"email": "user2@email.com"}' http://0.0.0.0:9091/certificates/<certificate-id>/transfers
```

where the `<certificate-id>` should be replaced with one of the certificate Ids that we saw in step 2.


5 - And finally to complete the transaction
```
curl -X PATCH -d '{"email": "user2@email.com", "status": "accepted"}' http://0.0.0.0:9091/certificates/<certificate-id>/transfers
```
Where the `<certificate-id>` should be replaced with one of the certificate Ids that we saw in step 2.

If no error was returned we can verify that the certificate was successfully transfered by repeting step 2. Joe should now have 0 certificates while Mary should have 2.

## API Endpoints
The API expected content type is JSON.

### Creating certificates
certificates can be created by existing users only.
Requsts must include a `X-User-Email` header containing the certificate owner email address.

Method: POST
Endpoint: /certificates

A request payload looks like:
```json
{
  "title": "my new certificate",
  "year": 2018,
  "note": "some notes about my certificate"
}
```

On success the application returns the cetificate that was created.
In case of an error the application will return an error containg the http status code and a message.

### Updating certificates
Existing certificates can be updated by specifying the fields that needs to be modified.
Note that attempting to update a transaction object will result in an error as transaction can only updated via the a certificate transfer.

Method: PATCH
Endpoint: /certificates/<the-certificate-id>

A request payload looks like:
```json
{
  "title": "my new certificate title",
  "year": 2018,
  "note": "new notes about my certificate"
}
```

An example of updating a certificate could look like:
```
curl -X PATCH -d '{"title" : "my new shiny title", "year": 2018, "notes": "new notes" }' http://0.0.0.0:9091/certificates/<the-certificate-id>

```

Certificate IDs and their time of creation cannot be modified directly.

Attempting to directly update the certificate ownerID or a transaction status will produce an error. Certificates ownership can only be updated using transactions.

### Deleting certificates
Existing certificates can be also removed. Once deleted a certificate cannot be recovered.

Method: DELETE
Endpoint: /certificates/<the-certificate-id>

An example of updating a certificate could look like:
```
curl -X DELETE http://0.0.0.0:9091/certificates/<the-certificate-id>

```

### Creating new users

Method: POST
Endpoint: /users

A request payload looks like:
```json
{
  "email": "test@email.com",
  "name": "test-user"
}
```

On success the application returns the user that was created.
In case of an error the application will return an error containing the http status code and a message.


### Listing certificates for a user
Certificates can be retrieved by specifying the owner ID in the URL.

Method: GET
Endpoint: /users/<userId>/certificates

An example of updating a certificate could look like:
```
curl http://0.0.0.0:9091/users/<userId>/certificates

```

The application will reponnd with a JSON array containing the certificates that belong to a user.


### Creating a new transaction

Method: POST
Endpoint: /certificates/:id/transfers

A request payload looks like:
```json
{
  "email": "user@email.com",
  "status": "pending"
}
```

Currently only the email address can be specified as the application will automatically set the transaction status to "pending".

The application will reponnd with a JSON object containing the certificates that belong to a user.
Errors will be returned when trying the create a new trasaction for a certificate that already has a pending transaction.


### Accepting a transaction
Certificate ownership can be updated only after a transaction has been accepted.

Method: PATCH
Endpoint: /users/<userId>/certificates

A request payload looks like:
```json
{
  "email": "user@email.com",
  "status": "accepted"
}
```

The application will reponnd with a JSON array containing the updated certificate.


## Test the app
The application codebase can be tested with
```
make test
```

The above command will also generate code coverage, accessible as an HTML file in the /artefacts folder.

## TODO/Nice to have
- ensure reading and writing to the store is thread safe
- user authentication
- better error handling
- CI for automated builds
- logging and monitoring
- A cli tool for allowing for runtime configuration
