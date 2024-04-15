options

{number | string} if its the same type of the primary key it will be treated as that

--crawl -c crawls to all the foreign keys that the program can find

-d depth of the crawl

tracer ad category 1

ad -> products -> categories

ad | products | categories

1 | 23 | 94

---

tracer ad category_specification 1 --full

ad -> products -> categories -> category_specification

just selects everything

---

tracer ad 1

ad

1

tracer ad 1 --full

id name type ...
