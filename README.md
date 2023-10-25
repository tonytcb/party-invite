### Objective

Your assignment is to implement a service that parses as a list of customers and returns their names based on location. Use Go and no framework.

### Brief

We have some customer records in a text file `./Data/customers.txt` -- one customer per line, JSON lines formatted. We want to invite any customer within 100km of our office for some food and drinks on us.

### Tasks

Write a web service that:

-   Has an endpoint that accepts a .txt file containting customers. See a sample file in `./Data`
    -   Read the full list of customers
    -   Output the names and user ids of matching customers (within 100km), sorted by User ID (ascending). The output should be in JSON.
-   No authentication is required
-   You can use the first formula from [this Wikipedia article](https://en.wikipedia.org/wiki/Great-circle_distance) to calculate distance. Don't forget, you'll need to convert degrees to radians. The GPS coordinates for our Dublin office are 53.339428, -6.257664. You can find the Customer list in `./Data`.

### Evaluation Criteria

-   **Go** best practices
-   We're looking for you to produce working code, with enough room to demonstrate how to structure components in a small program.
-   Poor answers will be in the form of one big function. It’s impossible to test anything smaller than the entire operation of the program, including reading from the input file. Errors are caught and ignored.
-   Good answers are well composed. Calculating distances and reading from a file are separate concerns. Classes or functions have clearly defined responsibilities. Test cases cover likely problems with input data.
-   It’s an excellent answer if we've learned something from reading the code.

### CodeSubmit

Please organize, design, test and document your code as if it were going into production - then push your changes to the master branch. After you have pushed your code, you may submit the assignment on the assignment page.

All the best and happy coding,

The sFOX Team