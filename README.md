# NYT_best_sellers
Initial work with Go, making GET calls, parsin jsons and visualizing the output


The main.go file makes a GET to NYT for a specific list. More information on the possible calls here: 

- https://developer.nytimes.com/docs/books-product/1/routes/lists.json/get

You will need to create a user, it is free to do so and straight forward.


Once you get your API key, you can modify the URL path, structure and the for loop prints needed for each field.

This output is then modeled as a table and finally will fill in the file output.html

In order to create a correct output.html we need to create a layout.html

The layout will depend on the struct you create for the response and how you parse the json (if you iterate every book, etc.)

The final output.html can be opened with any browser and will provide a tabular view of the output.
