#!/bin/bash

# Define the number of iterations
n=200

# Create a CSV file with headers
echo "name,sku,html_table_attributes_pdp" > output.csv

# Loop to call the API for different IDs
for ((i=1; i<=n; i++)); do
    # Call API and extract the name, sku, and html_table_attributes_pdp
    api_response=$(curl -s --location "https://spdp.abcdin.cl/prd/product?id=$((1175975+i))")

    name=$(echo "$api_response" | jq -r '.name')
    sku=$(echo "$api_response" | jq -r '.sku')
    html_table=$(echo "$api_response" | jq -r '.custom_attributes[] | select(.attribute_code=="html_table_atrributes_pdp") | .value')

    # Check if html_table is not empty
    if [ -n "$html_table" ]; then
        # Parse HTML and extract key-value pairs
        table_content=$(echo "$html_table" | pup 'tr json{}' | jq -r '.[] | .children | [.[0].text, .[1].text] | @tsv' | sed 's/\t/:/g')
        
        # Replace double quotes with single quotes
        table_content=${table_content//\"/\'}
        name=${name//\"/\'}
        sku=${sku//\"/\'}
        
        # Append data to the CSV file
        echo "\"$name\",\"$sku\",\"$table_content\"" >> output.csv
    else
        echo "Error: Unable to extract html_table_attributes_pdp for ID=$((1152637+i))"
    fi
done
