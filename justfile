requests := env_var_or_default("REQUESTS", "1000")

concurrency := env_var_or_default("CONCURRENCY", "50")

base_url := "http://localhost:8080"

rr:
    @echo "Sending requests to /rr endpoint..."
    hey -n {{requests}} -c {{concurrency}} "{{base_url}}/rr"

p2c:
    @echo "Sending requests to /p2c endpoint..."
    hey -n {{requests}} -c {{concurrency}} "{{base_url}}/p2c"

all: rr p2c
