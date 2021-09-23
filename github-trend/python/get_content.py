
import requests

GITHUB_TREANDING = "https://github.com/trending"

def get_github_trending() -> str:
    r = requests.get(GITHUB_TREANDING)
    if r.status_code==200:
        return r.text
    
def get_github_trending_local() -> str:
    with open("./test.html") as f:
        return f.read()