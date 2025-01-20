/*
Package raildata is a library to access NJ Transit's RailData API.

# API token management

The RailData API requires you to create a token using your username and password, and then use
that token for all operations. There is a limit to the number of tokens you can create per day,
so it is essential to manage the API token properly to avoid spurious token creations.

This library takes care of token management for you. It can receive a token to use throughout
the session, or it can create one by itself. It can also create a new token automatically when
the old token expires. When it gets a new token, it will call a function you provide so you can
save the token for later.

# Enriched API

Some RailData API methods return station names, while others return short names or station codes or
even "destination" names that don't match any existing station. The same happens with other types
of information, like line codes or names.

This library normalizes and enriches the API results so that, whenever you have a station, you get
its code and its name and its short name; whenever you have a line, you get its code and its name
and its official abbreviation.

Similarly, dates and times are represented as [time.Time] objects, delays and dwell intervals are
represented as [time.Duration], true/false and yes/no values are represented as booleans,
we have a special type for colors, and optional values are represented as pointers.

# Rate-limited functions

Some RailData API methods can only be called 5 or 10 times per day. This library splits them out
to a separate interface that you can get by calling the [Client.RateLimitedMethods] method. This makes
it clear to you, the programmer, that you should try to avoid calling those methods too often.

# NJ Transit developer credentials

In order to use this library, you need to visit https://developer.njtransit.com/registration/login
to request your NJ Transit developer API credentials for the RailData API.
*/
package raildata
