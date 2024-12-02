import requests

headers = {
    "apikey": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InRmc3R6bWR5enRka2hkc3FiaXBwIiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImlhdCI6MTczMzA0MDMwMiwiZXhwIjoyMDQ4NjE2MzAyfQ.Ww5AjZ-ZkGL5QxxPUDVcq8WzyAGOcEhhyY84glTlXqI",
    "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InRmc3R6bWR5enRka2hkc3FiaXBwIiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImlhdCI6MTczMzA0MDMwMiwiZXhwIjoyMDQ4NjE2MzAyfQ.Ww5AjZ-ZkGL5QxxPUDVcq8WzyAGOcEhhyY84glTlXqI"
}

req = requests.get("https://tfstzmdyztdkhdsqbipp.supabase.co/rest/v1/users?select=*", headers=headers)
print(req.json())
