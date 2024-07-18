# Caching

# Problem Statement :

## Simple Cache Implementation: 
detailed : (https://chatgpt.com/share/c265f813-f4ea-428e-bbb7-5986978ed104)

Implement a basic in-memory cache in Go using a map and implement set, get, and delete operations.
Add a simple expiration mechanism using a goroutine that periodically cleans up expired items.

### What is Caching?
Caching is a technique use to temporarily save frequently accessed data in high-speed memory like RAM to allow faster access.
This helps us to reduce the time and resource needed to retrieve data from slower medium or compute is again and helps in improving overall efficiency and performance of application.

Race Condition: [https://www.notion.so/Race-Condition-6c0f7cbd8d8c4d98b421e0772495312a](https://olivine-asparagus-094.notion.site/Race-Condition-6c0f7cbd8d8c4d98b421e0772495312a)

### How improper Caching implementation can lead to Race Condition?
As in current problem statement if we implement set , get and delete operations which are running concurrently and directly manipulating memory data (atleast having a write access)[Data Race], could lead to race-condition if not done properly . Race condition here means it can lead to data inconsistency , unpredictable behaviour , corrupt data etc.

Potential Race Conditions
1. Concurrent Writes and Reads:
    If one goroutine is writing to c.items while another is reading from it, the reader might see a partially written CacheItem or corrupt data.
2. Concurrent Writes:
    If multiple goroutines write to c.items simultaneously, the map's internal structure can be corrupted, leading to unpredictable behavior, crashes, or panics.
3. Concurrent Writes and Deletes:
    If one goroutine deletes an item while another is writing to it, the map might end up in an inconsistent state, and subsequent operations could fail or behave unexpectedly.

Refer : race-condition.go 
run :
```bash 
 go run -race race-condition.go
```
here race flag is used to seek any race condition appear in the code.
So here the list of all data race will appear like below example.

```vbnet
==================
WARNING: DATA RACE
Read at 0x00c0000140a0 by goroutine 6:
  main.(*Cache).Get()
      /path/to/main.go:29 +0x40
  main.main.func2()
      /path/to/main.go:60 +0x3a

Previous write at 0x00c0000140a0 by goroutine 5:
  main.(*Cache).Set()
      /path/to/main.go:23 +0x8b
  main.main.func1()
      /path/to/main.go:53 +0x3a
...
==================
Found 1 data race(s)
```
hence to avoid any race condition for the concurrent running threads , we use mutual exclusions using mutext locks to run thereads synchonouly and in ordered fashion.
Explanation
Mutex (mu sync.RWMutex): Ensures that only one goroutine can write to the cache at a time (Lock()), while allowing multiple concurrent reads (RLock()).
Locks and Unlocks: Protect access to the items map, ensuring that read and write operations do not interfere with each other.
By using the mutex, you prevent race conditions, ensuring the cache remains consistent and preventing data corruption or crashes.

### Why caching implementation should be thread-safe?
To prevent racce condition.

```bash
go run simple-cache.go 
```

### How to manage stale cache?
We are using [func cleanup()] go routine which will run periodically can clean up the map[cache] , means remove expired data. his ensures that the cache does not grow indefinitely with stale data, which can consume unnecessary memory and potentially degrade performance over time. 
Why a Cleanup Job is Needed
Memory Management:

Over time, without cleanup, the cache could accumulate a large number of expired items. This could lead to increased memory usage, potentially leading to memory exhaustion and application crashes.
Performance:

Keeping stale data in the cache could slow down cache operations such as searching for valid items. Regular cleanup helps maintain optimal performance by ensuring that the cache only contains relevant, non-expired data.
Resource Efficiency:

Cleaning up expired items helps in reclaiming memory resources that can be used for other parts of the application, ensuring efficient use of system resources.

### Why are we using channel done in the code ?
Channel are use to make sure main func is not exited before execution of set , get and delete thread . After complete execution of these threads only man funtion get exit.
In the provided code, the done channel is used to synchronize the completion of multiple goroutines with the main goroutine. This ensures that the main function waits for all spawned goroutines to complete their work before exiting.

```go
func main() {
    cache := NewCache()
    done := make(chan bool)

    // Goroutine 1: Set key1
    go func() {
        for i := 0; i < 1000; i++ {
            cache.Set("key1", "value1", 10*time.Second)
        }
        done <- true
    }()

    // Goroutine 2: Get key1
    go func() {
        for i := 0; i < 1000; i++ {
            cache.Get("key1")
        }
        done <- true
    }()

    // Goroutine 3: Delete key1
    go func() {
        for i := 0; i < 1000; i++ {
            cache.Delete("key1")
        }
        done <- true
    }()

    // Wait for all goroutines to complete
    <-done
    <-done
    <-done
}

```
Purpose and Usage of the done Channel
Synchronization:

The done channel is used to signal when each goroutine has completed its work. The main function waits for these signals before it terminates.
Coordination:

Each goroutine writes a true value to the done channel once it finishes executing its loop. This write operation indicates the completion of that particular goroutine's task.
Blocking Until Completion:

The main function blocks on <-done three times, which means it waits for three values to be received from the done channel. This ensures that all three goroutines have completed their work before the main function exits.
Why Use the done Channel?
Ensuring All Goroutines Complete:

Without the done channel, the main function might exit before the goroutines finish their execution. This would cause the program to terminate prematurely, potentially leaving operations incomplete.
Simple Synchronization Mechanism:

The done channel provides a simple and effective way to synchronize the main function with multiple goroutines, ensuring that all tasks are completed before the program exits.
Example Without the done Channel
If you remove the done channel synchronization, the main function might terminate before the goroutines complete, leading to unpredictable results:

```go
func main() {
    cache := NewCache()

    // Goroutine 1: Set key1
    go func() {
        for i := 0; i < 1000; i++ {
            cache.Set("key1", "value1", 10*time.Second)
        }
    }()

    // Goroutine 2: Get key1
    go func() {
        for i := 0; i < 1000; i++ {
            cache.Get("key1")
        }
    }()

    // Goroutine 3: Delete key1
    go func() {
        for i := 0; i < 1000; i++ {
            cache.Delete("key1")
        }
    }()

    // Without synchronization, the main function might exit immediately
    fmt.Println("Main function exiting")
}
```
In this case, the main function might print "Main function exiting" and terminate before the goroutines finish their work, causing incomplete operations.

Conclusion
The done channel is used to ensure that the main function waits for all the goroutines to complete their tasks before exiting. It provides a simple synchronization mechanism that coordinates the completion of concurrent operations in the program.

### Is goroutine if run periodically is said to be a cron job?
Comparison of Cleanup Goroutine and Cron Job
Purpose:

Cleanup Goroutine: Periodically removes expired items from the in-memory cache to free up memory and keep the cache efficient.
Cron Job: A scheduled task in Unix-like operating systems that runs at specified intervals to perform various maintenance tasks, such as backups, log rotation, and data cleanup.
Scheduling:

Cleanup Goroutine: Uses a time.Sleep() call within an infinite loop to run the cleanup task at regular intervals.
Cron Job: Uses the cron daemon with a crontab file to schedule tasks based on a specified time and date pattern.
Execution Environment:

Cleanup Goroutine: Runs within the application process, maintaining the in-memory cache directly.
Cron Job: Runs as a separate process initiated by the cron daemon, often used for system-wide maintenance tasks.
Implementation of Cleanup as a Goroutine
In the provided code, the cleanup process is implemented as a goroutine:

```go
// cleanup periodically removes expired items from the cache
func (c *Cache) cleanup() {
    for {
        time.Sleep(time.Minute) // Sleep for a minute before each cleanup
        c.mu.Lock() // Acquire a write lock to ensure exclusive access
        for key, item := range c.items { // Iterate through all items in the cache
            if time.Now().UnixNano() > item.Expiration { // Check if the item is expired
                delete(c.items, key) // Delete expired item
            }
        }
        c.mu.Unlock() // Release the lock after cleanup
    }
}
```
Similarity to Cron Jobs
Periodic Execution:
Just like a cron job, the cleanup goroutine runs periodically (every minute in this example) to perform its task.
Maintenance Task:
Both the cleanup goroutine and cron jobs are used for maintenance tasks that keep the system in good health.
Example of a Cron Job
To further illustrate the similarity, here's an example of a cron job that might be used to clean up a temporary directory every day at midnight:

Crontab Entry:

```bash
Copy code
0 0 * * * /usr/bin/find /tmp -type f -atime +7 -delete
```
This cron job runs the find command every day at midnight to delete files in /tmp that are older than 7 days.
Explanation:
0 0 * * *: Runs the job at midnight (00:00) every day.
/usr/bin/find /tmp -type f -atime +7 -delete: Finds and deletes files in /tmp that haven't been accessed in the last 7 days.
Conclusion
The cleanup goroutine in the cache implementation serves a similar purpose to a cron job by performing regular maintenance tasks to keep the cache efficient and prevent memory bloat. Both mechanisms ensure periodic execution of necessary tasks to maintain system health. Referring to the cleanup goroutine as a cron job within the context of an application is a valid analogy, emphasizing their shared goal of periodic maintenance.



Hands-on Exercises
Simple Cache Implementation:

Implement a basic in-memory cache in Go using a map and implement set, get, and delete operations.
Add a simple expiration mechanism using a goroutine that periodically cleans up expired items.
Redis Integration:

Set up a Redis server and integrate it into a simple Go application.
Implement caching for an API response using Redis.
Experiment with different data structures in Redis (e.g., strings, hashes, lists, sets).
Cache Invalidation:

Implement a mechanism to invalidate the cache when the underlying data changes. For example, when updating a database record, ensure the corresponding cache entry is invalidated.
Implement a write-through cache where writes are done to both the cache and the database.
Distributed Cache:

Set up a distributed cache using Redis in a cluster mode.
Implement a consistent hashing mechanism to distribute keys across multiple nodes.
Cache-aside Pattern:

Implement the cache-aside pattern where the application code explicitly loads data into the cache.
Use this pattern in a more complex application, such as an e-commerce site where product details are cached.
Write-back and Write-through Caches:

Implement a write-back cache where data is written to the cache first and then asynchronously written to the database.
Compare and contrast it with a write-through cache.
Benchmarking and Performance Testing:

Implement a benchmarking suite to measure the performance of different caching strategies.
Test the impact of cache hits and misses on application performance.
Recommended Resources
Books
"Designing Data-Intensive Applications" by Martin Kleppmann: This book covers various aspects of data storage, including caching mechanisms, and provides in-depth knowledge on building scalable systems.
"High Performance Browser Networking" by Ilya Grigorik: Although focused on networking, this book includes sections on optimizing web performance with caching.
Online Courses
Coursera: Cloud Computing Specialization: Offers courses on caching strategies as part of the cloud computing curriculum.
Udacity: Full Stack Web Developer Nanodegree: Includes caching and database optimization lessons.
Pluralsight: Offers various courses on Redis and caching strategies for developers.
Tutorials and Articles
Redis Official Documentation: Comprehensive guide and tutorials on using Redis for various use cases.
GeeksforGeeks and Medium: Both platforms have numerous articles and tutorials on implementing caching in different programming languages.
RealPython: Offers in-depth tutorials on using Redis and caching in Python applications, which can be adapted to Go.
Open Source Projects and Repositories
Go-Redis GitHub Repository: Explore examples and documentation for using Redis with Go.
Awesome Go: A curated list of Go frameworks, libraries, and software, including caching libraries.
Practical Projects
Build a URL Shortener: Implement a URL shortening service with Redis caching for fast redirects.
Develop a News Aggregator: Cache news articles fetched from various APIs to improve response times.
Create a Real-Time Chat Application: Use Redis for caching user sessions and message history.
Interview Preparation
LeetCode and HackerRank: Solve problems related to caching, such as LRU Cache implementations.
System Design Interviews: Practice designing systems with caching layers, including discussing trade-offs and optimizations.
By working on these exercises and utilizing the resources mentioned, you'll gain a solid understanding of caching mechanisms and be well-prepared for interviews.
