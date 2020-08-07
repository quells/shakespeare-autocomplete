# Shakespeare Autocomplete

Search through the collected works of Shakespeare (or any other large file).

Results are ordered by frequency and limited to the top 25 matches.

Non-word entities (like numbers) are ignored, although edge cases such as roman numerals are not.

## Usage

Build web server:

```
$ go build ./cmd/auto
```

Run web server:

```
$ ./auto shakespeare.txt
```

Query web server:

```
$curl 'http://localhost:5000/autocomplete?term=ha'

have
hath
had
hand
hast
has
ham
hands
hang
half
ha
happy
hard
hate
haste
having
hastings
hark
hair
hamlet
harry
hadst
harm
hail
happiness
```

See the documents in `artifacts` for example results for various queries.

If you run the executable from this directory, you can visit http://localhost:5000 to see a minimal website consuming this API.
