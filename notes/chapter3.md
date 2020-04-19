### Chapter 3 Notes

- Logs are at the heart of many distributed services, such as storage engines, message queues, version control, and replication and consensus algorithms.
- A log is an append-only sequence of records, split into a list of segments.
- An `active` segment is the one segment being actively written to.
- Each segment comprises a store file and an index file
  - The store file stores the record data, and is where records are continually appended to.
  - The index file stores an index of each record in the store file - it maps record offsets to their position in the store file.
  - To read a record given its offset, first get the entry from the index file for the record (which gives the position of the record in the store file), then read the record at that position in the store file.
  - The index file requires only two fields — the offset and stored position of the record. It is small enough to be memory-mapped.

- Build the log from bottom up, starting with the store and index files, then the segment, and finally the log.
- Terminology:
  - `Record`: the data stored in the log.
  - `Store`: the file where records are stored in.
  - `Index`: the file where index entries are stored in.
  - `Segment`: the abstraction that ties a store and an index together.
  - `Log`: the abstraction that ties all the segments together.

Running tests:
  ```bash
  cd internal/log
  go test -c && ./log.test
  ```

---

- Create a `store` struct, which is a simple wrapper around a file with two APIs to append and read bytes to and from the file.
- Write to the buffered writer instead of directly to the file to reduce the number of system calls and improve performance.

---

- Create an `index` struct, which is a simple wrapper around a physical file and a memory-mapped file.
- A graceful shutdown occurs when a service finishes its ongoing tasks, performs its processes to ensure there’s no data loss, and prepares for a restart.
- Handle ungraceful shutdowns by performing a sanity check when the service restarts. If there is corrupted data, rebuild the data or replicate the data from an uncorrupted source.

---

- The `segment` wraps the `index` and `store` types to coordinate operations across the two:
- When the log appends a record to the active segment, the segment needs to write the data to its store and add a new entry in the index.
- For reads, the segment needs to look up the entry from the index and then fetch the data from the store.
