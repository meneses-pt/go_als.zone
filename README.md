# go_als.zone

This is an API developed in GO, that acts as a read interface for the [goals.zone](https://goals.zone/) project

## Build

```bash
go build
```

## Run

Make sure you have an `.env` file with the following variables:
- `DB_HOST` - Database host address
- `DB_PORT` - Database port
- `DB_NAME` - Database name
- `DB_USER` - Database username
- `DB_PASSWORD` - Database password
- `HTTP_ADDR` - Address where to run the application
- `MEDIA_ROOT` - Media root URL for image paths in the DB
- `REDDIT_ROOT` - Reddit base URL~~~~

Run the executable created on the build stage.

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
[MIT](https://choosealicense.com/licenses/mit/)