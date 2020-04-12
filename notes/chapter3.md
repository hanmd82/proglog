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
  - Record: the data stored in the log.
  - Store: the file where records are stored in.
  - Index: the file where index entries are stored in.
  - Segment: the abstraction that ties a store and an index together.
  - Log: the abstraction that ties all the segments together.

- Create a `store` struct, which is a simple wrapper around a file with two APIs to append and read bytes to and from the file.
- Write to the buffered writer instead of directly to the file to reduce the number of system calls and improve performance.
- Run tests:
    ```bash
    cd internal/log
    go test -c && ./log.test
    ```
