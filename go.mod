module gatorcli

go 1.23.2

require internal/config v1.0.0

require (
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
	internal/rss v1.0.0
)

replace internal/config => ./internal/config

replace internal/rss => ./internal/rss
