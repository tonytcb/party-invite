### Objective

Your assignment is to implement a service that parses as a list of customers and returns their names based on location. Use Go and no framework.

### Brief

We have some customer records in a text file `./Data/customers.txt` -- one customer per line, JSON lines formatted. We want to invite any customer within 100km of our office for some food and drinks on us.

### Tasks

Write a web service that:

-   Has an endpoint that accepts a .txt file containing customers. See a sample file in `./Data`
    -   Read the full list of customers
    -   Output the names and user ids of matching customers (within 100km), sorted by User ID (ascending). The output should be in JSON.
-   No authentication is required
-   You can use the first formula from [this Wikipedia article](https://en.wikipedia.org/wiki/Great-circle_distance) to calculate distance. Don't forget, you'll need to convert degrees to radians. The GPS coordinates for our Dublin office are 53.339428, -6.257664. You can find the Customer list in `./Data`.

### CodeSubmit

Please organize, design, test and document your code as if it were going into production - then push your changes to the master branch. After you have pushed your code, you may submit the assignment on the assignment page.

All the best and happy coding,
