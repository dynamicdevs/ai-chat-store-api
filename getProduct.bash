#!/bin/bash

# Define the number of iterations
n=200

# Create a CSV file with headers
echo "name,sku,html_table_attributes_pdp,link,price" >> output.csv

process_id() {
    # Call API and extract the name, sku, and html_table_attributes_pdp
    api_response=$(curl -s --location "https://spdp.abcdin.cl/prd/product?id=$((1174465+$1))")

    name=$(echo "$api_response" | jq -r '.name')
    sku=$(echo "$api_response" | jq -r '.sku')
    html_table=$(echo "$api_response" | jq -r '.custom_attributes[] | select(.attribute_code=="html_table_atrributes_pdp") | .value')
    link=$(echo "$api_response" | jq -r '.custom_attributes[] | select(.attribute_code=="url_key") | .value')
    price=$(echo "$api_response" | jq -r '.price')

    # Check if html_table, link, and price are not empty
    if [ -n "$html_table" ] && [ -n "$link" ] && [ -n "$price" ]; then
        # Parse HTML and extract key-value pairs
        table_content=$(echo "$html_table" | pup 'tr json{}' | jq -r '.[] | .children | [.[0].text, .[1].text] | @tsv' | sed 's/\t/:/g')

        # Replace double quotes with single quotes
        table_content=${table_content//\"/\'}
        name=${name//\"/\'}
        sku=${sku//\"/\'}
        link=${link//\"/\'}
        price=${price//\"/\'}

        # Append data to the CSV file
        echo "\"$name\",\"$sku\",\"$table_content\",\"$link\",\"$price\"" >> output.csv
    else
        echo "Error: Unable to extract html_table_attributes_pdp, url_key, or price for ID=$((1175975+$1))"
    fi
}

# Loop to call the API for different IDs
for ((i=1; i<=n; i++)); do
    # Run the process_id function in the background
    process_id "$i" &
done

# Wait for all background jobs to complete
wait

echo "All data processed."
