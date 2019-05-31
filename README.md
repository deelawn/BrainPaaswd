# BrainPaaswd

### Endpoints
Refer to the Cloud_Services_Programming_Challenge.pdf document

### How to run it
If you have Docker:
- Run `docker build -t paaswd .` to build the container
- Run `docker run -p 8000:8000 -v $(pwd):/go/src/github.com/deelawn/BrainPaaswd --name paaswd_c --rm paaswd` 
to start the service
- It runs on port 8000 and uses the files in the .data directory by default

If you do not have Docker Option 1:
- Install Docker
- See above

If you do not have Docker Option 2:
- Download and install Go 1.12
- Run `go run main.go --help` to see possible arguments and default values
- Run again without `--help` to start the service


### How does it work?
- Exit if a provided resource path can't be read
- Read resources into a cache on startup
- For each request, retrieve the result from the cache and return the results or read from the file, 
update the cache, and return the results. This depends on if the cache is stale.

### Things that could be improved
These will be long form explanations as opposed to bullet points. These are all things that I would implemnt if I had 
the time to do so and, in some cases, run benchmarks to determine if performance gains warrant merging it to master.

#### Caching
I chose to write my own caching mechanism because I though it would be fun. One thing I did not address is the issue of 
timing out when trying to get a lock on the cache. It should create a context with a timeout and move on without reading
from the cache if it can't establish a lock within the time frame. There should also be a request based cache that
stores query results in the same way that UID and GIDs are indexed in the cache. Something else that could be done is to
create a new type of cache that implements storage.Cache; this new type of cache would be an in-memory sqlite instance.
This may increase performance and eliminate the need for caching by IDs if the ID columns on the tables are indexed.

#### Query Parsing
The current method of parsing query parameters is admittedly inefficient. Efficiency could be increased by implementing
a sqlite cache as expressed in the Caching section above. If it were to keep the existing caching mechanism, I would 
implement a generic query parser similar to something I created at my current company; it allows you to define a schema
for each service where each part of the schema defines a query parameter with properties such as name, type, accepts
multiple values, etc. Then reflection could be used to map model fields to query parameters and filtering could be done
over a single loop.  In fact, both services could use the same filtering logic. Additionally, it might be cool to
implement a query parsing pipeline using channels in order to concurrently filter, though the effort to achieve this
might eclipse the minimal performance gains from filtering on such a small amount of data.

#### File Sizes
If this services needed to be able to scale to handle massive amounts of data, then some changes would need to be made,
such as the implemention of a cache that stores selected resources while the majority of things are stored on disk; a
SQL database would be a solution to this problem.

#### Error Handling
Errors should be defined as format string constants somewhere so that the source value and error value can be
substituted when the error response is written to the HTTP writer.

#### Testing
Pretty decent test coverage was achieved just by executing the endpoint handler functions, but it could be better. Some
of the other packages should have tests written to test all public functions. For the services, tests should be written
to check that the error execution paths are behaving correctly. The `/users/{id}/groups` endpoint isn't tested
explicitly but is mostly covered by the `/groups/query?member={username}` endpoint tests. If I wanted to test the former
mentioned endpoint I would need to either spin up a server while running my tests or inject a mocked HTTP client into
the users service. Would have also been nice to write a few tests that change a files contents to ensure that the
file gets read again and cache updated on the next request.

#### CI
CI would have been cool to include where the tests get run every time I push to Github.

#### Badges
Badges are cool to add to a README. It would probably be a requirement if I had set up CI.

### Other notes
- `.scripts/coverage.sh` will run the tests and generate the code coverage report
- The `.coverage` directory conveniently contains the code coverage report in HTML format
- The `.data` directory contains sample `passwd` and `group` files
- The name of this repository is not a spelling mistake