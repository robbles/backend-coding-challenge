# Backend Coding Challenge

From [coveo/backend-coding-challenge](https://github.com/coveo/backend-coding-challenge)

Goals:

- Implement a basic solution in Python or Go
- Clean it up and post on Github
- Write a short blog post about it
    * Note: make sure it doesn't link directly to the challenge in any way to avoid spoiling it for companies that are currently using it for recruiting

## Ideas

- Use an in-memory structure and load the entire ~1MB file on startup

- Internal data structure could be a radix tree:
    * This gives the partial matching trivially
    * A limited breadth-first search could then be used when not at a leaf node to construct a list of matches
    * Depth first would give something like alphabetical order, whereas breadth first should result in "shortest first" order
    * Could each branch node store a "best child" or list of best children? That could offload the cost of building the partial list of completions to start up.
    * If there's only a few matches (e.g. longer query), iterating through the results and applying the Haversine formula for dstance would work
    * But what if there are lots of matches? Can't restrict up front without cutting out results that could be closer. Is there a time-optimized structure that could fix this?

- The latitude/longitude are "optional ... to help improve relative scores", so the system needs to rely primarily on matching by name, then by location. I assume that not matching a city just because it's too far away is unacceptable, so dividing results by geo-boxing first will likely not work.

- Am I supposed to pick a limit for the result set? Or is there no limit? I assume arbitrarily using 10 to limit any O(N) algorithm is probably fair.

- The obvious way to evaluate the correctness of the implementation is to compare against the "dumb" implementation: Manually checking every entry for a partial match, scoring every partial match by distance, then sorting and limiting the result set. Ideally a more efficient algorithm could get identical results without having to scan through the entire dataset.

- For scoring by distance efficiently, an R-tree might be a good approach. The query could go down to the most specific bounding box, then back up the chain of parent nodes as needed to get more results. e.g:
    * R-tree query hits a leaf node with 2 results
    * yield those results first, then move to parent node with 4 results
    * yield those results then move to parent node with 100 results
    * yield the closest 4, then return
- This would have an edge case of locations that are not close to many cities: the result set that was iterated over to produce enough results would be huge. e.g. in the middle of the Pacific Ocean.
- Could also produce subtly different results than the base implementation because of results on the border of a bounding box being included/excluded by the tree. Might be wrong about this, need to read description of R-tree balancing more closely.

- Starting point: implement just the radix tree approach, with scoring calculated as 1.0 - (len / max len)

## Example Cases

- query: "a", no lat/lng
    * 138 results match
    * yield the 10 shortest results ordered by length ascending
    * Ada, Adel, Ajax, Alma, Alma, Ames, Amos, Anna, Apex, Avon

- query: "van", no lat/lng
    * 5 results match
    * yield all results ordered by length ascending
    * Vandalia, Vandalia, Van Buren, Vancleave, Vancouver

- query: "van", lat/lng: 49.1886, -122.9384
    * 5 results match
    * yield all results ordered by length ascending
    * Vancouver, ...

