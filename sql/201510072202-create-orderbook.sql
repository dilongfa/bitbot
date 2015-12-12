create table orderbook (
	created_at timestamp,
	exchanger  text,
	pair       text,
	bids       text,
	asks       text,
	duration   int
) DISTRIBUTE BY APPEND (created_at);
