import requests
import time
import json
import random
from datetime import datetime, timezone

def generate_random_data():
    timestamp = datetime.utcnow().replace(tzinfo=timezone.utc).isoformat()
    temperature = round(random.uniform(20, 30), 1)
    dust = random.randint(50, 150)
    humidity = random.randint(60, 90)

    data = {
        "timestamp": timestamp,
        "temp": temperature,
        "dust": dust,
        "humidity": humidity
    }

    return data

def post_data_to_api(data, url):
    headers = {'Content-Type': 'application/json'}
    response = requests.post(url, data=json.dumps(data), headers=headers)

    if response.status_code == 200:
        print(f"Data posted successfully: {data}")
    else:
        print(f"Failed to post data. Status code: {response.status_code}, Response: {response.text}")

if __name__ == "__main__":
    api_url = "http://localhost:8080/api/sensor"

    while True:
        random_data = generate_random_data()
        post_data_to_api(random_data, api_url)

        time.sleep(random.randint(5, 15))

