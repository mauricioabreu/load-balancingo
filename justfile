requests := env_var_or_default("REQUESTS", "10000")
concurrency := env_var_or_default("CONCURRENCY", "2")
rate_limit := env_var_or_default("RATE_LIMIT", "50")

base_url := "http://localhost:8080"


# Load testing the /rr endpoint using the hey tool
rr:
    @echo "Sending requests to /rr endpoint..."
    hey -n {{requests}} -c {{concurrency}} -q {{rate_limit}} "{{base_url}}/rr"

# Load testing the /p2c endpoint using the hey tool
p2c:
    @echo "Sending requests to /p2c endpoint..."
    hey -n {{requests}} -c {{concurrency}} -q {{rate_limit}} "{{base_url}}/p2c"
