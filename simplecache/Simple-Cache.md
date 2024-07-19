
Making CLI interactive cache:

added set , get and delete commands 


The issue where the key is set but not found later can be due to several reasons, including scope and lifecycle of the cache instance. When you run ./main set ..., it creates a new instance of the cache, sets the key, and then exits. When you run ./main get ..., it creates another new instance of the cache, which starts empty, so it cannot find the key.

To persist the cache between runs, you need to store the cache data externally, such as in a file or a database. For simplicity, let's use a file-based approach to serialize and deserialize the cache.