# Vote Circle Service

Vote Circle Service.

Manage votes and circles.

- Manage votes
- Manage circles
- Gives the user votes for each circle
- count the votes up/down
- has the ranking list for each circle
- creates/updates circle
- has member list of each circle

## GQL

The gql resolvers as well as the models for them, 
will be created using the `github.com/99designs/gqlgen`
package.

### Resolver

Each package with his own schema has it's own
resolvers.


### Generate GQL resolvers

To generate GQL resolvers, `gqlgen` is used.

To generate the resolvers, against the schemas, run:
```bash
go run github.com/99designs/gqlgen generate
```
## Database

go-migrate is in use:
Download the CLI tool from here [migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

## Testing

To run all the tests from the project root:

```bash
 go test ./...
```

## Development & Production

### Development

For local development, either have a development.yml as
config file or inject the .env file into the current shell.

If you run the go app on your current machine, the development.yml
have to be in the project folder.

If you run locally with the docker image, this is not obligatory.